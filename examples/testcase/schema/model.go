package schema

import (
	"database/sql"
	t "time"

	s "github.com/si3nloong/sqlgen/sequel"
)

type LongText string

type A struct {
	s.Table   `sql:"Apple"`
	ID        string
	Text      LongText
	CreatedAt t.Time
}

type B struct {
	// go:staticcheck
	// n, p      schema.Name `sql:"Boy"`
	ID        string
	CreatedAt t.Time
}

type C struct {
	ID int64 `sql:",pk"`
}

type D struct {
	ID sql.NullString `sql:",pk"`
}
