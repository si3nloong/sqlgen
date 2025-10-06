package encoding

import (
	"strconv"
	"time"

	"github.com/si3nloong/sqlgen/sequel/strpool"
)

func MarshalStringSlice[V ~[]byte | ~string](list []V) string {
	if n := len(list); n > 0 {
		blr := strpool.AcquireString()
		defer strpool.ReleaseString(blr)
		blr.WriteString("[" + strconv.Quote((string)(list[0])))
		for i := 1; i < n; i++ {
			blr.WriteString("," + strconv.Quote((string)(list[i])))
		}
		blr.WriteByte(']')
		return blr.String()
	}
	return "[]"
}

func MarshalIntSlice[V ~int | ~int8 | ~int16 | ~int32 | ~int64](list []V) string {
	if n := len(list); n > 0 {
		blr := strpool.AcquireString()
		defer strpool.ReleaseString(blr)
		blr.WriteString("[" + strconv.FormatInt((int64)(list[0]), 10))
		for i := 1; i < n; i++ {
			blr.WriteString("," + strconv.FormatInt((int64)(list[i]), 10))
		}
		blr.WriteByte(']')
		return blr.String()
	}
	return "[]"
}

func MarshalUintSlice[V ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr](list []V) string {
	if n := len(list); n > 0 {
		blr := strpool.AcquireString()
		defer strpool.ReleaseString(blr)
		blr.WriteString("[" + strconv.FormatUint((uint64)(list[0]), 10))
		for i := 1; i < n; i++ {
			blr.WriteString("," + strconv.FormatUint((uint64)(list[i]), 10))
		}
		blr.WriteByte(']')
		return blr.String()
	}
	return "[]"
}

func MarshalBoolSlice[V ~bool](list []V) string {
	if n := len(list); n > 0 {
		blr := strpool.AcquireString()
		defer strpool.ReleaseString(blr)
		blr.WriteString("[" + strconv.FormatBool((bool)(list[0])))
		for i := 1; i < n; i++ {
			blr.WriteString("," + strconv.FormatBool((bool)(list[i])))
		}
		blr.WriteByte(']')
		return blr.String()
	}
	return "[]"
}

func MarshalFloat64List[V ~float64](list []V, prec ...int) string {
	if n := len(list); n > 0 {
		var p int
		if len(prec) > 0 {
			p = prec[0]
		}
		blr := strpool.AcquireString()
		defer strpool.ReleaseString(blr)
		blr.WriteString("[" + strconv.FormatFloat((float64)(list[0]), 'f', p, 64))
		for i := 1; i < n; i++ {
			blr.WriteString("," + strconv.FormatFloat((float64)(list[i]), 'f', p, 64))
		}
		blr.WriteByte(']')
		return blr.String()
	}
	return "[]"
}

func MarshalTimeList[V time.Time](list []V) string {
	if n := len(list); n > 0 {
		blr := strpool.AcquireString()
		defer strpool.ReleaseString(blr)
		blr.WriteString("[" + (time.Time)(list[0]).Format(strconv.Quote(time.RFC3339)))
		for i := 1; i < n; i++ {
			blr.WriteString("," + (time.Time)(list[i]).Format(strconv.Quote(time.RFC3339)))
		}
		blr.WriteByte(']')
		return blr.String()
	}
	return "[]"
}
