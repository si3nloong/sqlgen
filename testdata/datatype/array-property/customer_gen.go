package arrayproperty

func (Customer) Table() string {
	return "Customer"
}

func (Customer) Columns() []string {
    return []string{"id", "howOld", "Name", "Address", "Nicknames", "status"}
}

func (v Customer) Values() []any {
	return []any{v.ID, uint64(v.Age), v.Name, encoding.MarshalStringList(v.Address), encoding.MarshalStringList(v.Nicknames), string(v.Status)}
}

func (v *Customer) Addrs() []any {
	return []any{&v.ID, &v.Age, &v.Name, &v.Address, &v.Nicknames, &v.Status}
}