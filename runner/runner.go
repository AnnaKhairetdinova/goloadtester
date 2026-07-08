package runner

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/AnnaKhairetdinova/goloadtester/worker"
)

type Config struct {
	URL         string
	N           int
	Concurrency int
	Timeout     time.Duration
}

func Run(cfg Config) []worker.Result {
	var done atomic.Int64

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go showProgress(ctx, &done, cfg.N)

	jobs := make(chan struct{}, cfg.N)
	for i := 0; i < cfg.N; i++ {
		jobs <- struct{}{}
	}
	close(jobs)

	results := make(chan worker.Result, cfg.N)

	var wg sync.WaitGroup

	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				result := worker.Do(cfg.URL, cfg.Timeout)
				results <- result
				done.Add(1)
			}
		}()
	}

	wg.Wait()
	close(results)

	cancel()

	time.Sleep(50 * time.Millisecond)
	fmt.Println()

	var workerResults []worker.Result
	for res := range results {
		workerResults = append(workerResults, res)
	}

	return workerResults
}

func showProgress(ctx context.Context, done *atomic.Int64, total int) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			current := done.Load()
			percent := float64(current) / float64(total) * 100
			fmt.Printf("\rПрогресс: %d/%d (%.1f%%)  ", done, total, percent)

		case <-ctx.Done():
			fmt.Printf("\rПрогресс: %d/%d (100%%) - готово", total, total)
			return
		}
	}
}
