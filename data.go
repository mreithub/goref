package goref

// Data -- Reference counter Snapshot data
type Data struct {
	// Currently active invocations
	RefCount int32
	// Total number of invocations
	TotalCount int64

	// Total number of nanoseconds spent in that function
	TotalNsec int64

	// Computad field (TotalNsec/1000000), provided for convenience
	TotalMsec int64
	// Computed field (TotalMsec/TotalCount), provided for convenience
	AvgMsec float32
}

// Creates a Data object from an (internal) data object
//
// Copies all the duplicate fields over and calculates the convenience fields.
func newData(d *data) *Data {
	var avgMsec float64
	if d.totalCount > 0 {
		avgMsec = float64(d.totalNsec) / float64(1000000.*d.totalCount)
	}

	return &Data{
		RefCount:   d.refCount,
		TotalCount: d.totalCount,
		TotalNsec:  d.totalNsec,
		TotalMsec:  d.totalNsec / 1000000,
		AvgMsec:    float32(avgMsec),
	}
}
