package gitlab

import (
	"context"

	_gitlab "github.com/xanzy/go-gitlab"
)

func ListProjectMergeRequests(ctx context.Context, glab *_gitlab.Client, pid int64, opt _gitlab.ListProjectMergeRequestsOptions, yield func(p []*_gitlab.MergeRequest) bool) error {
	opt.ListOptions.Pagination = "keyset"
	if opt.ListOptions.OrderBy == "" {
		opt.ListOptions.OrderBy = "updated_at"
	}
	if opt.ListOptions.Sort == "" {
		opt.ListOptions.Sort = "desc"
	}

	options := []_gitlab.RequestOptionFunc{
		_gitlab.WithContext(ctx),
	}

	for {
		mrs, resp, err := glab.MergeRequests.ListProjectMergeRequests(int(pid), &opt, options...)
		if err != nil {
			return err
		}

		if !yield(mrs) {
			break
		}

		if resp.NextLink == "" {
			break
		}

		options = []_gitlab.RequestOptionFunc{
			_gitlab.WithContext(ctx),
			_gitlab.WithKeysetPaginationParameters(resp.NextLink),
		}
	}

	return nil
}
