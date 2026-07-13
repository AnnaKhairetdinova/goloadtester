package stats

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AnnaKhairetdinova/goloadtester/worker"
)

func TestTable_driven(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	res1 := worker.Result{
		StatusCode: 200,
		Duration:   50 * time.Millisecond,
		Err:        nil,
	}

	res2 := worker.Result{
		StatusCode: 200,
		Duration:   90 * time.Millisecond,
		Err:        nil,
	}

	res3 := worker.Result{
		StatusCode: 200,
		Duration:   99 * time.Millisecond,
		Err:        nil,
	}

	workerRes := []worker.Result{}
	workerRes = append(workerRes, res1, res2, res3)

	result := Aggregate(workerRes, 3*time.Second)

	if result.P50 < 49*time.Millisecond || result.P50 > 51*time.Millisecond {
		t.Errorf("P50 = %v, ожидалось ~50ms", result.P50)
	}

	if result.P90 < 89*time.Millisecond || result.P90 > 91*time.Millisecond {
		t.Errorf("P90 = %v, ожидалось ~90ms", result.P90)
	}

	if result.P99 < 98*time.Millisecond || result.P99 > 100*time.Millisecond {
		t.Errorf("P99 = %v, ожидалось ~99ms", result.P99)
	}

	mean := (50 + 90 + 99) * time.Millisecond / 3
	if result.Mean < mean-1*time.Millisecond || result.Mean > mean+1*time.Millisecond || result.Mean < mean {
		t.Errorf("Mean = %v, ожидалось ~%v", result.Mean, mean)
	}

	rps := float64(3) / 3.0
	if result.RPS < rps-0.01 || result.RPS > rps+0.01 {
		t.Errorf("RPS = %v, ожидалось ~%v", result.RPS, rps)
	}
}

func TestStats(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		result := Aggregate([]worker.Result{}, time.Second)

		if result.P50 != 0 || result.P90 != 0 || result.P99 != 0 {
			t.Errorf("Ожидались нулевые значения для пустого слайса")
		}

		if result.Mean != 0 {
			t.Errorf("Mean должен быть 0 для пустого слайса")
		}
	})

	t.Run("all errors", func(t *testing.T) {
		results := []worker.Result{
			{
				StatusCode: 0,
				Duration:   50 * time.Millisecond,
				Err:        fmt.Errorf("error 1"),
			},
			{
				StatusCode: 0,
				Duration:   90 * time.Millisecond,
				Err:        fmt.Errorf("error 2"),
			},
			{
				StatusCode: 0,
				Duration:   100 * time.Millisecond,
				Err:        fmt.Errorf("error 3"),
			},
		}

		result := Aggregate(results, 3*time.Second)

		if result.P50 == 0 {
			t.Errorf("P50 не может быть 0 при наличии результатов")
		}

		if result.Mean == 0 {
			t.Errorf("Mean не может быть 0 при наличии результатов")
		}
	})
}
