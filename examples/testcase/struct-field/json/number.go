package json

import (
	"encoding/json"
	"fmt"
	"strconv"
	"unsafe"
)

type number struct {
	v *json.Number
}

func Number(v *json.Number) *number {
	return &number{v: v}
}

func (n number) Scan(v any) error {
	switch vi := v.(type) {
	case []byte:
		if _, err := strconv.ParseFloat(unsafe.String(unsafe.SliceData(vi), len(vi)), 64); err != nil {
			return err
		}
		*n.v = json.Number(vi)
	case string:
		if _, err := strconv.ParseFloat(vi, 64); err != nil {
			return err
		}
		*n.v = json.Number(vi)
	case int64:
		*n.v = json.Number(strconv.FormatInt(vi, 10))
	case float64:
		*n.v = json.Number(strconv.FormatFloat(vi, 'f', 10, 64))
	default:
		return fmt.Errorf(`unsupported type %T for json.Number`, vi)
	}
	return nil
}

type JSON struct {
	Num      json.Number
	RawBytes json.RawMessage
}
