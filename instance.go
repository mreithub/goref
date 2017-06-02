package goref

import (
	"log"
	"time"
)

// Instance - Trackable instance
type Instance struct {
	parent    *GoRef
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
	var nsec int64
	if !i.startTime.IsZero() {
		// only measure time if startTime was set
		now = time.Now()
		nsec = now.Sub(i.startTime).Nanoseconds()
	}

	i.parent.do(evDeref, i.key, nsec)
	i.parent = nil // prevent double Deref()
}
