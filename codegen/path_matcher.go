package codegen

import (
	"regexp"
)

type Matcher interface {
	Match(v string) bool
}

type FileMatcher map[string]struct{}

func (f FileMatcher) Match(v string) bool {
	_, ok := f[v]
	return ok
}

type RegexpMatcher struct {
	*regexp.Regexp
}

func (r RegexpMatcher) Match(v string) bool {
	return r.MatchString(v)
}

type EmptyMatcher struct{}

func (*EmptyMatcher) Match(v string) bool {
	return true
}
