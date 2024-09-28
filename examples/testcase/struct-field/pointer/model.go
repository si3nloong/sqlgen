package pointer

import "time"

type Ptr struct {
	ID     int64 `sql:",pk,auto_increment"`
	Str    *string
	Bytes  *[]byte
	Bool   *bool
	Int    *int
	Int8   *int8
	Int16  *int16
	Int32  *int32
	Int64  *int64
	Uint   *uint
	Uint8  *uint8
	Uint16 *uint16
	Uint32 *uint32
	Uint64 *uint64
	F32    *float32
	F64    *float64
	Time   *time.Time `sql:",size:6"`
	Nested *nested
	*embeded
}

type nested struct {
	ID *int64
}

type embeded struct {
	EmbededTime *time.Time
}
