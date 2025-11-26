package readonly

type Model struct {
	A        string
	B        bool
	ReadOnly string `sql:",readonly"`
}
