package valuerproperty

func (B) Columns() []string {
    return []string{"ID", "Value", "N"}
}

func (m B) Values() []any {
	return []any{m.ID, any(m.Value), m.N}
}

func (m *B) Addrs() []any {
	return []any{&m.ID, &m.Value, &m.N}
}