package enum

type Enum int

const (
	success Enum = iota + 1_00
	failed
	otherStr longText = "very very very long text"
	cancelled
)

type uint16Enum uint16

const (
	zero uint16Enum = iota
	one
	two
	three
	four
	five
	six
	seven
	eight
	nine
	ten
)

type Custom struct {
	Str    longText  `sql:"text"`
	Enum   Enum      `sql:"e"`
	PtrStr *longText `sql:"ptr_str"`
	priv   int
	Status string `sql:"-"`
	Uint16 uint16Enum
	Num    uint16
}

type (
	smallText string
	longText  string
)
