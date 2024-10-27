package graphql

import (
	"strconv"
	"strings"
)

const (
	GlobalIdPrefix             = "gid://gitlab/"
	GlobalIdProjectPrefix      = GlobalIdPrefix + "Project/"
	GlobalIdGroupPrefix        = GlobalIdPrefix + "Group/"
	GlobalIdPipelinePrefix     = GlobalIdPrefix + "Ci::Pipeline/"
	GlobalIdJobBuildPrefix     = GlobalIdPrefix + "Ci::Build/"
	GlobalIdJobBridgePrefix    = GlobalIdPrefix + "Ci::Bridge/"
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
