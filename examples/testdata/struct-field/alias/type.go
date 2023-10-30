package alias

import (
	t "time"
)

type Empty struct{}

type customStr string

type aliasStr = customStr

type DT = t.Time

// This will exclude
type C struct {
	ID int64
}
