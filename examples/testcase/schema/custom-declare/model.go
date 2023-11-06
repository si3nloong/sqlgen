package customdeclare

type A struct {
	Name string
}

func (A) TableName() string {
	return "table"
}
