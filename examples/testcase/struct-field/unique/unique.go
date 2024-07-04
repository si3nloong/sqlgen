package unique

type User struct {
	Email     string `sql:",unique"`
	Age       uint8
	FirstName string `sql:",unique:n"`
	LastName  string `sql:",unique:n"`
}
