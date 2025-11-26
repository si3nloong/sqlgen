package readonly

// +sql:readonly
type Model struct {
	A        string
	B        bool
	ReadOnly string `sql:",readonly"`
}

func init() {
	m := Model{}
	m.A = ""
	m.B = false
}
