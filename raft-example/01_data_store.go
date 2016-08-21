package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"sync"

	"github.com/gyuho/db/pkg/fileutil"
)

type keyValue struct {
	Key string
	Val string
}

type dataStore struct {
	// shared with raftNode
	propc   chan []byte // propc to write proposals "TO"
	commitc chan []byte // commitc to receive ready-to-commit data "FROM"
	///////////////////////////////

	errc  chan error
	stopc chan struct{}
	donec chan struct{}

	mu    sync.RWMutex
	store map[string]string
}

func newDataStore(propc, commitc chan []byte) *dataStore {
	ds := &dataStore{
		propc:   propc,
		commitc: commitc,

		errc:  make(chan error),
		stopc: make(chan struct{}),
		donec: make(chan struct{}),

		store: make(map[string]string),
	}
	go ds.readCommit()
	return ds
}

func (ds *dataStore) get(key string) (string, bool) {
	ds.mu.RLock()
	v, ok := ds.store[key]
	ds.mu.RUnlock()
	return v, ok
}

func (ds *dataStore) propose(ctx context.Context, kv keyValue) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(kv); err != nil {
		ds.errc <- err
		return
	}
	data := buf.Bytes()

	for {
		select {
		case ds.propc <- data:
			return

		case <-ctx.Done():
			ds.errc <- ctx.Err()
			return

		case <-ds.donec:
			return
		}
	}
}

func (ds *dataStore) stop() {
	close(ds.stopc)
	<-ds.donec
}

func (ds *dataStore) readCommit() {
	for {
		select {
		case cc := <-ds.commitc:
			if len(cc) == 0 {
				continue
			}
			var kv keyValue
			if err := gob.NewDecoder(bytes.NewBuffer(cc)).Decode(&kv); err != nil {
				panic(err)
			}
			ds.mu.Lock()
			ds.store[kv.Key] = kv.Val
			ds.mu.Unlock()

		case err := <-ds.errc:
			if err != nil {
				panic(err)
			}

		case <-ds.stopc:
			close(ds.donec)
			return

		case <-ds.donec:
			return
		}
	}
}

func (ds *dataStore) saveSnapshot(fpath string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(ds.store); err != nil {
		ds.errc <- err
		return
	}

	if err := fileutil.WriteSync(fpath, buf.Bytes(), fileutil.PrivateFileMode); err != nil {
		ds.errc <- err
		return
	}
}

func (ds *dataStore) loadSnapshot(fpath string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	f, err := fileutil.OpenToRead(fpath)
	if err != nil {
		ds.errc <- err
		return
	}
	defer f.Close()

	var store map[string]string
	if err := gob.NewDecoder(f).Decode(&store); err != nil {
		ds.errc <- err
		return
	}

	ds.store = store
}
