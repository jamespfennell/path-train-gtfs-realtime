package main

import (
	"flag"
	"fmt"
	"github.com/jamespfennell/path-train-gtfs-realtime/feed"
	"github.com/jamespfennell/path-train-gtfs-realtime/server"
	"os"
	"time"
)

func main() {
	port := flag.Int("port", 8080, "the port to bind the HTTP server to")
	updatePeriod := flag.Duration("update_period", 5*time.Second, "how often to update the feed")
	timeoutPeriod := flag.Duration("timeout_period", 5*time.Second,
		"the maximum duration to wait for a response from the source API")
	useHTTPSourceAPI := flag.Bool("use_http_source_api", false,
		"use the HTTP source API instead of the default gRPC API")
	flag.Parse()

	f, err := feed.NewFeed(*updatePeriod, *timeoutPeriod, *useHTTPSourceAPI)
	if err != nil {
		fmt.Println("Failed to initialize feed:", err)
		os.Exit(1)
	}
	server.Run(*port, f)
}
