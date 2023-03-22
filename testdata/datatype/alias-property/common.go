package aliasproperty

import t "time"

type customStr string

type aliasStr = customStr

type AliasStruct struct {
	ID       int64 `sql:"-"`
	Name     aliasStr
	DateTime t.Time
}
