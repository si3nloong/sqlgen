package decode

import "time"

type LongText string

type Model struct {
	ID   uint8    `sql:",decode:github.com/si3nloong/sqlgen/examples/testencoding.UnmarshalAny"`
	Text LongText `sql:",decode:github.com/si3nloong/sqlgen/examples/testencoding.UnmarshalString"`
	T    time.Time
}
