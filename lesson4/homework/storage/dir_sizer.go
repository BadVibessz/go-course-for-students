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

func getAllFiles(ctx context.Context, d Dir, m *sync.Mutex) ([]File, error) {

	var totalFiles []File

	dirs, files, err := d.Ls(ctx)

	if err != nil {
		return totalFiles, err
	}

	totalFiles = append(totalFiles, files...)

	wg := sync.WaitGroup{}

	wg.Add(len(dirs))
	for _, dir := range dirs {

		go func(d Dir) {
			defer wg.Done()

			subFiles, err := getAllFiles(ctx, d, m)

			if err != nil {
				//return totalFiles, err
				// todo:
			}

			m.Lock()
			totalFiles = append(totalFiles, subFiles...)
			m.Unlock()
		}(dir)

	}

	wg.Wait()
	return totalFiles, nil
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {

	var totalSize int64
	var totalCount int64

	dirs, files, err := d.Ls(ctx)

	if err != nil {
		return Result{}, err
	}

	var totalFiles []File
	totalFiles = append(totalFiles, files...)

	m := sync.Mutex{}

	wg := sync.WaitGroup{}

	wg.Add(len(dirs))
	for _, dir := range dirs {
		go func(d Dir) {
			defer wg.Done()

			allFiles, err := getAllFiles(ctx, d, &m)

			if err != nil {
				// todo:
			}

			m.Lock()
			totalFiles = append(totalFiles, allFiles...)
			m.Unlock()

		}(dir)

		if err != nil {
			return Result{}, err
		}

	}

	wg.Wait()
	for _, f := range totalFiles {
		size, err := f.Stat(ctx)

		if err != nil {
			return Result{}, err
		}
		totalSize += size
		totalCount++
	}

	return Result{totalSize, totalCount}, nil
}
