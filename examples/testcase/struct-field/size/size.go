package size

import "time"

type Size struct {
	Str       string    `sql:",size:25"`
	Timestamp time.Time `sql:",size:6"`
	Time      time.Time
}
