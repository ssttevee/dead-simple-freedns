package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

func sync(ctx context.Context, prefix, token string) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%ssync.afraid.org/u/%s/", prefix, token), nil)
	if err != nil {
		log.Println("ERROR: failed to create request:", err.Error())
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("ERROR: failed to dispatch request:", err.Error())
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("ERROR: failed to read response body:", err.Error())
		fmt.Fprintln(os.Stderr, string(body))
		return
	}

	if res.StatusCode != http.StatusOK {
		log.Println("WARN: unexpected status code:", res.StatusCode)
		return
	}

	for line := range strings.SplitSeq(string(body), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		log.Println("INFO:", line)
	}
}

func main() {
	update4 := flag.Bool("4", true, "enable ipv4")
	update6 := flag.Bool("6", true, "enable ipv6")
	token := flag.String("token", "", "update token (required)")
	interval := flag.Duration("interval", 0, "update interval")
	flag.Parse()

	if !*update4 && !*update6 {
		log.Println("FATAL: both ipv4 and ipv6 are disabled")
		flag.Usage()
		os.Exit(125)
	}

	if *token == "" {
		log.Println("FATAL: token is empty")
		flag.Usage()
		os.Exit(125)
	}

	if *interval > 0 {
		log.Println("INFO: update interval is", *interval)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	for {
		if *update4 {
			sync(ctx, "", *token)
		}

		if *update6 {
			sync(ctx, "v6.", *token)
		}

		if *interval == 0 {
			return
		}

		select {
		case <-ctx.Done():
			log.Println("INFO: interrupt received")
			os.Exit(1)
		case <-time.After(*interval):
		}
	}
}
