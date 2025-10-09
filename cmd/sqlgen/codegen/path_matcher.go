package codegen

import (
	"regexp"
	"strings"
)

type Matcher interface {
	Match(v string) bool
}

type EmptyMatcher struct{}

func (*EmptyMatcher) Match(v string) bool {
	return true
}

type FileMatcher map[string]struct{}

func (f FileMatcher) Match(v string) bool {
	_, ok := f[v]
	return ok
}

type FolderMatcher string

func (f FolderMatcher) Match(v string) bool {
	return strings.HasPrefix(v, (string)(f))
}

type RegexMatcher struct {
	*regexp.Regexp
}

func (r *RegexMatcher) Match(v string) bool {
	return r.MatchString(v)
}
