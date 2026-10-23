package mysqldb

import (
	"reflect"
	"strings"
	"testing"

	"github.com/si3nloong/sqlgen/sequel"
)

func TestComparisonOperators(t *testing.T) {
	columns := []struct {
		name        string
		column      sequel.ColumnClause[string]
		placeholder string
		arg         any
	}{
		{
			name:        "primitive",
			column:      sequel.PrimitiveColumn("name", "stored value"),
			placeholder: "?",
			arg:         "Ada%",
		},
		{
			name:        "converted",
			column:      sequel.Column("name", "stored value", func(v string) any { return strings.ToUpper(v) }),
			placeholder: "?",
			arg:         "ADA%",
		},
		{
			name:        "SQL expression",
			column:      sequel.SQLColumn("name", "stored value", func(p string) string { return "LOWER(" + p + ")" }, func(v string) any { return strings.ToUpper(v) }),
			placeholder: "LOWER(?)",
			arg:         "ADA%",
		},
	}
	operators := []struct {
		name     string
		operator func(sequel.ColumnClause[string], string) sequel.WhereClause
		sql      string
	}{
		{
			name:     "Equal",
			operator: Equal[string],
			sql:      "=",
		},
		{
			name:     "NotEqual",
			operator: NotEqual[string],
			sql:      "<>",
		},
		{
			name:     "GreaterThan",
			operator: GreaterThan[string],
			sql:      ">",
		},
		{
			name:     "GreaterThanOrEqual",
			operator: GreaterThanOrEqual[string],
			sql:      ">=",
		},
		{
			name:     "LessThan",
			operator: LessThan[string],
			sql:      "<",
		},
		{
			name:     "LessThanOrEqual",
			operator: LessThanOrEqual[string],
			sql:      "<=",
		},
		{
			name:     "Like",
			operator: Like[string],
			sql:      "LIKE",
		},
		{
			name:     "NotLike",
			operator: NotLike[string],
			sql:      "NOT LIKE",
		},
	}
	for _, column := range columns {
		t.Run(column.name, func(t *testing.T) {
			for _, operator := range operators {
				t.Run(operator.name, func(t *testing.T) {
					assertOperator(t, operator.operator(column.column, "Ada%"), "name "+operator.sql+" "+column.placeholder, []any{column.arg})
				})
			}
			t.Run("Set", func(t *testing.T) {
				assertOperator(t, Set(column.column, "Ada%"), "name = "+column.placeholder, []any{column.arg})
			})
		})
	}
}

func TestRangeOperators(t *testing.T) {
	columns := []struct {
		name   string
		column sequel.ColumnClause[string]
		bounds string
		args   []any
	}{
		{
			name:   "primitive",
			column: sequel.PrimitiveColumn("name", "stored value"),
			bounds: "? AND ?",
			args:   []any{"ada", "bob"},
		},
		{
			name: "converted",
			column: sequel.Column("name", "stored value", func(v string) any {
				return strings.ToUpper(v)
			}),
			bounds: "? AND ?",
			args:   []any{"ADA", "BOB"},
		},
		{
			name: "SQL expression",
			column: sequel.SQLColumn("name", "stored value", func(p string) string {
				return "LOWER(" + p + ")"
			}, func(v string) any {
				return strings.ToUpper(v)
			}),
			bounds: "LOWER(?) AND LOWER(?)",
			args:   []any{"ADA", "BOB"},
		},
	}
	for _, column := range columns {
		t.Run(column.name, func(t *testing.T) {
			t.Run("Between", func(t *testing.T) {
				assertOperator(t, Between(column.column, "ada", "bob"), "name BETWEEN "+column.bounds, column.args)
			})
			t.Run("NotBetween", func(t *testing.T) {
				assertOperator(t, NotBetween(column.column, "ada", "bob"), "name NOT BETWEEN "+column.bounds, column.args)
			})
		})
	}
}

