package main

import (
	"flag"
	"fmt"
	"os"
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

	cfg := runner.Config{
		URL:         *url,
		N:           *n,
		Concurrency: *c,
		Timeout:     *timeout,
	}

	res, td := runner.Run(cfg)
	s := stats.Aggregate(res, td)

	report.Print(cfg, s)
}
