package goref

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// TODO tracking execution time here might cause performance issues (e.g. in virtualized environments gettimeofday() might be slow)
//   if that turns out to be the case, deactivate Data.totalNsec

// singleton GoRef instance
var instance = NewGoRef()

// Data -- RefCounter data
type Data struct {
	refCount   int32
	totalCount int64
	totalNsec  int64
}

// GoRef -- A simple, thread safe key-based reference counter that can be used for profiling your application
type GoRef struct {
	data map[string]*Data
	lock *sync.Mutex

	// linked list to old snapshots
	lastSnapshot *GoRef
}

// Instance - Trackable instance
type Instance struct {
	parent    *GoRef
	key       string
	startTime time.Time
}

// Deref -- Dereference an instance of 'key'
func (i Instance) Deref() {
	now := time.Now()
	data := i.parent.get(i.key)
	atomic.AddInt32(&data.refCount, -1)
	nsec := now.Sub(i.startTime).Nanoseconds()
	log.Print("nsec: ", nsec)
	atomic.AddInt64(&data.totalNsec, nsec)
}

// get -- Get the Data object for the specified key (or create it)
func (g *GoRef) get(key string) *Data {
	g.lock.Lock()
	defer g.lock.Unlock()

	rc, ok := g.data[key]
	if !ok {
		rc = &Data{}
		g.data[key] = rc
	}

	return rc
}

// Clone -- Returns a copy of the GoRef  (synchronously)
func (g *GoRef) Clone() GoRef {
	g.lock.Lock()
	defer g.lock.Unlock()

	data := map[string]*Data{}

	for key, d := range g.data {
		data[key] = &Data{
			refCount:   d.refCount,
			totalCount: d.totalCount,
			totalNsec:  d.totalNsec,
		}
	}

	// return a cloned GoRef instance
	return GoRef{
		data:         data,
		lock:         nil, // clones are (meant to be) read-only -> no need for locks
		lastSnapshot: nil, //
	}
}

// Get -- returns the refcounter Data for the specified key (or nil if not found)
func (g *GoRef) Get(key string) *Data {
	if g.lock != nil {
		// make sure this instance is readonly
		panic("GoRef: Called Get() on an active instance! call Clone() or TakeSnapshot() first!")
	}

	return g.data[key]
}

// Keys -- List all keys of this read-only instance
func (g *GoRef) Keys() []string {
	if g.lock != nil {
		panic("GoRef: Called Keys() on an active instance! call Clone() or TakeSnapshot() first!")
	}
	rc := make([]string, 0, len(g.data))

	for k := range g.data {
		rc = append(rc, k)
	}

	return rc
}

// Ref -- References an instance of 'key'
func (g *GoRef) Ref(key string) Instance {
	data := g.get(key)
	atomic.AddInt32(&data.refCount, 1)
	atomic.AddInt64(&data.totalCount, 1)

	return Instance{
		parent:    g,
		key:       key,
		startTime: time.Now(),
	}
}

// TakeSnapshot -- Clone the current GoRef instance and return
func (g *GoRef) TakeSnapshot() GoRef {
	old := g.lastSnapshot
	rc := g.Clone()
	rc.lastSnapshot = old
	g.lastSnapshot = &rc
	return rc
}

// NewGoRef -- GoRef constructor
func NewGoRef() *GoRef {
	return &GoRef{
		lock:         &sync.Mutex{},
		data:         map[string]*Data{},
		lastSnapshot: nil,
	}
}
