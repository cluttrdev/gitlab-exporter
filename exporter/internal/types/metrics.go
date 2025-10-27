package types

type Metric struct {
	Id  string
	Iid int64
	Job JobReference

	Name      string
	Labels    map[string]string
	Value     float64
	Timestamp int64
}
