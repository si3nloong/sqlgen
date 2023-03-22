package encoding

import (
	"strconv"
	"time"

	"github.com/valyala/bytebufferpool"
	"golang.org/x/exp/constraints"
)

type ValueType interface {
	[]byte | bool | float64 | uint64 | int64 | string | time.Time
}

type BaseType interface {
	~[]byte | ~bool | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64 | ~string | time.Time
}

func MarshalStringList[V ~[]byte | ~string](list []V) string {
	blr := bytebufferpool.Get()
	defer bytebufferpool.Put(blr)
	blr.WriteByte('[')
	for i, el := range list {
		if i > 0 {
			blr.WriteByte(',')
		}
		blr.WriteString(strconv.Quote(string(el)))
	}
	blr.WriteByte(']')
	return blr.String()
}

func MarshalIntList[V constraints.Integer](list []V) string {
	blr := bytebufferpool.Get()
	defer bytebufferpool.Put(blr)
	blr.WriteByte('[')
	for i, el := range list {
		if i > 0 {
			blr.WriteByte(',')
		}
		blr.WriteString(strconv.FormatInt(int64(el), 10))
	}
	blr.WriteByte(']')
	return blr.String()
}

func MarshalBoolList[V ~bool](list []V) string {
	blr := bytebufferpool.Get()
	defer bytebufferpool.Put(blr)
	blr.WriteByte('[')
	for i, el := range list {
		if i > 0 {
			blr.WriteByte(',')
		}
		blr.WriteString(strconv.FormatBool(bool(el)))
	}
	blr.WriteByte(']')
	return blr.String()
}

func MarshalFloatList[V constraints.Float](list []V) string {
	blr := bytebufferpool.Get()
	defer bytebufferpool.Put(blr)
	blr.WriteByte('[')
	for i, el := range list {
		if i > 0 {
			blr.WriteByte(',')
		}
		blr.WriteString(strconv.FormatFloat(float64(el), 'e', -1, 64))
	}
	blr.WriteByte(']')
	return blr.String()
}
