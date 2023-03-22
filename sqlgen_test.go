package sqlgen

import (
	"testing"
)

type a struct {
	Name string
}

func (a) Columns() []string {
	return []string{"Name"}
}

func (a a) Values() []any {
	return []any{a.Name}
}

func (a) FieldNames() []string {
	return []string{"Name"}
}

func (a *a) Addrs() []any {
	return []any{&a.Name}
}

func TestQuery(t *testing.T) {
	// QueryScan[a](context.Background(), nil, "INSERT INTO (?,?)", nil)
}

func TestQueryScan(t *testing.T) {
	// QueryScan[a](context.Background(), nil, "", nil)
}
