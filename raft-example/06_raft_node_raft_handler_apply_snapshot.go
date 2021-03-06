package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/gyuho/db/pkg/fileutil"
	"github.com/gyuho/db/pkg/types"
	"github.com/gyuho/db/raft"
	"github.com/gyuho/db/raft/raftpb"
)

// (etcd etcdserver.apply)
type apply struct {
	entriesToApply []raftpb.Entry  // (etcd etcdserver.apply.entries)
	snapshotToSave raftpb.Snapshot // (etcd etcdserver.apply.snapshot)
	applyDone      chan struct{}   // (etcd etcdserver.apply.raftDone)
}

// (etcd etcdserver.EtcdServer.applySnapshot)
func (rnd *raftNode) applySnapshot(ap *apply) {
	if raftpb.IsEmptySnapshot(ap.snapshotToSave) {
		return
	}

	logger.Infof("applying snapshot at index %d", rnd.snapshotIndex)
	defer logger.Infof("finished applying snapshot at index %d", rnd.snapshotIndex)

	if ap.snapshotToSave.Metadata.Index <= rnd.appliedIndex {
		logger.Panicf("snapshot index [%d] should > progress.appliedIndex [%d] + 1", ap.snapshotToSave.Metadata.Index, rnd.appliedIndex)
	}

	dbFilePath, err := rnd.storage.DBFilePath(ap.snapshotToSave.Metadata.Index)
	if err != nil {
		panic(err)
	}
	fpath := filepath.Join(rnd.snapDir, "db")
	if err = os.Rename(dbFilePath, fpath); err != nil {
		panic(err)
	}

	logger.Infof("loading snapshot at %q", fpath)
	f, err := fileutil.OpenToRead(fpath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := rnd.ds.loadSnapshot(f); err != nil {
		panic(err)
	}

	rnd.configState = ap.snapshotToSave.Metadata.ConfigState
	rnd.snapshotIndex = ap.snapshotToSave.Metadata.Index
	rnd.appliedIndex = ap.snapshotToSave.Metadata.Index
}

// (etcd etcdserver.EtcdServer.applyEntries,apply)
func (rnd *raftNode) applyEntries(ap *apply) {
	if len(ap.entriesToApply) == 0 {
		return
	}

	firstIdx := ap.entriesToApply[0].Index
	if firstIdx > rnd.appliedIndex+1 {
		logger.Panicf("first index of committed entry[%d] should <= progress.appliedIndex[%d] + 1", firstIdx, rnd.appliedIndex)
	}

	var ents []raftpb.Entry
	if rnd.appliedIndex-firstIdx+1 < uint64(len(ap.entriesToApply)) {
		ents = ap.entriesToApply[rnd.appliedIndex-firstIdx+1:]
	}
	if len(ents) == 0 {
		return
	}

	for i := range ents {
		switch ap.entriesToApply[i].Type {
		case raftpb.ENTRY_TYPE_NORMAL:
			if len(ap.entriesToApply[i].Data) == 0 {
				// ignore empty message
				break
			}
			select {
			case rnd.commitc <- ap.entriesToApply[i].Data:
			case <-rnd.stopc:
				return
			}

		case raftpb.ENTRY_TYPE_CONFIG_CHANGE:
			var cc raftpb.ConfigChange
			cc.Unmarshal(ap.entriesToApply[i].Data)
			rnd.configState = *rnd.node.ApplyConfigChange(cc)
			switch cc.Type {
			case raftpb.CONFIG_CHANGE_TYPE_ADD_NODE:
				if len(cc.Context) > 0 {
					rnd.transport.AddPeer(types.ID(cc.NodeID), []string{string(cc.Context)})
				}
			case raftpb.CONFIG_CHANGE_TYPE_REMOVE_NODE:
				if cc.NodeID == rnd.id {
					logger.Warningln("%s had already been removed!", types.ID(rnd.id))
					return
				}
				rnd.transport.RemovePeer(types.ID(cc.NodeID))
			}
		}

		// after commit, update appliedIndex
		rnd.mu.Lock()
		rnd.appliedIndex = ents[i].Index
		rnd.mu.Unlock()

		if ents[i].Index == rnd.lastIndex { // special nil commit to signal that replay has finished
			select {
			case rnd.commitc <- nil:
			case <-rnd.stopc:
				return
			}
		}
	}
}

const catchUpEntriesN = 10

// (etcd etcdserver.EtcdServer.snapshot)
func (rnd *raftNode) createSnapshot() {
	data, err := rnd.ds.createSnapshot()
	if err != nil {
		panic(err)
	}

	snap, err := rnd.storageMemory.CreateSnapshot(rnd.appliedIndex, &rnd.configState, data)
	if err != nil {
		// the snapshot was done asynchronously with the progress of raft.
		// raft might have already got a newer snapshot.
		if err == raft.ErrSnapOutOfDate {
			return
		}
		panic(err)
	}
	if err := rnd.storage.SaveSnap(snap); err != nil {
		panic(err)
	}
	logger.Infof("saved snapshot at index %d", snap.Metadata.Index)

	rnd.mu.Lock()
	compactIndex := uint64(1)
	if rnd.snapshotIndex > catchUpEntriesN {
		compactIndex = rnd.snapshotIndex - catchUpEntriesN
	}
	rnd.mu.Unlock()

	if err := rnd.storageMemory.Compact(compactIndex); err != nil {
		panic(err)
	}
	logger.Infof("saved snapshot at index %d", compactIndex)
}

func (rnd *raftNode) triggerSnapshot() {
	rnd.mu.Lock()
	if rnd.appliedIndex-rnd.snapshotIndex <= rnd.snapCount {
		rnd.mu.Unlock()
		return
	}
	rnd.mu.Unlock()

	logger.Infof("start snapshot [applied index: %d | last snapshot index: %d]", rnd.appliedIndex, rnd.snapshotIndex)
	rnd.createSnapshot()
	rnd.mu.Lock()
	rnd.snapshotIndex = rnd.appliedIndex
	rnd.mu.Unlock()
}

// (etcd etcdserver.EtcdServer.applyAll)
func (rnd *raftNode) applyAll(ap *apply) {
	rnd.applySnapshot(ap)
	rnd.applyEntries(ap)
	close(ap.applyDone)
	rnd.triggerSnapshot()
}

// (etcd etcdserver.raftNode.start, contrib.raftexample.raftNode.serveChannels)
func (rnd *raftNode) startRaftHandler() {
	snap, err := rnd.storageMemory.Snapshot()
	if err != nil {
		panic(err)
	}
	rnd.configState = snap.Metadata.ConfigState
	rnd.snapshotIndex = snap.Metadata.Index
	rnd.appliedIndex = snap.Metadata.Index

	defer rnd.storage.Close()

	ticker := time.NewTicker(time.Duration(rnd.electionTickN) * time.Millisecond)
	defer ticker.Stop()

	go rnd.handleProposal()

	for {
		select {
		case <-ticker.C:
			rnd.node.Tick()

		case rd := <-rnd.node.Ready():
			isLeader := false
			if rd.SoftState != nil && rd.SoftState.NodeState == raftpb.NODE_STATE_LEADER {
				isLeader = true
			}

			applyDone := make(chan struct{})
			go rnd.applyAll(&apply{
				entriesToApply: rd.EntriesToApply,
				snapshotToSave: rd.SnapshotToSave,
				applyDone:      applyDone,
			})

			// (Raft §10.2.1 Writing to the leader’s disk in parallel, p.141)
			// leader writes the new log entry to disk before replicating the entry
			// to its followers. Then, the followers write the entry to their disks.
			// Fortunately, the leader can write to its disk in parallel with replicating
			// to the followers and them writing to their disks.
			if isLeader {
				rnd.transport.Send(rd.MessagesToSend)
			}

			// etcdserver/raft.go: r.storage.Save(rd.HardState, rd.Entries)
			if err := rnd.storage.Save(rd.HardStateToSave, rd.EntriesToAppend); err != nil {
				panic(err)
			}

			if !raftpb.IsEmptySnapshot(rd.SnapshotToSave) {
				// etcdserver/raft.go: r.storage.SaveSnap(rd.Snapshot)
				if err := rnd.storage.SaveSnap(rd.SnapshotToSave); err != nil {
					panic(err)
				}

				// etcdserver/raft.go: r.raftStorage.ApplySnapshot(rd.Snapshot)
				if err := rnd.storageMemory.ApplySnapshot(rd.SnapshotToSave); err != nil {
					panic(err)
				}
			}

			// etcdserver/raft.go: r.raftStorage.Append(rd.Entries)
			if err := rnd.storageMemory.Append(rd.EntriesToAppend...); err != nil {
				panic(err)
			}

			if !isLeader {
				rnd.transport.Send(rd.MessagesToSend)
			}

			// wait for the raft routine to finish the disk writes before triggering a
			// snapshot. or applied index might be greater than the last index in raft
			// storage, since the raft routine might be slower than apply routine.
			<-applyDone

			// after commit, must call Advance
			// etcdserver/raft.go: r.Advance()
			rnd.node.Advance()

		case err := <-rnd.transport.Errc:
			rnd.errc <- err

			logger.Warningln("stopping %s;", types.ID(rnd.id), err)
			select {
			case <-rnd.stopc:
			default:
				rnd.stop()
			}
			return

		case <-rnd.stopc:
			return
		}
	}
}
