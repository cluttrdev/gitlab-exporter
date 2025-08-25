package graphql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

type TimeRangeOptions struct {
	UpdatedAfter  *time.Time
	UpdatedBefore *time.Time
}

const (
	GlobalIdPrefix             = "gid://gitlab/"
	GlobalIdMergeRequestPrefix = GlobalIdPrefix + "MergeRequest/"
	GlobalIdMilestonePrefix    = GlobalIdPrefix + "Milestone/"
	GlobalIdNotePrefix         = GlobalIdPrefix + "Note/"
	GlobalIdPipelinePrefix     = GlobalIdPrefix + "Ci::Pipeline/"
	GlobalIdProjectPrefix      = GlobalIdPrefix + "Project/"
	GlobalIdUserPrefix         = GlobalIdPrefix + "User/"
	GlobalIdIssuePrefix        = GlobalIdPrefix + "Issue/"
)

func FormatId(id int64, prefix string) string {
	return prefix + strconv.FormatInt(id, 10)
}

func ParseId(s string, prefix string) (int64, error) {
	return strconv.ParseInt(strings.TrimPrefix(s, prefix), 10, 64)
}

func parseNamespaceId(s string) (int64, error) {
	var s_ string

	prefixes := []string{
		GlobalIdPrefix + "Group/",
		GlobalIdPrefix + "Namespaces::UserNamespace/",
		GlobalIdPrefix + "Namespaces::ProjectNamespace/",
	}

	for _, prefix := range prefixes {
		s_ = strings.TrimPrefix(s, prefix)
		if len(s_) < len(s) {
			return strconv.ParseInt(s_, 10, 64)
		}
	}

	return strconv.ParseInt(s, 10, 64)
}

func parseJobId(s string) (int64, error) {
	var s_ string

	prefixes := []string{
		GlobalIdPrefix + "Ci::Build/",
		GlobalIdPrefix + "Ci::Bridge/",
		GlobalIdPrefix + "GenericCommitStatus/",
		GlobalIdPrefix + "CommitStatus/",
	}

	for _, prefix := range prefixes {
		s_ = strings.TrimPrefix(s, prefix)
		if len(s_) < len(s) {
			return strconv.ParseInt(s_, 10, 64)
		}
	}

	return strconv.ParseInt(s, 10, 64)
}

func valOrZero[T any](t *T) T {
	var v T
	if t != nil {
		v = *t
	}
	return v
}

func valOr[T any](t *T, v T) T {
	if t != nil {
		return *t
	}
	return v
}

func ptr[T any](v T) *T {
	return &v
}

func handleError(err error, query string, attrs ...slog.Attr) error {
	var errList gqlerror.List
	if errors.As(err, &errList) {
		attrs_ := make([]slog.Attr, 1, 1+len(attrs)+len(errList))

		const msg = "gqlerror"
		attrs_[0] = slog.String("query", query)

		attrs_ = append(attrs_, attrs...)

		for i, e := range errList {
			attrs_ = append(attrs_, slog.Group(
				fmt.Sprintf("error[%d]", i),
				slog.String("path", e.Path.String()),
				slog.String("message", e.Message),
			))
		}

		slog.LogAttrs(context.Background(), slog.LevelError, msg, attrs_...)
		return nil
	}
	return err
}
