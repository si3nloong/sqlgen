package primitivestruct

func (Primitive) Table() string {
	return "Base"
}

func (Primitive) Columns() []string {
	return []string{"Str", "Bytes", "Bool", "Int", "Int8", "Int16", "Int32", "Int64", "Uint", "Uint8", "Uint16", "Uint32", "Uint64", "F32", "F64", "Time"}
}

func (m Primitive) Values() []any {
	return []any{m.Str, m.Bytes, m.Bool, int64(m.Int), int64(m.Int8), int64(m.Int16), int64(m.Int32), m.Int64, uint64(m.Uint), uint64(m.Uint8), uint64(m.Uint16), uint64(m.Uint32), m.Uint64, float64(m.F32), m.F64, any(m.Time)}
}

func (m *Primitive) Addrs() []any {
	return []any{&m.Str, &m.Bytes, &m.Bool, &m.Int, &m.Int8, &m.Int16, &m.Int32, &m.Int64, &m.Uint, &m.Uint8, &m.Uint16, &m.Uint32, &m.Uint64, &m.F32, &m.F64, &m.Time}
}
