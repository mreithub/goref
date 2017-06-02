package goref

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleton(t *testing.T) {
	Ref("hello").Deref()
	ref := Ref("world")
	snap1 := GetSnapshot()
	ref.Deref()
	snap2 := GetSnapshot()

	// current state
	assert.Contains(t, instance.data, "hello")
	assert.Contains(t, instance.data, "world")
	d := instance.get("hello")
	assert.Equal(t, int32(0), d.active)
	assert.Equal(t, int64(1), d.count)
	d = instance.get("world")
	assert.Equal(t, int32(0), d.active)
	assert.Equal(t, int64(1), d.count)

	// reset instance
	Reset()
	Ref("bla").Deref()

	GetSnapshot() // synchronize

	assert.NotContains(t, instance.data, "hello")
	assert.Contains(t, instance.data, "bla")

	//
	// check Snapshot data after the fact
	//

	// snap1: Ref('hello'), Deref('hello'), Ref('world')
	d1 := snap1.Data["hello"]
	assert.Equal(t, int32(0), d1.Active)
	assert.Equal(t, int64(1), d1.Count)
	d2 := snap1.Data["world"]
	assert.Equal(t, int32(1), d2.Active)
	assert.Equal(t, int64(0), d2.Count)
	assert.Equal(t, int64(0), d2.USec)
	assert.Equal(t, 2, len(snap1.Data))

	// snap2: snap1 + Deref('world')
	d1 = snap2.Data["hello"]
	assert.Equal(t, int32(0), d1.Active)
	assert.Equal(t, int64(1), d1.Count)
	d2 = snap2.Data["world"]
	assert.Equal(t, int32(0), d2.Active)
	assert.Equal(t, int64(1), d2.Count)
	assert.True(t, d2.USec > 0)
	assert.Equal(t, 2, len(snap2.Data))
}
