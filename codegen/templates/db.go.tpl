{{- reserveImport "context" }}
{{- reserveImport "database/sql" }}
{{- reserveImport "strings" }}
{{- reserveImport "strconv" }}
{{- reserveImport "sync" }}
{{- reserveImport "github.com/si3nloong/sqlgen/sequel" }}
{{- reserveImport "github.com/si3nloong/sqlgen/sequel/strpool" }}

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
	stmt.WriteString("INSERT INTO " + dbName(model) + model.TableName() + " (" + strings.Join(columns, ",") + ") VALUES (")
	for i := range args {
		if i > 0 {
			stmt.WriteByte(',')
		}
		// argument always started from 1
		stmt.WriteString(wrapVar(i + 1))
	}
	{{- /* postgres */ -}}
	{{ if eq driver "postgres" }}
	stmt.WriteString(") RETURNING "+ strings.Join(model.Columns(), ",") +";")
	if err := sqlConn.QueryRowContext(ctx, stmt.String(), args...).Scan(model.Addrs()...); err != nil {
		return nil, err
	}
	return new(sequel.EmptyResult), nil
	{{ else }}
	stmt.WriteString(");")
	return sqlConn.ExecContext(ctx, stmt.String(), args...)
	{{ end }}
	{{ else -}}
	return sqlConn.ExecContext(ctx, "INSERT INTO "+ dbName(model) + model.TableName() +" ("+ strings.Join(columns, ",") +") VALUES ("+ strings.Repeat({{ quote (print "," varRune) }}, len(columns))[1:]+");", args...)
	{{ end -}}
}

