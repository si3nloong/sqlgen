package inlinenested

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type Str string

type NestedModel struct {
	Nested struct {
		Time    time.Time
		Decimal decimal.Decimal
		Bool    bool
	}
}

type NestedModelWithTag struct {
	Nested struct {
		Time    time.Time
		Decimal decimal.Decimal
		// DeepNested struct {
		// 	Str      string
		// 	Duration time.Duration
		// }
		Bool bool `json:"bool"`
		Str  Str  `json:"str"`
	}
}

type DeepNestedModel struct {
	Nested struct {
		Byte    sql.NullByte `sql:"byte" json:"byte"`
		Time    time.Time
		Decimal decimal.Decimal
		// DeepNested struct {
		// 	Str      string
		// 	Duration time.Duration
		// }
		Bool bool
	}
}
