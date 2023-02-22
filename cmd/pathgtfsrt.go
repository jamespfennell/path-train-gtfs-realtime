package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/jamespfennell/path-train-gtfs-realtime/server"
)

func main() {
	port := flag.Int("port", 8080, "the port to bind the HTTP server to")
	updatePeriod := flag.Duration("update_period", 5*time.Second, "how often to update the feed")
	timeoutPeriod := flag.Duration("timeout_period", 5*time.Second,
		"the maximum duration to wait for a response from the source API")
	useHTTPSourceAPI := flag.Bool("use_http_source_api", false,
		"use the HTTP source API instead of the default gRPC API")
	flag.Parse()
	if err := server.Run(context.Background(), server.RunArgs{
		Port:             *port,
		UpdatePeriod:     *updatePeriod,
		TimeoutPeriod:    *timeoutPeriod,
		UseHTTPSourceAPI: *useHTTPSourceAPI,
	}); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
