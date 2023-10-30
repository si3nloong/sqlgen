package alias

import (
	"database/sql"
	s "database/sql"
)

type AliasStruct struct {
	// sql.NullString
	a, B, c float64
	Empty
	pk
	// PtrHeader *aliasStr
	Header  aliasStr
	Raw     sql.RawBytes
	Text    customStr
	NullStr s.NullString
	model
	private string
}

type pk struct {
	ID int64 `sql:"Id,primary_key"`
}

type model struct {
	Created DT
	Updated DT
}
