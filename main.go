package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/AnnaKhairetdinova/goloadtester/runner"
)

func main() {
	url := flag.String("url", "", "URL для запроса")
	n := flag.Int("n", 10, "Количество запросов")
	c := flag.Int("c", 5, "Количество параллельных запросов")
	timeout := flag.Duration("timeout", 10*time.Second, "Время запроса")
	flag.Parse()

	cfg := runner.Config{
		URL:         *url,
		N:           *n,
		Concurrency: *c,
		Timeout:     *timeout,
	}

	if *url == "" {
		fmt.Println("Ошибка: флаг -url обязательный")
		flag.Usage()
		os.Exit(1)
	}

	if *c > *n {
		fmt.Println("Ошибка: -c не может быть больше -n")
		*c = *n
	}

	fmt.Printf("URL: %s\n", cfg.URL)
	fmt.Printf("n=%d, c=%d, timeout=%v\n", cfg.N, cfg.Concurrency, cfg.Timeout)

	result := runner.Run(cfg)
	fmt.Print(result[:11])
}
