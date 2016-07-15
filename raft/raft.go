package raft

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/gyuho/db/raft/raftpb"
)

// NoNodeID is a placeholder node ID, only used when there is no leader in the cluster,
// or to reset leader transfer.
const NoNodeID uint64 = 0

// raftNode represents Raft-algorithm-specific node.
//
// (etcd raft.raft)
type raftNode struct {
	id    uint64
	state raftpb.NODE_STATE

	leaderID      uint64               // (etcd raft.raft.lead)
	allProgresses map[uint64]*Progress // (etcd raft.raft.prs)

	// (etcd raft.raft.raftLog)
	storageRaftLog *storageRaftLog

	rand *rand.Rand

	// electionTimeoutTickNum is the number of ticks for election to time out.
	//
	// (etcd raft.raft.electionTimeout)
	electionTimeoutTickNum int

	// electionTimeoutElapsedTickNum is the number of ticks elapsed
	// since the tick reached the last election timeout.
	//
	// (etcd raft.raft.electionElapsed)
	electionTimeoutElapsedTickNum int

	// randomizedElectionTimeoutTickNum is the random number between
	// [electionTimeoutTickNum, 2 * electionTimeoutTickNum - 1], and gets reset
	// when raftNode state changes to follower or to candidate.
	//
	// (etcd raft.raft.randomizedElectionTimeoutTickNum)
	randomizedElectionTimeoutTickNum int

	// heartbeatTimeoutTickNum is the number of ticks for leader to send heartbeat to its followers.
	//
	// (etcd raft.raft.heartbeatTimeout)
	heartbeatTimeoutTickNum int

	// heartbeatTimeoutElapsedTickNum is the number of ticks elapsed
	// since the tick reached the last heartbeat timeout.
	//
	// (etcd raft.raft.heartbeatElapsed)
	heartbeatTimeoutElapsedTickNum int

	tickFunc func()
	stepFunc func(r *raftNode, msg raftpb.Message)

	maxEntryNumPerMsg uint64 // (etcd raft.raft.maxMsgSize)
	maxInflightMsgNum int    // (etcd raft.raft.maxInflight)

	// leaderCheckQuorum is true, then the leader steps down
	// if quorum of cluster is not active for an election timeout.
	//
	// 1. leaderCheckQuorum is true
	// AND
	// 2. election time out hasn't passed yet
	//
	// THEN leader last confirmed its leadership, guaranteed to have been
	// in contact with quorum within the election timeout, so it will ignore
	// the vote request from candidate.
	//
	// So it shouldn't increase its term either.
	//
	// (etcd raft.raft.checkQuorum)
	leaderCheckQuorum bool

	term      uint64          // (etcd raft.raft.Term)
	votedFor  uint64          // (etcd raft.raft.Vote)
	votedFrom map[uint64]bool // (etcd raft.raft.votes)

	// mailbox contains a slice of messages to be filtered and processed by each step method.
	//
	// (etcd raft.raft.msgs)
	mailbox []raftpb.Message

	// pendingConfigExist is true, then new configuration will be ignored,
	// in preference to the unapplied configuration.
	//
	// (etcd raft.raft.pendingConf)
	pendingConfigExist bool

	// (Raft §3.10  Leadership transfer extension, p.28)
	// leaderTransfereeID is the ID of the leader transfer target
	// when it's not zero.
	//
	// (etcd raft.raft.leadTransferee)
	leaderTransfereeID uint64

	// (etcd raft.raft.readState)
	leaderReadState LeaderReadState
}

// newRaftNode creates a new raftNode with the given Config.
func newRaftNode(c *Config) *raftNode {
	if err := c.validate(); err != nil {
		raftLogger.Panicf("invalid raft.Config %+v (%v)", c, err)
	}

	if c.Logger != nil { // set the Logger
		raftLogger.SetLogger(c.Logger)
	}

	rnd := &raftNode{
		id:    c.ID,
		state: raftpb.NODE_STATE_FOLLOWER, // 0

		leaderID:       NoNodeID,
		allProgresses:  make(map[uint64]*Progress),
		storageRaftLog: newStorageRaftLog(c.StorageStable),

		rand: rand.New(rand.NewSource(int64(c.ID))),

		electionTimeoutTickNum:  c.ElectionTickNum,
		heartbeatTimeoutTickNum: c.HeartbeatTimeoutTickNum,

		maxEntryNumPerMsg: c.MaxEntryNumPerMsg,
		maxInflightMsgNum: c.MaxInflightMsgNum,

		leaderCheckQuorum: c.LeaderCheckQuorum,

		leaderReadState: LeaderReadState{Index: uint64(0), RequestCtx: nil},
	}

	hardState, configState, err := c.StorageStable.GetState()
	if err != nil {
		raftLogger.Panicf("newRaftNode c.StorageStable.GetState error (%v)", err)
	}
	if !raftpb.IsEmptyHardState(hardState) {
		rnd.loadHardState(hardState)
	}

	peerIDs := c.allPeerIDs
	if len(configState.IDs) > 0 {
		if len(peerIDs) > 0 {
			raftLogger.Panicf("cannot specify peer IDs both in Config.allPeerIDs and configState.IDs")
		}
		// overwrite peerIDs
		peerIDs = configState.IDs
	}
	for _, id := range peerIDs {
		rnd.allProgresses[id] = &Progress{
			NextIndex: 1,
			inflights: newInflights(rnd.maxInflightMsgNum),
		}
	}

	if c.LastAppliedIndex > 0 {
		rnd.storageRaftLog.appliedTo(c.LastAppliedIndex)
	}

	rnd.becomeFollower(rnd.term, rnd.leaderID)

	var nodeSlice []string
	for _, id := range rnd.allNodeIDs() {
		nodeSlice = append(nodeSlice, fmt.Sprintf("%x", id))
	}

	raftLogger.Infof(`

	newRaftNode

	state     = %q
	id        = %x
	all nodes = %v

	leader check quorum = %v

	first index = %d
	last  index = %d

	term      = %d
	last term = %d

	committed index = %d
	applied   index = %d

`, rnd.state, rnd.id, strings.Join(nodeSlice, ", "), rnd.leaderCheckQuorum,
		rnd.storageRaftLog.firstIndex(), rnd.storageRaftLog.lastIndex(),
		rnd.term, rnd.storageRaftLog.lastTerm(),
		rnd.storageRaftLog.committedIndex, rnd.storageRaftLog.appliedIndex,
	)

	return rnd
}

