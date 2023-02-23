package server

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/jamespfennell/path-train-gtfs-realtime/feed"
	"github.com/jamespfennell/path-train-gtfs-realtime/monitoring"
)

//go:embed index.html
var indexHTMLPage string

type RunArgs struct {
	Port             int
	UpdatePeriod     time.Duration
	TimeoutPeriod    time.Duration
	UseHTTPSourceAPI bool
}

func Run(ctx context.Context, args RunArgs) error {
	f, err := feed.NewFeed(ctx, args.UpdatePeriod, args.TimeoutPeriod, args.UseHTTPSourceAPI, monitoring.RecordUpdate)
	if err != nil {
		return fmt.Errorf("failed to initialize feed: %s", err)
	}

	http.HandleFunc("/", rootHandler)
	http.Handle("/gtfsrt", monitoring.CountRequests(f.HttpHandler()))
	http.Handle("/metrics", monitoring.HttpHandler())

	return http.ListenAndServe(fmt.Sprintf(":%d", args.Port), nil)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index.html").Parse(indexHTMLPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, "data goes here")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
