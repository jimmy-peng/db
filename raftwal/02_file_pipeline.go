package raftwal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gyuho/db/pkg/fileutil"
)

// filePipeline pipelines disk space allocation.
type filePipeline struct {
	dir   string
	size  int64 // in bytes.
	count int

	lockedFileCh chan *fileutil.LockedFile
	errc         chan error
	donec        chan struct{}
}

// (etcd wal.newFilePipeline)
func newFilePipeline(dir string, size int64) *filePipeline {
	fp := &filePipeline{
		dir:          dir,
		size:         size,
		count:        0,
		lockedFileCh: make(chan *fileutil.LockedFile),
		errc:         make(chan error, 1),
		donec:        make(chan struct{}),
	}
	go fp.run()
	return fp
}

const extendFile = true

// (etcd wal.filePipeline.alloc)
func (fp *filePipeline) alloc() (f *fileutil.LockedFile, err error) {
	fpath := filepath.Join(fp.dir, fmt.Sprintf("%d.tmp", fp.count%2)) // to make it different than previous one
	if f, err = fileutil.OpenFileWithLock(fpath, os.O_WRONLY|os.O_CREATE, fileutil.PrivateFileMode); err != nil {
		return nil, err
	}

	if err = fileutil.Preallocate(f.File, fp.size, extendFile); err != nil {
		logger.Errorf("failed to allocate space when creating %q (%v)", fpath, err)
		f.Close()
		return nil, err
	}
	fp.count++
	return f, nil
}

// Open returns a fresh file ready for writes.
// Rename the file before calling this or duplicate Open
// will trigger file collisions.
//
// (etcd wal.filePipeline.Open)
func (fp *filePipeline) Open() (f *fileutil.LockedFile, err error) {
	select {
	case f = <-fp.lockedFileCh:
	case err = <-fp.errc:
	}
	return
}

// (etcd wal.filePipeline.run)
func (fp *filePipeline) run() {
	defer close(fp.errc)

	for {
		f, err := fp.alloc()
		if err != nil {
			fp.errc <- err
			return
		}

		select {
		case fp.lockedFileCh <- f:
		case <-fp.donec:
			os.Remove(f.Name())
			f.Close()
			return
		}
	}
}

// Close closes filePipeline.
//
// (etcd wal.filePipeline.Close)
func (fp *filePipeline) Close() error {
	close(fp.donec)
	return <-fp.errc
}
