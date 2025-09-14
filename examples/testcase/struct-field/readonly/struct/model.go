package readonly

import "sync"

type noCopy [0]sync.Mutex

type Model struct {
	noCopy
	A        string
	B        bool
	ReadOnly string `sql:",readonly"`
}

func init() {
	m := Model{}
	m.A = ""
	m.B = false
}
