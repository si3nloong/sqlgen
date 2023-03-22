package customprimitivestruct

func (Custom) Columns() []string {
    return []string{"text", "e", "Num"}
}

func (m Custom) Values() []any {
	return []any{string(m.Str), int64(m.Enum), uint64(m.Num)}
}

func (m *Custom) Addrs() []any {
	return []any{&m.Str, &m.Enum, &m.Num}
}