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

	// all the assertions are done after the fact (to make sure the different clones
	// keep their own copies of the Data)

	// final (current) state
	assert.Contains(t, g.data, "hello")
	assert.Contains(t, g.data, "world")
	d := g.get("hello")
	assert.Equal(t, int32(0), d.RefCount)
	assert.Equal(t, int64(2), d.TotalCount)
	assert.True(t, d.TotalNsec > 0)
	d = g.get("world")
	assert.Equal(t, int32(0), d.RefCount)
	assert.Equal(t, int64(1), d.TotalCount)
	assert.True(t, d.TotalNsec >= 100000000)

	// clone1: Ref('hello'), Deref('hello')
	keys := clone1.Keys()
	assert.Contains(t, keys, "hello")
	assert.NotContains(t, keys, "world")
	d1 := clone1.Get("hello")
	assert.Equal(t, int32(0), d1.RefCount)
	assert.Equal(t, int64(1), d1.TotalCount)
	assert.True(t, d1.TotalNsec > 0)
	assert.Equal(t, 1, len(clone1.GetData()))

	// clone2: clone1 + Ref('world'),  sleep(100ms)
	keys = clone2.Keys()
	assert.Contains(t, keys, "hello")
	assert.Contains(t, keys, "world")
	d2 := clone2.Get("world")
	assert.Equal(t, int32(1), d2.RefCount)
	assert.Equal(t, int64(1), d2.TotalCount)
	assert.Equal(t, int64(0), d2.TotalNsec)

	// clone3: clone2 + Deref('world'), Ref('hello')
	keys = clone3.Keys()
	assert.Contains(t, keys, "hello")
	assert.Contains(t, keys, "world")
	d3 := clone3.Get("world")
	assert.Equal(t, int32(0), d3.RefCount)
	assert.Equal(t, int64(1), d3.TotalCount)
	assert.True(t, d3.TotalNsec >= 100000000)
	assert.True(t, clone3.Get("hello").TotalNsec < 100000)
	assert.NotEqual(t, d1.TotalNsec, d3.TotalNsec)
}

func TestGetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			assert.Fail(t, "Expected a panic() in GoRef.Get()")
		}
	}()
	NewGoRef().Get("foo")
}

func TestGetDataPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			assert.Fail(t, "Expected a panic() in GoRef.GetData()")
		}
	}()
	NewGoRef().GetData()
}

func TestKeysPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			assert.Fail(t, "Expected a panic() in GoRef.Keys()")
		}
	}()
	NewGoRef().Keys()
}
