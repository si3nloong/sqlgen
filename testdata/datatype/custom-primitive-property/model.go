package customprimitivestruct

type enum int

const (
	success enum = iota
	failed
)

type Custom struct {
	Str    longText `sql:"text"`
	Enum   enum     `sql:"e"`
	priv   int
	Status string `sql:"-"`
	Num    uint16
}

type longText string
