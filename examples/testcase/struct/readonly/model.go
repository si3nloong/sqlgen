package readonly

import (
	"time"

	"github.com/shopspring/decimal"
)

// +sqlgen:readonly
type A struct {
	ID   string `sql:",pk"`
	Time time.Time
	Dec  decimal.Decimal
}

type B struct {
	ID int64 `sql:",pk"`
}
