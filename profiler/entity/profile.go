package entity

// profile specifies a profiles data (gzipped protobuf, json), and the types contained within it.
type Profile struct {
	// name indicates profile type and format (e.g. cpu.pprof, metrics.json)
	Name string
	Data []byte
}