// Insert is a helper function to insert multiple records.
func Insert[T sequel.TableColumnValuer[T]](ctx context.Context, sqlConn sequel.DB, data []T) (sql.Result, error) {
	n := len(data)
	if n == 0 {
		return new(sequel.EmptyResult), nil
	}

	var (
		model    = data[0]
		columns  = model.Columns()
		idx      = -1
		noOfCols = len(columns)
		args     = make([]any, 0, noOfCols*len(data))
	)

	switch vi := any(model).(type) {
	case sequel.Inserter:
		query := strings.Repeat(vi.InsertVarQuery()+",", len(data))
		return sqlConn.ExecContext(ctx, "INSERT INTO "+dbName(model)+model.TableName()+" ("+strings.Join(columns, ",")+") VALUES "+query[:len(query)-1]+";", args...)

	case sequel.AutoIncrKeyer:
		_, idx, _ = vi.PK()
		noOfCols--
		columns = append(columns[:idx], columns[idx+1:]...)
	}

	var stmt = strpool.AcquireString()
	defer strpool.ReleaseString(stmt)
	stmt.WriteString("INSERT INTO " + dbName(model) + model.TableName() + " (" + strings.Join(columns, ",") + ") VALUES ")
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
			stmt.WriteString(wrapVar(offset + 1 + j))
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

func UpsertOne[T sequel.KeyValuer[T], Ptr sequel.KeyValueScanner[T]](ctx context.Context, sqlConn sequel.DB, model Ptr, override bool, omittedFields ...string) (sql.Result, error) {
	var (
		{{ if eq driver "mysql" -}}
		_, idx, _ = model.PK()
		{{ else -}}
		pkName, idx, _ = model.PK()
		{{ end -}}
		args      = model.Values()
		columns   []string
	)
	switch vi := any(model).(type) {
	case sequel.SingleUpserter:
		return sqlConn.ExecContext(ctx, vi.UpsertOneStmt(), args...)

	case sequel.AutoIncrKeyer:
		// If it's an AUTO_INCREMENT primary key
		// We don't need to pass the value
		columns = model.Columns()
		columns = append(columns[:idx], columns[idx+1:]...)
		args = append(args[:idx], args[idx+1:]...)

	default:
		columns = model.Columns()
	}

	var stmt = strpool.AcquireString()
	defer strpool.ReleaseString(stmt)
	{{ if eq driver "mysql" -}}
	if !override {
		stmt.WriteString("INSERT IGNORE INTO " + dbName(model) + model.TableName() + " (" + strings.Join(columns, ",") + ") VALUES (" + strings.Repeat(",?", len(columns))[1:] + ");")
	} else {
		dict := make(map[string]struct{})
		for i := range omittedFields {
			dict[omittedFields[i]] = struct{}{}
		}
		stmt.WriteString("INSERT INTO " + dbName(model) + model.TableName() + " (" + strings.Join(columns, ",") + ") VALUES (" + strings.Repeat(",?", len(columns))[1:] + ") ON DUPLICATE KEY UPDATE ")
		columns = append(columns[:idx], columns[idx+1:]...)
		args = append(args[:idx], args[idx+1:]...)
		noOfCols := len(columns)
		for i := range columns {
			if _, ok := dict[columns[i]]; ok {
				continue
			}
			if i < noOfCols-1 {
				stmt.WriteString(columns[i] + " =VALUES(" + columns[i] + "),")
			} else {
				stmt.WriteString(columns[i] + " =VALUES(" + columns[i] + ")")
			}
		}
		clear(dict)
	}
	{{ else -}}
	noOfCols := len(columns)
	stmt.WriteString("INSERT INTO " + dbName(model) + model.TableName() + " (" + strings.Join(columns, ",") + ") VALUES ")
	for j := 0; j < noOfCols; j++ {
		if j > 0 {
			stmt.WriteByte(',')
		}
		stmt.WriteString(wrapVar(noOfCols + j + 1))
	}
	stmt.WriteString(" ON CONFLICT(" + pkName + ")")
	if override {
		dict := map[string]struct{}{pkName: {}}
		for i := range omittedFields {
			dict[omittedFields[i]] = struct{}{}
		}
		stmt.WriteString(" DO UPDATE SET ")
		for i := range columns {
			if _, ok := dict[columns[i]]; ok {
				continue
			}
			if i < noOfCols-1 {
				stmt.WriteString(columns[i] + " = EXCLUDED." + columns[i] + ",")
			} else {
				stmt.WriteString(columns[i] + " = EXCLUDED." + columns[i])
			}
		}
		stmt.WriteByte(';')
		clear(dict)
	} else {
		stmt.WriteString(" DO NOTHING;")
	}
	{{ end -}}
	return sqlConn.ExecContext(ctx, stmt.String(), args...)
}

// Upsert is a helper function to upsert multiple records.
func Upsert[T sequel.KeyValuer[T]](ctx context.Context, sqlConn sequel.DB, data []T, override bool, omittedFields ...string) (sql.Result, error) {
	n := len(data)
	if n == 0 {
		return new(sequel.EmptyResult), nil
	}

	var (
		model    = data[0]
		columns  = model.Columns()
		noOfCols = len(columns)
		args     = make([]any, 0, noOfCols*len(data))
	)

	pkName, idx, _ := model.PK()
	switch any(model).(type) {
	case sequel.AutoIncrKeyer:
		noOfCols--
		columns = append(columns[:idx], columns[idx+1:]...)
	}

	var stmt = strpool.AcquireString()
	defer strpool.ReleaseString(stmt)
	{{ if eq driver "mysql" -}}
	if override {
		stmt.WriteString("INSERT INTO " + dbName(model) + model.TableName() + " (" + strings.Join(columns, ",") + ") VALUES ")
	} else {
		stmt.WriteString("INSERT IGNORE INTO " + dbName(model) + model.TableName() + " (" + strings.Join(columns, ",") + ") VALUES ")
	}
	{{ else -}}
	stmt.WriteString("INSERT INTO " + dbName(model) + model.TableName() + " (" + strings.Join(columns, ",") + ") VALUES ")
	{{ end -}}
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
			stmt.WriteString(wrapVar(offset + 1 + j))
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
	{{ if eq driver "mysql" -}}
	if override {
		stmt.WriteString(" ON DUPLICATE KEY UPDATE ")
		dict := map[string]struct{}{pkName: {}}
		for i := range omittedFields {
			dict[omittedFields[i]] = struct{}{}
		}
		for i := range columns {
			if _, ok := dict[columns[i]]; ok {
				continue
			}
			if i < noOfCols-1 {
				stmt.WriteString(columns[i] + " =VALUES(" + columns[i] + "),")
			} else {
				stmt.WriteString(columns[i] + " =VALUES(" + columns[i] + ")")
			}
		}
		clear(dict)
	}
	stmt.WriteByte(';')
	{{ else -}}
	if dupKey, ok := any(model).(sequel.DuplicateKeyer); ok {
		stmt.WriteString(" ON CONFLICT(" + strings.Join(dupKey.OnDuplicateKey(), ",") + ")")
	} else {
		stmt.WriteString(" ON CONFLICT(" + pkName + ")")
	}
	if override {
		dict := map[string]struct{}{pkName: {}}
		for i := range omittedFields {
			dict[omittedFields[i]] = struct{}{}
		}
		stmt.WriteString(" DO UPDATE SET ")
		for i := range columns {
			if _, ok := dict[columns[i]]; ok {
				continue
			}
			if i < noOfCols-1 {
				stmt.WriteString(columns[i] + " = EXCLUDED." + columns[i] + ",")
			} else {
				stmt.WriteString(columns[i] + " = EXCLUDED." + columns[i])
			}
		}
		clear(dict)
		stmt.WriteByte(';')
	} else {
		stmt.WriteString(" DO NOTHING;")
	}
	{{ end -}}
	return sqlConn.ExecContext(ctx, stmt.String(), args...)
}

// FindByPK is to find single record using primary key.
func FindByPK[T sequel.KeyValuer[T], Ptr sequel.KeyValueScanner[T]](ctx context.Context, sqlConn sequel.DB, model Ptr) error {
	switch vi := any(model).(type) {
	case sequel.KeyFinder:
		_, _, pk := vi.PK()
		return sqlConn.QueryRowContext(ctx, vi.FindByPKStmt(), pk).Scan(model.Addrs()...)

	default:
		var (
			pkName, _, pk = model.PK()
			columns       = model.Columns()
		)
		return sqlConn.QueryRowContext(ctx, "SELECT "+strings.Join(columns, ",")+" FROM "+dbName(model)+model.TableName()+" WHERE "+pkName+" = {{ quoteVar 1 }} LIMIT 1;", pk).Scan(model.Addrs()...)
	}
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
		return sqlConn.ExecContext(ctx, "UPDATE "+dbName(model)+model.TableName()+" SET "+strings.Join(columns, " = {{ quoteVar 1 }},")+" = {{ quoteVar 1 }} WHERE "+pkName+" = {{ quoteVar 1 }};", append(values, pk)...)
		{{ else -}}
		stmt := strpool.AcquireString()
		defer strpool.ReleaseString(stmt)
		stmt.WriteString("UPDATE "+dbName(model)+model.TableName()+" SET ")
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
	return sqlConn.ExecContext(ctx, "DELETE FROM "+ dbName(model) + model.TableName() +" WHERE "+ pkName +" = {{ quoteVar 1 }};", pk)
}

type SelectStmt struct {
	Select    []string
	FromTable string
	Where     sequel.WhereClause
	OrderBy   []sequel.OrderByClause
	GroupBy	  []string
	Offset	  uint64
	Limit     uint16
}

type SQLStatement struct {
	Query	  string
	Arguments []any
}

func QueryStmt[T any, Ptr interface {
	*T
	sequel.Scanner[T]
}, Stmt interface{ 
	SelectStmt | SQLStatement
}](ctx context.Context, sqlConn sequel.DB, stmt Stmt) ([]T, error) {
	var (
		rows *sql.Rows
		err  error
	)

	switch vi := any(stmt).(type) {
	case SelectStmt:
		var (
			blr = AcquireStmt()
			v T
		)
		defer ReleaseStmt(blr)
		blr.WriteString("SELECT ")
		if len(vi.Select) > 0 {
			blr.WriteString(strings.Join(vi.Select, ","))
		} else {
			switch vj := any(v).(type) {
			case sequel.Columner:
				blr.WriteString(strings.Join(vj.Columns(), ","))
			default:
				blr.WriteByte('*')
			}
		}
		if vi.FromTable != "" {
			blr.WriteString(" FROM " + dbName(v) + vi.FromTable)
		} else {
			switch vj := any(v).(type) {
			case sequel.Tabler:
				blr.WriteString(" FROM " + dbName(v) + vj.TableName())
			default:
				return nil, fmt.Errorf("missing table name for model %T", v)
			}
		}
		if vi.Where != nil {
			blr.WriteString(" WHERE ")
			vi.Where(blr)
		}
		if len(vi.GroupBy) > 0 {
			blr.WriteString(" GROUP BY ")
			for i := range vi.GroupBy {
				if i > 0 {
					blr.WriteByte(',')
				}
				blr.WriteString(vi.GroupBy[i])
			}
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
		rows, err = sqlConn.QueryContext(ctx, blr.String(), blr.Args()...)
		ReleaseStmt(blr)

	case SQLStatement:
		rows, err = sqlConn.QueryContext(ctx, vi.Query, vi.Arguments...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
	Table	string
	Set		[]sequel.SetClause
	Where	sequel.WhereClause
	OrderBy []sequel.OrderByClause
	Limit	uint16
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

	var v T
	switch vi := any(stmt).(type) {
	case UpdateStmt:
		if vt, ok := any(v).(sequel.Tabler); ok {
			blr.WriteString("UPDATE " + dbName(v) + vt.TableName())
		} else {
			blr.WriteString("UPDATE " + dbName(v) + vi.Table)
		}
		if vi.Where != nil {
			blr.WriteString(" WHERE ")
			vi.Where(blr)
		}
		if len(vi.Set) > 0 {
			blr.WriteString(" SET ")
			for i := range vi.Set {
				if i > 0 {
					blr.WriteByte(',')
				}
				vi.Set[i](blr)
			}
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

	case DeleteStmt:
		if vt, ok := any(v).(sequel.Tabler); ok {
			blr.WriteString("DELETE FROM " + dbName(v) + vt.TableName())
		} else {
			blr.WriteString("DELETE FROM " + dbName(v) + vi.FromTable)
		}
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
	}
	blr.WriteByte(';')
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
	if stmt != nil {
		stmt.Reset()
		pool.Put(stmt)
	}
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
	s.WriteString(query+wrapVar(s.pos))
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

func dbName(model any) string {
	if v, ok := model.(sequel.DatabaseNamer); ok {
		return v.DatabaseName() + "."
	}
	return ""
}

{{ if not isStaticVar -}}
func wrapVar(i int) string {
	return {{ quote varRune }}+ strconv.Itoa(i)
}
{{ end }}
