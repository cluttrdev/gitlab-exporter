package gitlab

import "gitlab.com/gitlab-org/api/client-go"

var ErrNotFound error = gitlab.ErrNotFound

func Ptr[T any](v T) *T {
	return gitlab.Ptr(v)
}
