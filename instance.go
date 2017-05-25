package goref

import (
	"log"
	"sync/atomic"
	"time"
)

// Instance - Trackable instance
type Instance struct {
	parent    *GoRef
	data      *data
	key       string
	startTime time.Time
}

// Deref -- Dereference an instance of 'key'
func (i *Instance) Deref() {
	if i.parent == nil {
		log.Print("GoRef warning: possible double Deref()")
		return
	}

	var now time.Time
	if !i.startTime.IsZero() {
		// only measure time if startTime was set
		now = time.Now()
	}

	d := i.data
	if d == nil {
		d = i.parent.get(i.key)
	}
	atomic.AddInt32(&d.active, -1)
	atomic.AddInt64(&d.total, 1)
	if !now.IsZero() {
		nsec := now.Sub(i.startTime).Nanoseconds()
		atomic.AddInt64(&d.totalNsec, nsec)
	}

	i.parent = nil // prevent double Deref()
}
