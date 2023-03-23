package schema

func (A) Table() string {
	return "A"
}

func (A) Columns() []string {
    return []string{"ID", "CreatedAt"}
}

func (v A) Values() []any {
	return []any{v.ID, v.CreatedAt}
}

func (v *A) Addrs() []any {
	return []any{&v.ID, &v.CreatedAt}
}