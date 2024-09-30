package codegen

import (
	"database/sql/driver"
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
	name  string
	index []int
	// path     string
	paths    []string
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
}

func (b *tableInfo) GoName() string {
	return b.goName
}

func (b *tableInfo) DBName() string {
	return b.dbName
}

func (b *tableInfo) TableName() string {
	return b.tableName
}

func (b *tableInfo) PK() []string {
	return lo.Map(b.keys, func(c *columnInfo, _ int) string {
		return c.columnName
	})
}

func (b *tableInfo) Columns() []string {
	return lo.Map(b.columns, func(c *columnInfo, _ int) string {
		return c.columnName
	})
}

func (b *tableInfo) ColumnByIndex(i int) dialect.GoColumn {
	return b.columns[i]
}

func (b *tableInfo) RangeIndex(rangeFunc func(dialect.Index, int)) {
	var (
		optMap = make(map[string]*indexInfo)
		key    string
	)
	for _, col := range b.columns {
		switch {
		case col.hasOption(TagOptionIndex):
			key, _ = col.getOptionValue(TagOptionIndex)
		case col.hasOption(TagOptionUnique):
			key, _ = col.getOptionValue(TagOptionIndex)
		default:
			continue
		}

		if _, ok := optMap[key]; !ok {
			optMap[key] = &indexInfo{}
		}

		optMap[key].columns = append(optMap[key].columns, col.columnName)
	}
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

// Column info contain go actual type information
// Some of the default behaviour is not able to override, such as go size, go enum, go tags, go path, go name, go nullable
type columnInfo struct {
	goName     string
	goPaths    []string
	columnName string
	columnPos  int
	size       int
	t          types.Type
	enums      *enum
	tags       []goTag
	mapper     *dialect.ColumnType
}

func (c *columnInfo) Column() string {
	return c.columnName
}

func (c *columnInfo) ColumnName() string {
	return c.columnName
}

func (c *columnInfo) ColumnPos() int {
	return c.columnPos
}

func (c *columnInfo) GoName() string {
	return c.goName
}

func (c *columnInfo) GoPath() string {
	return strings.Join(lo.Map(c.goPaths, func(v string, _ int) string {
		if v[0] == '*' {
			return v[1:]
		}
		return v
	}), ".")
}

func (c *columnInfo) GoPaths() []string {
	var goPath string
	paths := []string{}
	for _, path := range c.goPaths {
		if path[0] == '*' {
			paths = append(paths, goPath+path[1:])
			goPath = ""
			continue
		}
		if goPath != "" {
			goPath += "." + path
		} else {
			goPath += path
		}
	}
	paths = append(paths, goPath)
	return paths
}

func (c *columnInfo) GoType() types.Type {
	return c.t
}

func (c *columnInfo) isPtr() bool {
	_, ok := c.t.(*types.Pointer)
	return ok
}

func (c *columnInfo) GoNullable() bool {
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

func (c *columnInfo) Implements(T *types.Interface) (wrongType bool) {
	_, wrongType = types.MissingMethod(c.t, T, true)
	return
}

func (c *columnInfo) hasOption(k string) bool {
	_, _, ok := lo.FindLastIndexOf(c.tags, func(v goTag) bool {
		return v.key == k
	})
	return ok
}

func (c *columnInfo) getOptionValue(k string) (string, bool) {
	tag, _, ok := lo.FindLastIndexOf(c.tags, func(v goTag) bool {
		return v.key == k
	})
	if ok {
		return tag.value, true
	}
	return "", false
}

func (c *columnInfo) DataType() string {
	return c.mapper.DataType(c)
}

func (c *columnInfo) Default() (driver.Value, bool) {
	return nil, false
}

func (c *columnInfo) Key() bool {
	_, ok1 := c.getOptionValue(TagOptionPKAlias)
	_, ok2 := c.getOptionValue(TagOptionFK)
	_, ok3 := c.getOptionValue(TagOptionPK)
	_, ok4 := c.getOptionValue(TagOptionAutoIncrement)
	return ok1 || ok2 || ok3 || ok4
}

func (c columnInfo) Name() string {
	return c.columnName
}

func (c *columnInfo) AutoIncr() bool {
	if _, ok := lo.Find(c.tags, func(v goTag) bool {
		return v.key == TagOptionAutoIncrement
	}); ok {
		return true
	}
	return false
}

func (c columnInfo) Size() int {
	return c.size
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
	columns []string
	unique  bool
}

func (i indexInfo) Columns() []string {
	return i.columns
}

func (i indexInfo) Unique() bool {
	return i.unique
}
