package types

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"unsafe"
)

type jsonMarshaler[T any] struct {
	v T
}

func JSONMarshaler[T any](addr T) driver.Valuer {
	return jsonMarshaler[T]{v: addr}
}

func (j jsonMarshaler[T]) Value() (driver.Value, error) {
	switch vj := any(j.v).(type) {
	case nil:
		return []byte(`null`), nil
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

type jsonUnmarshaler[T any, Ptr interface {
	*T
}] struct {
	v Ptr
}

func JSONUnmarshaler[T any, Ptr interface {
	*T
}](addr Ptr) sql.Scanner {
	return &jsonUnmarshaler[T, Ptr]{v: addr}
}

func (j *jsonUnmarshaler[T, Ptr]) Scan(v any) error {
	switch vj := any(j.v).(type) {
	case json.Unmarshaler:
		switch vi := v.(type) {
		case []byte:
			return vj.UnmarshalJSON(vi)
		case json.RawMessage:
			return vj.UnmarshalJSON(vi)
		case string:
			return vj.UnmarshalJSON(unsafe.Slice(unsafe.StringData(vi), len(vi)))
		}
	default:
		switch vi := v.(type) {
		case nil:
			j.v = nil
			return nil
		case []byte:
			return json.NewDecoder(bytes.NewBuffer(vi)).Decode(j.v)
		case json.RawMessage:
			return json.NewDecoder(bytes.NewBuffer(vi)).Decode(j.v)
		case string:
			return json.NewDecoder(bytes.NewBufferString(vi)).Decode(j.v)
		}
	}
	return fmt.Errorf(`sequel/types: invalid scan type for JSON, %T`, v)
}
