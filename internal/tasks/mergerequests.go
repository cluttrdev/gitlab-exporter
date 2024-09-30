package tasks

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-exporter/internal/exporter"
	"github.com/cluttrdev/gitlab-exporter/internal/types"
	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

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

type ExportProjectMergeRequestsOptions struct {
	ProjectID        int
	MergeRequestIIDs []int

	ExportNoteEvents bool
}

func ExportProjectMergeRequests(ctx context.Context, glab *gitlab.Client, exp *exporter.Exporter, opt ExportProjectMergeRequestsOptions) error {
	mergerequests := make([]*typespb.MergeRequest, 0, len(opt.MergeRequestIIDs))
	mrNoteEvents := make([]*typespb.MergeRequestNoteEvent, 0, len(opt.MergeRequestIIDs))

	_opt := gitlab.GetMergeRequestsOptions{}
	for _, iid := range opt.MergeRequestIIDs {
		mr, _, err := glab.MergeRequests.GetMergeRequest(opt.ProjectID, iid, &_opt, gitlab.WithContext(ctx))
		if err != nil {
			if errors.Is(err, context.Canceled) {
				break
			}
			slog.Error("error fetching merge request", "project_id", opt.ProjectID, "iid", iid, "error", err)
			continue
		}

		mergerequests = append(mergerequests, types.ConvertMergeRequest(mr))

		if opt.ExportNoteEvents {
			notes, err := FetchMergeRequestNotes(ctx, glab, opt.ProjectID, iid)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					break
				}
				slog.Error("error fetching merge request note events", "project_id", opt.ProjectID, "iid", iid)
				continue
			}

			for _, note := range notes {
				if ev := types.ConvertToMergeRequestNoteEvent(note); ev != nil {
					mrNoteEvents = append(mrNoteEvents, ev)
				}
			}
		}
	}

	if err := exp.ExportMergeRequests(ctx, mergerequests); err != nil {
		return fmt.Errorf("export merge requests: %w", err)
	}

	if err := exp.ExportMergeRequestNoteEvents(ctx, mrNoteEvents); err != nil {
		return fmt.Errorf("export merge request note events: %w", err)
	}

	return nil
}
