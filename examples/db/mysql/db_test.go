package mysqldb

import (
	"reflect"
	"testing"

	"github.com/si3nloong/sqlgen/sequel"
)

// These tests exercise SQL rendering without requiring a database connection.
func TestSelectStmt(t *testing.T) {
	id := sequel.PrimitiveColumn("id", int64(0))
	name := sequel.PrimitiveColumn("name", "")
	tests := []struct {
		name  string
		stmt  SelectStmt
		query string
		args  []any
	}{
		{
			name: "default columns",
			stmt: SelectStmt{
				FromTable: "users",
			},
			query: "SELECT * FROM users;",
		},
		{
			name: "explicit columns",
			stmt: SelectStmt{
				Select:    []string{"id", "name"},
				FromTable: "users",
			},
			query: "SELECT id,name FROM users;",
		},
		{
			name: "where",
			stmt: SelectStmt{
				FromTable: "users",
				Where:     Equal(id, int64(42)),
			},
			query: "SELECT * FROM users WHERE id = ?;",
			args:  []any{int64(42)},
		},
		{
			name: "group by",
			stmt: SelectStmt{
				Select:    []string{"name", "COUNT(*)"},
				FromTable: "users",
				GroupBy:   []string{"name"},
			},
			query: "SELECT name,COUNT(*) FROM users GROUP BY name;",
		},
		{
			name: "ascending first",
			stmt: SelectStmt{
				FromTable: "users",
				OrderBy:   []sequel.OrderByClause{Asc(name), Desc(id)},
			},
			query: "SELECT * FROM users ORDER BY name,id DESC;",
		},
		{
			name: "descending first",
			stmt: SelectStmt{
				FromTable: "users",
				OrderBy:   []sequel.OrderByClause{Desc(name), Asc(id)},
			},
			query: "SELECT * FROM users ORDER BY name DESC,id;",
		},
		{
			name: "limit",
			stmt: SelectStmt{
				FromTable: "users",
				Limit:     10,
			},
			query: "SELECT * FROM users LIMIT 10;",
		},
		{
			name: "limit and offset",
			stmt: SelectStmt{
				FromTable: "users",
				Limit:     10,
				Offset:    20,
			},
			query: "SELECT * FROM users LIMIT 10 OFFSET 20;",
		},
		{
			name: "maximum pagination values",
			stmt: SelectStmt{
				FromTable: "users",
				Limit:     65535,
				Offset:    ^uint64(0),
			},
			query: "SELECT * FROM users LIMIT 65535 OFFSET 18446744073709551615;",
		},
		{
			name: "all clauses",
			stmt: SelectStmt{
				Select:    []string{"id", "name"},
				FromTable: "users",
				Where:     And(GreaterThan(id, int64(10)), Or(Equal(name, "Ada"), Equal(name, "Bob"))),
				GroupBy:   []string{"id", "name"},
				OrderBy:   []sequel.OrderByClause{Desc(name), Asc(id)},
				Limit:     5,
				Offset:    10,
			},
			query: "SELECT id,name FROM users WHERE (id > ? AND (name = ? OR name = ?)) GROUP BY id,name ORDER BY name DESC,id LIMIT 5 OFFSET 10;",
			args:  []any{int64(10), "Ada", "Bob"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { assertStatement(t, tt.stmt, tt.query, tt.args) })
	}
	t.Run("missing table", func(t *testing.T) {
		assertMissingTable(t, SelectStmt{
			Where: Equal(id, int64(42)),
		}, "missing table name for SELECT statement")
	})
}

func TestUpdateStmt(t *testing.T) {
	id := sequel.PrimitiveColumn("id", int64(0))
	name := sequel.PrimitiveColumn("name", "")
	active := sequel.PrimitiveColumn("active", false)
	tests := []struct {
		name  string
		stmt  UpdateStmt
		query string
		args  []any
	}{
		{
			name: "single assignment",
			stmt: UpdateStmt{
				Table: "users",
				Set:   []sequel.SetClause{Set(name, "Ada")},
			},
			query: "UPDATE users SET name = ?;",
			args:  []any{"Ada"},
		},
		{
			name: "multiple assignments",
			stmt: UpdateStmt{
				Table: "users",
				Set:   []sequel.SetClause{Set(name, "Ada"), Set(active, true)},
			},
			query: "UPDATE users SET name = ?,active = ?;",
			args:  []any{"Ada", true},
		},
		{
			name: "assignments before where arguments",
			stmt: UpdateStmt{
				Table: "users",
				Set:   []sequel.SetClause{Set(name, "Ada"), Set(active, true)},
				Where: In(id, []int64{42, 43}),
			},
			query: "UPDATE users SET name = ?,active = ? WHERE id IN (?,?);",
			args:  []any{"Ada", true, int64(42), int64(43)},
		},
		{
			name: "ordered id,name DESC",
			stmt: UpdateStmt{
				Table:   "users",
				Set:     []sequel.SetClause{Set(active, false)},
				Where:   GreaterThan(id, int64(10)),
				OrderBy: []sequel.OrderByClause{Asc(id), Desc(sequel.PrimitiveColumn("name", ""))},
				Limit:   5,
			},
			query: "UPDATE users SET active = ? WHERE id > ? ORDER BY id,name DESC LIMIT 5;",
			args:  []any{false, int64(10)},
		},
		{
			name: "ordered id DESC,name",
			stmt: UpdateStmt{
				Table:   "users",
				Set:     []sequel.SetClause{Set(active, false)},
				Where:   GreaterThan(id, int64(10)),
				OrderBy: []sequel.OrderByClause{Desc(id), Asc(sequel.PrimitiveColumn("name", ""))},
				Limit:   5,
			},
			query: "UPDATE users SET active = ? WHERE id > ? ORDER BY id DESC,name ASC LIMIT 5;",
			args:  []any{false, int64(10)},
		},
		{
			name: "limit only",
			stmt: UpdateStmt{
				Table: "users",
				Set:   []sequel.SetClause{Set(active, false)},
				Limit: 65535,
			},
			query: "UPDATE users SET active = ? LIMIT 65535;",
			args:  []any{false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { assertStatement(t, tt.stmt, tt.query, tt.args) })
	}
	t.Run("missing table", func(t *testing.T) {
		assertMissingTable(t, UpdateStmt{
			Set:   []sequel.SetClause{Set(name, "Ada")},
			Where: Equal(id, int64(42)),
		}, "missing table name for UPDATE statement")
	})
}

func TestDeleteStmt(t *testing.T) {
	id := sequel.PrimitiveColumn("id", int64(0))
	tests := []struct {
		name  string
		stmt  DeleteStmt
		query string
		args  []any
	}{
		{
			name: "all rows",
			stmt: DeleteStmt{
				FromTable: "users",
			},
			query: "DELETE FROM users;",
		},
		{
			name: "where",
			stmt: DeleteStmt{
				FromTable: "users",
				Where:     Equal(id, int64(42)),
			},
			query: "DELETE FROM users WHERE id = ?;",
			args:  []any{int64(42)},
		},
		{
			name: "multiple arguments",
			stmt: DeleteStmt{
				FromTable: "users",
				Where:     Between(id, int64(10), int64(20)),
			},
			query: "DELETE FROM users WHERE id BETWEEN ? AND ?;",
			args:  []any{int64(10), int64(20)},
		},
		{
			name: "ordered id,name DESC",
			stmt: DeleteStmt{
				FromTable: "users",
				Where:     GreaterThan(id, int64(10)),
				OrderBy:   []sequel.OrderByClause{Asc(id), Desc(sequel.PrimitiveColumn("name", ""))},
				Limit:     5,
			},
			query: "DELETE FROM users WHERE id > ? ORDER BY id,name DESC LIMIT 5;",
			args:  []any{int64(10)},
		},
		{
			name: "ordered id DESC,name",
			stmt: DeleteStmt{
				FromTable: "users",
				Where:     GreaterThan(id, int64(10)),
				OrderBy:   []sequel.OrderByClause{Desc(id), Asc(sequel.PrimitiveColumn("name", ""))},
				Limit:     5,
			},
			query: "DELETE FROM users WHERE id > ? ORDER BY id DESC,name LIMIT 5;",
			args:  []any{int64(10)},
		},
		{
			name: "limit only",
			stmt: DeleteStmt{
				FromTable: "users",
				Limit:     65535,
			},
			query: "DELETE FROM users LIMIT 65535;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { assertStatement(t, tt.stmt, tt.query, tt.args) })
	}
	t.Run("missing table", func(t *testing.T) {
		assertMissingTable(t, DeleteStmt{
			Where: Equal(id, int64(42)),
		}, "missing table name for DELETE FROM statement")
	})
}

type statementWriter interface {
	WriteSQL(sequel.Stmt) error
}

func assertStatement(t *testing.T, statement statementWriter, wantQuery string, wantArgs []any) {
	t.Helper()
	var w SqlStmt
	if err := statement.WriteSQL(&w); err != nil {
		t.Fatalf("WriteSQL() error: %v", err)
	}
	if got := w.Query(); got != wantQuery {
		t.Errorf("Query() = %q, want %q", got, wantQuery)
	}
	if got := w.Args(); !reflect.DeepEqual(got, wantArgs) {
		t.Errorf("Args() = %#v, want %#v", got, wantArgs)
	}
}

func assertMissingTable(t *testing.T, statement statementWriter, wantError string) {
	t.Helper()
	var w SqlStmt
	// Validation must leave even an existing writer untouched.
	w.WriteString("prefix " + w.Var("existing"))
	query := w.Query()
	err := statement.WriteSQL(&w)
	if err == nil || err.Error() != wantError {
		t.Fatalf("WriteSQL() error = %v, want %q", err, wantError)
	}
	if w.Query() != query || !reflect.DeepEqual(w.Args(), []any{"existing"}) {
		t.Errorf("invalid statement modified writer: query = %q, args = %#v", w.Query(), w.Args())
	}
}
