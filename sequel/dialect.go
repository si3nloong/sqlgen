package sequel

import (
	"sync"
)

var (
	dialectMap     = new(sync.Map)
	defaultDialect = "mysql"
)

func RegisterDialect(name string, d Dialect) {
	if d == nil {
		panic("sqlgen: cannot register nil as dialect")
	}
	dialectMap.Store(name, d)
}

func GetDialect(name string) (Dialect, bool) {
	v, ok := dialectMap.Load(name)
	if !ok {
		return nil, ok
	}
	return v.(Dialect), ok
}

func DefaultDialect() Dialect {
	v, _ := dialectMap.Load(defaultDialect)
	return v.(Dialect)
}
