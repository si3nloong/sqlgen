package slice

type customStr string

type Slice struct {
	priv          int64
	ID            uint64 `sql:",pk,auto_increment"`
	BoolList      []bool
	StrList       []string
	CustomStrList []customStr
	IntList       []int
	Int8List      []int8
	Int16List     []int16
	Int32List     []int32
	Int64List     []int64
	UintList      []uint
	Uint8List     []uint8
	Uint16List    []uint16
	Uint32List    []uint32
	Uint64List    []uint64
	F32List       []float32
	F64List       []float64
	// TimeList      []t.Time
}
