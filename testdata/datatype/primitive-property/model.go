package primitivestruct

import "time"

type Primitive struct {
	Str string
	// Rune   rune
	// Char   byte
	Bytes  []byte
	Bool   bool
	Int    int
	Int8   int8
	Int16  int16
	Int32  int32
	Int64  int64
	Uint   uint
	Uint8  uint8
	Uint16 uint16
	Uint32 uint32
	Uint64 uint64
	F32    float32
	F64    float64
	Time   time.Time
}

func (m *Primitive) AddrsMap() []any {
	return []any{&m.Str, &m.Bytes, &m.Bool, &m.Int, &m.Int8, &m.Int16, &m.Int32, &m.Int64, &m.Uint, &m.Uint8, &m.Uint16, &m.Uint32, &m.Uint64, &m.F32, &m.F64, &m.Time}
}
