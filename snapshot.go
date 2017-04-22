package goref

import "time"

// Snapshot -- (Quasi-)read-only point-in-time copy of a GoRef instance
type Snapshot struct {
	// Snapshot/Clone data
	Data map[string]*Data

	// Snapshot name ("" for simple clones)
	Name string `json:",omitempty"`

	// Creation timestamp
	Ts time.Time

	// Previous Snapshot (or nil if this is the first one)
	Previous *Snapshot
}

// Keys -- List all keys of this read-only instance
func (s *Snapshot) Keys() []string {
	rc := make([]string, 0, len(s.Data))

	for k := range s.Data {
		rc = append(rc, k)
	}

	return rc
}
