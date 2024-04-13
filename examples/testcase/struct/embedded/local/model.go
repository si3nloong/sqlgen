package embedded

type a struct {
	ID   int64
	Name string
	Z    bool
}

type B struct {
	a
	ts
	// // FIXME: embedded will never overwrite parent property
	// Name string
}
