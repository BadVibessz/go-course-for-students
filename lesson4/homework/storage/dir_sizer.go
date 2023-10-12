package storage

import (
	"context"
	"sync"
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

	// TODO: add other fields as you wish
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{}
}

func getAllFiles(ctx context.Context, d Dir, m *sync.Mutex, errChan chan error) ([]File, error) {

	var totalFiles []File
	var err error

	wg := sync.WaitGroup{}

	dirs, files, err := d.Ls(ctx)
	if err != nil {
		return totalFiles, err
	}

	totalFiles = append(totalFiles, files...)
	for _, dir := range dirs {
		wg.Add(1)
		go func(d Dir) {
			defer wg.Done()

			select {
			case <-errChan:
				// error occurred
				return

			default:
				subFiles, err := getAllFiles(ctx, d, m, errChan)
				if err != nil {
					errChan <- err
					close(errChan)
					return
				} else {
					m.Lock()
					totalFiles = append(totalFiles, subFiles...)
					m.Unlock()
					return
				}
			}

		}(dir)

	}
	wg.Wait()

	err, ok := <-errChan
	if !ok && err != nil {
		return totalFiles, err
	}

	return totalFiles, nil
}

func worker(ctx context.Context, m *sync.Mutex, wg *sync.WaitGroup, dirCh chan Dir, allFiles []File) {
	defer wg.Done()

	// todo: use errorGroup
	select {
	case dir, ok := <-dirCh:

		if !ok {
			// todo handle
			return
		}

		dirs, files, err := dir.Ls(ctx)

		if err != nil {
			//errCh <- err
			return
		}

		// add new dirs into dir channel
		for _, dir := range dirs {
			dirCh <- dir
		}

		// add new files into resulting slice
		m.Lock()
		allFiles = append(allFiles, files...)
		m.Unlock()

	//case <-errCh:
	// todo error occurred

	default:
		return
	}

}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {

	var totalSize int64
	var totalCount int64

	var files []File

	//errChan := make(chan error, 1)

	m := sync.Mutex{}
	wg := sync.WaitGroup{}

	dirChan := make(chan Dir)
	dirChan <- d

	for i := 0; i < a.maxWorkersCount; i++ {
		wg.Add(1)
		go worker(ctx, &m, &wg, dirChan, files)
	}

	// listen for context cancellation
	wg.Add(1)
	go func() {
		defer wg.Done()

		select {
		case <-ctx.Done():
			close(dirChan)
			return
		}
	}()
	wg.Wait()

	for _, f := range files {
		size, err := f.Stat(ctx)

		if err != nil {
			return Result{}, err
		}
		totalSize += size
		totalCount++
	}

	return Result{totalSize, totalCount}, nil
}
