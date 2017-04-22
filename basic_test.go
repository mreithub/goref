package goref

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBasics(t *testing.T) {
	g := NewGoRef()

	g.Ref("hello").Deref()
	clone1 := g.Clone()
	ref := g.Ref("world")
	time.Sleep(100 * time.Millisecond)
	clone2 := g.Clone()
	ref.Deref()
	ref = g.Ref("hello")
	clone3 := g.Clone()
	ref.Deref()

	// final (current) state
	assert.Contains(t, g.data, "hello")
	assert.Contains(t, g.data, "world")
	d := g.get("hello")
	assert.Equal(t, int32(0), d.refCount)
	assert.Equal(t, int64(2), d.totalCount)
	assert.True(t, d.totalNsec > 0)
	d = g.get("world")
	assert.Equal(t, int32(0), d.refCount)
	assert.Equal(t, int64(1), d.totalCount)
	assert.True(t, d.totalNsec >= 100000000)

	// clone1: Ref('hello'), Deref('hello')
	keys := clone1.Keys()
	assert.Contains(t, keys, "hello")
	assert.NotContains(t, keys, "world")
	d1 := clone1.Get("hello")
	assert.Equal(t, int32(0), d1.refCount)
	assert.Equal(t, int64(1), d1.totalCount)
	assert.True(t, d1.totalNsec > 0)

	// clone2: clone1 + Ref('world'),  sleep(100ms)
	keys = clone2.Keys()
	assert.Contains(t, keys, "hello")
	assert.Contains(t, keys, "world")
	d2 := clone2.Get("world")
	assert.Equal(t, int32(1), d2.refCount)
	assert.Equal(t, int64(1), d2.totalCount)
	assert.Equal(t, int64(0), d2.totalNsec)

	// clone3: clone2 + Deref('world'), Ref('hello')
	keys = clone3.Keys()
	assert.Contains(t, keys, "hello")
	assert.Contains(t, keys, "world")
	d3 := clone3.Get("world")
	assert.Equal(t, int32(0), d3.refCount)
	assert.Equal(t, int64(1), d3.totalCount)
	assert.True(t, d3.totalNsec >= 100000000)
	assert.True(t, clone3.Get("hello").totalNsec < 100000)
	assert.NotEqual(t, d1.totalNsec, d3.totalNsec)
}
