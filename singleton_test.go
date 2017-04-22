package goref

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleton(t *testing.T) {
	Ref("hello").Deref()
	ref := Ref("world")
	clone := Clone()
	ref.Deref()
	snap := TakeSnapshot("snap")

	// current state
	assert.Contains(t, instance.data, "hello")
	assert.Contains(t, instance.data, "world")
	d := instance.get("hello")
	assert.Equal(t, int32(0), d.refCount)
	assert.Equal(t, int64(1), d.totalCount)
	d = instance.get("world")
	assert.Equal(t, int32(0), d.refCount)
	assert.Equal(t, int64(1), d.totalCount)
	assert.Equal(t, snap, instance.Snapshots())

	// reset instance
	Reset()
	assert.Empty(t, Snapshots())
	Ref("bla").Deref()

	assert.NotContains(t, instance.data, "hello")
	assert.Contains(t, instance.data, "bla")

	//
	// check Snapshot data after the fact
	//

	// clone: Ref('hello'), Deref('hello'), Ref('world')
	d1 := clone.Data["hello"]
	assert.Equal(t, int32(0), d1.RefCount)
	assert.Equal(t, int64(1), d1.TotalCount)
	assert.True(t, d1.TotalNsec > 0)
	d2 := clone.Data["world"]
	assert.Equal(t, int32(1), d2.RefCount)
	assert.Equal(t, int64(1), d2.TotalCount)
	assert.Equal(t, int64(0), d2.TotalNsec)
	assert.Equal(t, 2, len(clone.Data))

	// snap: clone + Deref('world')
	d1 = snap.Data["hello"]
	assert.Equal(t, int32(0), d1.RefCount)
	assert.Equal(t, int64(1), d1.TotalCount)
	assert.True(t, d1.TotalNsec > 0)
	d2 = snap.Data["world"]
	assert.Equal(t, int32(0), d2.RefCount)
	assert.Equal(t, int64(1), d2.TotalCount)
	assert.True(t, d2.TotalNsec > 0)
	assert.Equal(t, 2, len(snap.Data))
}
