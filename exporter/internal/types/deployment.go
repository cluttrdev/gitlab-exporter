package types

import (
	"time"
)

type EnvironmentReference struct {
	Id   int64
	Name string
	Tier string

	Project ProjectReference
}

type Environment struct {
	Id      int64
	Project ProjectReference

	CreatedAt *time.Time
	UpdatedAt *time.Time

	Name        string
	Slug        string
	Description string

	State string
	Tier  string
}

type Deployment struct {
	Id  int64
	Iid int64

	Job         JobReference
	Triggerer   UserReference
	Environment EnvironmentReference

	CreatedAt  *time.Time
	FinishedAt *time.Time
	UpdatedAt  *time.Time

	Status string
	Ref    string
	Sha    string
}
