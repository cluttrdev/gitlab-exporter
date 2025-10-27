package messages

import (
	"go.cluttr.dev/gitlab-exporter/exporter/internal/types"
	"go.cluttr.dev/gitlab-exporter/exporter/protobuf/typespb"
)

func NewUserReference(user types.UserReference) *typespb.UserReference {
	return &typespb.UserReference{
		Id:       user.Id,
		Username: user.Username,
		Name:     user.Name,
	}
}
