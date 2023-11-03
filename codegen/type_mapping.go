package codegen

var (
	typeMap = map[string]*Mapping{
		"string":     {"string(%v)", "github.com/si3nloong/sqlgen/sequel/types.String(%v)"},
		"[]byte":     {"string(%v)", "github.com/si3nloong/sqlgen/sequel/types.String(%v)"},
		"bool":       {"bool(%v)", "github.com/si3nloong/sqlgen/sequel/types.Bool(%v)"},
		"uint":       {"int64(%v)", "github.com/si3nloong/sqlgen/sequel/types.Integer(%v)"},
		"uint8":      {"int64(%v)", "github.com/si3nloong/sqlgen/sequel/types.Integer(%v)"},
		"uint16":     {"int64(%v)", "github.com/si3nloong/sqlgen/sequel/types.Integer(%v)"},
		"uint32":     {"int64(%v)", "github.com/si3nloong/sqlgen/sequel/types.Integer(%v)"},
		"uint64":     {"int64(%v)", "github.com/si3nloong/sqlgen/sequel/types.Integer(%v)"},
		"int":        {"int64(%v)", "github.com/si3nloong/sqlgen/sequel/types.Integer(%v)"},
		"int8":       {"int64(%v)", "github.com/si3nloong/sqlgen/sequel/types.Integer(%v)"},
		"int16":      {"int64(%v)", "github.com/si3nloong/sqlgen/sequel/types.Integer(%v)"},
		"int32":      {"int64(%v)", "github.com/si3nloong/sqlgen/sequel/types.Integer(%v)"},
		"int64":      {"int64(%v)", "github.com/si3nloong/sqlgen/sequel/types.Integer(%v)"},
		"float32":    {"float64(%v)", "github.com/si3nloong/sqlgen/sequel/types.Float(%v)"},
		"float64":    {"float64(%v)", "github.com/si3nloong/sqlgen/sequel/types.Float(%v)"},
		"time.Time":  {"time.Time(%v)", "(*time.Time)(%v)"},
		"*string":    {"github.com/si3nloong/sqlgen/sequel/types.String(%v)", "github.com/si3nloong/sqlgen/sequel/types.PtrOfString(%v)"},
		"*[]byte":    {"github.com/si3nloong/sqlgen/sequel/types.String(%v)", "github.com/si3nloong/sqlgen/sequel/types.PtrOfString(%v)"},
		"*bool":      {"github.com/si3nloong/sqlgen/sequel/types.Bool(%v)", "github.com/si3nloong/sqlgen/sequel/types.PtrOfBool(%v)"},
		"*uint":      {"github.com/si3nloong/sqlgen/sequel/types.Integer(%v)", "github.com/si3nloong/sqlgen/sequel/types.PtrOfInt(%v)"},
		"*uint8":     {"github.com/si3nloong/sqlgen/sequel/types.Integer(%v)", "github.com/si3nloong/sqlgen/sequel/types.PtrOfInt(%v)"},
		"*uint16":    {"github.com/si3nloong/sqlgen/sequel/types.Integer(%v)", "github.com/si3nloong/sqlgen/sequel/types.PtrOfInt(%v)"},
		"*uint32":    {"github.com/si3nloong/sqlgen/sequel/types.Integer(%v)", "github.com/si3nloong/sqlgen/sequel/types.PtrOfInt(%v)"},
		"*uint64":    {"github.com/si3nloong/sqlgen/sequel/types.Integer(%v)", "github.com/si3nloong/sqlgen/sequel/types.PtrOfInt(%v)"},
		"*int":       {"github.com/si3nloong/sqlgen/sequel/types.Integer(%v)", "github.com/si3nloong/sqlgen/sequel/types.PtrOfInt(%v)"},
		"*int8":      {"github.com/si3nloong/sqlgen/sequel/types.Integer(%v)", "github.com/si3nloong/sqlgen/sequel/types.PtrOfInt(%v)"},
		"*int16":     {"github.com/si3nloong/sqlgen/sequel/types.Integer(%v)", "github.com/si3nloong/sqlgen/sequel/types.PtrOfInt(%v)"},
		"*int32":     {"github.com/si3nloong/sqlgen/sequel/types.Integer(%v)", "github.com/si3nloong/sqlgen/sequel/types.PtrOfInt(%v)"},
		"*int64":     {"github.com/si3nloong/sqlgen/sequel/types.Integer(%v)", "github.com/si3nloong/sqlgen/sequel/types.PtrOfInt(%v)"},
		"*float32":   {"github.com/si3nloong/sqlgen/sequel/types.Float(%v)", "github.com/si3nloong/sqlgen/sequel/types.PtrOfFloat(%v)"},
		"*float64":   {"github.com/si3nloong/sqlgen/sequel/types.Float(%v)", "github.com/si3nloong/sqlgen/sequel/types.PtrOfFloat(%v)"},
		"*time.Time": {"github.com/si3nloong/sqlgen/sequel/types.Time(%v)", "github.com/si3nloong/sqlgen/sequel/types.PtrOfTime(%v)"},
		"[]string":   {"github.com/si3nloong/sqlgen/sequel/encoding.MarshalStringList(%v)", "github.com/si3nloong/sqlgen/sequel/types.StringList(%v)"},
		"[]bool":     {"github.com/si3nloong/sqlgen/sequel/encoding.MarshalBoolList(%v)", "github.com/si3nloong/sqlgen/sequel/types.BoolList(%v)"},
		"[]int":      {"github.com/si3nloong/sqlgen/sequel/encoding.MarshalSignedIntList(%v)", "github.com/si3nloong/sqlgen/sequel/types.IntList(%v)"},
		"[]int8":     {"github.com/si3nloong/sqlgen/sequel/encoding.MarshalSignedIntList(%v)", "github.com/si3nloong/sqlgen/sequel/types.IntList(%v)"},
		"[]int16":    {"github.com/si3nloong/sqlgen/sequel/encoding.MarshalSignedIntList(%v)", "github.com/si3nloong/sqlgen/sequel/types.IntList(%v)"},
		"[]int32":    {"github.com/si3nloong/sqlgen/sequel/encoding.MarshalSignedIntList(%v)", "github.com/si3nloong/sqlgen/sequel/types.IntList(%v)"},
		"[]int64":    {"github.com/si3nloong/sqlgen/sequel/encoding.MarshalSignedIntList(%v)", "github.com/si3nloong/sqlgen/sequel/types.IntList(%v)"},
		"[]uint":     {"github.com/si3nloong/sqlgen/sequel/encoding.MarshalUnsignedIntList(%v)", "github.com/si3nloong/sqlgen/sequel/types.UintList(%v)"},
		"[]uint8":    {"github.com/si3nloong/sqlgen/sequel/encoding.MarshalUnsignedIntList(%v)", "github.com/si3nloong/sqlgen/sequel/types.UintList(%v)"},
		"[]uint16":   {"github.com/si3nloong/sqlgen/sequel/encoding.MarshalUnsignedIntList(%v)", "github.com/si3nloong/sqlgen/sequel/types.UintList(%v)"},
		"[]uint32":   {"github.com/si3nloong/sqlgen/sequel/encoding.MarshalUnsignedIntList(%v)", "github.com/si3nloong/sqlgen/sequel/types.UintList(%v)"},
		"[]uint64":   {"github.com/si3nloong/sqlgen/sequel/encoding.MarshalUnsignedIntList(%v)", "github.com/si3nloong/sqlgen/sequel/types.UintList(%v)"},
		"[]float32":  {"github.com/si3nloong/sqlgen/sequel/encoding.MarshalFloatList(%v)", "github.com/si3nloong/sqlgen/sequel/types.FloatList(%v)"},
		"[]float64":  {"github.com/si3nloong/sqlgen/sequel/encoding.MarshalFloatList(%v)", "github.com/si3nloong/sqlgen/sequel/types.FloatList(%v)"},
	}
)

type Mapping struct {
	Encoder Expr
	Decoder Expr
}
