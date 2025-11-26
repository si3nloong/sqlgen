package inlinenested

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
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
		Bool    bool `json:"bool"`
		Str     Str  `json:"str"`
	}
}

type embedded struct {
	Flag  bool
	int64 int
}
type DeepNestedModel struct {
	Nested struct {
		Byte       sql.NullByte `sql:"byte" json:"byte"`
		Time       time.Time
		Decimal    decimal.Decimal
		number     float64
		embedded   `json:",inline"`
		DeepNested struct {
			time.Time
			Str      string
			Duration time.Duration
			Lvl3     struct {
				Float float64 `json:"float_64" xml:"float_64"`
				Bool  bool
				UUID  uuid.UUID `json:"uuid" xml:"uuid"`
			}
		}
		flag bool
	}
}
