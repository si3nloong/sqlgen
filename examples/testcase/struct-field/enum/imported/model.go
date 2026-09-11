package imported

import (
	"reflect"
	"time"
)

type ImportedEnum struct {
	Weekday time.Weekday
	Kind    reflect.Kind
}
