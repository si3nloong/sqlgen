package alias

import (
	s "database/sql"
	t "time"
)

type customStr string

type aliasStr = customStr

type DT = t.Time

type AliasStruct struct {
	Header  aliasStr
	Text    customStr
	Created DT
	NullStr s.NullString
}