// (etcd raft.raft.quorum)
func (rnd *raftNode) quorum() int {
	return len(rnd.allProgresses)/2 + 1
}

// (etcd raft.raft.resetRandomizedElectionTimeout)
func (rnd *raftNode) randomizeElectionTickTimeout() {
	// [electiontimeout, 2 * electiontimeout - 1].
	rnd.randomizedElectionTimeoutTickNum = rnd.electionTimeoutTickNum + rnd.rand.Intn(rnd.electionTimeoutTickNum)
}

// (etcd raft.raft.pastElectionTimeout)
func (rnd *raftNode) pastElectionTimeout() bool {
	return rnd.electionTimeoutElapsedTickNum >= rnd.randomizedElectionTimeoutTickNum
}

// (etcd raft.raft.abortLeaderTransfer)
func (rnd *raftNode) stopLeaderTransfer() {
	rnd.leaderTransfereeID = NoNodeID
}

// (etcd raft.raft.resetPendingConf)
func (rnd *raftNode) resetPendingConfigExist() {
	rnd.pendingConfigExist = false
}

// (etcd raft.raft.reset)
func (rnd *raftNode) resetWithTerm(term uint64) {
	if rnd.term != term {
		rnd.term = term
		rnd.votedFor = NoNodeID
	}

	rnd.leaderID = NoNodeID
	rnd.votedFrom = make(map[uint64]bool)

	rnd.electionTimeoutElapsedTickNum = 0
	rnd.heartbeatTimeoutElapsedTickNum = 0
	rnd.randomizeElectionTickTimeout()

	rnd.stopLeaderTransfer()

	rnd.pendingConfigExist = false

	for id := range rnd.allProgresses {
		rnd.allProgresses[id] = &Progress{
			// NextIndex is the starting index of entries for next replication.
			NextIndex: rnd.storageRaftLog.lastIndex() + 1,
			inflights: newInflights(rnd.maxInflightMsgNum),
		}

		if id == rnd.id {
			// MatchIndex is the highest known matched entry index of this node.
			rnd.allProgresses[id].MatchIndex = rnd.storageRaftLog.lastIndex()
		}
	}
}

// (etcd raft.raft.setProgress)
func (rnd *raftNode) updateProgress(id, matchIndex, nextIndex uint64) {
	rnd.allProgresses[id] = &Progress{
		MatchIndex: matchIndex,
		NextIndex:  nextIndex,
		inflights:  newInflights(rnd.maxInflightMsgNum),
	}
}

// (etcd raft.raft.delProgress)
func (rnd *raftNode) deleteProgress(id uint64) {
	delete(rnd.allProgresses, id)
}

// (etcd raft.raft.softState)
func (rnd *raftNode) softState() *raftpb.SoftState {
	return &raftpb.SoftState{
		NodeState: rnd.state,
		LeaderID:  rnd.leaderID,
	}
}

// (etcd raft.raft.hardState)
func (rnd *raftNode) hardState() raftpb.HardState {
	return raftpb.HardState{
		VotedFor:       rnd.votedFor,
		CommittedIndex: rnd.storageRaftLog.committedIndex,
		Term:           rnd.term,
	}
}

// (etcd raft.raft.loadState)
func (rnd *raftNode) loadHardState(state raftpb.HardState) {
	if state.CommittedIndex < rnd.storageRaftLog.committedIndex || state.CommittedIndex > rnd.storageRaftLog.lastIndex() {
		raftLogger.Panicf("HardState of %x has committed index %d out of range [%d, %d]",
			rnd.id, state.CommittedIndex, rnd.storageRaftLog.committedIndex, rnd.storageRaftLog.lastIndex())
	}

	rnd.votedFor = state.VotedFor
	rnd.storageRaftLog.committedIndex = state.CommittedIndex
	rnd.term = state.Term
}
