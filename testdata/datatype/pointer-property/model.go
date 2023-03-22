package pointerproperty

import "time"

type Ptr struct {
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
	Time   *time.Time
}
