package readonly

import (
	"time"

	"github.com/shopspring/decimal"
)

// +sql:readonly
type A struct {
	ID   string `sql:",pk"`
	Time time.Time
	Dec  decimal.Decimal
}

// +sql:ignore
type B struct {
	ID int64 `sql:",pk"`
}

type private struct {
	A
	B
}
