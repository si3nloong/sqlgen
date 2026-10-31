package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/si3nloong/sqlgen/sequel"
)

// A real database/sql row pipeline with controlled results and captured SQL.
type recorder struct {
	query string
	args  []any
	count int
	calls int
}

func (r *recorder) ExecContext(_ context.Context, q string, args ...any) (sql.Result, error) {
	r.query = q
	r.args = args
	r.calls++
	return driver.RowsAffected(1), nil
}
func (r *recorder) QueryContext(ctx context.Context, q string, args ...any) (*sql.Rows, error) {
	db := sql.OpenDB(connector{r})
	rows, err := db.QueryContext(ctx, q, args...)
	db.Close()
	return rows, err
}
func (r *recorder) QueryRowContext(ctx context.Context, q string, args ...any) *sql.Row {
	db := sql.OpenDB(connector{r})
	row := db.QueryRowContext(ctx, q, args...)
	db.Close()
	return row
}

type connector struct{ r *recorder }

func (c connector) Connect(context.Context) (driver.Conn, error) { return connection{c.r}, nil }
func (c connector) Driver() driver.Driver                        { return testDriver{} }

type testDriver struct{}

func (testDriver) Open(string) (driver.Conn, error) { return nil, errors.New("unused") }

type connection struct{ r *recorder }

func (connection) Prepare(string) (driver.Stmt, error) { return nil, errors.New("unused") }
func (connection) Close() error                        { return nil }
func (connection) Begin() (driver.Tx, error)           { return nil, errors.New("unused") }
func (c connection) QueryContext(_ context.Context, q string, args []driver.NamedValue) (driver.Rows, error) {
	c.r.query = q
	c.r.calls++
	c.r.args = nil
	for _, a := range args {
		c.r.args = append(c.r.args, a.Value)
	}
	return &testRows{remaining: c.r.count}, nil
}

type testRows struct{ remaining int }

func (*testRows) Columns() []string { return []string{"id"} }
func (*testRows) Close() error      { return nil }
func (r *testRows) Next(dest []driver.Value) error {
	if r.remaining == 0 {
		return io.EOF
	}
	dest[0] = int64(r.remaining)
	r.remaining--
	return nil
}

type model struct{ ID int64 }

func (model) TableName() string        { return "items" }
func (model) Columns() []string        { return []string{"id"} }
func (m model) Values() []any          { return []any{m.ID} }
func (model) HasPK()                   {}
func (m model) PK() (string, int, any) { return "id", 0, m.ID }
func (m *model) Addrs() []any          { return []any{&m.ID} }

type noTable struct{ ID int64 }

func (m *noTable) Addrs() []any { return []any{&m.ID} }

func TestEmptyVars(t *testing.T) {
	s := AcquireStmt()
	defer ReleaseStmt(s)
	first := s.Var(1)
	if got := s.Vars(nil); got != "()" {
		t.Fatalf("empty Vars: %q", got)
	}
	if !reflect.DeepEqual(s.Args(), []any{1}) {
		t.Fatal(s.Args())
	}
	if next := s.Var(2); first != "[[ quoteVar 1 ]]" || next != "[[ quoteVar 2 ]]" {
		t.Fatalf("placeholders: %s %s", first, next)
	}
}
func TestPaginationMaximumLimit(t *testing.T) {
	for _, prev := range []bool{false, true} {
		r := &recorder{count: 65536}
		p := Paginate[model](PaginateStmt{Limit: 65535})
		seq := p.Next(context.Background(), r)
		if prev {
			seq = p.Prev(context.Background(), r)
		}
		calls := 0
		seq(func(data []model, err error) bool {
			calls++
			if err != nil {
				t.Fatal(err)
			}
			if len(data) != 65535 {
				t.Fatalf("prev=%v: got %d rows", prev, len(data))
			}
			return false
		})
		if calls != 1 || !strings.HasSuffix(r.query, " LIMIT 65536;") {
			t.Fatalf("calls=%d SQL=%s", calls, r.query)
		}
	}
}
type orderColumn struct { name string; ascending bool }
func (o orderColumn) ColumnName() string { return o.name }
func (o orderColumn) Asc() bool { return o.ascending }
func (o orderColumn) Desc() bool { return !o.ascending }

