package encoding

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"unsafe"
)

type jsonValuer[T any] struct {
	v T
}

func JSONValue[T any](addr T) driver.Valuer {
	return jsonValuer[T]{v: addr}
}

func (j jsonValuer[T]) Value() (driver.Value, error) {
	switch vj := any(j.v).(type) {
	case nil:
		return nil, nil
	case json.RawMessage:
		return []byte(vj)[:], nil
	case json.Marshaler:
		return vj.MarshalJSON()
	default:
		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(j.v); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
}

type jsonScanner[T any, Addr addrOrPtr[T]] struct {
	v Addr
}

func JSONScanner[T any, Addr interface {
	*T
}](addr Addr) sql.Scanner {
	return &jsonScanner[T, Addr]{v: addr}
}

func (j *jsonScanner[T, Addr]) Scan(v any) error {
	switch vi := v.(type) {
	case nil:
		switch vj := any(j.v).(type) {
		case json.Unmarshaler:
			return vj.UnmarshalJSON(nil)
		case **T:
			*(**T)(unsafe.Pointer(j.v)) = nil
			return nil
		case *T:
			var v T
			*(*T)(unsafe.Pointer(j.v)) = v
			return nil
		default:
			panic("unreachable")
		}
	case string:
		return json.NewDecoder(bytes.NewBufferString(vi)).Decode(j.v)
	case []byte:
		return json.NewDecoder(bytes.NewBuffer(vi)).Decode(j.v)
	case json.RawMessage:
		return json.NewDecoder(bytes.NewBuffer(vi)).Decode(j.v)
	default:
		return fmt.Errorf(`sequel/types: invalid scan type for JSON, %T`, v)
	}
}
