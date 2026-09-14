package pk

import "time"

type Color int

const (
	Red Color = iota
	White
)

type Car struct {
	ID        PK `sql:",pk"`
	No        string
	Color     Color
	Year      *int
	ManucDate time.Time
}
