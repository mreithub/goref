package goref

import "time"

// Snapshot -- point-in-time copy of a GoRef instance
type Snapshot struct {
	// Snapshot data
	Data map[string]Data

	// Creation timestamp
	Ts time.Time
}

// Keys -- List all keys of this read-only instance
func (s *Snapshot) Keys() []string {
	rc := make([]string, 0, len(s.Data))

	for k := range s.Data {
		rc = append(rc, k)
	}

	return rc
}
