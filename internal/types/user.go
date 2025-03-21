package types

import "go.cluttr.dev/gitlab-exporter/protobuf/typespb"

type UserReference struct {
	Id       int64
	Username string
	Name     string
}

func ConvertUserReference(user UserReference) *typespb.UserReference {
	return &typespb.UserReference{
		Id:       user.Id,
		Username: user.Username,
		Name:     user.Name,
	}
}
