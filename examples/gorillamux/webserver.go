package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mreithub/goref"
)

func indexHTML(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<h1>Index</h1>
  <a href="/delayed.html">delayed.html</a><br />
  <a href="/goref.json">goref.json</a>`))
}

func delayedHTML(w http.ResponseWriter, r *http.Request) {
	time.Sleep(200 * time.Millisecond)
	msg := fmt.Sprintf("The time is %s", time.Now().String())
	w.Write([]byte(msg))
}

func gorefJSON(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(goref.GetSnapshot().Data)

	w.Header().Add("Content-type", "application/json")
	w.Write(data)
}

func trackRequests(router *mux.Router) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to find the matching HTTP route (we'll use that as GoRef key)
		var match mux.RouteMatch
		if router.Match(r, &match) {
			path, _ := match.Route.GetPathTemplate()
			path = fmt.Sprintf("%s %s", r.Method, path)

			ref := goref.Ref(path)
			router.ServeHTTP(w, r)
			ref.Deref()
		} else {
			// No route found (-> 404 error)
			router.ServeHTTP(w, r)
		}
	})
}

func main() {
	var r = mux.NewRouter()

	r.HandleFunc("/", indexHTML)
	r.HandleFunc("/delayed.html", delayedHTML)
	r.HandleFunc("/goref.json", gorefJSON)

	var handler = handlers.LoggingHandler(os.Stdout, trackRequests(r))
	log.Fatal(http.ListenAndServe("localhost:1234", handler))
}
