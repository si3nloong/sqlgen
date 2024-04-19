{{- reserveImport "context" }}
{{- reserveImport "database/sql" }}
{{- reserveImport "strings" }}
{{- reserveImport "strconv" }}
{{- reserveImport "sync" }}
{{- reserveImport "github.com/si3nloong/sqlgen/sequel" }}
{{- reserveImport "github.com/si3nloong/sqlgen/sequel/strpool" }}

const _getTableSQL = "SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = {{ var 1 }} LIMIT 1;"

{{ $dialect := dialectVar }}
func InsertOne[T sequel.TableColumnValuer[T], Ptr interface {
	sequel.TableColumnValuer[T]
	sequel.Scanner[T]
}](ctx context.Context, db sequel.DB, v Ptr) (sql.Result, error) {
	args := v.Values()
	switch vi := any(v).(type) {
	case sequel.SingleInserter:
		switch vk := vi.(type) {
		case sequel.AutoIncrKeyer:
			_, idx, _ := vk.PK()
			args = append(args[:idx], args[idx+1:]...)
		}
		return db.ExecContext(ctx, vi.InsertOneStmt(), args...)
	}

	columns := v.Columns()
	switch vi := any(v).(type) {
	case sequel.AutoIncrKeyer:
		// If it's a auto increment primary key
		// We don't need to pass the value
		_, idx, _ := vi.PK()
		columns = append(columns[:idx], columns[idx+1:]...)
		args = append(args[:idx], args[idx+1:]...)
	}
	{{ if not $dialect.IsVarSame -}}
	stmt := strpool.AcquireString()
	defer strpool.ReleaseString(stmt)
	stmt.WriteString("INSERT INTO " + v.TableName() + " (" + strings.Join(columns, ",") + ") VALUES (")
	for i := range args {
		if i > 0 {
			stmt.WriteByte(',')
		}
		stmt.WriteString(wrapVar(i))
	}
	stmt.WriteString(");")
	return db.ExecContext(ctx, stmt.String(), args...)
	{{ else -}}
	return db.ExecContext(ctx, "INSERT INTO "+ v.TableName() +" ("+ strings.Join(columns, ",") +") VALUES ("+ strings.Repeat({{ quote (print "," $dialect.Var) }}, len(columns))[1:]+")", args...)
	{{ end -}}
}

// InsertInto is a helper function to insert your records.
func InsertInto[T sequel.TableColumnValuer[T]](ctx context.Context, db sequel.DB, data []T) (sql.Result, error) {
	n := len(data)
	if n == 0 {
		return new(sequel.EmptyResult), nil
	}

	var (
		model    T
		columns  = model.Columns()
		idx      = -1
		noOfCols = len(columns)
		args     = make([]any, 0, noOfCols*len(data))
	)

	switch vi := any(model).(type) {
	case sequel.AutoIncrKeyer:
		_, idx, _ = vi.PK()
		noOfCols--
		columns = append(columns[:idx], columns[idx+1:]...)

	case sequel.Inserter:
		query := strings.Repeat(vi.InsertVarQuery(), len(data))
		return db.ExecContext(ctx, "INSERT INTO "+model.TableName()+" ("+strings.Join(columns, ",")+") VALUES "+query[:len(query)-1]+";", args...)
	}

	var stmt = strpool.AcquireString()
	defer strpool.ReleaseString(stmt)
	stmt.WriteString("INSERT INTO " + model.TableName() + " (" + strings.Join(columns, ",") + ") VALUES ")
	for i := range data {
		if i > 0 {
			stmt.WriteString(",(")
		} else {
			stmt.WriteByte('(')
		}
		{{ if not $dialect.IsVarSame -}}
		offset := noOfCols * i
		{{ end -}}
		for j := 0; j < noOfCols; j++ {
			if j > 0 {
				stmt.WriteByte(',')
			}
			{{ if $dialect.IsVarSame -}}
			stmt.WriteString({{ quote $dialect.Var }})
			{{ else -}}
			stmt.WriteString(wrapVar(offset + j))
			{{ end -}}
		}
		if idx > -1 {
			values := data[i].Values()
			values = append(values[:idx], values[idx+1:]...)
			args = append(args, values...)
		} else {
			args = append(args, data[i].Values()...)
		}
		stmt.WriteByte(')')
	}
	stmt.WriteByte(';')
	return db.ExecContext(ctx, stmt.String(), args...)
}