func TestMembershipOperators(t *testing.T) {
	columns := []struct {
		name   string
		column sequel.ColumnClause[string]
		args   []any
	}{
		{
			name:   "primitive",
			column: sequel.PrimitiveColumn("name", "stored value"),
			args:   []any{"ada", "bob"},
		},
		{
			name: "converted",
			column: sequel.Column("name", "stored value", func(v string) any {
				return strings.ToUpper(v)
			}),
			args: []any{"ADA", "BOB"},
		},
		{
			name: "SQL column",
			column: sequel.SQLColumn("name", "stored value", func(p string) string {
				return p
			}, func(v string) any {
				return strings.ToUpper(v)
			}),
			args: []any{"ADA", "BOB"},
		},
	}
	operators := []struct {
		name     string
		operator func(sequel.ColumnClause[string], []string) sequel.WhereClause
		sql      string
	}{
		{
			name:     "In",
			operator: In[string],
			sql:      " IN ",
		},
		{
			name:     "NotIn",
			operator: NotIn[string],
			sql:      " NOT IN ",
		},
	}
	for _, column := range columns {
		t.Run(column.name, func(t *testing.T) {
			for _, operator := range operators {
				t.Run(operator.name, func(t *testing.T) {
					t.Run("multiple values", func(t *testing.T) {
						assertOperator(t, operator.operator(column.column, []string{"ada", "bob"}), "name"+operator.sql+"(?,?)", column.args)
					})
					t.Run("single value", func(t *testing.T) {
						assertOperator(t, operator.operator(column.column, []string{"ada"}), "name"+operator.sql+"(?)", column.args[:1])
					})
					// Empty lists render literally; callers decide how to handle them for their dialect.
					t.Run("nil values", func(t *testing.T) {
						assertOperator(t, operator.operator(column.column, nil), "name"+operator.sql+"()", nil)
					})
					t.Run("empty values", func(t *testing.T) {
						assertOperator(t, operator.operator(column.column, []string{}), "name"+operator.sql+"()", nil)
					})
				})
			}
		})
	}
}

func TestLogicalOperators(t *testing.T) {
	id := sequel.PrimitiveColumn("id", int64(0))
	tests := []struct {
		name   string
		clause sequel.WhereClause
		query  string
		args   []any
	}{
		{
			name:   "empty And",
			clause: And(),
		},
		{
			name:   "empty Or",
			clause: Or(),
		},
		{
			name:   "single And",
			clause: And(Equal(id, int64(1))),
			query:  "(id = ?)",
			args:   []any{int64(1)},
		},
		{
			name:   "single Or",
			clause: Or(Equal(id, int64(1))),
			query:  "(id = ?)",
			args:   []any{int64(1)},
		},
		{
			name:   "multiple And",
			clause: And(GreaterThan(id, int64(1)), LessThan(id, int64(10)), NotEqual(id, int64(5))),
			query:  "(id > ? AND id < ? AND id <> ?)",
			args:   []any{int64(1), int64(10), int64(5)},
		},
		{
			name:   "multiple Or",
			clause: Or(Equal(id, int64(1)), Equal(id, int64(2)), Equal(id, int64(3))),
			query:  "(id = ? OR id = ? OR id = ?)",
			args:   []any{int64(1), int64(2), int64(3)},
		},
		{
			name:   "nested groups",
			clause: And(Or(Equal(id, int64(1)), Equal(id, int64(2))), Between(id, int64(1), int64(10))),
			query:  "((id = ? OR id = ?) AND id BETWEEN ? AND ?)",
			args:   []any{int64(1), int64(2), int64(1), int64(10)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertOperator(t, tt.clause, tt.query, tt.args)
		})
	}
}

func TestNullOperators(t *testing.T) {
	column := sequel.PrimitiveColumn("name", "stored value")
	t.Run("IsNull", func(t *testing.T) {
		assertOperator(t, IsNull(column), "name IS NULL", nil)
	})
	t.Run("IsNotNull", func(t *testing.T) {
		assertOperator(t, IsNotNull(column), "name IS NOT NULL", nil)
	})
}

func TestOrderOperators(t *testing.T) {
	column := sequel.PrimitiveColumn("name", "stored value")
	tests := []struct {
		name      string
		clause    sequel.OrderByClause
		ascending bool
	}{
		{
			name:      "Asc",
			clause:    Asc(column),
			ascending: true,
		},
		{
			name:      "Desc",
			clause:    Desc(column),
			ascending: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.clause.ColumnName(); got != "name" {
				t.Errorf("ColumnName() = %q, want name", got)
			}
			if got := tt.clause.Asc(); got != tt.ascending {
				t.Errorf("Asc() = %v, want %v", got, tt.ascending)
			}
			if got := tt.clause.Desc(); got != !tt.ascending {
				t.Errorf("Desc() = %v, want %v", got, !tt.ascending)
			}
		})
	}
}

func assertOperator(t *testing.T, clause func(sequel.StmtWriter), wantQuery string, wantArgs []any) {
	t.Helper()
	var w SqlStmt
	clause(&w)
	if got := w.Query(); got != wantQuery {
		t.Errorf("Query() = %q, want %q", got, wantQuery)
	}
	if got := w.Args(); !reflect.DeepEqual(got, wantArgs) {
		t.Errorf("Args() = %#v, want %#v", got, wantArgs)
	}
}
