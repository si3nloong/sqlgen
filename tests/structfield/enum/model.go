package enum

type Enum int

const (
	success Enum = iota
	failed
)

type Custom struct {
	Str    longText `sql:"text"`
	Enum   Enum     `sql:"e"`
	priv   int
	Status string `sql:"-"`
	Num    uint16
}

type longText string
