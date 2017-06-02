package goref

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBasics(t *testing.T) {
	g := NewGoRef()

	r := g.Ref("hello")
	time.Sleep(3 * time.Microsecond)
	r.Deref()

	clone1 := g.GetSnapshot()
	ref := g.Ref("world")
	time.Sleep(100 * time.Millisecond)
	clone2 := g.GetSnapshot()
	ref.Deref()
	ref = g.Ref("hello")
	clone3 := g.GetSnapshot()
	ref.Deref()

	// all the assertions are done after the fact (to make sure the different clones
	// keep their own copies of the Data)

	g.GetSnapshot() // wait for run() to catch up

	// final (current) state
	assert.Contains(t, g.data, "hello")
	assert.Contains(t, g.data, "world")
	d := g.get("hello")
	assert.Equal(t, int32(0), d.active)
	assert.Equal(t, int64(2), d.count)
	assert.True(t, d.nsec > 0)
	d = g.get("world")
	assert.Equal(t, int32(0), d.active)
	assert.Equal(t, int64(1), d.count)
	assert.True(t, d.nsec >= 100000000)

	// clone1: Ref('hello'), Deref('hello')
	keys := clone1.Keys()
	assert.Contains(t, keys, "hello")
	assert.NotContains(t, keys, "world")
	d1 := clone1.Data["hello"]
	assert.Equal(t, int32(0), d1.Active)
	assert.Equal(t, int64(1), d1.Count)
	assert.True(t, d1.USec > 0)
	assert.Equal(t, 1, len(clone1.Data))

	// clone2: clone1 + Ref('world'),  sleep(100ms)
	keys = clone2.Keys()
	assert.Contains(t, keys, "hello")
	assert.Contains(t, keys, "world")
	d2 := clone2.Data["world"]
	assert.Equal(t, int32(1), d2.Active)
	assert.Equal(t, int64(0), d2.Count)
	assert.Equal(t, int64(0), d2.USec)

	// clone3: clone2 + Deref('world'), Ref('hello')
	keys = clone3.Keys()
	assert.Contains(t, keys, "hello")
	assert.Contains(t, keys, "world")
	d3 := clone3.Data["world"]
	assert.Equal(t, int32(0), d3.Active)
	assert.Equal(t, int64(1), d3.Count)
	assert.True(t, d3.USec >= 100000)
	assert.True(t, clone3.Data["hello"].USec < 100)
	assert.NotEqual(t, d1.USec, d3.USec)
}
