package server

import (
	"fmt"
	"github.com/jamespfennell/path-train-gtfs-realtime/feed"
	"github.com/jamespfennell/path-train-gtfs-realtime/monitoring"
	"html/template"
	"log"
	"net/http"
)

func Run(port int, f *feed.Feed, monitor *monitoring.Monitor) {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/feed/", feedHandlerFactory(f))
	http.HandleFunc("/status/", htmlHandlerFactory(monitor))
	http.HandleFunc("/status/json", jsonHandlerFactory(monitor))
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

func feedHandlerFactory(f *feed.Feed) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: set content type
		_, err := w.Write(f.Get())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func jsonHandlerFactory(m *monitoring.Monitor) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := m.WriteJson(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func htmlHandlerFactory(m *monitoring.Monitor) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := m.WriteHTML(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
