package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type Result struct {
	url    string
	status string
	code   int
	time   time.Duration
}

func main() {
	argsWithoutProg := os.Args[1:]

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if len(argsWithoutProg) == 0 {
		fmt.Println("Please provide URLs e.g: https://google.com")
		return
	}

	// results := make(chan Result, len(argsWithoutProg))
	var wg sync.WaitGroup

	for _, v := range argsWithoutProg {
		wg.Add(1)
		go makeConnection(v, &wg, ctx)
	}

	wg.Wait()
}

func makeConnection(url string, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	start := time.Now()
	// resp, err := http.Get(url)
	req,err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Println("Malformed url")
	}
	resp, err := http.DefaultClient.Do(req)

	if errors.Is(err, context.DeadlineExceeded) {
		res := Result{
			url:    url,
			status: "Connection Timed out",
			code:   500,
			time:   time.Since(start),
		}
		fmt.Printf("url: %v, result: %v, statusCode: %v, (%v)\n", res.url, res.status, res.code, res.time)
		return
	} else if err != nil {
		res := Result{
			url:    url,
			status: err.Error(),
			code:   500,
			time:   time.Since(start),
		}
		fmt.Printf("url: %v, result: %v, statusCode: %v, (%v)\n", res.url, res.status, res.code, res.time)
		return
	}
	defer resp.Body.Close()
	res := Result{
		url:    url,
		status: "Success",
		code:   resp.StatusCode,
		time:   time.Since(start),
	}
	fmt.Printf("url: %v, result: %v, statusCode: %v, (%v)\n", res.url, res.status, res.code, res.time)
}
