package types

import "github.com/cluttrdev/gitlab-exporter/protobuf/typespb"

type UserReference struct {
	Id       int64
	Username string
}

func ConvertUserReference(user UserReference) *typespb.UserReference {
	return &typespb.UserReference{
		Id: user.Id,
	}
}