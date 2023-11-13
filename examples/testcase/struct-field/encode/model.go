package encode

import "time"

type LongText string

type Model struct {
	ID   uint8    `sql:",encode:github.com/si3nloong/sqlgen/examples/testencoding.MarshalAny"`
	Text LongText `sql:",encode:github.com/si3nloong/sqlgen/examples/testencoding.MarshalGenericString"`
	T    time.Time
}
