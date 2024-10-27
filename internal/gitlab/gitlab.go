package gitlab

import "github.com/xanzy/go-gitlab"

func Ptr[T any](v T) *T {
	return gitlab.Ptr(v)
}
