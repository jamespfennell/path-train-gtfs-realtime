package server

import (
	"fmt"
	"github.com/jamespfennell/path-train-gtfs-realtime/feed"
	"github.com/jamespfennell/path-train-gtfs-realtime/monitoring"
	"html/template"
	"log"
	"net/http"
)

func Run(port int, f *feed.Feed) {
	go monitoring.Listen(f.AddUpdateBroadcaster())

	http.HandleFunc("/", rootHandler)
	http.Handle("/gtfsrt", monitoring.CountRequests(f.HttpHandler()))
	http.Handle("/metrics", monitoring.HttpHandler())

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
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
