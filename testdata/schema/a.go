package schema

import (
	"time"

	"github.com/si3nloong/sqlgen/sql/schema"
)

type A struct {
	schema.Name `sql:"StructA"`
	ID          string
	CreatedAt   time.Time
}
