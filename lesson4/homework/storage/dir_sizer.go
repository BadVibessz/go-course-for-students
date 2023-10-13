package storage

import (
	"context"
	"sync"
	"sync/atomic"
)

// Result represents the Size function result
type Result struct {
	// Total Size of File objects
	Size int64
	// Count is a count of File objects processed
	Count int64
}

type DirSizer interface {
	// Size calculate a size of given Dir, receive a ctx and the root Dir instance
	// will return Result or error if happened
	Size(ctx context.Context, d Dir) (Result, error)
}

// sizer implement the DirSizer interface
type sizer struct {
	// maxWorkersCount number of workers for asynchronous run
	maxWorkersCount int
	countMadeDir    int64
	countDir        int64
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{
		maxWorkersCount: 10,
		countMadeDir:    0,
		countDir:        0,
	}
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {

	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}

	dirCh := make(chan Dir)
	endCh := make(chan any, 10)
	errCh := make(chan error, 1)

	var size int64
	var count int64

	var files []File
	for i := 1; i <= a.maxWorkersCount; i++ {
		wg.Add(1)
		go a.worker(ctx, wg, mutex, files, dirCh, endCh, errCh, &size, &count)
	}

	dirCh <- d
	atomic.AddInt64(&a.countDir, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			a.close(dirCh, endCh, errCh)
			return
		case <-endCh:
			a.close(dirCh, endCh, errCh)
			return
		}
	}()
	wg.Wait()

	for _, f := range files {

		sz, err := f.Stat(ctx)
		if err != nil {
			return Result{1, 1}, err
		}

		size += sz
		count++
	}

	if err := <-errCh; err != nil {
		return Result{}, err
	}
	return Result{size, count}, nil
}

func (a *sizer) close(dirCh chan Dir, endCh chan any, errCh chan error) {
	close(dirCh)
	close(endCh)
	close(errCh)
}

func (a *sizer) worker(ctx context.Context, wg *sync.WaitGroup, mutex *sync.Mutex, allFiles []File,
	dirCh chan Dir, endCh chan any, errCh chan error, sz *int64, ct *int64) {
	defer wg.Done()

	for dir := range dirCh {
		dirSlice, fileSlice, err := dir.Ls(ctx)
		if err != nil {
			errCh <- err
			endCh <- struct{}{}
			return
		}

		atomic.AddInt64(&a.countDir, int64(len(dirSlice)))

		for _, dir := range dirSlice {
			dirCh <- dir
		}

		for _, f := range fileSlice {
			size, err := f.Stat(ctx)
			if err != nil {
				endCh <- err
				errCh <- err

				return
			}

			atomic.AddInt64(sz, size)
			atomic.AddInt64(ct, 1)
		}

		//TODO: I FUCKING DO NOT UNDERSTAND WHY ALLFILES STAYS NIL
		//mutex.Lock()
		//allFiles = append(allFiles, fileSlice...)
		//mutex.Unlock()

		//atomic.AddInt64(sz, int64(len(fileSlice)))

		atomic.AddInt64(&a.countMadeDir, 1)

		if atomic.LoadInt64(&a.countMadeDir) == atomic.LoadInt64(&a.countDir) {
			endCh <- struct{}{}
			return
		}
	}
}
