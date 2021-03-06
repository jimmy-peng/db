package raft

import "github.com/gyuho/db/raft/raftpb"

// Ready represents entries and messages that are ready to read,
// ready to save to stable storage, ready to commit, ready to be
// sent to other peers. Ready is point-in-time state of a Node.
//
// All fields in Ready are read-only.
//
// (etcd raft.Ready)
type Ready struct {
	// SoftState provides state that is useful for logging and debugging.
	// The state is volatile and does not need to be persisted to the WAL.
	//
	// SoftState is nil, if there is no update.
	//
	// (etcd raft.Ready.*SoftState)
	SoftState *raftpb.SoftState

	// HardStateToSave is the current state of the Node to be saved in stable storage
	// BEFORE messages are sent out.
	//
	// HardStateToSave is raftpb.EmptyHardState, if there is no update.
	//
	// (etcd raft.Ready.raftpb.HardState)
	HardStateToSave raftpb.HardState

	// SnapshotToSave specifies the Snapshot to save to stable storage.
	// Only leader can send Snapshot.
	//
	// (etcd raft.Ready.Snapshot)
	SnapshotToSave raftpb.Snapshot

	// EntriesToAppend specifies the entries to save to stable storage
	// BEFORE messages are sent out.
	//
	// (etcd raft.Ready.Entries)
	EntriesToAppend []raftpb.Entry

	// EntriesToApply are committed entries, which have already been
	// saved in stable storage, so ready to apply.
	//
	// (etcd raft.Ready.CommittedEntries)
	EntriesToApply []raftpb.Entry

	// MessagesToSend is outbound messages to be sent AFTER EntriesToAppend are committed
	// to the stable storage. If it contains raftpb.LEADER_SNAPSHOT_REQUEST, the application
	// MUST report back to Raft when the snapshot has been received or has failed, by calling
	// ReportSnapshot.
	MessagesToSend []raftpb.Message

	// ReadStates are updated when Raft receives raftpb.MESSAGE_TYPE_TRIGGER_READ_INDEX,
	// only valid for the requested read-request.
	// ReadState is used to serve linearized read-only quorum-get requests without going
	// through Raft log appends, when the Node's applied index is greater than the index in ReadState.
	ReadStates []ReadState
}

// ContainsUpdates returns true if Ready contains any updates.
//
// (etcd raft.Ready.containsUpdates)
func (rd Ready) ContainsUpdates() bool {
	return rd.SoftState != nil ||
		!raftpb.IsEmptyHardState(rd.HardStateToSave) ||
		!raftpb.IsEmptySnapshot(rd.SnapshotToSave) ||
		len(rd.EntriesToAppend) > 0 ||
		len(rd.EntriesToApply) > 0 ||
		len(rd.MessagesToSend) > 0 ||
		len(rd.ReadStates) > 0
}

// (etcd raft.newReady)
func newReady(rnd *raftNode, prevSoftState *raftpb.SoftState, prevHardState raftpb.HardState) Ready {
	rd := Ready{
		EntriesToAppend: rnd.storageRaftLog.unstableEntries(),
		EntriesToApply:  rnd.storageRaftLog.nextEntriesToApply(),
		MessagesToSend:  rnd.mailbox,
	}

	if softState := rnd.softState(); !softState.Equal(prevSoftState) {
		rd.SoftState = softState
	}

	if hardState := rnd.hardState(); !hardState.Equal(prevHardState) {
		rd.HardStateToSave = hardState
	}

	if rnd.storageRaftLog.storageUnstable.snapshot != nil {
		rd.SnapshotToSave = *rnd.storageRaftLog.storageUnstable.snapshot
	}

	if len(rnd.readStates) != 0 {
		rd.ReadStates = rnd.readStates
	}
	return rd
}
