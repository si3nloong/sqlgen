package schema

import (
	"database/sql"
	t "time"
)

type A struct {
	// schema.Name `sql:"Apple"`
	ID        string
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
