package executor

import (
	"context"
)

type (
	In  <-chan any
	Out = In
)

type Stage func(in In) (out Out)

func ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {
	out := make(chan any)

	go func() {
		defer close(out)

		stageOut := in
		for _, stage := range stages {
			stageOut = stage(stageOut)
		}

		for {
			select {
			case val, open := <-stageOut:
				if open {
					out <- val
				} else {
					return
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}
