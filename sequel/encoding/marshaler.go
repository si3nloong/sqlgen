package encoding

import (
	"strconv"
	"time"

	"github.com/si3nloong/sqlgen/sequel/strpool"
	"golang.org/x/exp/constraints"
)

func MarshalStringList[V ~[]byte | ~string](list []V) string {
	blr := strpool.AcquireString()
	defer strpool.ReleaseString(blr)
	blr.WriteByte('[')
	for i := range list {
		if i > 0 {
			blr.WriteByte(',')
		}
		blr.WriteString(strconv.Quote(string(list[i])))
	}
	blr.WriteByte(']')
	return blr.String()
}

func MarshalSignedIntList[V constraints.Signed](list []V) string {
	blr := strpool.AcquireString()
	defer strpool.ReleaseString(blr)
	blr.WriteByte('[')
	for i := range list {
		if i > 0 {
			blr.WriteByte(',')
		}
		blr.WriteString(strconv.FormatInt(int64(list[i]), 10))
	}
	blr.WriteByte(']')
	return blr.String()
}

func MarshalUnsignedIntList[V constraints.Unsigned](list []V) string {
	blr := strpool.AcquireString()
	defer strpool.ReleaseString(blr)
	blr.WriteByte('[')
	for i := range list {
		if i > 0 {
			blr.WriteByte(',')
		}
		blr.WriteString(strconv.FormatUint(uint64(list[i]), 10))
	}
	blr.WriteByte(']')
	return blr.String()
}

func MarshalBoolList[V ~bool](list []V) string {
	blr := strpool.AcquireString()
	defer strpool.ReleaseString(blr)
	blr.WriteByte('[')
	for i := range list {
		if i > 0 {
			blr.WriteByte(',')
		}
		if list[i] {
			blr.WriteByte('1')
		} else {
			blr.WriteByte('0')
		}
	}
	blr.WriteByte(']')
	return blr.String()
}

func MarshalFloatList[V constraints.Float](list []V, precision ...int) string {
	blr := strpool.AcquireString()
	defer strpool.ReleaseString(blr)
	var prec = -1
	if len(precision) > 0 {
		prec = precision[0]
	}
	blr.WriteByte('[')
	for i := range list {
		if i > 0 {
			blr.WriteByte(',')
		}
		blr.WriteString(strconv.FormatFloat(float64(list[i]), 'f', prec, 64))
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
			blr.WriteByte(',')
		}
		blr.WriteString(time.Time(list[i]).Format(strconv.Quote(time.RFC3339)))
	}
	blr.WriteByte(']')
	return blr.String()
}
