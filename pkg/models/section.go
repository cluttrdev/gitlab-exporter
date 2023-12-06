package models

import (
	"time"
)

type Section struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Job  struct {
		ID     int64  `json:"id"`
		Name   string `json:"name"`
		Status string `json:"status"`
	} `json:"job"`
	Pipeline struct {
		ID        int64  `json:"id"`
		ProjectID int64  `json:"project_id"`
		Ref       string `json:"ref"`
		Sha       string `json:"sha"`
		Status    string `json:"status"`
	} `json:"pipeline"`
	StartedAt  *time.Time `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	Duration   float64    `json:"duration"`
}
