package models

type JobMetric struct {
	Name      string            `json:"name"`
	Labels    map[string]string `json:"labels"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"`
	Job       struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
}
