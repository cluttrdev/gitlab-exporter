package messages

import (
	"go.cluttr.dev/gitlab-exporter/exporter/internal/types"
	"go.cluttr.dev/gitlab-exporter/protobuf/typespb"
)

func NewUserReference(user types.UserReference) *typespb.UserReference {
	return &typespb.UserReference{
		Id:       user.Id,
		Username: user.Username,
		Name:     user.Name,
	}
}

func NewUserReferences(users []types.UserReference) []*typespb.UserReference {
	if users == nil {
		return nil
	}
	us := make([]*typespb.UserReference, 0, len(users))
	for _, user := range users {
		us = append(us, NewUserReference(user))
	}
	return us
}
