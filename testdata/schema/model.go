package schema

import (
	"database/sql"
	"time"

	s "github.com/si3nloong/sqlgen/sql/schema"
)

type A struct {
	s.Name    `sql:"Apple"`
	ID        string
	CreatedAt time.Time
}

type B struct {
	// go:staticcheck
	n, p      s.Name `sql:"Boy"`
	ID        string
	CreatedAt time.Time
}

type C struct {
	ID int64 `sql:",pk"`
}

type D struct {
	ID sql.NullString `sql:",pk"`
}
