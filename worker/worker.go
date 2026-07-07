package worker

import (
	"net/http"
	"time"
)

type Result struct {
	StatusCode int
	Duration   time.Duration
	Err        error
}

func Do(url string, timeout time.Duration) Result {
	client := &http.Client{
		Timeout: timeout,
	}

	start := time.Now()
	resp, err := client.Get(url)
	duration := time.Since(start)

	if err != nil {
		return Result{
			StatusCode: 0,
			Duration:   duration,
			Err:        err,
		}
	}

	defer resp.Body.Close()

	return Result{
		StatusCode: resp.StatusCode,
		Duration:   duration,
		Err:        nil,
	}
}
