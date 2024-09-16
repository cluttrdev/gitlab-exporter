package tasks

import (
	"context"

	gitlab "github.com/xanzy/go-gitlab"
)

func FetchProjectMergeRequest(ctx context.Context, glab *gitlab.Client, pid int64, iid int64) (*gitlab.MergeRequest, error) {
	opt := gitlab.GetMergeRequestsOptions{}

	mr, _, err := glab.MergeRequests.GetMergeRequest(int(pid), int(iid), &opt, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return mr, nil
}

func FetchMergeRequestNotes(ctx context.Context, glab *gitlab.Client, pid interface{}, iid int) ([]*gitlab.Note, error) {
	opts := gitlab.ListMergeRequestNotesOptions{
		ListOptions: gitlab.ListOptions{
			Pagination: "keyset",
			PerPage:    100,
			OrderBy:    "created_at",
			Sort:       "desc",
		},
	}

	options := []gitlab.RequestOptionFunc{
		gitlab.WithContext(ctx),
	}

	var notes []*gitlab.Note
	for {
		ns, resp, err := glab.Notes.ListMergeRequestNotes(pid, iid, &opts, options...)
		if err != nil {
			return nil, err
		}

		notes = append(notes, ns...)

		if resp.NextLink == "" {
			break
		}

		options = []gitlab.RequestOptionFunc{
			gitlab.WithContext(ctx),
			gitlab.WithKeysetPaginationParameters(resp.NextLink),
		}
	}

	return notes, nil
}
