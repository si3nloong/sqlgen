package postgres

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"strings"
	"unsafe"

	"github.com/samber/lo"
	"github.com/si3nloong/sqlgen/sequel"
)

//go:generate stringer --type indexType --linecomment
type indexType uint8

const (
	bTree indexType = iota // BTREE
	hash                   // HASH
	brin                   // BRIN
)

type tableDefinition struct {
	keys []string
	cols []*columnDefinition
	idxs []*indexDefinition
}

func (s *tableDefinition) PK() (sequel.TablePK, bool) {
	if len(s.keys) == 0 {
		return nil, false
	}
	return nil, true
}

func (s *tableDefinition) Columns() []string {
	return lo.Map(s.cols, func(v *columnDefinition, _ int) string {
		return v.name
	})
}

func (s *tableDefinition) Indexes() []string {
	return lo.Map(s.idxs, func(v *indexDefinition, _ int) string {
		return v.Name()
	})
}

func (s *tableDefinition) Column(i int) sequel.TableColumn {
	return s.cols[i]
}

func (s *tableDefinition) Index(i int) sequel.TableIndex {
	return s.idxs[i]
}

type columnDefinition struct {
	name         string
	length       int64
	dataType     string
	nullable     bool
	comment      string
	defaultValue any
	check        sql.RawBytes
}

func (c *columnDefinition) Name() string {
	return c.name
}

func (c *columnDefinition) Length() int64 {
	return c.length
}

func (c *columnDefinition) DataType() string {
	return c.dataType
}

func (c *columnDefinition) Nullable() bool {
	return c.nullable
}

func (c *columnDefinition) Comment() string {
	return c.comment
}

func (c *columnDefinition) Default() any {
	return c.defaultValue
}

func (c *columnDefinition) Definition() string {
	str := c.name + " " + c.dataType
	if !c.nullable {
		str += " NOT NULL"
	}
	if c.defaultValue != nil {
		str += " DEFAULT " + format(c.defaultValue)
	}
	if len(c.check) != 0 {
		str += " " + unsafe.String(unsafe.SliceData(c.check), len(c.check))
	}
	return str
}

type indexDefinition struct {
	cols      []string
	indexType indexType
}

func (i *indexDefinition) Name() string {
	str := strings.Join(i.cols, ",")
	hash := md5.Sum(unsafe.Slice(unsafe.StringData(str), len(str)))
	return hex.EncodeToString(hash[:])
}

func (i *indexDefinition) Type() string {
	return i.indexType.String()
}

func (i *indexDefinition) Columns() []string {
	return i.cols
}

func (i *indexDefinition) Definition() string {
	return "CONSTRAINT " + i.Name() + " " + i.Type() + " (" + strings.Join(i.cols, ",") + ")"
}
