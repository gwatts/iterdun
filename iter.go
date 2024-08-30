package iterdun

import (
	"context"
	"iter"
	"sync"
)

// Parallel iterates over multiple iterators of the same type, pulling from each in
// parallel and merging the result into a single resulting iterator.  The iteration order
// is non-deterministic.
//
// Iteration will stop if the supplied context is cancelled.
func Parallel[E any](ctx context.Context, iters ...iter.Seq[E]) iter.Seq[E] {
	return func(yield func(item E) bool) {
		var wg sync.WaitGroup
		wg.Add(len(iters))

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		ch := make(chan E)
		for _, i := range iters {
			go func() {
				defer wg.Done()
				for item := range i {
					select {
					case ch <- item:
					case <-ctx.Done():
						return
					}
				}
			}()
		}

		go func() {
			wg.Wait()
			close(ch)
		}()

		for item := range ch {
			if !yield(item) {
				return
			}
		}
	}
}
