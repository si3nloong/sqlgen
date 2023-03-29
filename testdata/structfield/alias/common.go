package alias

import (
	s "database/sql"
	t "time"
)

type customStr string

type aliasStr = customStr

type DT t.Time

type AliasStruct struct {
	ID      int64 `sql:",pk"`
	Header  aliasStr
	Text    customStr
	Created DT
	NullStr s.NullString
}
