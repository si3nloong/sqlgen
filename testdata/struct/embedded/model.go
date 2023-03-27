package embedded

type a struct {
	ID   int64
	Name string
	Z    bool
}

type B struct {
	a
	Name string
}