// FindByPK is to find single record using primary key.
func FindByPK[T sequel.KeyValuer[T], Ptr sequel.KeyValueScanner[T]](ctx context.Context, db sequel.DB, v Ptr) error {
	switch vi := any(v).(type) {
	case sequel.KeyFinder:
		_, _, pk := vi.PK()
		return db.QueryRowContext(ctx, vi.FindByPKStmt(), pk).Scan(v.Addrs()...)
	}

	var (
		pkName, _, pk = v.PK()
		columns       = v.Columns()
	)
	return db.QueryRowContext(ctx, "SELECT "+strings.Join(columns, ",")+" FROM "+v.TableName()+" WHERE "+pkName+" = {{ var 1 }} LIMIT 1;", pk).Scan(v.Addrs()...)
}

// UpdateByPK is to update single record using primary key.
func UpdateByPK[T sequel.KeyValuer[T]](ctx context.Context, db sequel.DB, v T) (sql.Result, error) {
	var (
		pkName, idx, pk = v.PK()
		columns, values = v.Columns(), v.Values()
	)
	switch vi := any(v).(type) {
	case sequel.KeyUpdater:
		values = append(values[:idx], append(values[idx+1:], pk)...)
		return db.ExecContext(ctx, vi.UpdateByPKStmt(), values...)

	default:
		columns = append(columns[:idx], columns[idx+1:]...)
		values = append(values[:idx], values[idx+1:]...)
		{{ if $dialect.IsVarSame -}}
		return db.ExecContext(ctx, "UPDATE "+v.TableName()+" SET "+strings.Join(columns, " = {{ var 1 }},")+" = {{ var 1 }} WHERE "+pkName+" = {{ var 1 }};", append(values, pk)...)
		{{ else -}}
		stmt := strpool.AcquireString()
		defer strpool.ReleaseString(stmt)
		stmt.WriteString("UPDATE "+ v.TableName()+ " SET ")
		for idx := range columns {
			if idx > 0 {
				stmt.WriteByte(',')
			}
			stmt.WriteString(columns[idx] +" = "+ wrapVar(idx + 1))
		}
		stmt.WriteString(" WHERE "+ pkName +" = "+ wrapVar(len(columns) + 2)+ ";")
		return db.ExecContext(ctx, stmt.String(), append(values, pk)...)
		{{ end -}}

	}
}

// DeleteByPK is to update single record using primary key.
func DeleteByPK[T sequel.KeyValuer[T]](ctx context.Context, db sequel.DB, v T) (sql.Result, error) {
	pkName, _, pk := v.PK()
	return db.ExecContext(ctx, "DELETE FROM "+v.TableName()+" WHERE "+pkName+" = {{ var 1 }};", pk)
}

// Migrate is to create or alter the table based on the defined schemas.
func Migrate[T sequel.Migrator](ctx context.Context, db sequel.DB) error {
	var (
		v           T
		table       string
		tableExists bool
	)
	tableName, _ := strconv.Unquote(v.TableName())
	if err := db.QueryRowContext(ctx, _getTableSQL, tableName).Scan(&table); err != nil {
		tableExists = false
	} else {
		tableExists = true
	}
	if tableExists {
		if _, err := db.ExecContext(ctx, v.AlterTableStmt()); err != nil {
			return err
		}
		return nil
	}
	if _, err := db.ExecContext(ctx, v.CreateTableStmt()); err != nil {
		return err
	}
	return nil
}

type SelectStmt struct {
	Select    []string
	FromTable string
	Where     sequel.WhereClause
	OrderBy   []sequel.OrderByClause
	Offset	  uint64
	Limit     uint16
}

