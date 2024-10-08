package codegen

import (
	"go/types"
	"reflect"
	"strings"

	"github.com/samber/lo"
	"github.com/si3nloong/sqlgen/codegen/dialect"
)

type structType struct {
	name   string
	t      types.Type
	fields []*structFieldType
}

type structFieldType struct {
	name     string
	index    []int
	path     string
	t        types.Type
	enums    *enum
	exported bool
	embedded bool
	tag      reflect.StructTag
}

type enum struct {
	typeName string
	values   []*enumValue
}

type enumValue struct {
	name  string
	value string
}

type tableInfo struct {
	goName      string
	dbName      string
	tableName   string
	t           types.Type
	autoIncrKey *columnInfo
	keys        []*columnInfo
	columns     []*columnInfo
	indexes     []*indexInfo
}

func (b *tableInfo) GoName() string {
	return b.goName
}

func (b *tableInfo) DatabaseName() string {
	return b.dbName
}

func (b *tableInfo) TableName() string {
	return b.tableName
}

func (b *tableInfo) Keys() []string {
	return lo.Map(b.keys, func(c *columnInfo, _ int) string {
		return c.columnName
	})
}
func (b *tableInfo) Columns() []string {
	return lo.Map(b.columns, func(c *columnInfo, _ int) string {
		return c.columnName
	})
}
func (b *tableInfo) Indexes() []string {
	return lo.Map(b.indexes, func(c *indexInfo, _ int) string {
		return strings.Join(c.columns, ",")
	})
}
func (b *tableInfo) Implements(T *types.Interface) (*types.Func, bool) {
	return types.MissingMethod(b.t, T, true)
}
func (b *tableInfo) PtrImplements(T *types.Interface) (*types.Func, bool) {
	return types.MissingMethod(types.NewPointer(b.t), T, true)
}
func (b *tableInfo) colsWithoutAutoIncrPK() []*columnInfo {
	return lo.Filter(b.columns, func(v *columnInfo, _ int) bool {
		return !v.AutoIncr()
	})
}

// Mean table has only pk
func (b tableInfo) hasNoColsExceptPK() bool {
	return len(b.keys) == len(b.columns)
}

func (b tableInfo) hasNoColsExceptAutoPK() bool {
	return b.autoIncrKey != nil && len(b.columns) == 1 &&
		b.autoIncrKey == b.columns[0]
}

type goTag struct {
	key   string
	value string
}

type columnInfo struct {
	goName     string
	goPath     string
	columnName string
	columnPos  int
	size       int
	t          types.Type
	enums      *enum
	tags       []goTag
	mapper     *dialect.ColumnType
}

func (c *columnInfo) hasOption(k string) bool {
	_, _, ok := lo.FindLastIndexOf(c.tags, func(v goTag) bool {
		return v.key == k
	})
	return ok
}

func (c *columnInfo) getOption(k string) (string, bool) {
	tag, _, ok := lo.FindLastIndexOf(c.tags, func(v goTag) bool {
		return v.key == k
	})
	if ok {
		return tag.value, true
	}
	return "", false
}

func (c *columnInfo) Nullable() bool {
	switch c.t.(type) {
	case *types.Pointer,
		*types.Map,
		*types.Chan,
		*types.Interface,
		*types.Slice:
		return true
	default:
		return false
	}
}

func (c columnInfo) GoName() string {
	return c.goName
}

func (c columnInfo) GoPath() string {
	return c.goPath
}

func (c *columnInfo) Type() types.Type {
	return c.t
}

func (c *columnInfo) AutoIncr() bool {
	if _, ok := lo.Find(c.tags, func(v goTag) bool {
		return v.key == TagOptionAutoIncrement
	}); ok {
		return true
	}
	return false
}

func (c columnInfo) ColumnName() string {
	return c.columnName
}

func (c *columnInfo) ColumnPos() int {
	return c.columnPos
}

func (c *columnInfo) Implements(T *types.Interface) (wrongType bool) {
	_, wrongType = types.MissingMethod(c.t, T, true)
	return
}

func (i *columnInfo) sqlValuer() (func(string) string, bool) {
	if i.mapper == nil || i.mapper.SQLValuer == "" {
		return nil, false
	}
	return func(column string) string {
		return strings.Replace(i.mapper.SQLValuer, "{{.}}", column, 1)
	}, true
}

func (i columnInfo) sqlScanner() (func(string) string, bool) {
	if i.mapper == nil || i.mapper.SQLScanner == "" {
		return nil, false
	}
	return func(column string) string {
		return strings.Replace(i.mapper.SQLScanner, "{{.}}", column, 1)
	}, true
}

type indexInfo struct {
	columns   []string
	indexType string
}

func (i indexInfo) Columns() []string {
	return i.columns
}

func (i indexInfo) Type() string {
	return i.indexType
}
