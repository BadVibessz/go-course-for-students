package storage

import (
	"context"
	"golang.org/x/sync/errgroup"
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

	errGroup, ctx := errgroup.WithContext(ctx)

	dirs, files, err := d.Ls(ctx)

	if err != nil {
		return totalFiles, err
	}

	totalFiles = append(totalFiles, files...)
	for _, dir := range dirs {
		errGroup.Go(
			func() error {

				subFiles, err := getAllFiles(ctx, dir, m)

				if err != nil {
					return err
				}

				m.Lock()
				totalFiles = append(totalFiles, subFiles...)
				m.Unlock()

				return nil
			})

	}

	if err := errGroup.Wait(); err != nil {
		return totalFiles, err
	}

	return totalFiles, nil
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {

	var totalSize int64
	var totalCount int64

	m := sync.Mutex{}

	totalFiles, err := getAllFiles(ctx, d, &m)
	if err != nil {
		return Result{}, err
	}

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
