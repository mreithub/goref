package goref

import "time"

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
	data map[string]*data

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
			//case msgStop:
			//	return // TODO find a
		default:
			panic("unsupported GoRef event type")
		}
	}
}

// GetSnapshot -- Creates and returns a deep copy of the current state
func (g *GoRef) GetSnapshot() Snapshot {
	g.do(evSnapshot, "", 0)
	return <-g.snapshotChannel
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
	rc := &GoRef{
		data:            map[string]*data{},
		evChannel:       make(chan event, 100),
		snapshotChannel: make(chan Snapshot, 5),
	}
	go rc.run()

	return rc
}
