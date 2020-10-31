package stats

// Meta represents process metadata. This will be not changed
// as long as the process continues.
type Meta struct {
	Username   string
	Cmmand     string
	GoMaxProcs int
	NumCPU     int
}
