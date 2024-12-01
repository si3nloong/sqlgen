package encoding

import (
	"strconv"
	"time"

	"github.com/si3nloong/sqlgen/sequel/strpool"
	"golang.org/x/exp/constraints"
)

func MarshalStringSlice[V ~[]byte | ~string](list []V, enclose ...[2]byte) string {
	blr := strpool.AcquireString()
	defer strpool.ReleaseString(blr)
	if len(enclose) > 0 {
		blr.WriteByte(enclose[0][0])
	} else {
		blr.WriteByte('[')
	}
	for i := range list {
		if i > 0 {
			blr.WriteByte(',')
		}
		blr.WriteString(strconv.Quote((string)(list[i])))
	}
	if len(enclose) > 0 {
		blr.WriteByte(enclose[0][1])
	} else {
		blr.WriteByte(']')
	}
	return blr.String()
}

func MarshalIntSlice[V constraints.Signed](list []V) string {
	blr := strpool.AcquireString()
	defer strpool.ReleaseString(blr)
	blr.WriteByte('[')
	for i := range list {
		if i > 0 {
			blr.WriteString("," + strconv.FormatInt((int64)(list[i]), 10))
		} else {
			blr.WriteString(strconv.FormatInt((int64)(list[i]), 10))
		}
	}
	blr.WriteByte(']')
	return blr.String()
}

func MarshalUintSlice[V constraints.Unsigned](list []V) string {
	blr := strpool.AcquireString()
	defer strpool.ReleaseString(blr)
	blr.WriteByte('[')
	for i := range list {
		if i > 0 {
			blr.WriteString("," + strconv.FormatUint((uint64)(list[i]), 10))
		} else {
			blr.WriteString(strconv.FormatUint((uint64)(list[i]), 10))
		}
	}
	blr.WriteByte(']')
	return blr.String()
}

func MarshalBoolSlice[V ~bool](list []V) string {
	blr := strpool.AcquireString()
	defer strpool.ReleaseString(blr)
	blr.WriteByte('[')
	for i := range list {
		if i > 0 {
			blr.WriteByte(',')
		}
		if list[i] {
			blr.WriteString("true")
		} else {
			blr.WriteString("false")
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
	for i := range list {
		if i > 0 {
			blr.WriteString("," + strconv.FormatFloat((float64)(list[i]), 'f', p, 64))
		} else {
			blr.WriteString(strconv.FormatFloat((float64)(list[i]), 'f', p, 64))
		}
	}
	blr.WriteByte(']')
	return blr.String()
}

func MarshalTimeList[V time.Time](list []V) string {
	blr := strpool.AcquireString()
	defer strpool.ReleaseString(blr)
	blr.WriteByte('[')
	for i := range list {
		if i > 0 {
			blr.WriteString("," + (time.Time)(list[i]).Format(strconv.Quote(time.RFC3339)))
		} else {
			blr.WriteString((time.Time)(list[i]).Format(strconv.Quote(time.RFC3339)))
		}
	}
	blr.WriteByte(']')
	return blr.String()
}
