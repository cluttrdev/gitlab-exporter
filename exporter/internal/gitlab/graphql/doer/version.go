package doer

import (
	"fmt"
	"strconv"
	"strings"
)

// ServerVersion is a parsed semver-like version string.
type ServerVersion struct {
	Major int
	Minor int
	Patch int

	raw string
}

// ParseServerVersion parses a dotted-decimal version string into a
// ServerVersion.
//
// Empty strings or inputs without a leading numeric component return an
// error. Each of the three numeric components is parsed from its leading
// run of ASCII digits; missing components default to zero.
func ParseServerVersion(s string) (ServerVersion, error) {
	if s == "" {
		return ServerVersion{}, fmt.Errorf("empty string")
	}

	v := ServerVersion{raw: s}

	parts := strings.SplitN(s, ".", 3)

	major, err := strconv.Atoi(leadingDigits(parts[0]))
	if err != nil || leadingDigits(parts[0]) == "" {
		return ServerVersion{}, fmt.Errorf("invalid major version %q: %q", s, parts[0])
	}
	v.Major = major

	if len(parts) >= 2 {
		minorDigits := leadingDigits(parts[1])
		if minorDigits == "" {
			return ServerVersion{}, fmt.Errorf("invalid minor version %q: %q", s, parts[1])
		}
		minor, err := strconv.Atoi(minorDigits)
		if err != nil {
			return ServerVersion{}, fmt.Errorf("minor %q: %w", parts[1], err)
		}
		v.Minor = minor
	}

	if len(parts) >= 3 {
		patchDigits := leadingDigits(parts[2])
		// Patch may be entirely non-numeric (e.g. "-pre" alone) - leave at zero.
		if patchDigits != "" {
			patch, err := strconv.Atoi(patchDigits)
			if err != nil {
				return ServerVersion{}, fmt.Errorf("patch %q: %w", parts[2], err)
			}
			v.Patch = patch
		}
	}

	return v, nil
}

// compareServerVersion returns -1, 0, or 1 if a is less than, equal to, or
// greater than b. Comparison is lexicographic across (Major, Minor, Patch);
// the raw field is ignored.
func compareServerVersion(a, b ServerVersion) int {
	if c := cmpInt(a.Major, b.Major); c != 0 {
		return c
	}
	if c := cmpInt(a.Minor, b.Minor); c != 0 {
		return c
	}
	return cmpInt(a.Patch, b.Patch)
}

// complexityHonorsConditionalFragments is the first GitLab version whose
// query-complexity analyzer respects @include / @skip on fragment spreads and
// inline fragments.
//
// graphql-ruby's analysis visitor sets its "skipping" flag from a node's
// directives, but until v2.3.0 it only did so in on_field. on_fragment_spread
// and on_inline_fragment gained the same handling in 2.3.0, so on older
// releases QueryComplexity charges for conditional fragments the server will
// never resolve. GitLab 16.11 ships 2.2.5 (affected); 17.0 ships 2.3.2 (fixed).
const complexityHonorsConditionalFragments = "17.0"

// Gate carries a fixed server version against which @minVersion / @maxVersion
// constraints are evaluated. Construct one via NewGate.
type Gate struct {
	server ServerVersion
}

// MaterializesConditionals reports whether the client must evaluate
// @include / @skip itself before sending a query, because this server's
// complexity analyzer would otherwise charge for conditional fragments it
// does not resolve. See complexityHonorsConditionalFragments.
func (g *Gate) MaterializesConditionals() bool {
	if g == nil {
		return false
	}
	honors, err := g.Allows(complexityHonorsConditionalFragments, "")
	if err != nil {
		// The bound is a compile-time constant, so a parse failure here is
		// impossible. Fall back to leaving the directives for the server.
		return false
	}
	return !honors
}

// NewGate parses a server version string and returns a Gate suitable for use
// by the AST rewriter and the VersionGatedDoer.
func NewGate(server string) (*Gate, error) {
	v, err := ParseServerVersion(server)
	if err != nil {
		return nil, fmt.Errorf("parse server version: %w", err)
	}
	return &Gate{server: v}, nil
}

// Allows reports whether the gate's server version falls within the inclusive
// range [minV, maxV]. Empty bounds mean "no constraint" on that side. An
// error is returned if either non-empty bound is malformed.
//
// Semantics:
//
//	minV == ""           -> no lower bound
//	maxV == ""           -> no upper bound
//	server >= minV       -> lower bound satisfied
//	server <= maxV       -> upper bound satisfied
//
// Both bounds are evaluated using compareServerVersion (Major.Minor.Patch
// only); pre-release suffixes are ignored.
func (g *Gate) Allows(minV, maxV string) (bool, error) {
	if g == nil {
		return false, fmt.Errorf("gate is nil")
	}

	if minV != "" {
		min, err := ParseServerVersion(minV)
		if err != nil {
			return false, fmt.Errorf("parse minVersion %q: %w", minV, err)
		}
		if compareServerVersion(g.server, min) < 0 {
			return false, nil
		}
	}

	if maxV != "" {
		max, err := ParseServerVersion(maxV)
		if err != nil {
			return false, fmt.Errorf("parse maxVersion %q: %w", maxV, err)
		}
		if compareServerVersion(g.server, max) > 0 {
			return false, nil
		}
	}

	return true, nil
}

// leadingDigits returns the longest prefix of s consisting of ASCII digits.
func leadingDigits(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return s[:i]
		}
	}
	return s
}

func cmpInt(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}