func QueryStmt[T any, Ptr interface {
	*T
	sequel.Scanner[T]
}, Stmt interface{ SelectStmt }](ctx context.Context, dbConn sequel.DB, stmt Stmt) ([]T, error) {
	blr := AcquireStmt()
	defer ReleaseStmt(blr)

	switch vi := any(stmt).(type) {
	case SelectStmt:
		blr.WriteString("SELECT " + strings.Join(vi.Select, ",") + " FROM " + vi.FromTable)
		if vi.Where != nil {
			blr.WriteString(" WHERE ")
			vi.Where(blr)
		}
		if len(vi.OrderBy) > 0 {
			blr.WriteString(" ORDER BY ")
			for i := range vi.OrderBy {
				if i > 0 {
					blr.WriteByte(',')
				}
				vi.OrderBy[i](blr)
			}
		}
		if vi.Limit > 0 {
			blr.WriteString(" LIMIT " + strconv.FormatUint(uint64(vi.Limit), 10))
		}
		if vi.Offset > 0 {
			blr.WriteString(" OFFSET " + strconv.FormatUint(vi.Offset, 10))
		}
		blr.WriteByte(';')
	}

	rows, err := dbConn.QueryContext(ctx, blr.String(), blr.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ReleaseStmt(blr)

	var result []T
	for rows.Next() {
		var v T
		if err := rows.Scan(Ptr(&v).Addrs()...); err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return result, nil
}

type UpdateStmt struct {
	FromTable string
	Set       []string
	Where     sequel.WhereClause
	OrderBy   []sequel.OrderByClause
	Limit     uint16
}

type DeleteStmt struct {
	FromTable string
	Where     sequel.WhereClause
	OrderBy   []sequel.OrderByClause
	Limit     uint16
}

func ExecStmt[T any, Stmt interface {
	UpdateStmt | DeleteStmt
}](ctx context.Context, dbConn sequel.DB, stmt Stmt) (sql.Result, error) {
	blr := AcquireStmt()
	defer ReleaseStmt(blr)

	switch vi := any(stmt).(type) {
	case UpdateStmt:
		blr.WriteString("UPDATE " + vi.FromTable)
		if vi.Where != nil {
			blr.WriteString(" WHERE ")
			vi.Where(blr)
		}
		if len(vi.Set) > 0 {
			blr.WriteString(" SET ")
		}
		if len(vi.OrderBy) > 0 {
			blr.WriteString(" ORDER BY ")
			for i := range vi.OrderBy {
				if i > 0 {
					blr.WriteByte(',')
				}
				vi.OrderBy[i](blr)
			}
		}
		if vi.Limit > 0 {
			blr.WriteString(" LIMIT " + strconv.FormatUint(uint64(vi.Limit), 10))
		}
		blr.WriteByte(';')

	case DeleteStmt:
		blr.WriteString("DELETE FROM " + vi.FromTable)
		if vi.Where != nil {
			blr.WriteString(" WHERE ")
			vi.Where(blr)
		}
		if len(vi.OrderBy) > 0 {
			blr.WriteString(" ORDER BY ")
			for i := range vi.OrderBy {
				if i > 0 {
					blr.WriteByte(',')
				}
				vi.OrderBy[i](blr)
			}
		}
		if vi.Limit > 0 {
			blr.WriteString(" LIMIT " + strconv.FormatUint(uint64(vi.Limit), 10))
		}
		blr.WriteByte(';')
	}
	return dbConn.ExecContext(ctx, blr.String(), blr.Args()...)
}

var (
	pool = sync.Pool{
		New: func() any {
			return new(sqlStmt)
		},
	}
)

func AcquireStmt() sequel.Stmt {
	return pool.Get().(*sqlStmt)
}

func ReleaseStmt(stmt sequel.Stmt) {
	stmt.Reset()
	pool.Put(stmt)
}

type sqlStmt struct {
	strings.Builder
	pos  int
	args []any
}

func (s *sqlStmt) Var(query string, value any) {
	s.pos++
	{{ if $dialect.IsVarSame -}}
	s.WriteString(query+"?")
	{{ else -}}
	s.WriteString(wrapVar(s.pos))
	{{ end -}}
	s.args = append(s.args, value)
}

func (s *sqlStmt) Vars(query string, values []any) {
	s.WriteString(query)
	noOfLen := len(values)
	{{ if $dialect.IsVarSame -}}
	s.WriteString("(" + strings.Repeat(",?", noOfLen)[1:] + ")")
	{{ else -}}
	s.WriteByte('(')
	i := s.pos
	s.pos += noOfLen
	for ; i < s.pos; i++ {
		s.WriteString(wrapVar(i + 1))
	}
	s.WriteByte(')')
	{{ end -}}
	s.args = append(s.args, values...)
}

func (s sqlStmt) Args() []any {
	return s.args
}

func (s *sqlStmt) Reset() {
	s.args = nil
	s.pos = 0
	s.Builder.Reset()
}

func QuoteIdentifier(v string) string {
	return string({{ quoteChar }}) + v + string({{ quoteChar }})
}

{{ if not $dialect.IsVarSame -}}
func wrapVar(i int) string {
	return "{{ $dialect.Var }}"+ strconv.Itoa(i)
}
{{ end }}
