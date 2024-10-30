package graphql

import (
	"strconv"
	"strings"
)

const (
	GlobalIdPrefix             = "gid://gitlab/"
	GlobalIdProjectPrefix      = GlobalIdPrefix + "Project/"
	GlobalIdPipelinePrefix     = GlobalIdPrefix + "Ci::Pipeline/"
	GlobalIdMergeRequestPrefix = GlobalIdPrefix + "MergeRequest/"
	GlobalIdMilestonePrefix    = GlobalIdPrefix + "Milestone/"
	GlobalIdNotePrefix         = GlobalIdPrefix + "Note/"
	GlobalIdUserPrefix         = GlobalIdPrefix + "User/"
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

func ptr[T any](v T) *T {
	return &v
}
