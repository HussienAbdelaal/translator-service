package utils

import (
	"context"
	"sync"
)

// WorkerFunc defines the signature of a function that could block the main thread
// and should be run in a goroutine.
type WorkerFunc func(ctx context.Context, input string) (string, error)

// DoInParallelFailFast concurrently executes WorkerFunc on a slice of strings while preserving order.
// It fails fast if any WorkerFunc fails.
func DoInParallelFailFast(ctx context.Context, inputs []string, workerFn WorkerFunc) ([]string, error) {
	results := make([]string, len(inputs))

	var (
		wg     sync.WaitGroup
		once   sync.Once
		errCh  = make(chan error, 1)
		cancel context.CancelFunc
	)

	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	for i, input := range inputs {
		wg.Add(1)
		go func(i int, text string) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return // skip if canceled
			default:
			}

			result, err := workerFn(ctx, text)
			if err != nil {
				once.Do(func() {
					errCh <- err
					cancel()
				})
				return
			}
			results[i] = result
		}(i, input)
	}

	wg.Wait()

	select {
	case err := <-errCh:
		return nil, err
	default:
		return results, nil
	}
}
