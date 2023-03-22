package arrayproperty

func (Customer) Columns() []string {
	return []string{"id", "Name", "Address", "Nicknames", "Status"}
}

func (m Customer) Values() []any {
	return []any{m.ID, m.Name, m.Address, m.Nicknames, string(m.Status)}
}

func (m *Customer) Addrs() []any {
	return []any{&m.ID, &m.Name, &m.Address, &m.Nicknames, &m.Status}
}
