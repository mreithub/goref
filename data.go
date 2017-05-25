package goref

// Data -- Reference counter Snapshot data
type Data struct {
	// Currently active invocations
	Active int32 `json:"active"`

	// Total number of (finished) invocations
	Total int64 `json:"total"`

	// Total time spent (in nanoseconds)
	TotalNsec int64 `json:"totalNsec"`

	// Computad field (TotalNsec/1000000), provided for convenience
	TotalMsec int64 `json:"totalMsec"`

	// Computed field (TotalMsec/TotalCount), provided for convenience
	AvgMsec float32 `json:"avgMsec"`
}

// Creates a Data object from an (internal) data object
//
// Copies all the duplicate fields over and calculates the convenience fields.
func newData(d *data) *Data {
	var avgMsec float64
	if d.total > 0 {
		avgMsec = float64(d.totalNsec) / float64(1000000.*d.total)
	}

	return &Data{
		Active:    d.active,
		Total:     d.total,
		TotalNsec: d.totalNsec,
		TotalMsec: d.totalNsec / 1000000,
		AvgMsec:   float32(avgMsec),
	}
}
