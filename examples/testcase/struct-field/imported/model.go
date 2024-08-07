package imported

import (
	"database/sql"

	"github.com/google/uuid"
)

type Model struct {
	Str      sql.NullString
	Bool     sql.NullBool
	RawBytes sql.RawBytes
	Int16    sql.NullInt16
	Int32    sql.NullInt32
	Int64    sql.NullInt64
	Float64  sql.NullFloat64
	Time     sql.NullTime
}

type Some struct {
	ID uuid.UUID
}
