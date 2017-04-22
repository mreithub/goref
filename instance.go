package goref

import (
	"sync/atomic"
	"time"
)

// Instance - Trackable instance
type Instance struct {
	parent    *GoRef
	key       string
	startTime time.Time
}

// Deref -- Dereference an instance of 'key'
func (i Instance) Deref() {
	var now time.Time
	if !i.startTime.IsZero() {
		// only measure time if startTime was set
		now = time.Now()
	}

	data := i.parent.get(i.key)
	atomic.AddInt32(&data.refCount, -1)
	if !now.IsZero() {
		nsec := now.Sub(i.startTime).Nanoseconds()
		atomic.AddInt64(&data.totalNsec, nsec)
	}
}
