package goref

import "time"

// Data -- Reference counter Snapshot data
type Data struct {
	// Currently active invocations
	Active int32 `json:"active"`

	// Total number of (finished) invocations
	Count int64 `json:"count"`

	// Total time spent
	Duration time.Duration `json:"duration"`

	// Computed average run time in msec, provided for convenience
	AvgMsec float32 `json:"avgMsec"`
}

// Fills a Data object with the values from an (internal) data object
//
// Copies all the duplicate fields over and calculates the convenience fields.
func newData(src *data) Data {
	var avgMsec float64
	if src.count > 0 {
		avgMsec = float64(src.nsec) / float64(1000000.*src.count)
	}

	return Data{
		Active:   src.active,
		Count:    src.count,
		Duration: time.Duration(src.nsec) * time.Nanosecond,
		AvgMsec:  float32(avgMsec),
	}
}
