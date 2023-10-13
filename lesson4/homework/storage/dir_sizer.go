package storage

import (
	"context"
	"fmt"
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
	err             error
	dirCh           chan Dir
	endCh           chan struct{}
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{
		maxWorkersCount: 10,
		countMadeDir:    0,
		countDir:        0,
		dirCh:           make(chan Dir),
		endCh:           make(chan struct{}, 10),
	}
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	res := Result{Size: 0, Count: 0}
	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}

	for i := 1; i <= a.maxWorkersCount; i++ {
		wg.Add(1)
		go a.dirProcessor(ctx, wg, mutex, &res)
	}

	a.dirCh <- d
	atomic.AddInt64(&a.countDir, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			a.close()
			return
		case <-a.endCh:
			a.close()
			return
		}
	}()
	wg.Wait()

	if a.err != nil {
		return Result{}, a.err
	}
	return res, nil
}

func (a *sizer) close() {
	close(a.dirCh)
	close(a.endCh)
}

func (a *sizer) dirProcessor(ctx context.Context, wg *sync.WaitGroup, mutex *sync.Mutex, res *Result) {
	defer wg.Done()

	for dir := range a.dirCh {
		dirSlice, fileSlice, err := dir.Ls(ctx)
		if err != nil {
			mutex.Lock()
			a.err = fmt.Errorf("dirProcessor : %w", err)
			mutex.Unlock()
			a.endCh <- struct{}{}
			return
		}
		atomic.AddInt64(&a.countDir, int64(len(dirSlice)))
		for _, dir := range dirSlice {
			a.dirCh <- dir
		}
		for _, file := range fileSlice {
			size, err := file.Stat(ctx)
			if err != nil {
				mutex.Lock()
				a.err = fmt.Errorf("fileProcessor : %w", err)
				mutex.Unlock()
				a.endCh <- struct{}{}
				return
			}

			atomic.AddInt64(&res.Count, 1)
			atomic.AddInt64(&res.Size, size)
		}
		atomic.AddInt64(&a.countMadeDir, 1)

		if atomic.LoadInt64(&a.countMadeDir) == atomic.LoadInt64(&a.countDir) {
			a.endCh <- struct{}{}
			return
		}
	}
}
