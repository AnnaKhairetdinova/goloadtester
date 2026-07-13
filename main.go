package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AnnaKhairetdinova/goloadtester/report"
	"github.com/AnnaKhairetdinova/goloadtester/runner"
	"github.com/AnnaKhairetdinova/goloadtester/stats"
)

func main() {
	url := flag.String("url", "", "URL для запроса")
	n := flag.Int("n", 10, "Количество запросов")
	c := flag.Int("c", 5, "Количество параллельных запросов")
	timeout := flag.Duration("timeout", 10*time.Second, "Время запроса")
	flag.Parse()

	if *url == "" {
		fmt.Println("Ошибка: флаг -url обязательный")
		flag.Usage()
		os.Exit(1)
	}

	if *c > *n {
		fmt.Println("Ошибка: -c не может быть больше -n")
		*c = *n
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := runner.Config{
		URL:         *url,
		N:           *n,
		Concurrency: *c,
		Timeout:     *timeout,
	}

	start := time.Now()
	res := runner.Run(ctx, cfg)
	duration := time.Since(start)

	if ctx.Err() != nil {
		return
	}

	if len(res) == 0 {
		fmt.Println("Нет данных")
		os.Exit(1)
	}

	s := stats.Aggregate(res, duration)

	report.Print(cfg, s)
}
