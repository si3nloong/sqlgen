package testdata

import "database/sql/driver"

func (Common) Columns() []string {
    return []string{"i64", "Name", "Yes", "str", "Alias", "Int", "Uint", "Flag", "f32bit", "f64bit", "StrList", "Enum", "t"}
}

func (m Common) Values() []any {
	return []any{m.Int64, m.Name, m.Yes, string(m.Cstr), string(m.Alias), int64(m.Int), uint64(m.Uint), m.Flag, float64(m.F32), m.F64, m.StrList, int64(m.Enum), any(m.T)}
}

func (m Common) Addrs() []any {
	return []any{&m.Int64, &m.Name, &m.Yes, &m.Cstr, &m.Alias, &m.Int, &m.Uint, &m.Flag, &m.F32, &m.F64, &m.StrList, &m.Enum, &m.T}
}