{{- reserveImport "context" }}
{{- reserveImport "database/sql" }}
{{- reserveImport "strings" }}
{{- reserveImport "strconv" }}
{{- reserveImport "sync" }}
{{- reserveImport "github.com/si3nloong/sqlgen/sequel" }}
{{- reserveImport "github.com/si3nloong/sqlgen/sequel/strpool" }}

const _getTableSQL = "SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = {{ quoteVar 1 }} LIMIT 1;"

func InsertOne[T sequel.TableColumnValuer[T], Ptr interface {
	sequel.TableColumnValuer[T]
	sequel.Scanner[T]
}](ctx context.Context, sqlConn sequel.DB, model Ptr) (sql.Result, error) {
	var (
		args = model.Values()
		columns []string
	)
	switch vi := any(model).(type) {
	case sequel.SingleInserter:
		switch vk := vi.(type) {
		case sequel.AutoIncrKeyer:
			_, idx, _ := vk.PK()
			args = append(args[:idx], args[idx+1:]...)
		}
		return sqlConn.ExecContext(ctx, vi.InsertOneStmt(), args...)

	case sequel.AutoIncrKeyer:
		// If it's an AUTO_INCREMENT primary key
		// We don't need to pass the value
		_, idx, _ := vi.PK()
		columns = model.Columns()
		columns = append(columns[:idx], columns[idx+1:]...)
		args = append(args[:idx], args[idx+1:]...)

	default:
		columns = model.Columns()
	}

	{{ if not isStaticVar -}}
	stmt := strpool.AcquireString()
	defer strpool.ReleaseString(stmt)
	stmt.WriteString("INSERT INTO " + model.TableName() + " (" + strings.Join(columns, ",") + ") VALUES (")
	for i := range args {
		if i > 0 {
			stmt.WriteByte(',')
		}
		stmt.WriteString(wrapVar(i))
	}
	stmt.WriteString(");")
	return sqlConn.ExecContext(ctx, stmt.String(), args...)
	{{ else -}}
	return sqlConn.ExecContext(ctx, "INSERT INTO "+ model.TableName() +" ("+ strings.Join(columns, ",") +") VALUES ("+ strings.Repeat({{ quote (print "," varRune) }}, len(columns))[1:]+")", args...)
	{{ end -}}
}

// InsertInto is a helper function to insert your records.
func InsertInto[T sequel.TableColumnValuer[T]](ctx context.Context, sqlConn sequel.DB, data []T) (sql.Result, error) {
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
		query := strings.Repeat(vi.InsertVarQuery()+",", len(data))
		return sqlConn.ExecContext(ctx, "INSERT INTO "+model.TableName()+" ("+strings.Join(columns, ",")+") VALUES "+query[:len(query)-1]+";", args...)
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
		{{ if not isStaticVar -}}
		offset := noOfCols * i
		{{ end -}}
		for j := 0; j < noOfCols; j++ {
			if j > 0 {
				stmt.WriteByte(',')
			}
			{{ if isStaticVar -}}
			stmt.WriteString({{ quote varRune }})
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
	return sqlConn.ExecContext(ctx, stmt.String(), args...)
}

// FindByPK is to find single record using primary key.
func FindByPK[T sequel.KeyValuer[T], Ptr sequel.KeyValueScanner[T]](ctx context.Context, sqlConn sequel.DB, model Ptr) error {
	switch vi := any(model).(type) {
	case sequel.KeyFinder:
		_, _, pk := vi.PK()
		return sqlConn.QueryRowContext(ctx, vi.FindByPKStmt(), pk).Scan(model.Addrs()...)
	}

	var (
		pkName, _, pk = model.PK()
		columns       = model.Columns()
	)
	return sqlConn.QueryRowContext(ctx, "SELECT "+strings.Join(columns, ",")+" FROM "+model.TableName()+" WHERE "+pkName+" = {{ quoteVar 1 }} LIMIT 1;", pk).Scan(model.Addrs()...)
}

// UpdateByPK is to update single record using primary key.
func UpdateByPK[T sequel.KeyValuer[T]](ctx context.Context, sqlConn sequel.DB, model T) (sql.Result, error) {
	var (
		pkName, idx, pk = model.PK()
		columns, values = model.Columns(), model.Values()
	)
	switch vi := any(model).(type) {
	case sequel.KeyUpdater:
		values = append(values[:idx], append(values[idx+1:], pk)...)
		return sqlConn.ExecContext(ctx, vi.UpdateByPKStmt(), values...)

	default:
		columns = append(columns[:idx], columns[idx+1:]...)
		values = append(values[:idx], values[idx+1:]...)
		{{ if isStaticVar -}}
		return sqlConn.ExecContext(ctx, "UPDATE "+model.TableName()+" SET "+strings.Join(columns, " = {{ quoteVar 1 }},")+" = {{ quoteVar 1 }} WHERE "+pkName+" = {{ quoteVar 1 }};", append(values, pk)...)
		{{ else -}}
		stmt := strpool.AcquireString()
		defer strpool.ReleaseString(stmt)
		stmt.WriteString("UPDATE "+ model.TableName()+ " SET ")
		for idx := range columns {
			if idx > 0 {
				stmt.WriteByte(',')
			}
			stmt.WriteString(columns[idx] +" = "+ wrapVar(idx + 1))
		}
		stmt.WriteString(" WHERE "+ pkName +" = "+ wrapVar(len(columns) + 2)+ ";")
		return sqlConn.ExecContext(ctx, stmt.String(), append(values, pk)...)
		{{ end -}}
	}
}

// DeleteByPK is to update single record using primary key.
func DeleteByPK[T sequel.KeyValuer[T]](ctx context.Context, sqlConn sequel.DB, model T) (sql.Result, error) {
	pkName, _, pk := model.PK()
	return sqlConn.ExecContext(ctx, "DELETE FROM "+model.TableName()+" WHERE "+pkName+" = {{ quoteVar 1 }};", pk)
}

// Migrate is to create or alter the table based on the defined schemas.
func Migrate[T sequel.Migrator](ctx context.Context, sqlConn sequel.DB) error {
	var (
		v           T
		table       string
		tableExists bool
	)
	tableName, _ := strconv.Unquote(v.TableName())
	if err := sqlConn.QueryRowContext(ctx, _getTableSQL, tableName).Scan(&table); err != nil {
		tableExists = false
	} else {
		tableExists = true
	}
	if tableExists {
		if _, err := sqlConn.ExecContext(ctx, v.AlterTableStmt()); err != nil {
			return err
		}
		return nil
	}
	if _, err := sqlConn.ExecContext(ctx, v.CreateTableStmt()); err != nil {
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
}, Stmt interface{ SelectStmt }](ctx context.Context, sqlConn sequel.DB, stmt Stmt) ([]T, error) {
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

	rows, err := sqlConn.QueryContext(ctx, blr.String(), blr.Args()...)
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
}](ctx context.Context, sqlConn sequel.DB, stmt Stmt) (sql.Result, error) {
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
	return sqlConn.ExecContext(ctx, blr.String(), blr.Args()...)
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
	{{ if isStaticVar -}}
	s.WriteString(query+"?")
	{{ else -}}
	s.WriteString(wrapVar(s.pos))
	{{ end -}}
	s.args = append(s.args, value)
}

func (s *sqlStmt) Vars(query string, values []any) {
	s.WriteString(query)
	noOfLen := len(values)
	{{ if isStaticVar -}}
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

{{ if not isStaticVar -}}
func wrapVar(i int) string {
	return {{ quote varRune }}+ strconv.Itoa(i)
}
{{ end }}
