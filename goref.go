package goref

import (
	"sync"
	"sync/atomic"
	"time"
)

// TODO tracking execution time might cause performance issues (e.g. in virtualized environments gettimeofday() might be slow)
//   if that turns out to be the case, deactivate Data.TotalNsec

// GoRef -- A simple, thread safe key-based reference counter that can be used for profiling your application
type GoRef struct {
	data map[string]*Data
	lock *sync.Mutex

	snapshots *Snapshot
}

// get -- Get the Data object for the specified key (or create it) - thread safe
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

// Clone -- Returns a Snapshot of the GoRef  (synchronously)
func (g *GoRef) Clone() *Snapshot {
	g.lock.Lock()
	defer g.lock.Unlock()

	data := map[string]*Data{}

	for key, d := range g.data {
		data[key] = &Data{
			RefCount:   d.RefCount,
			TotalCount: d.TotalCount,
			TotalNsec:  d.TotalNsec,
		}
	}

	// return a cloned GoRef instance
	return &Snapshot{
		Data: data,
		Ts:   time.Now(),
	}
}

// Ref -- References an instance of 'key'
func (g *GoRef) Ref(key string) Instance {
	data := g.get(key)
	atomic.AddInt32(&data.RefCount, 1)
	atomic.AddInt64(&data.TotalCount, 1)

	return Instance{
		parent:    g,
		key:       key,
		startTime: time.Now(),
	}
}

// Snapshots -- Linked list of Snapshots (in reverse order)
func (g *GoRef) Snapshots() *Snapshot {
	// We assume here that pointer access is atomic (to avoid locking the Mutex)
	return g.snapshots
}

// TakeSnapshot -- Clone the current GoRef instance and return the new snapshot
func (g *GoRef) TakeSnapshot(name string) *Snapshot {
	// prepends the snapshot to the list
	rc := g.Clone()
	rc.Name = name

	g.lock.Lock()
	defer g.lock.Unlock()

	rc.Previous = g.snapshots
	g.snapshots = rc

	return rc
}

// NewGoRef -- GoRef constructor
func NewGoRef() *GoRef {
	return &GoRef{
		lock: &sync.Mutex{},
		data: map[string]*Data{},
	}
}