func TestSelectOneStmtWriteSQL(t *testing.T) {
	var w SqlStmt
	if err := (SelectOneStmt{}).WriteSQL(&w); err == nil || w.Query() != "" {
		t.Fatal("missing table must fail before writing")
	}
	if err := (SelectOneStmt{FromTable: "items"}).WriteSQL(&w); err != nil {
		t.Fatal(err)
	}
	if w.Query() != "SELECT * FROM items LIMIT 1;" {
		t.Fatal(w.Query())
	}
	w.Reset()
	s := SelectOneStmt{
		Select: []string{"id", "name"}, FromTable: "items",
		Where: func(w sequel.StmtWriter) { w.WriteString("id = " + w.Var(42)) },
		GroupBy: []string{"id", "name"},
		OrderBy: []sequel.OrderByClause{orderColumn{"id", false}, orderColumn{"name", true}},
	}
	if err := s.WriteSQL(&w); err != nil { t.Fatal(err) }
	want := "SELECT id,name FROM items WHERE id = [[ quoteVar 1 ]] GROUP BY id,name ORDER BY id DESC,name LIMIT 1;"
	if w.Query() != want || !reflect.DeepEqual(w.Args(), []any{42}) {
		t.Fatalf("got %s %v, want %s", w.Query(), w.Args(), want)
	}
}

func TestQueryOneStmt(t *testing.T) {
	r := &recorder{count: 1}
	where := func(w sequel.StmtWriter) { w.WriteString("id = " + w.Var(int64(1))) }
	got, err := QueryOneStmt[model](context.Background(), r, func(model) SelectOneStmt { return SelectOneStmt{Where: where, GroupBy: []string{"id"}} })
	if err != nil || got.ID != 1 {
		t.Fatalf("%v %v", got, err)
	}
	var expected SqlStmt
	if err := (SelectStmt{Select: []string{"id"}, FromTable: "items", Where: where, GroupBy: []string{"id"}, Limit: 1}).WriteSQL(&expected); err != nil {
		t.Fatal(err)
	}
	if r.query != expected.Query() || !reflect.DeepEqual(r.args, expected.Args()) {
		t.Fatalf("got %s %v, want %s %v", r.query, r.args, expected.Query(), expected.Args())
	}
	before := r.calls
	if _, err := QueryOneStmt[noTable](context.Background(), r, func(noTable) SelectOneStmt { return SelectOneStmt{} }); err == nil || r.calls != before {
		t.Fatal("missing table should fail before querying")
	}
}
func TestFormatPreservesArguments(t *testing.T) {
	s := AcquireStmt()
	defer ReleaseStmt(s)
	s.WriteString("id = " + s.Var(int64(42)))
	for i := 0; i < 2; i++ {
		if got := fmt.Sprintf("%+v", s); got != "id = 42" {
			t.Fatal(got)
		}
	}
	if !reflect.DeepEqual(s.Args(), []any{int64(42)}) {
		t.Fatal(s.Args())
	}
}

[[ if ne driver "postgres" ]]
type autoModel struct {
	ID     int64
	Name   string
	Active bool
}

func (autoModel) TableName() string        { return "items" }
func (autoModel) Columns() []string        { return []string{"id", "name", "active"} }
func (m autoModel) Values() []any          { return []any{m.ID, m.Name, m.Active} }
func (autoModel) HasPK()                   {}
func (autoModel) IsAutoIncr()              {}
func (m autoModel) PK() (string, int, any) { return "id", 0, m.ID }

func TestBulkInsertPlaceholders(t *testing.T) {
	r := new(recorder)
	_, err := Insert(context.Background(), r, []autoModel{{Name: "Ada", Active: true}, {Name: "Bob", Active: false}})
	if err != nil {
		t.Fatal(err)
	}
	if r.query != "INSERT INTO items (name,active) VALUES (?,?),(?,?);" || !reflect.DeepEqual(r.args, []any{"Ada", true, "Bob", false}) {
		t.Fatalf("%s %v", r.query, r.args)
	}

}

type composite struct {
	A, B       int
	Name, Note string
}

