package gitlab

import "gitlab.com/gitlab-org/api/client-go"

func Ptr[T any](v T) *T {
	return gitlab.Ptr(v)
}
