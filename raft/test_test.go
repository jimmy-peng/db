package raft

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/gyuho/db/raft/raftpb"
)

// (etcd raft.stateMachine)
type stateMachine interface {
	Step(msg raftpb.Message) error
	readResetMailbox() []raftpb.Message
}

func (rnd *raftNode) readResetMailbox() []raftpb.Message {
	msgs := rnd.mailbox
	rnd.mailbox = make([]raftpb.Message, 0)
	return msgs
}

type blackHole struct{}

func (blackHole) Step(raftpb.Message) error          { return nil }
func (blackHole) readResetMailbox() []raftpb.Message { return nil }

var noOpBlackHole = &blackHole{}

type connection struct {
	from, to uint64
}

// fakeNetwork simulates network message passing for Raft tests.
//
// (etcd raft.network)
type fakeNetwork struct {
	allStateMachines         map[uint64]stateMachine
	allStableStorageInMemory map[uint64]*StorageStableInMemory

	allDroppedConnection  map[connection]float64
	allIgnoredMessageType map[raftpb.MESSAGE_TYPE]bool
}

// (etcd raft.idsBySize)
func generateIDs(n int) []uint64 {
	ids := make([]uint64, n)

	for i := 0; i < n; i++ {
		ids[i] = (uint64(i) << 56) | (uint64(time.Now().Add(time.Minute).UnixNano()))
	}
	return ids
}

func Test_generateIDs(t *testing.T) {
	ids := generateIDs(10)
	var prevID uint64
	for i, id := range ids {
		if i == 0 {
			prevID = id
			continue
		}
		fmt.Printf("generated %x\n", id)
		if id == prevID {
			t.Fatalf("#%d: expected %x != %x", i, prevID, id)
		}

		id = prevID
	}
}

// (etcd raft.newNetwork)
func newFakeNetwork(machines ...stateMachine) *fakeNetwork {
	peerIDs := generateIDs(len(machines))

	allStateMachines := make(map[uint64]stateMachine)
	allStableStorageInMemory := make(map[uint64]*StorageStableInMemory)

	for i := range machines {
		id := peerIDs[i]
		switch v := machines[i].(type) {
		case nil:
			allStateMachines[id] = newRaftNode(&Config{
				ID:         id,
				allPeerIDs: peerIDs,

				ElectionTickNum:         10,
				HeartbeatTimeoutTickNum: 1,
				LeaderCheckQuorum:       false,
				StorageStable:           NewStorageStableInMemory(),
				MaxEntryNumPerMsg:       0,
				MaxInflightMsgNum:       256,
				LastAppliedIndex:        0,
			})
			allStableStorageInMemory[id] = NewStorageStableInMemory()

		case *raftNode:
			v.id = id
			v.allProgresses = make(map[uint64]*Progress)
			for _, pid := range peerIDs {
				v.allProgresses[pid] = &Progress{}
			}
			v.resetWithTerm(0)
			allStateMachines[id] = v

		case *blackHole:
			allStateMachines[id] = v

		default:
			raftLogger.Panicf("unknown state machine type: %T", v)
		}
	}

	return &fakeNetwork{
		allStateMachines:         allStateMachines,
		allStableStorageInMemory: allStableStorageInMemory,

		allDroppedConnection:  make(map[connection]float64),
		allIgnoredMessageType: make(map[raftpb.MESSAGE_TYPE]bool),
	}
}

// (etcd raft.network.filter)
func (fn *fakeNetwork) filter(msgs ...raftpb.Message) []raftpb.Message {
	var filtered []raftpb.Message
	for _, msg := range msgs {
		if fn.allIgnoredMessageType[msg.Type] {
			continue
		}

		switch msg.Type {
		case raftpb.MESSAGE_TYPE_INTERNAL_TRIGGER_FOLLOWER_OR_CANDIDATE_TO_START_CAMPAIGN:
			raftLogger.Panicf("%q never goes over network", msg.Type)

		default:
			percentage := fn.allDroppedConnection[connection{from: msg.From, to: msg.To}]
			if rand.Float64() < percentage {
				continue // skip append
			}
		}

		filtered = append(filtered, msg)
	}

	return filtered
}

// (etcd raft.network.send)
func (fn *fakeNetwork) sendAndStepFirstMessage(msgs ...raftpb.Message) {
	if len(msgs) > 0 {
		firstMsg := msgs[0]
		machine := fn.allStateMachines[firstMsg.To]
		machine.Step(firstMsg)

		msgs = append(msgs[1:], fn.filter(machine.readResetMailbox()...)...)
	}
}

// (etcd raft.network.recover)
func (fn *fakeNetwork) recoverAll() {
	fn.allDroppedConnection = make(map[connection]float64)
	fn.allIgnoredMessageType = make(map[raftpb.MESSAGE_TYPE]bool)
}

// (etcd raft.network.drop)
func (fn *fakeNetwork) dropConnectionByPercentage(from, to uint64, percentage float64) {
	fn.allDroppedConnection[connection{from, to}] = percentage
}

// (etcd raft.network.cut)
func (fn *fakeNetwork) cutConnection(id1, id2 uint64) {
	fn.allDroppedConnection[connection{id1, id2}] = 1
	fn.allDroppedConnection[connection{id2, id1}] = 1
}

// (etcd raft.network.isolate)
func (fn *fakeNetwork) isolate(id uint64) {
	for sid := range fn.allStateMachines {
		if id != sid {
			fn.cutConnection(id, sid)
		}
	}
}

// (etcd raft.network.ignore)
func (fn *fakeNetwork) ignoreMessageType(tp raftpb.MESSAGE_TYPE) {
	fn.allIgnoredMessageType[tp] = true
}
