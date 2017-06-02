# GoRef - Simple (and fast) go-style invocation tracker

GoRef is a small [Go][golang] package which implements a simple key-based
method invocation counter and timing profiler.

It can be used to:
- track execution time of your functions/goroutines
- find bottlenecks in your code
- Check if your goroutines exit properly
- Track calls to your HTTP endpoints (and their execution time) - see below

To access the internal profiling data, use `GetSnapshot()`.
It'll ask the worker goroutine to create a deep copy of the GoRef's instance current state.

GoRef's code is thread safe. It uses a messaging channel read by a single worker goroutine
which does the heavy lifting.  
Calls to `Ref()` and `Deref()` are asynchronous
(that asynchronousity doesn't affect time measurement though).  


### Getting started

Download the package, e.g.:

    go get github.com/mreithub/goref

Add the following snippet to each function (or goroutine) you want to track
(and replace 'foo' with your own key names).

```go
ref := goref.Ref("foo"); defer ref.Deref()
```

The above snippet uses `GoRef` in singleton mode. But you can also create your
own `GoRef` instances (and e.g. use different ones in different parts of your
application):

```go
g := goref.NewGoRef()

// and then instead of the above snippet:
ref := g.Ref("foo"); defer ref.Deref()
```


At any point in time you can call `GetSnapshot()` to obtain a deep copy of the measurements.


### Example (excerpt from [webserver.go](examples/webserver.go)):

This example shows how to use GoRef in your web applications.  
Here it tracks all web handler invocations.

Have a look at the usage documentation at [godoc.org][godoc].

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

	data, _ := json.Marshal(goref.GetSnapshot().Data)

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
    "active": 0,
    "count": 5,
    "usec": 12,
    "avgMsec": 0.0024592
  },
  "/goref.json": {
    "Active": 1,
    "count": 9,
    "usec": 547,
    "avgMsec": 0.060820557
  },
  "/delayed.html": {
    "active": 0,
    "count": 2,
    "usec": 412555,
    "avgMsec": 206.27777
  }
}
```

- `active`: the number of currently active instances
- `count`: number of (finished) instances (doesn't include the `active` ones yet)
- `usec`: total time spent in that function (in microseconds)
- `avgMsec`: calculated average (`usec/(1000*total)`)

### Using [`gorilla-mux`][gorillamux]

If you're using [gorilla-mux][gorillamux], there's a simple way to
add GoRef to your project:

```go
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
```

and in your main function something like:

```go
var router = mux.NewRouter()
// add your routes here using router.HandleFunc() and the like
var addr = ":8080"
var handler = handlers.LoggingHandler(os.Stdout, trackRequests(router))
log.Fatal(http.ListenAndServe(addr, handler))
```

[golang]: https://golang.org/
[godoc]: https://godoc.org/github.com/mreithub/goref
[gorillamux]: https://github.com/gorilla/mux
