package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mreithub/goref"
)

func indexHTML(w http.ResponseWriter, r *http.Request) {
	ref := goref.Ref("/")
	defer ref.Deref()

	w.Write([]byte(`<h1>Index</h1>
  <a href="/delayed.html">delayed.html</a><br />
  <a href="/goref.json">goref.json</a>`))
}

func delayedHTML(w http.ResponseWriter, r *http.Request) {
	ref := goref.Ref("/hello.html")
	defer ref.Deref()

	time.Sleep(200 * time.Millisecond)
	msg := fmt.Sprintf("The time is %s", time.Now().String())
	w.Write([]byte(msg))
}

func gorefJSON(w http.ResponseWriter, r *http.Request) {
	ref := goref.Ref("/goref.json")
	defer ref.Deref()

	data, _ := json.Marshal(goref.Clone().Data)

	w.Header().Add("Content-type", "application/json")
	w.Write(data)
}

func main() {
	http.HandleFunc("/", indexHTML)
	http.HandleFunc("/delayed.html", delayedHTML)
	http.HandleFunc("/goref.json", gorefJSON)

	http.ListenAndServe("localhost:1234", nil)
}
