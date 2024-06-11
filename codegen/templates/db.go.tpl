{{- reserveImport "context" }}
{{- reserveImport "database/sql" }}
{{- reserveImport "strings" }}
{{- reserveImport "strconv" }}
{{- reserveImport "sync" }}
{{- reserveImport "github.com/si3nloong/sqlgen/sequel" }}
{{- reserveImport "github.com/si3nloong/sqlgen/sequel/strpool" }}

type autoIncrKeyInserter interface {
	sequel.AutoIncrKeyer
	sequel.SingleInserter
}

type primaryKeyInserter interface {
	sequel.PrimaryKeyer
	sequel.SingleInserter
}
{{ if eq driver "postgres" -}}
{{- /* postgres */ -}}
func InsertOne[T sequel.TableColumnValuer[T], Ptr interface {
	sequel.TableColumnValuer[T]
	sequel.Scanner[T]
}](ctx context.Context, sqlConn sequel.DB, model Ptr) error {
	switch v := any(model).(type) {
	case autoIncrKeyInserter:
		_, idx, _ := v.PK()
		values := model.Values()
		values = append(values[:idx], values[idx+1:]...)
		return sqlConn.QueryRowContext(ctx, v.InsertOneStmt(), values...).Scan(model.Addrs()...)
	case primaryKeyInserter:
		return sqlConn.QueryRowContext(ctx, v.InsertOneStmt(), model.Values()...).Scan(model.Addrs()...)
	default:
		columns, values := model.Columns(), model.Values()
		stmt := strpool.AcquireString()
		defer strpool.ReleaseString(stmt)
		stmt.WriteString("INSERT INTO " + dbName(model) + model.TableName() + " (" + strings.Join(columns, ",") + ") VALUES (")
		for i := range values {
			if i > 0 {
				stmt.WriteString(","+wrapVar(i + 1))
			} else {
				// argument always started from 1
				stmt.WriteString(wrapVar(i + 1))
			}
		}
		stmt.WriteString(") RETURNING "+ strings.Join(columns, ",") +";")
		return sqlConn.QueryRowContext(ctx, stmt.String(), values...).Scan(model.Addrs()...)
	}
}	
{{ else }}
func InsertOne[T sequel.TableColumnValuer[T], Ptr interface {
	sequel.TableColumnValuer[T]
	sequel.Scanner[T]
}](ctx context.Context, sqlConn sequel.DB, model Ptr) (sql.Result, error) {
	switch v := any(model).(type) {
	case autoIncrKeyInserter:
		_, idx, _ := v.PK()
		values := model.Values()
		values = append(values[:idx], values[idx+1:]...)
		result, err := sqlConn.ExecContext(ctx, v.InsertOneStmt(), values...)
		if err != nil {
			return nil, err
		}
		i64, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}
		switch v := model.Addrs()[idx].(type) {
		case *int64:
			*v = i64
		case sql.Scanner:
			if err := v.Scan(i64); err != nil {
				return nil, err
			}
		default:
			return nil, errors.New(`sqlgen: invalid auto increment data type`)
		}
		return result, nil
	case primaryKeyInserter:
		return sqlConn.ExecContext(ctx, v.InsertOneStmt(), model.Values()...)
	default:
		columns, values := model.Columns(), model.Values()
		return sqlConn.ExecContext(ctx, "INSERT INTO "+dbName(model)+model.TableName()+" ("+strings.Join(columns, ",")+") VALUES ("+strings.Repeat(",?", len(columns))[1:]+");", values...)
	}
}
{{ end }}

