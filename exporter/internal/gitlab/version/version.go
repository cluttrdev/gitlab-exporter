package version

import (
	"fmt"
	"strconv"
	"strings"
)

// GitLabVersion represents a parsed GitLab version.
type GitLabVersion struct {
	Major int
	Minor int
	Patch int
	Raw   string
}

// IsZero reports whether v has not been set. A zero value indicates that
// version detection (or fallback assignment) has not yet occurred.
func (v GitLabVersion) IsZero() bool {
	return v.Major == 0 && v.Minor == 0 && v.Patch == 0 && v.Raw == ""
}

// Parse parses a GitLab version string into a GitLabVersion.
//
// Parsing is tolerant: anything after the second numeric component is
// captured into Patch on a best-effort basis (the leading digits are
// parsed; any suffix like "-ee", "-pre", or unexpected segments is
// discarded for the numeric fields but preserved in Raw).
//
// Examples of accepted inputs:
//
//	"17"          -> {Major: 17}
//	"17.0"        -> {Major: 17, Minor: 0}
//	"17.4.2"      -> {Major: 17, Minor: 4, Patch: 2}
//	"17.4.2-ee"   -> {Major: 17, Minor: 4, Patch: 2}
//	"16.10.0-pre" -> {Major: 16, Minor: 10, Patch: 0}
//
// Empty strings and inputs without a leading numeric component return an
// error.
func Parse(s string) (GitLabVersion, error) {
	if s == "" {
		return GitLabVersion{}, fmt.Errorf("version: parse: empty string")
	}

	v := GitLabVersion{Raw: s}

	// Split on '.' for the numeric components; the third component may
	// carry a non-numeric suffix (e.g. "2-ee") that we trim before parsing.
	parts := strings.SplitN(s, ".", 3)

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return GitLabVersion{}, fmt.Errorf("version: parse: major %q: %w", parts[0], err)
	}
	v.Major = major

	if len(parts) >= 2 {
		minor, err := strconv.Atoi(leadingDigits(parts[1]))
		if err != nil {
			return GitLabVersion{}, fmt.Errorf("version: parse: minor %q: %w", parts[1], err)
		}
		v.Minor = minor
	}

	if len(parts) >= 3 {
		// Patch can carry suffixes like "-ee" or "-pre"; take leading
		// digits only. Empty leading digits (e.g. "0-pre" -> "0", but
		// "-pre" alone -> "") are tolerated and leave Patch at zero.
		digits := leadingDigits(parts[2])
		if digits != "" {
			patch, err := strconv.Atoi(digits)
			if err != nil {
				return GitLabVersion{}, fmt.Errorf("version: parse: patch %q: %w", parts[2], err)
			}
			v.Patch = patch
		}
	}

	return v, nil
}

// Compare returns an integer comparing two GitLabVersions. The Raw field is
// ignored.
func Compare(a, b GitLabVersion) int {
	if c := compareInt(a.Major, b.Major); c != 0 {
		return c
	} else if c := compareInt(a.Minor, b.Minor); c != 0 {
		return c
	} else {
		return compareInt(a.Patch, b.Patch)
	}
}

// leadingDigits returns the longest prefix of s consisting of ASCII
// digits.
func leadingDigits(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return s[:i]
		}
	}
	return s
}

func compareInt(a, b int) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	} else {
		return 0
	}
}
