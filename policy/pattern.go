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
	// PatternAll matches everything.
	PatternAll    = NewPattern(globPrefix + "*")
	PatternLatest = NewPattern(globPrefix + "latest")
)

// Pattern provides an interface to match image tags.
//
// While the string representation aims to prefix patterns with
// their respective type, it can also be omitted and defaults then
// to glob matching.
type Pattern interface {
	Matches(tag string) bool
	// String returns the prefixed string representation.
	String() string
	// ImageNewerFunc returns a function to compare image newness.
	ImageNewerFunc() image.SortLessFunc
	// Valid returns true if the pattern is considered valid.
	Valid() bool
}

type GlobPattern string

// SemverPattern matches by semantic versioning.
// See https://semver.org/
type SemverPattern struct {
	pattern     string // pattern without prefix
	constraints *semver.Constraints
}

// NewPattern instantiates a Pattern according to the prefix
// it finds. The prefix is optional and defaults to `glob`.
//
// `semver:*` matches only tags that are valid according to
// semantic versioning while `glob:*` matches every single tag.
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

func (g GlobPattern) ImageNewerFunc() image.SortLessFunc {
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

	// Allow `*` as match-all for valid semver tags
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

func (s SemverPattern) ImageNewerFunc() image.SortLessFunc {
	return image.BySemverTagDesc
}

func (s SemverPattern) Valid() bool {
	return s.constraints != nil
}