// Insert is a helper function to insert multiple records.
func Insert[T sequel.TableColumnValuer[T]](ctx context.Context, sqlConn sequel.DB, data []T) (sql.Result, error) {
	noOfData := len(data)
	if noOfData == 0 {
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

{{ if eq driver "postgres" -}}
{{- /* postgres */ -}}
type UpsertOptions struct {
	override		bool
	omitFields		[]string
	onDuplicateKeys []string
}
type UpsertOption func(*UpsertOptions)

func WithOverride(override bool) UpsertOption {
	return func(opt *UpsertOptions) {
		opt.override = override
	}
}

func WithOmitFields(fields []string) UpsertOption {
	return func(opt *UpsertOptions) {
		opt.omitFields = fields
	}
}

func WithDuplicateKeys(keys []string) UpsertOption {
	return func(opt *UpsertOptions) {
		opt.onDuplicateKeys = keys
	}
}

func UpsertOne[T sequel.KeyValuer[T], Ptr sequel.KeyValueScanner[T]](ctx context.Context, sqlConn sequel.DB, model Ptr, opts ...UpsertOption) error {
	var opt UpsertOptions
	for i := range opts {
		opts[i](&opt)
	}
	switch any(model).(type) {
	case sequel.Keyer:
		columns := model.Columns()
		stmt := strpool.AcquireString()
		noOfCols := len(columns)
		defer strpool.ReleaseString(stmt)
		stmt.WriteString("INSERT INTO " + dbName(model) + model.TableName() + " (" + strings.Join(columns, ",") + ") VALUES (")
		for i := 0; i < noOfCols; i++ {
			if i > 0 {
				stmt.WriteString("," + wrapVar(i))
			} else {
				stmt.WriteString(wrapVar(i))
			}
		}
		if len(opt.onDuplicateKeys) > 0 {
			stmt.WriteString(") ON CONFLICT(" + strings.Join(opt.onDuplicateKeys, ",") + ")")
		} else {
			switch vi := any(model).(type) {
			case sequel.PrimaryKeyer:
				pkName, _, _ := vi.PK()
				opt.omitFields = append(opt.omitFields, pkName)
				stmt.WriteString(") ON CONFLICT(" + pkName + ")")
			case sequel.CompositeKeyer:
				names, _, _ := vi.CompositeKey()
				opt.omitFields = append(opt.omitFields, names...)
				stmt.WriteString(") ON CONFLICT(" + strings.Join(names, ",") + ")")
			default:
				panic("unreachable")
			}
		}
		if opt.override {
			omitDict := make(map[string]struct{})
			for i := range opt.omitFields {
				omitDict[opt.omitFields[i]] = struct{}{}
			}
			for i := 0; i < noOfCols; i++ {
				if _, ok := omitDict[columns[i]]; ok {
					continue
				}
				if i < noOfCols-1 {
					stmt.WriteString(columns[i] + " = EXCLUDED." + columns[i] + ",")
				} else {
					stmt.WriteString(columns[i] + " = EXCLUDED." + columns[i])
				}
			}
			clear(omitDict)
			stmt.WriteString(" RETURNING " + strings.Join(columns, ",") + ";")
		} else {
			stmt.WriteString(" DO NOTHING;")
		}
		return sqlConn.QueryRowContext(ctx, stmt.String(), model.Values()...).Scan(model.Addrs()...)
	}
	return nil
}
{{ else }}
func UpsertOne[T sequel.KeyValuer[T], Ptr sequel.KeyValueScanner[T]](ctx context.Context, sqlConn sequel.DB, model Ptr, override bool, omittedFields ...string) (sql.Result, error) {
	switch v := any(model).(type) {
	case sequel.PrimaryKeyer:
		pkName, idx, _ := v.PK()
		columns := model.Columns()
		stmt := strpool.AcquireString()
		defer strpool.ReleaseString(stmt)
		if !override {
			stmt.WriteString("INSERT IGNORE INTO " + dbName(model) + model.TableName() + " (" + strings.Join(columns, ",") + ") VALUES (" + strings.Repeat(",?", len(columns))[1:] + ");")
		} else {
			omitDict := map[string]struct{}{pkName: {}}
			for i := range omittedFields {
				omitDict[omittedFields[i]] = struct{}{}
			}
			noOfCols := len(columns)
			stmt.WriteString("INSERT INTO " + dbName(model) + model.TableName() + " (" + strings.Join(columns, ",") + ") VALUES (" + strings.Repeat(",?", noOfCols)[1:] + ") ON DUPLICATE KEY UPDATE ")
			for i := range columns {
				if _, ok := omitDict[columns[i]]; ok {
					continue
				}
				if i < noOfCols-1 {
					stmt.WriteString(columns[i] + " =VALUES(" + columns[i] + "),")
				} else {
					stmt.WriteString(columns[i] + " =VALUES(" + columns[i] + ")")
				}
			}
			clear(omitDict)
		}
		if _, ok := any(model).(sequel.AutoIncrKeyer); ok {
			result, err := sqlConn.ExecContext(ctx, stmt.String(), model.Values()...)
			if err != nil {
				return nil, err
			}
			i64, err := result.LastInsertId()
			if err != nil {
				return nil, err
			}
			switch v := model.Addrs()[idx].(type) {
			case *int64:
				*v = i64
			case sql.Scanner:
				if err := v.Scan(i64); err != nil {
					return nil, err
				}
			default:
				return nil, errors.New(`sqlgen: invalid auto increment data type`)
			}
			return result, nil
		}
		return sqlConn.ExecContext(ctx, stmt.String(), model.Values()...)
	case sequel.CompositeKeyer:
		names, idxs, _ := v.CompositeKey()
		columns := model.Columns()
		stmt := strpool.AcquireString()
		defer strpool.ReleaseString(stmt)
		if !override {
			stmt.WriteString("INSERT IGNORE INTO " + dbName(model) + model.TableName() + " (" + strings.Join(columns, ",") + ") VALUES (" + strings.Repeat(",?", len(columns))[1:] + ");")
		} else {
			dict := make(map[string]struct{})
			for i := range append(names, omittedFields...) {
				dict[omittedFields[i]] = struct{}{}
			}
			noOfCols := len(columns)
			stmt.WriteString("INSERT INTO " + dbName(model) + model.TableName() + " (" + strings.Join(columns, ",") + ") VALUES (" + strings.Repeat(",?", noOfCols)[1:] + ") ON DUPLICATE KEY UPDATE ")
			// Exclude composite key, don't update it
			for i := len(idxs) - 1; i >= 0; i-- {
				columns = append(columns[:idxs[i]], columns[idxs[i]+1:]...)
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
		return sqlConn.ExecContext(ctx, stmt.String(), model.Values()...)
	default:
		panic("unreachable")
	}
}
{{ end }}

// Upsert is a helper function to upsert multiple records.
func Upsert[T sequel.KeyValuer[T], Ptr sequel.Scanner[T]](ctx context.Context, sqlConn sequel.DB, data []T, override bool, omittedFields ...string) (sql.Result, error) {
	noOfData := len(data)
	if noOfData == 0 {
		return new(sequel.EmptyResult), nil
	}

	var (
		model    = data[0]
		columns  = model.Columns()
		noOfCols = len(columns)
		args     = make([]any, 0, noOfCols*noOfData)
	)

	var stmt = strpool.AcquireString()
	defer strpool.ReleaseString(stmt)
	switch v := any(model).(type) {
	case sequel.AutoIncrKeyer:
		pkName, idx, _ := v.PK()
		omittedFields = append(omittedFields, pkName)
		// Don't include auto increment primary key on INSERT
		columns = append(columns[:idx], columns[idx+1:]...)
		if override {
			stmt.WriteString("INSERT INTO " + dbName(model) + model.TableName() + " (" + strings.Join(columns, ",") + ") VALUES ")
		} else {
			stmt.WriteString("INSERT IGNORE INTO " + dbName(model) + model.TableName() + " (" + strings.Join(columns, ",") + ") VALUES ")
		}
		placeholder := ",(" + strings.Repeat(",?", noOfCols)[1:] + ")"
		for i := 0; i < noOfData; i++ {
			values := data[i].Values()
			values = append(values[:idx], values[idx+1:]...)
			args = append(args, values...)
			if i > 0 {
				stmt.WriteString(placeholder)
			} else {
				stmt.WriteString(placeholder[1:])
			}
		}
	case sequel.PrimaryKeyer:
		pkName, _, _ := v.PK()
		omittedFields = append(omittedFields, pkName)
		if override {
			stmt.WriteString("INSERT INTO " + dbName(model) + model.TableName() + " (" + strings.Join(columns, ",") + ") VALUES ")
		} else {
			stmt.WriteString("INSERT IGNORE INTO " + dbName(model) + model.TableName() + " (" + strings.Join(columns, ",") + ") VALUES ")
		}
		placeholder := ",(" + strings.Repeat(",?", noOfCols)[1:] + ")"
		for i := 0; i < noOfData; i++ {
			args = append(args, data[i].Values()...)
			if i > 0 {
				stmt.WriteString(placeholder)
			} else {
				stmt.WriteString(placeholder[1:])
			}
		}
	case sequel.CompositeKeyer:
		if override {
			stmt.WriteString("INSERT INTO " + dbName(model) + model.TableName() + " (" + strings.Join(columns, ",") + ") VALUES ")
		} else {
			stmt.WriteString("INSERT IGNORE INTO " + dbName(model) + model.TableName() + " (" + strings.Join(columns, ",") + ") VALUES ")
		}
		placeholder := ",(" + strings.Repeat(",?", noOfCols)[1:] + ")"
		for i := 0; i < noOfData; i++ {
			args = append(args, data[i].Values()...)
			if i > 0 {
				stmt.WriteString(placeholder)
			} else {
				stmt.WriteString(placeholder[1:])
			}
		}
		_, idxs, _ := v.CompositeKey()
		// Exclude primary key, don't update it
		for i := len(idxs) - 1; i >= 0; i-- {
			columns = append(columns[:idxs[i]], columns[idxs[i]+1:]...)
		}
	}
	if override {
		stmt.WriteString(" ON DUPLICATE KEY UPDATE ")
		/* don't update primary key when we do upsert */
		omitDict := map[string]struct{}{}
		for i := range omittedFields {
			omitDict[omittedFields[i]] = struct{}{}
		}
		for i := range columns {
			if _, ok := omitDict[columns[i]]; ok {
				continue
			}
			if i < noOfCols-1 {
				stmt.WriteString(columns[i] + " =VALUES(" + columns[i] + "),")
			} else {
				stmt.WriteString(columns[i] + " =VALUES(" + columns[i] + ")")
			}
		}
		clear(omitDict)
	}
	stmt.WriteByte(';')
	return sqlConn.ExecContext(ctx, stmt.String(), args...)
}

type primaryKeyFinder interface {
	sequel.PrimaryKeyer
	sequel.KeyFinder
}

// FindByPK is to find single record using primary key.
func FindByPK[T sequel.KeyValuer[T], Ptr sequel.KeyValueScanner[T]](ctx context.Context, sqlConn sequel.DB, model Ptr) error {
	switch v := any(model).(type) {
	case primaryKeyFinder:
		_, _, pk := v.PK()
		return sqlConn.QueryRowContext(ctx, v.FindByPKStmt(), pk).Scan(model.Addrs()...)
	case sequel.PrimaryKeyer:
		columns := model.Columns()
		pkName, _, pk := v.PK()
		return sqlConn.QueryRowContext(ctx, "SELECT "+strings.Join(columns, ",")+" FROM "+dbName(model)+model.TableName()+" WHERE "+pkName+" = {{ quoteVar 1 }} LIMIT 1;", pk).Scan(model.Addrs()...)
	case sequel.CompositeKeyer:
		columns := model.Columns()
		names, _, keys := v.CompositeKey()
		{{ if isStaticVar -}}
		return sqlConn.QueryRowContext(ctx, "SELECT "+strings.Join(columns, ",")+" FROM "+dbName(model)+model.TableName()+" WHERE "+strings.Join(names, " = {{ quoteVar 1 }} AND ")+" = {{ quoteVar 1 }} LIMIT 1;", keys...).Scan(model.Addrs()...)
		{{ else -}}
		stmt := strpool.AcquireString()
		defer strpool.ReleaseString(stmt)
		stmt.WriteString("SELECT "+strings.Join(columns, ",")+" FROM "+dbName(model)+model.TableName()+" WHERE ")
		max := len(names)
		for i := 1; i <= max; i++ {
			if i > 1 {
				stmt.WriteString(" AND "+ names[i]+" = "+ wrapVar(i))
			} else {
				stmt.WriteString(names[i]+" = "+ wrapVar(i))
			}
		}
		stmt.WriteString(" LIMIT 1;")
		return sqlConn.QueryRowContext(ctx, stmt.String(), keys...).Scan(model.Addrs()...)
		{{ end -}}
	default:
		panic("unreachable")
	}
}

type primaryKeyUpdater interface {
	sequel.PrimaryKeyer
	sequel.KeyUpdater
}

// UpdateByPK is to update single record using primary key.
func UpdateByPK[T sequel.KeyValuer[T]](ctx context.Context, sqlConn sequel.DB, model T) (sql.Result, error) {
	switch v := any(model).(type) {
	case primaryKeyUpdater:
		_, pkIdx, pk := v.PK()
		values := model.Values()
		values = append(values[:pkIdx], append(values[pkIdx+1:], pk)...)
		return sqlConn.ExecContext(ctx, v.UpdateByPKStmt(), values...)
	case sequel.PrimaryKeyer:
		pkName, pkIdx, pk := v.PK()
		values := model.Values()
		columns := model.Columns()
		values = append(values[:pkIdx], append(values[pkIdx+1:], pk)...)
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
	default:
		panic("unreachable")
	}
}

type primaryKeyDeleter interface {
	sequel.PrimaryKeyer
	sequel.KeyDeleter
}

// DeleteByPK is to update single record using primary key.
func DeleteByPK[T sequel.KeyValuer[T]](ctx context.Context, sqlConn sequel.DB, model T) (sql.Result, error) {
	switch v := any(model).(type) {
	case primaryKeyDeleter:
		_, _, pk := v.PK()
		return sqlConn.ExecContext(ctx, v.DeleteByPKStmt(), pk)
	case sequel.PrimaryKeyer:
		pkName, _, pk := v.PK()
		return sqlConn.ExecContext(ctx, "DELETE FROM "+dbName(model)+model.TableName()+" WHERE "+pkName+" = {{ quoteVar 1 }};", pk)
	case sequel.CompositeKeyer:
		names, _, keys := v.CompositeKey()
		{{ if isStaticVar -}}
		return sqlConn.ExecContext(ctx, "DELETE FROM "+dbName(model)+model.TableName()+" WHERE "+strings.Join(names, " = ? AND ")+" = ?;", keys...)
		{{ else -}}
		stmt := strpool.AcquireString()
		defer strpool.ReleaseString(stmt)
		stmt.WriteString("DELETE FROM "+dbName(model)+model.TableName()+" WHERE ")
		max := len(names)
		for i := 1; i <= max; i++ {
			if i == 1 {
				stmt.WriteString(names[i]+" = "+ wrapVar(i))
			} else {
				stmt.WriteString(" AND "+ names[i]+" = "+ wrapVar(i))
			}
		}
		stmt.WriteByte(';')
		return sqlConn.ExecContext(ctx, stmt.String(), keys...)
		{{ end -}}
	default:
		panic("unreachable")
	}
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
