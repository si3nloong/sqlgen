package testdata

import (
	"database/sql/driver"
	"strings"
	"sync"
	"time"
)

func (c *Common) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c Common) String() string {
	return strings.Join(c.StrList, ",")
}

type enum int

const (
	success enum = iota
	failed
)

// type Map[K comparable, V any] = map[K]V

type Empty struct{}

type No[T any] struct {
	sync.RWMutex
}

func (n No[T]) Something() {}

func (No[T]) Something2() {}

type List[T any] = []T

type customStr string

// var _ (driver.Valuer) = (*customStr)(nil)

func (c customStr) Value() (driver.Value, error) {
	return nil, nil
}

type aliasStr = customStr

type Common struct {
	sync.RWMutex `sql:"-"`
	priv         int
	Int64        int64     `sql:"i64"`
	Name, Yes    string    ``
	Cstr         customStr `sql:"str"`
	Alias        aliasStr
	Int          int
	Uint         uint
	Flag         bool
	r            rune
	// c            complex128
	// b            byte
	Ignore string `sql:"-"`
	// Bytes        []byte
	F32 float32 `sql:"f32bit"`
	F64 float64 `sql:"f64bit"`
	// a            any
	// Interface    interface{} `sql:"interface"`
	// inter        interface{}
	// Arr      [2]string `sql:"arr"`
	StrList []string
	// IntList  []int
	// BoolList []bool
	List List[any]
	Enum enum
	// Raw  json.RawMessage
	// Map          map[string]string
	T time.Time `sql:"t"`
}

var (
	_ driver.Valuer = (*Common)(nil)
)
