package goref

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnapshots(t *testing.T) {
	g := NewGoRef()

	assert.Empty(t, g.Snapshots())

	g.Ref("hello").Deref()

	s1 := g.TakeSnapshot("s1")
	assert.Equal(t, s1.Name, "s1")
	assert.Equal(t, []string{"hello"}, s1.Keys())
	assert.Equal(t, s1, g.Snapshots())

	ref := g.Ref("hello")
	s2 := g.TakeSnapshot("s2")
	ref.Deref()

	//test snapshot data
	snap2 := g.Snapshots()
	snap1 := snap2.Previous
	assert.Equal(t, int32(1), snap2.Get("hello").RefCount)
	assert.Equal(t, int64(2), snap2.Get("hello").TotalCount)

	assert.Equal(t, int32(0), snap1.Get("hello").RefCount)
	assert.Equal(t, int64(1), snap1.Get("hello").TotalCount)

	// test linked list (and order)
	assert.Equal(t, s2, snap2)
	assert.Equal(t, s1, snap1)
	assert.Equal(t, s2.Previous, s1)
	assert.Nil(t, s1.Previous)
}
