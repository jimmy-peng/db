package rafthttp

import (
	"net/http"
	"sync"
	"time"

	"github.com/gyuho/db/pkg/probing"
	"github.com/gyuho/db/pkg/tlsutil"
	"github.com/gyuho/db/pkg/types"
	"github.com/gyuho/db/raft/raftpb"
	"github.com/gyuho/db/raftsnap"
)

// Transporter defines rafthttp transport layer.
//
// (etcd rafthttp.Transporter)
type Transporter interface {
	// Start starts transporter.
	// Start must be called first.
	//
	// (etcd rafthttp.Transporter.Start)
	Start() error

	// Stop closes all connections and stops the transporter.
	//
	// (etcd rafthttp.Transporter.Stop)
	Stop()

	// HTTPHandler returns http.Handler with '/raft' prefix.
	//
	// (etcd rafthttp.Transporter.Handler)
	HTTPHandler() http.Handler

	// Send sends messages to its remote peers.
	//
	// (etcd rafthttp.Transporter.Send)
	Send(msgs []raftpb.Message)

	// SendSnapshot sends snapshot to its remote peers.
	//
	// (etcd rafthttp.Transporter.SendSnapshot)
	SendSnapshot(msg raftsnap.Message)

	// AddPeer adds a peer with given peer URLs to the transport.
	//
	// (etcd rafthttp.Transporter.AddPeer)
	AddPeer(peerID types.ID, urls []string)

	// RemovePeer removes the peer with the given ID.
	//
	// (etcd rafthttp.Transporter.RemovePeer)
	RemovePeer(peerID types.ID)

	// RemoveAllPeers removes all existing peers in transporter.
	//
	// (etcd rafthttp.Transporter.RemoveAllPeers)
	RemoveAllPeers()

	// UpdatePeer updates the peer with given ID and URLs.
	//
	// (etcd rafthttp.Transporter.UpdatePeer)
	UpdatePeer(peerID types.ID, urls []string)

	// ActiveSince returns the time that the connection with the peer became active.
	// If the connection is currently inactive, it returns zero time.
	//
	// (etcd rafthttp.Transporter.ActiveSince)
	ActiveSince(peerID types.ID) time.Time

	// AddPeerRemote adds a remote peer with given URLs.
	//
	// (etcd rafthttp.Transporter.AddRemote)
	AddPeerRemote(peerID types.ID, urls []string)
}

// Transport implements Transporter. It sends and receives raft messages to/from peers.
//
// (etcd rafthttp.Transport)
type Transport struct {
	TLSInfo     tlsutil.TLSInfo
	DialTimeout time.Duration

	Sender    types.ID   // (etcd rafthttp.Transport.ID)
	ClusterID types.ID   // (etcd rafthttp.Transport.ClusterID)
	PeerURLs  types.URLs // (etcd rafthttp.Transport.URLs)

	Raft            Raft                  // (etcd rafthttp.Transport.Raft)
	RaftSnapshotter *raftsnap.Snapshotter // (etcd rafthttp.Transport.Snapshotter)

	Errc chan error // (etcd rafthttp.Transport.ErrorC)

	streamRoundTripper         http.RoundTripper // (etcd rafthttp.Transport.streamRt)
	pipelineRoundTripper       http.RoundTripper // (etcd rafthttp.Transport.pipelineRt)
	pipelineRoundTripperProber probing.Prober    // (etcd rafthttp.Transport.prober)

	mu          sync.RWMutex
	peers       map[types.ID]Peer        // (etcd rafthttp.Transport.peers)
	peerRemotes map[types.ID]*peerRemote // (etcd rafthttp.Transport.remotes)
}

func (tr *Transport) Start() error {
	var err error
	tr.streamRoundTripper, err = NewStreamRoundTripper(tr.TLSInfo, tr.DialTimeout)
	if err != nil {
		return err
	}

	// allow more idle connections
	tr.pipelineRoundTripper, err = NewRoundTripper(tr.TLSInfo, tr.DialTimeout)
	if err != nil {
		return err
	}

	tr.pipelineRoundTripperProber = probing.NewProber(tr.pipelineRoundTripper)

	tr.peers = make(map[types.ID]Peer)
	tr.peerRemotes = make(map[types.ID]*peerRemote)
	return nil
}

// CutPeer drops messages to the specified peer.
func (tr *Transport) CutPeer(peerID types.ID) {
	tr.mu.RLock()
	p, pok := tr.peers[peerID]
	r, rok := tr.peerRemotes[peerID]
	tr.mu.RUnlock()

	if pok {
		p.(Pausable).Pause()
	}
	if rok {
		r.Pause()
	}
}

// MendPeer recovers the dropping message to the given peer.
func (tr *Transport) MendPeer(peerID types.ID) {
	tr.mu.RLock()
	p, pok := tr.peers[peerID]
	r, rok := tr.peerRemotes[peerID]
	tr.mu.RUnlock()

	if pok {
		p.(Pausable).Resume()
	}
	if rok {
		r.Resume()
	}
}

//
//
// ↓ need peerPipeline, streamReader, Peer, peer, peerRemote implementation
//
//
//
