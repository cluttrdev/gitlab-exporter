package types

import (
	"time"
)

type Section struct {
	Id  int64
	Job JobReference

	Name       string
	StartedAt  *time.Time
	FinishedAt *time.Time
	Duration   time.Duration
}
