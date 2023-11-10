package aliasstruct

import (
	"database/sql"

	"cloud.google.com/go/civil"
)

// This should generate
type A civil.DateTime // Type Definition

func (a A) TableName() string {
	return ""
}

// This shouldn't generate
type B = civil.DateTime

// This should generate
type C sql.NullString
