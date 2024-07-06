package array

const fixed = 10

type Str string

type Array struct {
	Tuple     [2]byte
	Runes     [4]rune
	Bytes     [10]byte
	FixedSize [fixed]byte
	Str       [100]Str
}
