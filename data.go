package goref

// Data -- Reference counter Snapshot data
type Data struct {
	RefCount   int32
	TotalCount int64
	TotalNsec  int64
	TotalMsec  int64
	AvgMsec    float32
}

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
