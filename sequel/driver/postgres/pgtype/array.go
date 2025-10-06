package pgtype

import (
	"bytes"
	"cmp"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

func arrayScan[T cmp.Ordered, Arr interface{ ~[]T }](a *Arr, src any, t string) error {
	switch src := src.(type) {
	case []byte:
		return scanBytes(a, src, t)
	case string:
		return scanBytes(a, ([]byte)(src), t)
	case nil:
		*a = nil
		return nil
	}
	return fmt.Errorf("pgtype: cannot convert %T to %s", src, t)
}

func intArrayValue[T ~int | ~int8 | ~int16 | ~int32 | ~int64](a []T) (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '{'

		b = strconv.AppendInt(b, (int64)(a[0]), 10)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendInt(b, (int64)(a[i]), 10)
		}

		return string(append(b, '}')), nil
	}
	return "{}", nil
}

func uintArrayValue[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr](a []T) (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '{'

		b = strconv.AppendUint(b, (uint64)(a[0]), 10)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendUint(b, (uint64)(a[i]), 10)
		}

		return string(append(b, '}')), nil
	}
	return "{}", nil
}

// parseArray extracts the dimensions and elements of an array represented in
// text format. Only representations emitted by the backend are supported.
// Notably, whitespace around brackets and delimiters is significant, and NULL
// is case-sensitive.
//
// See http://www.postgresql.org/docs/current/static/arrays.html#ARRAYS-IO
func parseArray(src, del []byte) (dims []int, elems [][]byte, err error) {
	var depth, i int

	if len(src) < 1 || src[0] != '{' {
		return nil, nil, fmt.Errorf("pgtype: unable to parse array; expected %q at offset %d", '{', 0)
	}

Open:
	for i < len(src) {
		switch src[i] {
		case '{':
			depth++
			i++
		case '}':
			elems = make([][]byte, 0)
			goto Close
		default:
			break Open
		}
	}
	dims = make([]int, i)

Element:
	for i < len(src) {
		switch src[i] {
		case '{':
			if depth == len(dims) {
				break Element
			}
			depth++
			dims[depth-1] = 0
			i++
		case '"':
			var elem = []byte{}
			var escape bool
			for i++; i < len(src); i++ {
				if escape {
					elem = append(elem, src[i])
					escape = false
				} else {
					switch src[i] {
					default:
						elem = append(elem, src[i])
					case '\\':
						escape = true
					case '"':
						elems = append(elems, elem)
						i++
						break Element
					}
				}
			}
		default:
			for start := i; i < len(src); i++ {
				if bytes.HasPrefix(src[i:], del) || src[i] == '}' {
					elem := src[start:i]
					if len(elem) == 0 {
						return nil, nil, fmt.Errorf("pgtype: unable to parse array; unexpected %q at offset %d", src[i], i)
					}
					if bytes.Equal(elem, []byte("NULL")) {
						elem = nil
					}
					elems = append(elems, elem)
					break Element
				}
			}
		}
	}

	for i < len(src) {
		if bytes.HasPrefix(src[i:], del) && depth > 0 {
			dims[depth-1]++
			i += len(del)
			goto Element
		} else if src[i] == '}' && depth > 0 {
			dims[depth-1]++
			depth--
			i++
		} else {
			return nil, nil, fmt.Errorf("pgtype: unable to parse array; unexpected %q at offset %d", src[i], i)
		}
	}

Close:
	for i < len(src) {
		if src[i] == '}' && depth > 0 {
			depth--
			i++
		} else {
			return nil, nil, fmt.Errorf("pgtype: unable to parse array; unexpected %q at offset %d", src[i], i)
		}
	}
	if depth > 0 {
		err = fmt.Errorf("pgtype: unable to parse array; expected %q at offset %d", '}', i)
	}
	if err == nil {
		for _, d := range dims {
			if (len(elems) % d) != 0 {
				err = fmt.Errorf("pgtype: multidimensional arrays must have elements with matching dimensions")
			}
		}
	}
	return
}

func scanLinearArray(src, del []byte, typ string) (elems [][]byte, err error) {
	dims, elems, err := parseArray(src, del)
	if err != nil {
		return nil, err
	}
	if len(dims) > 1 {
		return nil, fmt.Errorf("pgtype: cannot convert ARRAY%s to %s", strings.Replace(fmt.Sprint(dims), " ", "][", -1), typ)
	}
	return elems, err
}

func appendArrayQuotedBytes(b, v []byte) []byte {
	b = append(b, '"')
	for {
		i := bytes.IndexAny(v, `"\`)
		if i < 0 {
			b = append(b, v...)
			break
		}
		if i > 0 {
			b = append(b, v[:i]...)
		}
		b = append(b, '\\', v[i])
		v = v[i+1:]
	}
	return append(b, '"')
}

// func appendValue(b []byte, v driver.Value) ([]byte, error) {
// 	return append(b, encode(nil, v, 0)...), nil
// }
