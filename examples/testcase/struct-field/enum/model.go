package enum

type Enum int

const (
	success Enum = iota
	failed
	otherStr longText = "very very very long text"
	cancelled
)

type Custom struct {
	Str    longText `sql:"text"`
	Enum   Enum     `sql:"e"`
	priv   int
	Status string `sql:"-"`
	Num    uint16
}

type (
	smallText string
	longText  string
)
