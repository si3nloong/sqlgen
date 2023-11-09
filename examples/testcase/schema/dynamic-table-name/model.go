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
