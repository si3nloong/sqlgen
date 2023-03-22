package aliasproperty

func (AliasStruct) Columns() []string {
	return []string{"Name", "DateTime"}
}

func (m AliasStruct) Values() []any {
	return []any{string(m.Name), m.DateTime}
}

func (m *AliasStruct) Addrs() []any {
	return []any{&m.Name, &m.DateTime}
}
