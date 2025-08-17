package encoding

import (
	"strconv"
	"time"

	"github.com/si3nloong/sqlgen/sequel/strpool"
	"golang.org/x/exp/constraints"
)

func MarshalStringSlice[V ~[]byte | ~string](list []V) string {
	blr := strpool.AcquireString()
	defer strpool.ReleaseString(blr)
	blr.WriteByte('[')
	if len(list) > 0 {
		blr.WriteString(strconv.Quote((string)(list[0])))
		for i := 1; i < len(list); i++ {
			blr.WriteString("," + strconv.Quote((string)(list[i])))
		}
	}
	blr.WriteByte(']')
	return blr.String()
}

func MarshalIntSlice[V constraints.Signed](list []V) string {
	blr := strpool.AcquireString()
	defer strpool.ReleaseString(blr)
	blr.WriteByte('[')
	if len(list) > 0 {
		blr.WriteString(strconv.FormatInt((int64)(list[0]), 10))
		for i := 1; i < len(list); i++ {
			blr.WriteString("," + strconv.FormatInt((int64)(list[i]), 10))
		}
	}
	blr.WriteByte(']')
	return blr.String()
}

func MarshalUintSlice[V constraints.Unsigned](list []V) string {
	blr := strpool.AcquireString()
	defer strpool.ReleaseString(blr)
	blr.WriteByte('[')
	if len(list) > 0 {
		blr.WriteString(strconv.FormatUint((uint64)(list[0]), 10))
		for i := 1; i < len(list); i++ {
			blr.WriteString("," + strconv.FormatUint((uint64)(list[i]), 10))
		}
	}
	blr.WriteByte(']')
	return blr.String()
}

func MarshalBoolSlice[V ~bool](list []V) string {
	blr := strpool.AcquireString()
	defer strpool.ReleaseString(blr)
	blr.WriteByte('[')
	if len(list) > 0 {
		blr.WriteString(strconv.FormatBool((bool)(list[0])))
		for i := 1; i < len(list); i++ {
			blr.WriteString("," + strconv.FormatBool((bool)(list[i])))
		}
	}
	blr.WriteByte(']')
	return blr.String()
}

func MarshalFloat64List[V ~float64](list []V, prec ...int) string {
	var p int
	if len(prec) > 0 {
		p = prec[0]
	}
	blr := strpool.AcquireString()
	defer strpool.ReleaseString(blr)
	blr.WriteByte('[')
	if len(list) > 0 {
		blr.WriteString(strconv.FormatFloat((float64)(list[0]), 'f', p, 64))
		for i := 1; i < len(list); i++ {
			blr.WriteString("," + strconv.FormatFloat((float64)(list[i]), 'f', p, 64))
		}
	}
	blr.WriteByte(']')
	return blr.String()
}

func MarshalTimeList[V time.Time](list []V) string {
	blr := strpool.AcquireString()
	defer strpool.ReleaseString(blr)
	blr.WriteByte('[')
	if len(list) > 0 {
		blr.WriteString((time.Time)(list[0]).Format(strconv.Quote(time.RFC3339)))
		for i := 1; i < len(list); i++ {
			blr.WriteString("," + (time.Time)(list[i]).Format(strconv.Quote(time.RFC3339)))
		}
	}
	blr.WriteByte(']')
	return blr.String()
}
