package goref

import (
	"sync"
	"time"
)

// TODO tracking execution time might cause performance issues (e.g. in virtualized environments gettimeofday() might be slow)
//   if that turns out to be the case, deactivate Data.TotalNsec

// data -- internal GoRef data structure
type data struct {
	// currently active invocations
	active int32
	// number of finished invocations
	count int64
	// time spent in those invocations (in nanoseconds)
	nsec int64
}

// event types (for internal communication):
const (
	// stop the goroutine handling this GoRef instance
	evStop = iota
	// resets this GoRef instance
	evReset = iota
	// Takes a snapshot and sends it to snapshotChannel
	evSnapshot = iota
	// increments a ref counter
	evRef = iota
	// decrements a ref counter (and updates the total count + time)
	evDeref = iota
)

type event struct {
	typ  int
	key  string
	nsec int64
}

// GoRef -- A simple, go-style key-based reference counter that can be used for profiling your application (main class)
type GoRef struct {
	name   string
	parent *GoRef

	data map[string]*data

	_children map[string]*GoRef
	childLock sync.Mutex

	evChannel       chan event
	snapshotChannel chan Snapshot
}

func (g *GoRef) do(evType int, key string, nsec int64) {
	g.evChannel <- event{
		typ:  evType,
		key:  key,
		nsec: nsec,
	}
}

// get -- Get the Data object for the specified key (or create it) - thread safe
func (g *GoRef) get(key string) *data {
	rc, ok := g.data[key]
	if !ok {
		rc = &data{}
		g.data[key] = rc
	}

	return rc
}

// Ref -- References an instance of 'key'
func (g *GoRef) Ref(key string) *Instance {
	g.do(evRef, key, 0)

	return &Instance{
		parent:    g,
		key:       key,
		startTime: time.Now(),
	}
}

func (g *GoRef) run() {
	for msg := range g.evChannel {
		//log.Print("~~goref: ", msg)
		switch msg.typ {
		case evRef:
			g.get(msg.key).active++
			break
		case evDeref:
			d := g.get(msg.key)
			d.active--
			d.count++
			d.nsec += msg.nsec
			break
		case evSnapshot:
			g.takeSnapshot()
			break
		case evReset:
			g.data = map[string]*data{}
			break
		case evStop:
			return // TODO stop this GoRef instance safely
		default:
			panic("unsupported GoRef event type")
		}
	}
}

// GetChild -- Gets (or creates) a specific child instance (recursively)
func (g *GoRef) GetChild(path ...string) *GoRef {
	if len(path) == 0 {
		return g
	}

	firstSegment := path[0]

	var child *GoRef
	{ // keep the lock as short as possible
		g.childLock.Lock()
		defer g.childLock.Unlock()

		var ok bool
		child, ok = g._children[firstSegment]
		if !ok {
			// create a new child transparently
			child = newGoRef(firstSegment, g)
			g._children[firstSegment] = child
		}
	}

	return child.GetChild(path[1:]...)
}

// GetChildren -- Creates a point-in-time copy of this GoRef instance's children
func (g *GoRef) GetChildren() map[string]*GoRef {
	g.childLock.Lock()
	defer g.childLock.Unlock()

	// simply copy all entries
	var rc = make(map[string]*GoRef, len(g._children))
	for name, child := range g._children {
		rc[name] = child
	}

	return rc
}

// GetParent -- Get the parent of this GoRef instance (will return nil for root instances)
func (g *GoRef) GetParent() *GoRef {
	// g.parent is immutable -> no locking necessary
	return g.parent
}

// GetPath -- Get this GoRef instance's path (i.e. its parents' and its own name)
//
// Root instances have empty names, all the others have the name you give them
// when creating them with GetChild().
//
// To get a single string path, you can use strings.Join()
//
// ```go
// strings.Join(g.GetPath(), "/")
// ```
func (g *GoRef) GetPath() []string {
	var rc []string
	// this method needs no thread safety mechanisms.
	// A GoRef's name and parent are immutable.
	if g.parent != nil {
		rc = append(g.parent.GetPath(), g.name)
	}

	return rc
}

// GetSnapshot -- Creates and returns a deep copy of the current state (including child instance states)
func (g *GoRef) GetSnapshot() Snapshot {
	g.do(evSnapshot, "", 0)

	// get child snapshots while we wait
	children := g.GetChildren()
	childData := make(map[string]Snapshot, len(children))

	for name, child := range children {
		childData[name] = child.GetSnapshot()
	}

	rc := <-g.snapshotChannel
	rc.Children = childData
	return rc
}

// takeSnapshot -- internal (-> thread-unsafe) method taking a deep copy of the current state and sending it to snapshotChannel
func (g *GoRef) takeSnapshot() {
	// copy entries
	data := make(map[string]Data, len(g.data))
	for key, d := range g.data {
		data[key] = newData(d)
	}

	// send Snapshot
	g.snapshotChannel <- Snapshot{
		Data: data,
		Ts:   time.Now(),
	}
}

// Reset -- Resets this GoRef instance to its initial state
func (g *GoRef) Reset() {
	g.do(evReset, "", 0)
}

// NewGoRef -- Construct a new root-level GoRef instance
func NewGoRef() *GoRef {
	return newGoRef("", nil)
}

func newGoRef(name string, parent *GoRef) *GoRef {
	rc := &GoRef{
		name:            name,
		parent:          parent,
		data:            map[string]*data{},
		_children:       map[string]*GoRef{},
		evChannel:       make(chan event, 100),
		snapshotChannel: make(chan Snapshot, 5),
	}
	go rc.run()

	return rc
}
