package policy

import (
	"github.com/Masterminds/semver"
	"github.com/ryanuber/go-glob"
	"github.com/weaveworks/flux/image"
	"strings"
)

const (
	globPrefix   = "glob:"
	semverPrefix = "semver:"
)

var (
	// PatternAll matches everything. To the contrary of `semver:*`
	// which only matches tags that are valid semvers.
	PatternAll    = NewPattern("glob:*")
	PatternLatest = NewPattern("glob:latest")
)

// Pattern provides an interface to match arbitrary strings.
//
// While the string representation aims to prefix patterns with
// their respective type, it can also be omitted and defaults then
// to glob matching.
type Pattern interface {
	Matches(tag string) bool
	// Prefixed string representation
	String() string
	ImageLess() image.SortLessFunc
	Valid() bool
}

type GlobPattern string

type SemverPattern struct {
	// Pattern without prefix
	pattern     string
	constraints *semver.Constraints
}

func NewPattern(pattern string) Pattern {
	if strings.HasPrefix(pattern, semverPrefix) {
		pattern = strings.TrimPrefix(pattern, semverPrefix)
		c, _ := semver.NewConstraint(pattern)
		return SemverPattern{pattern, c}
	}
	return GlobPattern(strings.TrimPrefix(pattern, globPrefix))
}

func (g GlobPattern) Matches(tag string) bool {
	return glob.Glob(string(g), tag)
}

func (g GlobPattern) String() string {
	return globPrefix + string(g)
}

func (g GlobPattern) ImageLess() image.SortLessFunc {
	return image.ByCreatedDesc
}

func (g GlobPattern) Valid() bool {
	return true
}

func (s SemverPattern) Matches(tag string) bool {
	v, err := semver.NewVersion(tag)
	if err != nil {
		return false
	}

	// Allow `*` as match-all for valid semvers
	if s.pattern == semverPrefix+"*" {
		return true
	}
	if s.constraints == nil {
		// Invalid constraints match anything
		return true
	}
	return s.constraints.Check(v)
}

func (s SemverPattern) String() string {
	return semverPrefix + s.pattern
}

func (s SemverPattern) ImageLess() image.SortLessFunc {
	return image.BySemverTagDesc
}

func (s SemverPattern) Valid() bool {
	return s.constraints != nil
}
