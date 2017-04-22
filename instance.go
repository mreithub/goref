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
	now := time.Now()
	data := i.parent.get(i.key)
	atomic.AddInt32(&data.RefCount, -1)
	nsec := now.Sub(i.startTime).Nanoseconds()
	atomic.AddInt64(&data.TotalNsec, nsec)
}
