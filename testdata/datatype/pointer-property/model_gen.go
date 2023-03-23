package pointerproperty

func (Ptr) Table() string {
	return "BasePtr"
}

func (Ptr) Columns() []string {
	return []string{"Str", "Bytes", "Bool", "Int", "Int8", "Int16", "Int32", "Int64", "Uint", "Uint8", "Uint16", "Uint32", "Uint64", "F32", "F64", "Time"}
}

func (m Ptr) Values() []any {
	return []any{any(m.Str), any(m.Bytes), any(m.Bool), any(m.Int), any(m.Int8), any(m.Int16), any(m.Int32), any(m.Int64), any(m.Uint), any(m.Uint8), any(m.Uint16), any(m.Uint32), any(m.Uint64), any(m.F32), any(m.F64), any(m.Time)}
}

func (m *Ptr) Addrs() []any {
	return []any{&m.Str, &m.Bytes, &m.Bool, &m.Int, &m.Int8, &m.Int16, &m.Int32, &m.Int64, &m.Uint, &m.Uint8, &m.Uint16, &m.Uint32, &m.Uint64, &m.F32, &m.F64, &m.Time}
}
