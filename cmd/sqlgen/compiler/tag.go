package compiler

import "strings"

const (
	TagOptionAutoIncrement = "auto_increment"
	TagOptionPKAlias       = "pk"
	TagOptionPK            = "primary_key"
	TagOptionFKAlias       = "fk"
	TagOptionFK            = "foreign_key"
	TagOptionUnsigned      = "unsigned"
	TagOptionReadonly      = "readonly"
	TagOptionSize          = "size"
	TagOptionUnique        = "unique"
	TagOptionIndex         = "index"
)

type structTag struct {
	// values []string
	pos map[string]int
}

func (s *structTag) hasOpts(keys ...string) bool {
	for _, k := range keys {
		if _, ok := s.pos[k]; ok {
			return true
		}
	}
	return false
}

func parseTag(tagVal []string) *structTag {
	tag := structTag{pos: make(map[string]int)}
	for _, v := range tagVal {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		paths := strings.SplitN(v, ":", 2)
		tag.pos[paths[0]] = 0
	}
	return &tag
}
