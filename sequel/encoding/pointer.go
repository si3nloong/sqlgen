package encoding

import (
	"database/sql"
	"encoding"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"time"
	"unsafe"
)

func PtrScanner[T any, Ptr interface {
	*T
	sql.Scanner
}, Addr *Ptr](addr Addr) sql.Scanner {
	return &addrOfPtr[T, Ptr, Addr]{v: addr}
}

type addrOfPtr[T any, Ptr interface {
	*T
	sql.Scanner
}, Addr *Ptr] struct {
	v Addr
}

func (p *addrOfPtr[T, Ptr, Addr]) Scan(v any) error {
	switch vt := v.(type) {
	case nil:
		*(*Ptr)(unsafe.Pointer(p.v)) = nil
		return nil
	case string:
		switch vi := any(*p.v).(type) {
		case sql.Scanner:
			return vi.Scan(vt)
		case encoding.TextUnmarshaler:
			return vi.UnmarshalText(unsafe.Slice(unsafe.StringData(vt), len(vt)))
		case encoding.BinaryUnmarshaler:
			return vi.UnmarshalBinary(unsafe.Slice(unsafe.StringData(vt), len(vt)))
		default:
			return fmt.Errorf(`sql/encoding: unsupported %T`, vi)
		}
	case []byte:
		switch vi := any(*p.v).(type) {
		case sql.Scanner:
			return vi.Scan(vt)
		case encoding.TextUnmarshaler:
			return vi.UnmarshalText(vt)
		case encoding.BinaryUnmarshaler:
			return vi.UnmarshalBinary(vt)
		default:
			return fmt.Errorf(`sql/encoding: unsupported %T`, vi)
		}
	case bool:
		switch vi := any(*p.v).(type) {
		case sql.Scanner:
			return vi.Scan(vt)
		case encoding.TextUnmarshaler:
			val := strconv.FormatBool(vt)
			return vi.UnmarshalText(unsafe.Slice(unsafe.StringData(val), len(val)))
		case encoding.BinaryUnmarshaler:
			if vt {
				return vi.UnmarshalBinary([]byte{1})
			}
			return vi.UnmarshalBinary([]byte{0})
		default:
			return fmt.Errorf(`sql/encoding: unsupported %T`, vi)
		}
	case int64:
		switch vi := any(*p.v).(type) {
		case sql.Scanner:
			return vi.Scan(vt)
		case encoding.TextUnmarshaler:
			val := strconv.FormatInt(vt, 10)
			return vi.UnmarshalText(unsafe.Slice(unsafe.StringData(val), len(val)))
		case encoding.BinaryUnmarshaler:
			b := make([]byte, 8)
			binary.LittleEndian.PutUint64(b, uint64(vt))
			return vi.UnmarshalBinary(b)
		default:
			return fmt.Errorf(`sql/encoding: unsupported %T`, vi)
		}
	case float64:
		switch vi := any(*p.v).(type) {
		case sql.Scanner:
			return vi.Scan(vt)
		case encoding.TextUnmarshaler:
			b := make([]byte, 8)
			binary.LittleEndian.PutUint64(b, math.Float64bits(vt))
			return vi.UnmarshalText(b)
		case encoding.BinaryUnmarshaler:
			b := make([]byte, 8)
			binary.LittleEndian.PutUint64(b, math.Float64bits(vt))
			return vi.UnmarshalBinary(b)
		default:
			return fmt.Errorf(`sql/encoding: unsupported %T`, vi)
		}
	case time.Time:
		switch vi := any(*p.v).(type) {
		case sql.Scanner:
			return vi.Scan(vt)
		case encoding.TextUnmarshaler:
			b, err := vt.MarshalText()
			if err != nil {
				return err
			}
			return vi.UnmarshalText(b)
		case encoding.BinaryUnmarshaler:
			b, err := vt.MarshalBinary()
			if err != nil {
				return err
			}
			return vi.UnmarshalBinary(b)
		default:
			return fmt.Errorf(`sql/encoding: unsupported %T`, vi)
		}
	default:
		return fmt.Errorf(`sql/encoding: unsupported %T`, vt)
	}
}
