package report

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/AnnaKhairetdinova/goloadtester/runner"
	"github.com/AnnaKhairetdinova/goloadtester/stats"
)

func Print(cfg runner.Config, s stats.Stats) {
	fmt.Fprintf(os.Stdout, "URL %s\n", cfg.URL)
	fmt.Fprintf(os.Stdout, "Запросов %d\n", cfg.N)
	fmt.Fprintf(os.Stdout, "Конкурентность %d\n", cfg.Concurrency)
	fmt.Fprintf(os.Stdout, "Таймаут %s\n\n", cfg.Timeout)

	fmt.Fprintf(os.Stdout, "Всего времени %s\n", beautifulDuration(s.TotalDuration))
	fmt.Fprintf(os.Stdout, "RPS %.2f\n\n", s.RPS)

	fmt.Fprintf(os.Stdout, "Латентность:\n")
	fmt.Fprintf(os.Stdout, "Среднее %s\n", beautifulDuration(s.Mean))
	fmt.Fprintf(os.Stdout, "Минимум %s\n", beautifulDuration(s.Min))
	fmt.Fprintf(os.Stdout, "Максимум %s\n", beautifulDuration(s.Max))
	fmt.Fprintf(os.Stdout, "p50 %s\n", beautifulDuration(s.P50))
	fmt.Fprintf(os.Stdout, "p90 %s\n", beautifulDuration(s.P90))
	fmt.Fprintf(os.Stdout, "p99 %s\n\n", beautifulDuration(s.P99))

	fmt.Fprintf(os.Stdout, "Статус-коды:\n")
	// структура и сортировка для вывода статус-кодов
	type pair struct {
		key   int
		value int
	}

	pairs := make([]pair, len(s.StatusCodes))
	for k, v := range s.StatusCodes {
		pairs = append(pairs, pair{k, v})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].value > pairs[j].value
	})

	for _, p := range pairs {
		if p.key != 0 {
			percentP := float64(p.value) / float64(cfg.N) * 100
			fmt.Fprintf(os.Stdout, "%d: %d (%.1f%%)\n", p.key, p.value, percentP)
		}
	}

	percentErr := float64(s.Errors) / float64(cfg.N) * 100
	fmt.Fprintf(os.Stdout, "Ошибки: %d (%.1f%%)\n", s.Errors, percentErr)
}

func beautifulDuration(t time.Duration) string {
	switch {
	case t < time.Millisecond:
		return fmt.Sprintf("%d µs", t.Microseconds())
	case t < time.Second:
		return fmt.Sprintf("%d ms", t.Milliseconds())
	default:
		return fmt.Sprintf("%.2f s", t.Seconds())
	}
}
