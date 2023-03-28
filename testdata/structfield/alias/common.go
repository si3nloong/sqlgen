package alias

import t "time"

type customStr string

type aliasStr = customStr

type DT t.Time

type AliasStruct struct {
	ID      int64 `sql:"-"`
	Header  aliasStr
	Text    customStr
	Created DT
}
