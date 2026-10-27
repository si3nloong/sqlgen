package encoding

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

func JSONValue[T any](value T) driver.Valuer {
	return jsonValuer[T]{v: value}
}

type jsonValuer[T any] struct {
	v T
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

func JSONScanner[T any, Addr interface{ *T }](addr Addr) sql.Scanner {
	return &jsonScanner[T, Addr]{v: addr}
}

type jsonScanner[T any, Addr interface{ *T }] struct {
	v Addr
}

func (j *jsonScanner[T, Addr]) Scan(v any) error {
	switch vi := v.(type) {
	case nil:
		switch vj := any(j.v).(type) {
		case json.Unmarshaler:
			return vj.UnmarshalJSON(nil)
		default:
			var zero T
			*j.v = zero
			return nil
		}
	case string:
		return json.NewDecoder(bytes.NewBufferString(vi)).Decode(j.v)
	case []byte:
		return json.NewDecoder(bytes.NewBuffer(vi)).Decode(j.v)
	case json.RawMessage:
		return json.NewDecoder(bytes.NewBuffer(vi)).Decode(j.v)
	default:
		return fmt.Errorf(`sequel/encoding: invalid scan type for JSON, %T`, v)
	}
}
