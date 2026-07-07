package runner

import (
	"sync"
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
			}
		}()
	}

	wg.Wait()
	close(results)

	var workerResults []worker.Result
	for res := range results {
		workerResults = append(workerResults, res)
	}

	return workerResults
}
