# GoRef - Simple (and thread safe) golang reference counter

GoRef is a small [Go][golang] package which implements a simple key-based
invocation counter.
This can be used to check that your code behaves the way you expect
(e.g. whether you're properly cleaning up resources you're using, that all your
goroutines exit as expected, ...)

It also tracks execution time, so it might help you to find bottlenecks in your
web application.

And there's snapshot support. It allows you to create Snapshot copies of the
ref counter (e.g. periodically or at certain key points in your application)
which help to further narrow down potential issues.

### Getting started

Download the package, e.g.:

    go get github.com/mreithub/goref

Add the following snippet to each function (or goroutine) you want to track
(and replace 'foo' with your own key names).

```go
goref.Ref("foo"); defer ref.Deref()
```

The above snippet uses `GoRef` in singleton mode. But you can also create your
own `GoRef` instances (and e.g. use different ones in different parts of your
application):

```go
g := goref.NewGoRef()

// and then instead of the above snippet:
g.Ref("foo"); defer ref.Deref()
```


At any point in time you can call `Clone()` to obtain a copy of the current state
or `TakeSnapshot(name string)` to create point-in-time snapshots.


### Example (excerpt from [webserver.go](examples/webserver.go)):

This example shows how to use GoRef in your web applications.  
Here it tracks all web handler invocations.

```go
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
```

Run it with

    go run examples/webserver.go

and browse to http://localhost:1234/

After accessing each page a couple of times `/goref.json` should look something
like this:

```json
{
  "/": {
    "RefCount": 0,
    "TotalCount": 5,
    "TotalNsec": 12296,
    "TotalMsec": 0,
    "AvgMsec": 0.0024592
  },
  "/goref.json": {
    "RefCount": 1,
    "TotalCount": 9,
    "TotalNsec": 547385,
    "TotalMsec": 0,
    "AvgMsec": 0.060820557
  },
  "/delayed.html": {
    "RefCount": 0,
    "TotalCount": 2,
    "TotalNsec": 412555528,
    "TotalMsec": 412,
    "AvgMsec": 206.27777
  }
}
```

Internally, GoRef calculates with nanosecond precision. The `TotalMsec` and `AvgMsec`
fields are provided for convenience.

[golang]: https://golang.org/
