package codegen

import (
	"go/types"
	"path/filepath"
	"strings"
)

var (
	typeMap = map[string]Mapping{
		"string":     {"string", "github.com/si3nloong/sqlgen/sql/types.String"},
		"[]byte":     {"string", "github.com/si3nloong/sqlgen/sql/types.String"},
		"bool":       {"bool", "github.com/si3nloong/sqlgen/sql/types.Bool"},
		"uint":       {"int64", "github.com/si3nloong/sqlgen/sql/types.Integer"},
		"uint8":      {"int64", "github.com/si3nloong/sqlgen/sql/types.Integer"},
		"uint16":     {"int64", "github.com/si3nloong/sqlgen/sql/types.Integer"},
		"uint32":     {"int64", "github.com/si3nloong/sqlgen/sql/types.Integer"},
		"uint64":     {"int64", "github.com/si3nloong/sqlgen/sql/types.Integer"},
		"int":        {"int64", "github.com/si3nloong/sqlgen/sql/types.Integer"},
		"int8":       {"int64", "github.com/si3nloong/sqlgen/sql/types.Integer"},
		"int16":      {"int64", "github.com/si3nloong/sqlgen/sql/types.Integer"},
		"int32":      {"int64", "github.com/si3nloong/sqlgen/sql/types.Integer"},
		"int64":      {"int64", "github.com/si3nloong/sqlgen/sql/types.Integer"},
		"float32":    {"float64", "github.com/si3nloong/sqlgen/sql/types.Float"},
		"float64":    {"float64", "github.com/si3nloong/sqlgen/sql/types.Float"},
		"*string":    {"github.com/si3nloong/sqlgen/sql/types.String", "github.com/si3nloong/sqlgen/sql/types.PtrOfString"},
		"*[]byte":    {"github.com/si3nloong/sqlgen/sql/types.String", "github.com/si3nloong/sqlgen/sql/types.PtrOfString"},
		"*bool":      {"github.com/si3nloong/sqlgen/sql/types.Bool", "github.com/si3nloong/sqlgen/sql/types.PtrOfBool"},
		"*uint":      {"github.com/si3nloong/sqlgen/sql/types.Integer", "github.com/si3nloong/sqlgen/sql/types.PtrOfInt"},
		"*uint8":     {"github.com/si3nloong/sqlgen/sql/types.Integer", "github.com/si3nloong/sqlgen/sql/types.PtrOfInt"},
		"*uint16":    {"github.com/si3nloong/sqlgen/sql/types.Integer", "github.com/si3nloong/sqlgen/sql/types.PtrOfInt"},
		"*uint32":    {"github.com/si3nloong/sqlgen/sql/types.Integer", "github.com/si3nloong/sqlgen/sql/types.PtrOfInt"},
		"*uint64":    {"github.com/si3nloong/sqlgen/sql/types.Integer", "github.com/si3nloong/sqlgen/sql/types.PtrOfInt"},
		"*int":       {"github.com/si3nloong/sqlgen/sql/types.Integer", "github.com/si3nloong/sqlgen/sql/types.PtrOfInt"},
		"*int8":      {"github.com/si3nloong/sqlgen/sql/types.Integer", "github.com/si3nloong/sqlgen/sql/types.PtrOfInt"},
		"*int16":     {"github.com/si3nloong/sqlgen/sql/types.Integer", "github.com/si3nloong/sqlgen/sql/types.PtrOfInt"},
		"*int32":     {"github.com/si3nloong/sqlgen/sql/types.Integer", "github.com/si3nloong/sqlgen/sql/types.PtrOfInt"},
		"*int64":     {"github.com/si3nloong/sqlgen/sql/types.Integer", "github.com/si3nloong/sqlgen/sql/types.PtrOfInt"},
		"*float32":   {"github.com/si3nloong/sqlgen/sql/types.Float", "github.com/si3nloong/sqlgen/sql/types.PtrOfFloat"},
		"*float64":   {"github.com/si3nloong/sqlgen/sql/types.Float", "github.com/si3nloong/sqlgen/sql/types.PtrOfFloat"},
		"*time.Time": {"github.com/si3nloong/sqlgen/sql/types.Time", "github.com/si3nloong/sqlgen/sql/types.PtrOfTime"},
	}
)

type Codec string

func (c Codec) IsPkgFunc() (*types.Package, string, bool) {
	pkg := string(c)
	idx := strings.LastIndexByte(pkg, '.')
	if idx > 0 {
		path := pkg[:idx]
		cb := pkg[idx+1:]
		return types.NewPackage(path, filepath.Base(path)), cb, true
	}
	return nil, "", false
}

func (c Codec) CastOrInvoke(pkg *Package, v string) string {
	if p, invoke, ok := c.IsPkgFunc(); ok {
		p, _ = pkg.Import(p)
		return p.Name() + "." + invoke + "(" + v + ")"
	}
	return string(c) + "(" + v + ")"
}

type Mapping struct {
	Encoder Codec
	Decoder Codec
}
