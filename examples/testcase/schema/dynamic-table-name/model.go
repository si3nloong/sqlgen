package tabler

type Model struct {
	Name string
}

func (m Model) TableName() string {
	if m.Name == "" {
		return "unknown"
	}
	return m.Name
}

type A struct {
	ID   int64 `sql:",pk"`
	Name string
}

func (m A) TableName() string {
	if m.Name == "" {
		return "apple"
	}
	return m.Name
}