func (composite) TableName() string { return "items" }
func (composite) Columns() []string { return []string{"a", "b", "name", "note"} }
func (m composite) Values() []any   { return []any{m.A, m.B, m.Name, m.Note} }
func (composite) HasPK()            {}
func (m composite) CompositeKey() ([]string, []int, []any) {
	return []string{"a", "b"}, []int{0, 1}, []any{m.A, m.B}
}
func (m *composite) Addrs() []any { return []any{&m.A, &m.B, &m.Name, &m.Note} }
func TestCompositeUpsert(t *testing.T) {
	m := composite{1, 2, "Ada", "note"}
	for _, omitted := range [][]string{nil, {"name"}, {"note"}} {
		r := new(recorder)
		if _, err := UpsertOne(context.Background(), r, &m, true, omitted...); err != nil {
			t.Fatal(err)
		}
		updates := "name = VALUES(name),note = VALUES(note)"
		if len(omitted) > 0 {
			if omitted[0] == "name" {
				updates = "note = VALUES(note)"
			} else {
				updates = "name = VALUES(name)"
			}
		}
		want := "INSERT INTO items (a,b,name,note) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE " + updates
		if r.query != want || !reflect.DeepEqual(r.args, m.Values()) {
			t.Fatalf("got %s %v, want %s", r.query, r.args, want)
		}
	}
	r := new(recorder)
	if _, err := UpsertOne(context.Background(), r, &m, true, "name", "note"); err == nil || r.calls != 0 {
		t.Fatal("expected no-updatable-columns error")
	}
}

[[ end ]]

type plainStringWriter struct { io.StringWriter }

func TestwriteJoin(t *testing.T) {
	for _, values := range [][]string{nil, {}, {""}, {"id"}, {"id", "name"}, {"", "name", ""}, {"名字", "id"}} {
		var builder strings.Builder
		var stmt SqlStmt
		var plain strings.Builder
		for _, w := range []io.StringWriter{&builder, &stmt, plainStringWriter{&plain}} {
			w.WriteString("prefix:")
			writeJoin(w, values)
		}
		want := "prefix:" + strings.Join(values, ",")
		if builder.String() != want || stmt.Query() != want || plain.String() != want {
			t.Fatalf("values=%q: %q / %q / %q; want %q", values, builder.String(), stmt.Query(), plain.String(), want)
		}
	}
}

var sqlSink string

func BenchmarkwriteJoin(b *testing.B) {
	cols := []string{"id", "name", "email", "created_at", "updated_at", "status", "account_id", "deleted_at"}
	b.Run("JoinThenWrite", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			var w strings.Builder
			w.WriteString("SELECT ")
			w.WriteString(strings.Join(cols, ","))
			w.WriteString(" FROM items;")
			sqlSink = w.String()
		}
	})
	b.Run("writeJoin", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			var w strings.Builder
			w.WriteString("SELECT ")
			writeJoin(&w, cols)
			w.WriteString(" FROM items;")
			sqlSink = w.String()
		}
	})
}

[[ if eq driver "postgres" ]]
type postgresUpsertModel struct { model }
func (postgresUpsertModel) Columns() []string { return []string{"id", "name", "note", "active"} }
func (m postgresUpsertModel) Values() []any { return []any{m.ID, "Ada", "note", true} }
func (postgresUpsertModel) SQLInsertColumns() string { return "id,name,note,active" }
func (m postgresUpsertModel) SQLInsertValues() []any { return m.Values() }
func (postgresUpsertModel) SQLInsertPlaceholders(int) string { return "($1,$2,$3,$4)" }

func TestPostgresUpsertSeparators(t *testing.T) {
	for _, tc := range []struct { omit []string; updates string }{
		{nil, "name = EXCLUDED.name,note = EXCLUDED.note,active = EXCLUDED.active"},
		{[]string{"name"}, "note = EXCLUDED.note,active = EXCLUDED.active"},
		{[]string{"note"}, "name = EXCLUDED.name,active = EXCLUDED.active"},
		{[]string{"active"}, "name = EXCLUDED.name,note = EXCLUDED.note"},
		{[]string{"name", "note"}, "active = EXCLUDED.active"},
	} {
		r := &recorder{count: 1}
		m := postgresUpsertModel{}
		if err := UpsertOne(context.Background(), r, &m, WithOmitFields(tc.omit...)); err != nil { t.Fatal(err) }
		want := ") DO UPDATE SET " + tc.updates + " RETURNING id,name,note,active;"
		if !strings.HasSuffix(r.query, want) { t.Fatalf("omitted %v: got %s, want suffix %s", tc.omit, r.query, want) }
	}
}
[[ end ]]
