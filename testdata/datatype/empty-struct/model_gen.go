package emptystruct

func (empty) Columns() []string {
    return []string{}
}

func (m empty) Values() []any {
	return []any{}
}

func (m *empty) Addrs() []any {
	return []any{}
}