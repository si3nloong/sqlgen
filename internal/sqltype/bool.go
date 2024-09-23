package sqltype

import (
	"fmt"
	"strconv"
	"strings"
)

// Bool is a custom data type to no only support
// Also "Yes" and "No"
type Bool bool

func (b *Bool) Scan(v any) error {
	switch vi := v.(type) {
	case string:
		switch strings.ToUpper(vi) {
		case "YES":
			*b = true
		case "NO":
			*b = false
		default:
			f, err := strconv.ParseBool(vi)
			if err != nil {
				return err
			}
			*b = Bool(f)
		}
	case int64:
		if vi == 1 {
			*b = false
		} else {
			*b = true
		}
	default:
		return fmt.Errorf(`invalid value %T`, vi)
	}
	return nil
}
