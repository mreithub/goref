package goref

import (
	"fmt"
	"testing"
	"time"
)

var nsec int64
var snap Snapshot

// BenchmarkMeasureTime -- Measures how long measuring the time takes (using time.Now() and Time.Sub())
func BenchmarkMeasureTime(b *testing.B) {
	for n := 0; n < b.N; n++ {
		start := time.Now()
		end := time.Now()
		nsec = end.Sub(start).Nanoseconds()
	}

}

// BenchmarkRefDeref -- Measures how long an empty Ref().Deref() call takes
func BenchmarkRefDeref(b *testing.B) {
	g := NewGoRef()

	for n := 0; n < b.N; n++ {
		g.Ref("hello").Deref()
	}
	//snap := g.Clone()
	//j, _ := json.Marshal(snap.Data)
	//log.Printf("data: %s", j)
}

// BenchmarkRefDeref -- Measures how long an empty Ref().Deref() call takes (doing the Deref() in a defer statement)
func BenchmarkRefDerefDeferred(b *testing.B) {
	g := NewGoRef()

	for n := 0; n < b.N; n++ {
		r := g.Ref("hello")
		defer r.Deref()
	}
	//snap := g.Clone()
	//j, _ := json.Marshal(snap.Data)
	//log.Printf("data: %s", j)
}

// benchmarkGetSnapshot -- Measure how long it takes to create a deep copy of the snapshot data
func benchmarkGetSnapshot(count int, b *testing.B) {
	// setup
	g := NewGoRef()
	for n := 0; n < count; n++ {
		g.Ref(fmt.Sprintf("ref%d", n)).Deref()
	}

	for n := 0; n < b.N; n++ {
		snap = *g.Clone()
	}
}

func BenchmarkGetSnapshot100(b *testing.B) {
	benchmarkGetSnapshot(100, b)
}

func BenchmarkGetSnapshot1000(b *testing.B) {
	benchmarkGetSnapshot(1000, b)
}
