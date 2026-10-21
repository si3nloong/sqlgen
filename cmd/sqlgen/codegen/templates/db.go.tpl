{{- reserveImport "context" }}
{{- reserveImport "database/sql" }}
{{- reserveImport "database/sql/driver" }}
{{- reserveImport "strings" }}
{{- reserveImport "io" }}
{{- reserveImport "strconv" }}
{{- reserveImport "fmt" }}
{{- reserveImport "sync" }}
{{- reserveImport "unsafe" }}
{{- reserveImport "github.com/si3nloong/sqlgen/sequel" }}
{{- reserveImport "github.com/si3nloong/sqlgen/sequel/strpool" }}
{{- reserveImport "github.com/si3nloong/sqlgen/sequel/driver/postgres/pgutil" }}

{{ if eq driver "postgres" -}}
{{- /* postgres */ -}}
func InsertOne[T sequel.ColumnValuer, Ptr interface {
	sequel.ColumnValuer
	sequel.PtrScanner[T]
}](ctx context.Context, db sequel.DB, model Ptr) error {
	switch v := any(model).(type) {
	case sequel.SingleInserter:
		query, args := v.InsertOneStmt()
		return db.QueryRowContext(ctx, query, args...).Scan(model.Addrs()...)
	case sequel.Inserter:
		stmt := strpool.AcquireString()
		stmt.WriteString("INSERT INTO ") 
		stmt.WriteString(DbTable(model))
		stmt.WriteString(v.SQLInsertColumns())
		stmt.WriteString(" VALUES ")
		stmt.WriteString(v.SQLInsertPlaceholders(0))
		stmt.WriteString(" RETURNING ")
		writej(stmt, TableColumns(model))
		stmt.WriteString(";")
		row := db.QueryRowContext(ctx, stmt.String(), v.SQLInsertValues()...)
		strpool.ReleaseString(stmt) // Cleanup
		return row.Scan(model.Addrs()...)
	default:
		columns, args := TableColumns(model), model.Values()
		columnSQL := strings.Join(columns, ",")
		stmt := strpool.AcquireString()
		stmt.WriteString("INSERT INTO ") 
		stmt.WriteString(DbTable(model))
		stmt.WriteString(" (")
		stmt.WriteString(columnSQL)
		stmt.WriteString(") VALUES ({{ quoteVar 1 }}")
		noOfArgs := len(args)
		for i := 1; i < noOfArgs; i++ {
			stmt.WriteString(",")
			stmt.WriteString(wrapVar(i+1))
		}
		stmt.WriteString(") RETURNING ")
		stmt.WriteString(columnSQL)
		stmt.WriteString(";")
		row := db.QueryRowContext(ctx, stmt.String(), args...)
		strpool.ReleaseString(stmt) // Cleanup
		args = nil // Cleanup
		return row.Scan(model.Addrs()...)
	}
}	
{{ else }}
type autoIncrKeyInserter interface {
	sequel.AutoIncrKeyer
	sequel.SingleInserter
	ScanAutoIncr(int64) error
}

func InsertOne[T sequel.ColumnValuer, Ptr interface {
	sequel.ColumnValuer
	sequel.PtrScanner[T]
}](ctx context.Context, db sequel.DB, model Ptr) (sql.Result, error) {
	switch v := any(model).(type) {
	case autoIncrKeyInserter:
		query, args := v.InsertOneStmt()
		result, err := db.ExecContext(ctx, query, args...)
		if err != nil {
			return nil, err
		}
		i64, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}
		return result, v.ScanAutoIncr(i64)
	case sequel.SingleInserter:
		query, args := v.InsertOneStmt()
		return db.ExecContext(ctx, query, args...)
	case sequel.Inserter:
		stmt := strpool.AcquireString()
		stmt.WriteString("INSERT INTO ") 
		stmt.WriteString(DbTable(model))
		stmt.WriteString(v.SQLInsertColumns())
		stmt.WriteString(" VALUES ")
		stmt.WriteString(v.SQLInsertPlaceholders(0))
		stmt.WriteString(";")
		result, err := db.ExecContext(ctx, stmt.String(), v.SQLInsertValues()...)
		strpool.ReleaseString(stmt) // Cleanup
		return result, err
	default:
		columns, args := TableColumns(model), model.Values()
		stmt := strpool.AcquireString()
		stmt.WriteString("INSERT INTO ")
		stmt.WriteString(DbTable(model))
		stmt.WriteString(" (")
		writej(stmt, columns)
		stmt.WriteString(") VALUES (")
		stmt.WriteString(strings.Repeat("{{ quoteVar 1 }},", len(columns) - 1))
		stmt.WriteString("{{ quoteVar 1 }});")
		result, err := db.ExecContext(ctx, stmt.String(), args...)
		strpool.ReleaseString(stmt) // Cleanup
		args = nil // Cleanup
		return result, err
	}
}
{{ end }}

{{ if eq driver "postgres" -}}
{{- /* postgres */ -}}
// Insert is a helper function to insert multiple records.
func Insert[T sequel.ColumnValuer, Ptr sequel.PtrScanner[T]](ctx context.Context, db sequel.DB, data []T) (sql.Result, error) {
	noOfData := len(data)
	if noOfData == 0 {
		return new(sequel.EmptyResult), nil
	}

	model := data[0]
	columns := TableColumns(model)
	stmt := strpool.AcquireString()

	var (
		rows *sql.Rows
		err  error
	)
	switch v := any(model).(type) {
	case sequel.Inserter:
		noOfCols := len(columns)
		args := make([]any, 0, noOfCols*noOfData)
		stmt.WriteString("INSERT INTO ")
		stmt.WriteString(DbTable(model))
		stmt.WriteString(" (")
		stmt.WriteString(v.SQLInsertColumns())
		stmt.WriteString(") VALUES ")
		stmt.WriteString(v.SQLInsertPlaceholders(0))
		args = append(args, data[0].Values()...)
		for i := 1; i < noOfData; i++ {
			stmt.WriteString(",")
			stmt.WriteString(v.SQLInsertPlaceholders(i))
			args = append(args, data[i].Values()...)
		}
		stmt.WriteString(" RETURNING ")
		writej(stmt, columns)
		stmt.WriteString(";")
		rows, err = db.QueryContext(ctx, stmt.String(), args...)
		args = nil // Cleanup
	default:
		noOfCols := len(columns)
		args := make([]any, 0, noOfCols*noOfData)
		stmt.WriteString("INSERT INTO ")
		stmt.WriteString(DbTable(model))
		stmt.WriteString(" (")
		writej(stmt, columns)
		stmt.WriteString(") VALUES ({{ quoteVar 1 }}")
		for i := 1; i < noOfCols; i++ {
			stmt.WriteString(",")
			stmt.WriteString(wrapVar(i+1))
		}
		stmt.WriteString(")")
		args = append(args, data[0].Values()...)
		for i := 1; i < noOfData; i++ {
			stmt.WriteString(",(")
			stmt.WriteString(wrapVar(i*noOfCols + 1))
			for j := 1; j < noOfCols; j++ {
				stmt.WriteString(",")
				stmt.WriteString(wrapVar(i*noOfCols + j + 1))
			}
			stmt.WriteString(")")
			args = append(args, data[i].Values()...)
		}
		stmt.WriteString(" RETURNING ")
		writej(stmt, columns)
		stmt.WriteString(";")
		rows, err = db.QueryContext(ctx, stmt.String(), args...)
		args = nil // Cleanup
	}
	strpool.ReleaseString(stmt) // Cleanup
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var count int64
	for rows.Next() {
		if err := rows.Scan(Ptr(&data[count]).Addrs()...); err != nil {
			return nil, err
		}
		count++
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sequel.NewRowsAffectedResult(count), nil
}
{{ else }}
// Insert is a helper function to insert multiple records.
// Models must have at least one insertable column after excluding an auto-increment key.
func Insert[T sequel.ColumnValuer](ctx context.Context, db sequel.DB, data []T) (sql.Result, error) {
	noOfData := len(data)
	if noOfData == 0 {
		return new(sequel.EmptyResult), nil
	}

	model := data[0]
	columns := model.Columns()
	stmt := strpool.AcquireString()

	switch v := any(model).(type) {
	case sequel.AutoIncrKeyer:
		_, idx, _ := v.PK()
		columns = append(columns[:idx], columns[idx+1:]...)
		noOfColumns := len(columns)
		args := make([]any, 0, noOfColumns*noOfData)
		placeholder := strings.Repeat(",{{ quoteVar 1 }}", noOfColumns)
		placeholder = "(" + placeholder[1:] + ")"
		stmt.WriteString("INSERT INTO ")
		stmt.WriteString(DbTable(model))
		stmt.WriteString(" (")
		writej(stmt, columns)
		stmt.WriteString(") VALUES ")
		stmt.WriteString(placeholder)
		values := data[0].Values()
		values = append(values[:idx], values[idx+1:]...)
		args = append(args, values...)
		for i := 1; i < noOfData; i++ {
			stmt.WriteString(",")
			stmt.WriteString(placeholder)
			values = data[i].Values()
			values = append(values[:idx], values[idx+1:]...)
			args = append(args, values...)
		}
		stmt.WriteString(";")
		result, err := db.ExecContext(ctx, stmt.String(), args...)
		strpool.ReleaseString(stmt) // Cleanup
		args = nil // Cleanup
		return result, err
	default:
		noOfColumns := len(columns)
		args := make([]any, 0, noOfColumns*noOfData)
		placeholder := "(" + strings.Repeat(",{{ quoteVar 1 }}", noOfColumns)[1:] + ")"
		stmt.WriteString("INSERT INTO ")
		stmt.WriteString(DbTable(model))
		stmt.WriteString(" (")
		writej(stmt, columns)
		stmt.WriteString(") VALUES ")
		stmt.WriteString(placeholder)
		args = append(args, data[0].Values()...)
		for i := 1; i < noOfData; i++ {
			stmt.WriteString(",")
			stmt.WriteString(placeholder)
			args = append(args, data[i].Values()...)
		}
		stmt.WriteString(";")
		result, err := db.ExecContext(ctx, stmt.String(), args...)
		strpool.ReleaseString(stmt) // Cleanup
		args = nil // Cleanup
		return result, err
	}
}
{{ end }}

{{ if eq driver "postgres" -}}
{{- /* postgres */ -}}
type duplicateKeyer interface {
	DuplicateKeys() []string
}
type setOnUpsertOmitter interface {
	OnUpsertOmitColumns() []string
}
type upsertOpts struct {
	doNothing     bool
	omitFields    []string
	duplicateKeys []string
}
type UpsertOption func(*upsertOpts)

func WithDoNothing(doNothing bool) UpsertOption {
	return func(opt *upsertOpts) {
		opt.doNothing = doNothing
	}
}

func WithOmitFields(fields ...string) UpsertOption {
	return func(opt *upsertOpts) {
		opt.omitFields = fields
	}
}

func WithDuplicateKeys(keys ...string) UpsertOption {
	return func(opt *upsertOpts) {
		opt.duplicateKeys = keys
	}
}

func UpsertOne[T sequel.KeyValuer, Ptr interface{
	sequel.KeyValueScanner[T]
	sequel.Inserter
}](ctx context.Context, db sequel.DB, model Ptr, opts ...UpsertOption) error {
	var opt upsertOpts
	for i := range opts {
		opts[i](&opt)
	}

	var (
		stmt     = strpool.AcquireString()
		columns  = model.Columns()
		noOfCols = len(columns)
		values   = model.Values()
	)
	defer strpool.ReleaseString(stmt)

	switch v := any(model).(type) {
	case sequel.AutoIncrKeyer:
		pkName, idx, _ := v.PK()
		opt.omitFields = append(opt.omitFields, pkName)
		columns = append(columns[:idx], columns[idx+1:]...)
		values = append(values[:idx], values[idx+1:]...)
		// Don't include auto increment primary key on INSERT
		stmt.WriteString("INSERT INTO ")
		stmt.WriteString(DbTable(model))
		stmt.WriteString(" (")
		writej(stmt, columns)
		stmt.WriteString(") VALUES ")
	case sequel.PrimaryKeyer:
		pkName, _, _ := v.PK()
		opt.omitFields = append(opt.omitFields, pkName)
		stmt.WriteString("INSERT INTO ")
		stmt.WriteString(DbTable(model))
		stmt.WriteString(" (")
		writej(stmt, columns)
		stmt.WriteString(") VALUES ")
	case sequel.CompositeKeyer:
		keyNames, _, _ := v.CompositeKey()
		opt.omitFields = append(opt.omitFields, keyNames...)
		stmt.WriteString("INSERT INTO ")
		stmt.WriteString(DbTable(model))
		stmt.WriteString(" (")
		writej(stmt, columns)
		stmt.WriteString(") VALUES ")
	default:
		panic("unreachable")
	}
	stmt.WriteString(model.SQLInsertPlaceholders(0))
	stmt.WriteString(" ON CONFLICT(")
	if len(opt.duplicateKeys) > 0 {
		writej(stmt, opt.duplicateKeys)
	} else {
		switch vi := any(model).(type) {
		case duplicateKeyer:
			keys := vi.DuplicateKeys()
			opt.omitFields = append(opt.omitFields, keys...)
			writej(stmt, keys)
		case sequel.PrimaryKeyer:
			pkName, _, _ := vi.PK()
			stmt.WriteString(pkName)
		case sequel.CompositeKeyer:
			names, _, _ := vi.CompositeKey()
			writej(stmt, names)
		default:
			panic("unreachable")
		}
	}
	if opt.doNothing {
		stmt.WriteString(") DO NOTHING")
	} else {
		stmt.WriteString(") DO UPDATE SET ")
		omitDict := make(map[string]struct{})
		noOfCols = len(columns)
		if omitter, ok := any(model).(setOnUpsertOmitter); ok {
			opt.omitFields = append(opt.omitFields, omitter.OnUpsertOmitColumns()...)
		}
		for i := range opt.omitFields {
			omitDict[opt.omitFields[i]] = struct{}{}
		}
		i := 0
		// Write the first included column without a leading comma.
		for ; i < noOfCols; i++ {
			if _, ok := omitDict[columns[i]]; ok {
				continue
			}
			stmt.WriteString(columns[i])
			stmt.WriteString(" = EXCLUDED.")
			stmt.WriteString(columns[i])
			i++
			break
		}
		for ; i < noOfCols; i++ {
			if _, ok := omitDict[columns[i]]; ok {
				continue
			}
			stmt.WriteString(",")
			stmt.WriteString(columns[i])
			stmt.WriteString(" = EXCLUDED.")
			stmt.WriteString(columns[i])
		}
		clear(omitDict)
	}
	stmt.WriteString(" RETURNING ")
	writej(stmt, TableColumns(model))
	stmt.WriteString(";")
	return db.QueryRowContext(ctx, stmt.String(), values...).Scan(model.Addrs()...)
}

// Upsert is a helper function to upsert multiple records.
func Upsert[T interface {
	sequel.Keyer
	sequel.Inserter
}, Ptr sequel.PtrScanner[T]](ctx context.Context, db sequel.DB, data []T, opts ...UpsertOption) (sql.Result, error) {
	noOfData := len(data)
	if noOfData == 0 {
		return new(sequel.EmptyResult), nil
	}

	var (
		model    = data[0]
		columns  = model.Columns()
		noOfCols = len(columns)
		args     = make([]any, 0, noOfCols*noOfData)
		opt      upsertOpts
	)

	for i := range opts {
		opts[i](&opt)
	}

	var stmt = strpool.AcquireString()
	switch v := any(model).(type) {
	case sequel.AutoIncrKeyer:
		pkName, idx, _ := v.PK()
		opt.omitFields = append(opt.omitFields, pkName)
		columns = append(columns[:idx], columns[idx+1:]...)
		// Don't include auto increment primary key on INSERT
		stmt.WriteString("INSERT INTO ")
		stmt.WriteString(DbTable(model))
		stmt.WriteString(" (")
		stmt.WriteString(model.SQLInsertColumns())
		stmt.WriteString(") VALUES ")
		stmt.WriteString(model.SQLInsertPlaceholders(0))
		values := data[0].Values()
		values = append(values[:idx], values[idx+1:]...)
		args = append(args, values...)
		for i := 1; i < noOfData; i++ {
			stmt.WriteString(",")
			stmt.WriteString(model.SQLInsertPlaceholders(i))
			values := data[i].Values()
			values = append(values[:idx], values[idx+1:]...)
			args = append(args, values...)
		}
	case sequel.PrimaryKeyer:
		pkName, _, _ := v.PK()
		opt.omitFields = append(opt.omitFields, pkName)
		stmt.WriteString("INSERT INTO ")
		stmt.WriteString(DbTable(model))
		stmt.WriteString(" (")
		stmt.WriteString(model.SQLInsertColumns())
		stmt.WriteString(") VALUES ")
		stmt.WriteString(model.SQLInsertPlaceholders(0))
		args = append(args, data[0].Values()...)
		for i := 1; i < noOfData; i++ {
			stmt.WriteString(",")
			stmt.WriteString(model.SQLInsertPlaceholders(i))
			args = append(args, data[i].Values()...)
		}
	case sequel.CompositeKeyer:
		keyNames, _, _ := v.CompositeKey()
		opt.omitFields = append(opt.omitFields, keyNames...)
		stmt.WriteString("INSERT INTO ")
		stmt.WriteString(DbTable(model))
		stmt.WriteString(" (")
		stmt.WriteString(model.SQLInsertColumns())
		stmt.WriteString(") VALUES ")
		stmt.WriteString(model.SQLInsertPlaceholders(0))
		args = append(args, data[0].Values()...)
		for i := 1; i < noOfData; i++ {
			stmt.WriteString(",")
			stmt.WriteString(model.SQLInsertPlaceholders(i))
			args = append(args, data[i].Values()...)
		}
	default:
		panic("unreachable")
	}
	stmt.WriteString(" ON CONFLICT(")
	if len(opt.duplicateKeys) > 0 {
		writej(stmt, opt.duplicateKeys)
	} else {
		switch vi := any(model).(type) {
		case duplicateKeyer:
			keys := vi.DuplicateKeys()
			opt.omitFields = append(opt.omitFields, keys...)
			writej(stmt, keys)
		case sequel.PrimaryKeyer:
			pkName, _, _ := vi.PK()
			stmt.WriteString(pkName)
		case sequel.CompositeKeyer:
			names, _, _ := vi.CompositeKey()
			// // Exclude primary key, don't update it
			// for i := len(idxs) - 1; i >= 0; i-- {
			// 	columns = append(columns[:idxs[i]], columns[idxs[i]+1:]...)
			// }
			writej(stmt, names)
		default:
			panic("unreachable")
		}
	}
	if opt.doNothing {
		stmt.WriteString(") DO NOTHING")
	} else {
		stmt.WriteString(") DO UPDATE SET ")
		omitDict := make(map[string]struct{})
		if omitter, ok := any(model).(setOnUpsertOmitter); ok {
			opt.omitFields = append(opt.omitFields, omitter.OnUpsertOmitColumns()...)
		}
		columns = model.Columns()
		for i := range opt.omitFields {
			omitDict[opt.omitFields[i]] = struct{}{}
		}
		first := true
		for i := 0; i < noOfCols; i++ {
			if _, ok := omitDict[columns[i]]; ok {
				continue
			}
			if !first {
				stmt.WriteString(",")
			}
			stmt.WriteString(columns[i])
			stmt.WriteString(" = EXCLUDED.")
			stmt.WriteString(columns[i])
			first = false
		}
		clear(omitDict)
	}
	stmt.WriteString(" RETURNING ")
	writej(stmt, TableColumns(model))
	stmt.WriteString(";")
	rows, err := db.QueryContext(ctx, stmt.String(), args...)
	strpool.ReleaseString(stmt) // Deallocate statement
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var count int64
	for rows.Next() {
		if err := rows.Scan(Ptr(&data[count]).Addrs()...); err != nil {
			return nil, err
		}
		count++
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sequel.NewRowsAffectedResult(count), nil
}
{{ else }}
func UpsertOne[T sequel.KeyValuer, Ptr sequel.KeyValueScanner[T]](ctx context.Context, db sequel.DB, model Ptr, override bool, omittedFields ...string) (sql.Result, error) {
	stmt := strpool.AcquireString()
	defer strpool.ReleaseString(stmt)
	switch v := any(model).(type) {
	case sequel.AutoIncrKeyer:
		pkName, idx, _ := v.PK()
		columns := model.Columns()
		if !override {
			stmt.WriteString("INSERT IGNORE INTO ")
			stmt.WriteString(DbTable(model))
			stmt.WriteString(" (")
			writej(stmt, columns)
			stmt.WriteString(") VALUES (")
			stmt.WriteString(strings.Repeat(",?", len(columns))[1:] + ");")
		} else {
			omitDict := map[string]struct{}{pkName: {}}
			for i := range omittedFields {
				omitDict[omittedFields[i]] = struct{}{}
			}
			noOfCols := len(columns)
			stmt.WriteString("INSERT INTO ")
			stmt.WriteString(DbTable(model))
			stmt.WriteString(" (")
			writej(stmt, columns)
			stmt.WriteString(") VALUES (")
			stmt.WriteString(strings.Repeat(",?", noOfCols)[1:])
			stmt.WriteString(") ON DUPLICATE KEY UPDATE ")
			first := true
			for i := range columns {
				if _, ok := omitDict[columns[i]]; ok {
					continue
				}
				if first {
					stmt.WriteString(columns[i])
					stmt.WriteString("=VALUES(")
					stmt.WriteString(columns[i])
					stmt.WriteString(")")
				} else {
					stmt.WriteString(",")
					stmt.WriteString(columns[i])
					stmt.WriteString("=VALUES(")
					stmt.WriteString(columns[i])
					stmt.WriteString(")")
				}
				first = false
			}
			clear(omitDict)
		}
		result, err := db.ExecContext(ctx, stmt.String(), model.Values()...)
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
			return nil, fmt.Errorf(`invalid auto increment data type %T`, v)
		}
		return result, nil
	case sequel.PrimaryKeyer:
		pkName, _, _ := v.PK()
		columns := model.Columns()
		if !override {
			stmt.WriteString("INSERT IGNORE INTO ")
			stmt.WriteString(DbTable(model))
			stmt.WriteString(" (")
			writej(stmt, columns)
			stmt.WriteString(") VALUES (")
			stmt.WriteString(strings.Repeat(",?", len(columns))[1:])
			stmt.WriteString(");")
		} else {
			omitDict := map[string]struct{}{pkName: {}}
			for i := range omittedFields {
				omitDict[omittedFields[i]] = struct{}{}
			}
			noOfCols := len(columns)
			stmt.WriteString("INSERT INTO ")
			stmt.WriteString(DbTable(model))
			stmt.WriteString(" (")
			writej(stmt, columns)
			stmt.WriteString(") VALUES (")
			stmt.WriteString(strings.Repeat(",?", noOfCols)[1:])
			stmt.WriteString(") ON DUPLICATE KEY UPDATE ")
			first := true
			for i := range columns {
				if _, ok := omitDict[columns[i]]; ok {
					continue
				}
				if first {
					stmt.WriteString(columns[i])
					stmt.WriteString(" = VALUES(")
					stmt.WriteString(columns[i])
					stmt.WriteString(")")
				} else {
					stmt.WriteString(",")
					stmt.WriteString(columns[i])
					stmt.WriteString(" = VALUES(")
					stmt.WriteString(columns[i])
					stmt.WriteString(")")
				}
				first = false
			}
			clear(omitDict)
		}
		return db.ExecContext(ctx, stmt.String(), model.Values()...)
	case sequel.CompositeKeyer:
		names, _, _ := v.CompositeKey()
		columns := model.Columns()
		if !override {
			stmt.WriteString("INSERT IGNORE INTO " + DbTable(model) + " (")
			writej(stmt, columns)
			stmt.WriteString(") VALUES (" + strings.Repeat(",?", len(columns))[1:] + ");")
		} else {
			dict := make(map[string]struct{}, len(names)+len(omittedFields))
			for _, name := range names {
				dict[name] = struct{}{}
			}
			for _, name := range omittedFields {
				dict[name] = struct{}{}
			}
			noOfCols := len(columns)
			stmt.WriteString("INSERT INTO " + DbTable(model) + " (")
			writej(stmt, columns)
			stmt.WriteString(") VALUES (" + strings.Repeat(",?", noOfCols)[1:] + ") ON DUPLICATE KEY UPDATE ")
			first := true
			for _, column := range columns {
				if _, ok := dict[column]; ok {
					continue
				}
				if !first {
					stmt.WriteString(",")
				}
				stmt.WriteString(column)
				stmt.WriteString(" = VALUES(")
				stmt.WriteString(column)
				stmt.WriteString(")")
				first = false
			}
			if first {
				return nil, fmt.Errorf("no columns to update for model %T", model)
			}
			clear(dict)
		}
		return db.ExecContext(ctx, stmt.String(), model.Values()...)
	default:
		panic("unreachable")
	}
}

// Upsert is a helper function to upsert multiple records.
func Upsert[T interface {
	sequel.Keyer
	sequel.Inserter
}, Ptr sequel.PtrScanner[T]](ctx context.Context, db sequel.DB, data []T, override bool, omittedFields ...string) (sql.Result, error) {
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
			stmt.WriteString("INSERT INTO ")
			stmt.WriteString(DbTable(model))
			stmt.WriteString(" (")
			writej(stmt, columns)
			stmt.WriteString(") VALUES ")
		} else {
			stmt.WriteString("INSERT IGNORE INTO ")
			stmt.WriteString(DbTable(model))
			stmt.WriteString(" (")
			writej(stmt, columns)
			stmt.WriteString(") VALUES ")
		}
		placeholder := ",(" + strings.Repeat(",?", noOfCols)[1:] + ")"
		stmt.WriteString(placeholder[1:])
		values := data[0].Values()
		values = append(values[:idx], values[idx+1:]...)
		args = append(args, values...)
		for i := 1; i < noOfData; i++ {
			stmt.WriteString(placeholder)
			values := data[i].Values()
			values = append(values[:idx], values[idx+1:]...)
			args = append(args, values...)
		}
	case sequel.PrimaryKeyer:
		pkName, _, _ := v.PK()
		omittedFields = append(omittedFields, pkName)
		if override {
			stmt.WriteString("INSERT INTO ")
			stmt.WriteString(DbTable(model))
			stmt.WriteString(" (")
			writej(stmt, columns)
			stmt.WriteString(") VALUES ")
		} else {
			stmt.WriteString("INSERT IGNORE INTO ")
			stmt.WriteString(DbTable(model))
			stmt.WriteString(" (")
			writej(stmt, columns)
			stmt.WriteString(") VALUES ")
		}
		stmt.WriteString(model.SQLInsertPlaceholders(0))
		args = append(args, data[0].Values()...)
		for i := 1; i < noOfData; i++ {
			stmt.WriteString(",")
			stmt.WriteString(model.SQLInsertPlaceholders(i))
			args = append(args, data[i].Values()...)
		}
	case sequel.CompositeKeyer:
		if override {
			stmt.WriteString("INSERT INTO ")
			stmt.WriteString(DbTable(model))
			stmt.WriteString(" (")
			writej(stmt, columns)
			stmt.WriteString(") VALUES ")
		} else {
			stmt.WriteString("INSERT IGNORE INTO ")
			stmt.WriteString(DbTable(model))
			stmt.WriteString(" (")
			writej(stmt, columns)
			stmt.WriteString(") VALUES ")
		}
		stmt.WriteString(model.SQLInsertPlaceholders(0))
		args = append(args, data[0].Values()...)
		for i := 1; i < noOfData; i++ {
			stmt.WriteString(",")
			stmt.WriteString(model.SQLInsertPlaceholders(i))
			args = append(args, data[i].Values()...)
		}
		_, idxs, _ := v.CompositeKey()
		// Exclude primary key, don't update it
		for i := len(idxs) - 1; i >= 0; i-- {
			columns = append(columns[:idxs[i]], columns[idxs[i]+1:]...)
		}
	}
	if override {
		stmt.WriteString(" ON DUPLICATE KEY UPDATE ")
		// Don't update primary key when we do upsert
		omitDict := map[string]struct{}{}
		for i := range omittedFields {
			omitDict[omittedFields[i]] = struct{}{}
		}
		noOfCols = len(columns)
		first := true
		for i := range columns {
			if _, ok := omitDict[columns[i]]; ok {
				continue
			}
			if first {
				stmt.WriteString(columns[i])
				stmt.WriteString("=VALUES(")
				stmt.WriteString(columns[i])
				stmt.WriteString(")")
			} else {
				stmt.WriteString(",")
				stmt.WriteString(columns[i])
				stmt.WriteString("=VALUES(")
				stmt.WriteString(columns[i])
				stmt.WriteString(")")
			}
			first = false
		}
		clear(omitDict)
	}
	stmt.WriteString(";")
	return db.ExecContext(ctx, stmt.String(), args...)
}
{{ end }}

// FindByPK is to find single record using primary key.
func FindByPK[T sequel.KeyScanner, Ptr sequel.KeyPtrScanner[T]](ctx context.Context, db sequel.DB, model Ptr) error {
	switch v := any(model).(type) {
	case sequel.KeyFinder:
		query, args := v.FindOneByPKStmt()
		return db.QueryRowContext(ctx, query, args...).Scan(model.Addrs()...)
	case sequel.PrimaryKeyer:
		pkName, _, pk := v.PK()
		stmt := strpool.AcquireString()
		stmt.WriteString("SELECT ")
		writej(stmt, TableColumns(model))
		stmt.WriteString(" FROM ")
		stmt.WriteString(DbTable(model))
		stmt.WriteString(" WHERE ")
		stmt.WriteString(pkName)
		stmt.WriteString(" = {{ quoteVar 1 }} LIMIT 1;")
		row := db.QueryRowContext(ctx, stmt.String(), pk)
		strpool.ReleaseString(stmt)
		return row.Scan(model.Addrs()...)
	case sequel.CompositeKeyer:
		keyNames, _, keys := v.CompositeKey()
		{{ if isStaticVar -}}
		stmt := strpool.AcquireString()
		stmt.WriteString("SELECT ")
		writej(stmt, TableColumns(model))
		stmt.WriteString(" FROM ")
		stmt.WriteString(DbTable(model))
		stmt.WriteString(" WHERE (")
		writej(stmt, keyNames)
		stmt.WriteString(") = ({{ quoteVar 1 }}")
		stmt.WriteString(strings.Repeat(",{{ quoteVar 1 }}", len(keyNames) - 1))
		stmt.WriteString(") LIMIT 1;")
		row := db.QueryRowContext(ctx, stmt.String(), keys...)
		strpool.ReleaseString(stmt)
		return row.Scan(model.Addrs()...)
		{{ else -}}
		stmt := strpool.AcquireString()
		stmt.WriteString("SELECT ")
		writej(stmt, TableColumns(model))
		stmt.WriteString(" FROM ")
		stmt.WriteString(DbTable(model))
		stmt.WriteString(" WHERE (")
		writej(stmt, keyNames)
		stmt.WriteString(") = ({{ quoteVar 1 }}")
		noOfKeys := len(keyNames)
		for i := 1; i < noOfKeys; i++ {
			stmt.WriteString(",")
			stmt.WriteString(wrapVar(i+1))
		}
		stmt.WriteString(") LIMIT 1;")
		row := db.QueryRowContext(ctx, stmt.String(), keys...)
		strpool.ReleaseString(stmt)
		return row.Scan(model.Addrs()...)
		{{ end -}}
	default:
		panic("unreachable")
	}
}

// UpdateByPK is to update single record using primary key.
func UpdateByPK[T sequel.KeyValuer](ctx context.Context, db sequel.DB, model T) (sql.Result, error) {
	switch v := any(model).(type) {
	case sequel.KeyUpdater:
		query, args := v.UpdateOneByPKStmt()
		return db.ExecContext(ctx, query, args...)
	case sequel.PrimaryKeyer:
		pkName, pkIdx, pk := v.PK()
		columns := model.Columns()
		columns = append(columns[:pkIdx], columns[pkIdx+1:]...)
		values := model.Values()
		values = append(values[:pkIdx], append(values[pkIdx+1:], pk)...)
		{{ if isStaticVar -}}
		return db.ExecContext(ctx, "UPDATE "+DbTable(model)+" SET "+strings.Join(columns, " = {{ quoteVar 1 }},")+" = {{ quoteVar 1 }} WHERE "+pkName+" = {{ quoteVar 1 }};", append(values, pk)...)
		{{ else -}}
		stmt := strpool.AcquireString()
		stmt.WriteString("UPDATE ")
		stmt.WriteString(DbTable(model))
		stmt.WriteString(" SET ")
		stmt.WriteString(columns[0])
		stmt.WriteString(" = {{ quoteVar 1 }}")
		noOfColumn := len(columns)
		for i := 1; i < noOfColumn; i++ {
			stmt.WriteString(",")
			stmt.WriteString(columns[i])
			stmt.WriteString(" = ")
			stmt.WriteString(wrapVar(i + 1))
		}
		stmt.WriteString(" WHERE ")
		stmt.WriteString(pkName)
		stmt.WriteString(" = ")
		stmt.WriteString(wrapVar(noOfColumn + 2))
		stmt.WriteString(";")
		result, err := db.ExecContext(ctx, stmt.String(), append(values, pk)...)
		strpool.ReleaseString(stmt)
		return result, err
		{{ end -}}
	case sequel.CompositeKeyer:
		columns := model.Columns()
		keyNames, keyIdxs, keys := v.CompositeKey()
		values := model.Values()
		sort.Ints(keyIdxs)
		for i := len(keyIdxs) - 1; i >= 0; i-- {
			columns = append(columns[:keyIdxs[i]], columns[keyIdxs[i]+1:]...)
			values = append(values[:keyIdxs[i]], values[keyIdxs[i]+1:]...)
		}
		noOfColumn := len(columns)
		args := make([]any, 0, noOfColumn+len(keys))
		{{ if isStaticVar -}}
		query := "UPDATE " + DbTable(model) + " SET " +strings.Join(columns, " = {{ quoteVar 1 }},")+" = {{ quoteVar 1 }} WHERE (" + strings.Join(keyNames, ",") + ") = ({{ quoteVar 1 }}" + strings.Repeat(",{{ quoteVar 1 }}", len(keyNames)-1)+");"
		return db.ExecContext(ctx, query, append(args, keys...)...)
		{{ else -}}
		stmt := strpool.AcquireString()
		defer strpool.ReleaseString(stmt)
		stmt.WriteString("UPDATE ")
		stmt.WriteString(DbTable(model))
		stmt.WriteString(" SET ")
		stmt.WriteString(columns[0])
		stmt.WriteString(" = {{ quoteVar 1 }}")
		for i := 1; i < noOfColumn; i++ {
			stmt.WriteString(",")
			stmt.WriteString(columns[i])
			stmt.WriteString(" = ")
			stmt.WriteString(wrapVar(i + 1))
		}
		stmt.WriteString(" WHERE (")
		writej(stmt, keyNames)
		stmt.WriteString(") = (")
		stmt.WriteString(wrapVar(noOfColumn + 1))
		noOfKey := len(keyNames)
		for i := 1; i < noOfKey; i++ {
			stmt.WriteString(",")
			stmt.WriteString(columns[i])
			stmt.WriteString(" = ")
			stmt.WriteString(wrapVar(noOfColumn + i + 1))
		}
		stmt.WriteString(");")
		args = append(args, values...)
		return db.ExecContext(ctx, stmt.String(), append(args, keys...)...)
		{{ end -}}
	default:
		panic("unreachable")
	}
}

// DeleteByPK is to update single record using primary key.
func DeleteByPK[T sequel.KeyValuer](ctx context.Context, db sequel.DB, model T) (sql.Result, error) {
	switch v := any(model).(type) {
	case sequel.KeyDeleter:
		query, args := v.DeleteOneByPKStmt()
		return db.ExecContext(ctx, query, args...)
	case sequel.PrimaryKeyer:
		pkName, _, pk := v.PK()
		stmt := strpool.AcquireString()
		stmt.WriteString("DELETE FROM ")
		stmt.WriteString(DbTable(model))
		stmt.WriteString(" WHERE ")
		stmt.WriteString(pkName)
		stmt.WriteString(" = {{ quoteVar 1 }};")
		result, err := db.ExecContext(ctx, stmt.String(), pk)
		strpool.ReleaseString(stmt)
		return result, err
	case sequel.CompositeKeyer:
		keyNames, _, keys := v.CompositeKey()
		{{ if isStaticVar -}}
		stmt := strpool.AcquireString()
		stmt.WriteString("DELETE FROM ")
		stmt.WriteString(DbTable(model))
		stmt.WriteString(" WHERE (")
		writej(stmt, keyNames)
		stmt.WriteString(") = ({{ quoteVar 1 }}")
		stmt.WriteString(strings.Repeat(",{{ quoteVar 1 }}", len(keyNames)-1))
		stmt.WriteString(");")
		result, err := db.ExecContext(ctx, stmt.String(), keys...)
		strpool.ReleaseString(stmt)
		return result, err
		{{ else -}}
		stmt := strpool.AcquireString()
		stmt.WriteString("DELETE FROM ")
		stmt.WriteString(DbTable(model))
		stmt.WriteString(" WHERE (")
		writej(stmt, keyNames)
		stmt.WriteString(") = ({{ quoteVar 1 }}")
		noOfKey := len(keyNames)
		for i := 1; i < noOfKey; i++ {
			stmt.WriteString(",")
			stmt.WriteString(wrapVar(i+1))
		}
		stmt.WriteString(");")
		result, err := db.ExecContext(ctx, stmt.String(), keys...)
		strpool.ReleaseString(stmt) // Clean up
		keys = nil // Clean up
		return result, err
		{{ end -}}
	default:
		panic("unreachable")
	}
}

type LockMode string

func (l LockMode) LockMode() string {
	return (string)(l)
}

const (
{{ if eq driver "postgres" -}}
	LockForUpdate 	   LockMode = "FOR UPDATE"
	LockForNoKeyUpdate LockMode = "FOR NO KEY UPDATE"
	LockForShare       LockMode = "FOR SHARE"
	LockForKeyShare    LockMode = "FOR KEY SHARE"
{{ else }}
	LockForUpdate LockMode = "FOR UPDATE"
	LockForShare  LockMode = "FOR SHARE"
{{- end }}
)

func LockByPK[T sequel.KeyScanner, Ptr sequel.KeyPtrScanner[T]](ctx context.Context, tx *sql.Tx, model Ptr, locker sequel.RowLevelLocker) error {
	stmt := strpool.AcquireString()
	switch v := any(model).(type) {
	case sequel.PrimaryKeyer:
		pkName, _, pk := v.PK()
		stmt.WriteString("SELECT ")
		writej(stmt, TableColumns(model))
		stmt.WriteString(" FROM ")
		stmt.WriteString(DbTable(model))
		stmt.WriteString(" WHERE ")
		stmt.WriteString(pkName)
		stmt.WriteString(" = {{ quoteVar 1 }} LIMIT 1 ")
		stmt.WriteString(locker.LockMode())
		stmt.WriteString(";")
		row := tx.QueryRowContext(ctx, stmt.String(), pk)
		strpool.ReleaseString(stmt)
		return row.Scan(model.Addrs()...)
	case sequel.CompositeKeyer:
		keyNames, _, keys := v.CompositeKey()
		{{ if isStaticVar -}}
		stmt.WriteString("SELECT ")
		writej(stmt, TableColumns(model))
		stmt.WriteString(" FROM ")
		stmt.WriteString(DbTable(model))
		stmt.WriteString(" WHERE ")
		stmt.WriteString("(")
		writej(stmt, keyNames)
		stmt.WriteString(") = ({{ quoteVar 1 }}"+strings.Repeat(",{{ quoteVar 1 }}", len(keyNames)-1)+") LIMIT 1 ")
		stmt.WriteString(locker.LockMode())
		stmt.WriteString(";")
		row := tx.QueryRowContext(ctx, stmt.String(), keys...)
		strpool.ReleaseString(stmt)
		return row.Scan(model.Addrs()...)
		{{ else -}}
		stmt.WriteString("SELECT ")
		writej(stmt, TableColumns(model))
		stmt.WriteString(" FROM ")
		stmt.WriteString(DbTable(model))
		stmt.WriteString(" WHERE ")
		stmt.WriteString("(")
		writej(stmt, keyNames)
		stmt.WriteString(") = ({{ quoteVar 1 }}")
		noOfKey := len(keyNames)
		for i := 1; i < noOfKey; i++ {
			stmt.WriteString(",")
			stmt.WriteString(wrapVar(i+1))
		}
		stmt.WriteString(") LIMIT 1 ")
		stmt.WriteString(locker.LockMode())
		stmt.WriteString(";")
		row := tx.QueryRowContext(ctx, stmt.String(), keys...)
		strpool.ReleaseString(stmt)
		return row.Scan(model.Addrs()...)
		{{ end -}}
	default:
		panic("unreachable")
	}
}

type PaginateStmt struct {
	Select    []string
	FromTable string
	Where     sequel.WhereClause
	OrderBy   []sequel.OrderByClause
	Limit     uint16
}

func Paginate[T sequel.KeyValuer, Ptr sequel.KeyValueScanner[T]](stmt PaginateStmt) *Pager[T, Ptr] {
	// Set default limit
	if stmt.Limit == 0 {
		stmt.Limit = 100
	}
	return &Pager[T, Ptr]{
		stmt: &stmt,
	}
}

type Pager[T sequel.KeyValuer, Ptr sequel.KeyValueScanner[T]] struct {
	stmt *PaginateStmt
}

func (r *Pager[T, Ptr]) Prev(ctx context.Context, db sequel.DB, cursor ...T) iter.Seq2[[]T, error] {
	return func(yield func([]T, error) bool) {
		var (
			v         T
			hasCursor bool
			maxLimit  = int(r.stmt.Limit)
			rows 	  *sql.Rows
			err 	  error
			noOfData  int
			data 	  = make([]T, 0, maxLimit+1) // Add one to limit to find prev cursor
		)
		if len(cursor) > 0 {
			v = cursor[0]
			hasCursor = true
		}

		stmt := AcquireStmt()
		defer ReleaseStmt(stmt)

		for {
			if hasCursor {
				if err = FindByPK(ctx, db, Ptr(&v)); err != nil {
					yield(nil, err)
					return
				}
			}

			stmt.WriteString("SELECT ")
			writej(stmt, TableColumns(v))
			stmt.WriteString(" FROM ")
			stmt.WriteString(DbTable(v))

			if r.stmt.Where != nil || len(r.stmt.OrderBy) > 0 || hasCursor {
				stmt.WriteString(" WHERE ")
				if r.stmt.Where != nil {
					r.stmt.Where(stmt)
				}
				if n := len(r.stmt.OrderBy); n > 0 {
					colDict := make(map[string]int)
					columns := v.Columns()
					values := v.Values()
					for i := range columns {
						colDict[columns[i]] = i
					}

					buf := strpool.AcquireString()
					for i := 0; i < n; i++ {
						orderedColumn := r.stmt.OrderBy[i]
						columnName := orderedColumn.ColumnName()
						value := values[colDict[columnName]]
						if orderedColumn.Asc() {
							// If ascending
							stmt.WriteString(columnName)
							stmt.WriteString(" <= ")
							stmt.WriteString(stmt.Var(value))
							stmt.WriteString(" AND ")
							buf.WriteString(columnName)
							buf.WriteString(" < ")
							buf.WriteString(stmt.Var(value))
							buf.WriteString(" OR ")
						} else {
							// If descending
							stmt.WriteString(columnName)
							stmt.WriteString(" >= ")
							stmt.WriteString(stmt.Var(value))
							stmt.WriteString(" AND ")
							buf.WriteString(columnName)
							buf.WriteString(" > ")
							buf.WriteString(stmt.Var(value))
							buf.WriteString(" OR ")
						}
					}
					clear(colDict)
					stmt.WriteString("(")
					stmt.WriteString(buf.String())
					strpool.ReleaseString(buf)
				}

				if hasCursor {
					// Check the primary key and compare value
					switch vi := any(v).(type) {
					case sequel.CompositeKeyer:
						pkNames, _, vals := vi.CompositeKey()
						stmt.WriteString("((")
						writej(stmt, pkNames)
						stmt.WriteString(") <= ")
						stmt.WriteString(stmt.Vars(vals))
						stmt.WriteString(")")
					case sequel.PrimaryKeyer:
						pkName, _, val := vi.PK()
						stmt.WriteString("(")
						stmt.WriteString(pkName)
						stmt.WriteString(" <= ")
						stmt.WriteString(stmt.Var(val))
						stmt.WriteString(")")
					default:
						panic("unreachable")
					}
				}
			}

			stmt.WriteString(" ORDER BY ")
			if len(r.stmt.OrderBy) > 0 {
				for i := range r.stmt.OrderBy {
					if r.stmt.OrderBy[i].Asc() {
						stmt.WriteString(r.stmt.OrderBy[i].ColumnName())
						stmt.WriteString(" ASC,")
					} else {
						stmt.WriteString(r.stmt.OrderBy[i].ColumnName())
						stmt.WriteString(" DESC,")
					}
				}
				var (
					suffix  = " DESC"
					lastCol = r.stmt.OrderBy[len(r.stmt.OrderBy)-1]
				)
				if lastCol.Asc() {
					suffix = " ASC"
				}
				switch vi := any(v).(type) {
				case sequel.CompositeKeyer:
					pkNames, _, _ := vi.CompositeKey()
					stmt.WriteString("(")
					writej(stmt, pkNames)
					stmt.WriteString(")")
					stmt.WriteString(suffix)
				case sequel.PrimaryKeyer:
					pkName, _, _ := vi.PK()
					stmt.WriteString(pkName)
					stmt.WriteString(suffix)
				default:
					panic("unreachable")
				}
			} else {
				// If there is no order by clause,
				// descending is the default order
				switch vi := any(v).(type) {
				case sequel.CompositeKeyer:
					pkNames, _, _ := vi.CompositeKey()
					stmt.WriteString("(")
					writej(stmt, pkNames)
					stmt.WriteString(") DESC")
				case sequel.PrimaryKeyer:
					pkName, _, _ := vi.PK()
					stmt.WriteString(pkName)
					stmt.WriteString(" DESC")
				default:
					panic("unreachable")
				}
			}
			// Add one to limit to find next cursor
			stmt.WriteString(" LIMIT ")
			stmt.WriteString(strconv.Itoa(maxLimit+1))
			stmt.WriteString(";")

			rows, err = db.QueryContext(ctx, stmt.Query(), stmt.Args()...)
			stmt.Reset()
			if err != nil {
				yield(nil, err)
				return
			}

			for rows.Next() {
				var v T
				if err := rows.Scan(Ptr(&v).Addrs()...); err != nil {
					rows.Close()
					yield(nil, err)
					return
				}
				data = append(data, v)
			}
			if err = rows.Close(); err != nil {
				yield(nil, err)
				return
			}
			if err = rows.Err(); err != nil {
				yield(nil, err)
				return
			}

			noOfData = len(data)
			if noOfData <= maxLimit {
				yield(data, nil)
				return
			}

			if !yield(data[:noOfData-1], nil) {
				return
			}

			v = data[noOfData-1] // Set previous cursor
			data = data[:0] // Reset result
			hasCursor = true
		}
	}
}

func (r *Pager[T, Ptr]) Next(ctx context.Context, db sequel.DB, cursor ...T) iter.Seq2[[]T, error] {
	return func(yield func([]T, error) bool) {
		var (
			v         T
			hasCursor bool
			maxLimit  = int(r.stmt.Limit)
			rows 	  *sql.Rows
			err		  error
			noOfData  int
			data 	  = make([]T, 0, maxLimit+1) // Add one to limit to find next cursor
		)
		if len(cursor) > 0 {
			v = cursor[0]
			hasCursor = true
		}

		stmt := AcquireStmt()
		defer ReleaseStmt(stmt)

		for {
			if hasCursor {
				if err = FindByPK(ctx, db, Ptr(&v)); err != nil {
					yield(nil, err)
					return
				}
			}

			stmt.WriteString("SELECT ")
			writej(stmt, TableColumns(v))
			stmt.WriteString(" FROM ")
			stmt.WriteString(DbTable(v))
			if r.stmt.Where != nil || len(r.stmt.OrderBy) > 0 || hasCursor {
				stmt.WriteString(" WHERE ")
				if r.stmt.Where != nil {
					r.stmt.Where(stmt)
				}

				if len(r.stmt.OrderBy) > 0 {
					colDict := make(map[string]int)
					columns := v.Columns()
					values := v.Values()
					for i := range columns {
						colDict[columns[i]] = i
					}

					for i, orderBy := range r.stmt.OrderBy {
						colName := orderBy.ColumnName()
						val := values[colDict[colName]]
						if i > 0 {
							stmt.WriteString(" AND ")
						}
						if orderBy.Asc() {
							// If ascending
							stmt.WriteString(colName)
							stmt.WriteString(" >= ")
							stmt.WriteString(stmt.Var(val))
						} else {
							// If descending
							stmt.WriteString(colName)
							stmt.WriteString(" <= ")
							stmt.WriteString(stmt.Var(val))
						}
					}
					stmt.WriteString(" AND (")
					for i, orderBy := range r.stmt.OrderBy {
						colName := orderBy.ColumnName()
						val := values[colDict[colName]]
						if i > 0 {
							stmt.WriteString(" OR ")
						}
						if orderBy.Asc() {
							// If ascending
							stmt.WriteString(colName)
							stmt.WriteString(" > ")
							stmt.WriteString(stmt.Var(val))
						} else {
							// If descending
							stmt.WriteString(colName)
							stmt.WriteString(" < ")
							stmt.WriteString(stmt.Var(val))
						}
					}
					clear(colDict)
					stmt.WriteString(" OR ")
				}

				if hasCursor {
					// Check the primary key and compare value
					switch vi := any(v).(type) {
					case sequel.CompositeKeyer:
						pkNames, _, vals := vi.CompositeKey()
						stmt.WriteString("((")
						writej(stmt, pkNames)
						stmt.WriteString(") >= ")
						stmt.WriteString(stmt.Vars(vals))
						stmt.WriteString(")")
					case sequel.PrimaryKeyer:
						pkName, _, val := vi.PK()
						stmt.WriteString("(")
						stmt.WriteString(pkName)
						stmt.WriteString(" >= ")
						stmt.WriteString(stmt.Var(val))
						stmt.WriteString(")")
					default:
						panic("unreachable")
					}
				}
			}

			stmt.WriteString(" ORDER BY ")
			if len(r.stmt.OrderBy) > 0 {
				for i := range r.stmt.OrderBy {
					if r.stmt.OrderBy[i].Asc() {
						stmt.WriteString(r.stmt.OrderBy[i].ColumnName())
						stmt.WriteString(" ASC,")
					} else {
						stmt.WriteString(r.stmt.OrderBy[i].ColumnName())
						stmt.WriteString(" DESC,")
					}
				}
				var (
					suffix  = " ASC"
					lastCol = r.stmt.OrderBy[len(r.stmt.OrderBy)-1]
				)
				if !lastCol.Asc() {
					suffix = " DESC"
				}
				switch vi := any(v).(type) {
				case sequel.CompositeKeyer:
					pkNames, _, _ := vi.CompositeKey()
					stmt.WriteString("(")
					writej(stmt, pkNames)
					stmt.WriteString(")")
					stmt.WriteString(suffix)
				case sequel.PrimaryKeyer:
					pkName, _, _ := vi.PK()
					stmt.WriteString(pkName)
					stmt.WriteString(suffix)
				default:
					panic("unreachable")
				}
			} else {
				// If there is no order by clause,
				// ascending is the default order
				switch vi := any(v).(type) {
				case sequel.CompositeKeyer:
					pkNames, _, _ := vi.CompositeKey()
					stmt.WriteString("(")
					writej(stmt, pkNames)
					stmt.WriteString(") ASC")
				case sequel.PrimaryKeyer:
					pkName, _, _ := vi.PK()
					stmt.WriteString(pkName)
					stmt.WriteString(" ASC")
				default:
					panic("unreachable")
				}
			}
			stmt.WriteString(" LIMIT ")
			stmt.WriteString(strconv.Itoa(maxLimit+1))
			stmt.WriteString(";")

			rows, err = db.QueryContext(ctx, stmt.Query(), stmt.Args()...)
			stmt.Reset()
			if err != nil {
				yield(nil, err)
				return
			}

			for rows.Next() {
				var v T
				if err = rows.Scan(Ptr(&v).Addrs()...); err != nil {
					rows.Close()
					yield(nil, err)
					return
				}
				data = append(data, v)
			}
			if err = rows.Close(); err != nil {
				yield(nil, err)
				return
			}
			if err = rows.Err(); err != nil {
				yield(nil, err)
				return
			}

			noOfData = len(data)
			if noOfData <= maxLimit {
				yield(data, nil)
				return
			}

			if !yield(data[:noOfData-1], nil) {
				return
			}

			v = data[noOfData-1] // Set next cursor
			data = data[:0] // Reset result
			hasCursor = true
		}
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

// WriteSQL writes a complete SELECT statement. FromTable must be resolved by the caller.
func (s SelectStmt) WriteSQL(w sequel.Stmt) error {
	if s.FromTable == "" {
		return fmt.Errorf("missing table name for SELECT statement")
	}
	w.WriteString("SELECT ")
	if len(s.Select) > 0 {
		writej(w, s.Select)
	} else {
		w.WriteString("*")
	}
	w.WriteString(" FROM ")
	w.WriteString(s.FromTable)
	if s.Where != nil {
		w.WriteString(" WHERE ")
		s.Where(w)
	}
	if n := len(s.GroupBy); n > 0 {
		w.WriteString(" GROUP BY ")
		w.WriteString(s.GroupBy[0])
		for i := 1; i < n; i++ {
			w.WriteString(",")
			w.WriteString(s.GroupBy[i])
		}
	}
	if n := len(s.OrderBy); n > 0 {
		w.WriteString(" ORDER BY ")
		w.WriteString(s.OrderBy[0].ColumnName())
		if !s.OrderBy[0].Asc() {
			w.WriteString(" DESC")
		}
		for i := 1; i < n; i++ {
			w.WriteString(",")
			w.WriteString(s.OrderBy[i].ColumnName())
			if !s.OrderBy[i].Asc() {
				w.WriteString(" DESC")
			}
		}
	}
	if s.Limit > 0 {
		w.WriteString(" LIMIT ")
		w.WriteString(strconv.Itoa((int)(s.Limit)))
	}
	if s.Offset > 0 {
		w.WriteString(" OFFSET ")
		w.WriteString(strconv.FormatUint(s.Offset, 10))
	}
	w.WriteString(";")
	return nil
}

func QueryStmt[T any, Ptr sequel.PtrScanner[T]](ctx context.Context, db sequel.DB, stmtFunc func(T) SelectStmt) ([]T, error) {
	var v T
	selectStmt := stmtFunc(v)
	stmt := AcquireStmt()
	if len(selectStmt.Select) == 0 {
		if cols, ok := any(v).(sequel.Columner); ok {
			selectStmt.Select = TableColumns(cols)
		}
	}
	if selectStmt.FromTable != "" {
		selectStmt.FromTable = dbName(v) + selectStmt.FromTable
	} else if table, ok := any(v).(sequel.Tabler); ok {
		selectStmt.FromTable = DbTable(table)
	} else {
		ReleaseStmt(stmt)
		return nil, fmt.Errorf("missing table name for model %T", v)
	}
	if err := selectStmt.WriteSQL(stmt); err != nil {
		ReleaseStmt(stmt)
		return nil, err
	}
	rows, err := db.QueryContext(ctx, stmt.Query(), stmt.Args()...)
	ReleaseStmt(stmt)
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
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

type SelectOneStmt struct {
	Select    []string
	FromTable string
	Where     sequel.WhereClause
	OrderBy   []sequel.OrderByClause
	GroupBy   []string
}

// WriteSQL writes a SELECT statement limited to one row. FromTable must be resolved by the caller.
func (s SelectOneStmt) WriteSQL(w sequel.Stmt) error {
	if s.FromTable == "" {
		return fmt.Errorf("missing table name for SELECT statement")
	}
	w.WriteString("SELECT ")
	if len(s.Select) > 0 {
		writej(w, s.Select)
	} else {
		w.WriteString("*")
	}
	w.WriteString(" FROM ")
	w.WriteString(s.FromTable)
	if s.Where != nil {
		w.WriteString(" WHERE ")
		s.Where(w)
	}
	if n := len(s.GroupBy); n > 0 {
		w.WriteString(" GROUP BY ")
		w.WriteString(s.GroupBy[0])
		for i := 1; i < n; i++ {
			w.WriteString(",")
			w.WriteString(s.GroupBy[i])
		}
	}
	if n := len(s.OrderBy); n > 0 {
		w.WriteString(" ORDER BY ")
		w.WriteString(s.OrderBy[0].ColumnName())
		if !s.OrderBy[0].Asc() {
			w.WriteString(" DESC")
		}
		for i := 1; i < n; i++ {
			w.WriteString(",")
			w.WriteString(s.OrderBy[i].ColumnName())
			if !s.OrderBy[i].Asc() {
				w.WriteString(" DESC")
			}
		}
	}
	w.WriteString(" LIMIT 1;")
	return nil
}

func QueryOneStmt[T any, Ptr sequel.PtrScanner[T]](ctx context.Context, db sequel.DB, stmtFunc func(T) SelectOneStmt) (Ptr, error) {
	var v T
	selectOneStmt := stmtFunc(v)
	if len(selectOneStmt.Select) == 0 {
		if cols, ok := any(v).(sequel.Columner); ok {
			selectOneStmt.Select = TableColumns(cols)
		}
	}
	if selectOneStmt.FromTable != "" {
		selectOneStmt.FromTable = dbName(v) + selectOneStmt.FromTable
	} else if table, ok := any(v).(sequel.Tabler); ok {
		selectOneStmt.FromTable = DbTable(table)
	} else {
		return nil, fmt.Errorf("missing table name for model %T", v)
	}
	stmt := AcquireStmt()
	defer ReleaseStmt(stmt)
	if err := selectOneStmt.WriteSQL(stmt); err != nil {
		return nil, err
	}
	row := db.QueryRowContext(ctx, stmt.Query(), stmt.Args()...)
	if err := row.Scan(Ptr(&v).Addrs()...); err != nil {
		return nil, err
	}
	return &v, nil
}

type UpdateStmt struct {
	Table	string
	Set		[]sequel.SetClause
	Where	sequel.WhereClause
	OrderBy []sequel.OrderByClause
	{{ if ne driver "postgres" -}}
	{{- /* postgres */ -}}
	Limit   uint16
	{{ end -}}
}

type DeleteStmt struct {
	FromTable string
	Where     sequel.WhereClause
	OrderBy   []sequel.OrderByClause
	{{ if ne driver "postgres" -}}
	{{- /* postgres */ -}}
	Limit     uint16
	{{ end -}}
}

// WriteSQL writes a complete UPDATE statement. Table must be resolved by the caller.
func (s UpdateStmt) WriteSQL(w sequel.Stmt) error {
	if s.Table == "" {
		return fmt.Errorf("missing table name for UPDATE statement")
	}
	w.WriteString("UPDATE ")
	w.WriteString(s.Table)
	if len(s.Set) > 0 {
		w.WriteString(" SET ")
		s.Set[0](w)
		for i := 1; i < len(s.Set); i++ {
			w.WriteString(",")
			s.Set[i](w)
		}
	}
	if s.Where != nil {
		w.WriteString(" WHERE ")
		s.Where(w)
	}
	if n := len(s.OrderBy); n > 0 {
		if s.OrderBy[0].Asc() {
			w.WriteString(" ORDER BY ")
			w.WriteString(s.OrderBy[0].ColumnName())
		} else {
			w.WriteString(" ORDER BY ")
			w.WriteString(s.OrderBy[0].ColumnName())
			w.WriteString(" DESC")
		}
		for i := 1; i < n; i++ {
			if s.OrderBy[i].Asc() {
				w.WriteString(",")
				w.WriteString(s.OrderBy[i].ColumnName())
				w.WriteString(" ASC")
			} else {
				w.WriteString(",")
				w.WriteString(s.OrderBy[i].ColumnName())
				w.WriteString(" DESC")
			}
		}
	}
	{{ if ne driver "postgres" -}}
	{{- /* postgres */ -}}
	if s.Limit > 0 {
		w.WriteString(" LIMIT ")
		w.WriteString(strconv.Itoa((int)(s.Limit)))
	}
	{{ end -}}
	w.WriteString(";")
	return nil
}

// WriteSQL writes a complete DELETE FROM statement. FromTable must be resolved by the caller.
func (s DeleteStmt) WriteSQL(w sequel.Stmt) error {
	if s.FromTable == "" {
		return fmt.Errorf("missing table name for DELETE FROM statement")
	}
	w.WriteString("DELETE FROM ")
	w.WriteString(s.FromTable)
	if s.Where != nil {
		w.WriteString(" WHERE ")
		s.Where(w)
	}
	if n := len(s.OrderBy); n > 0 {
		w.WriteString(" ORDER BY ")
		w.WriteString(s.OrderBy[0].ColumnName())
		if !s.OrderBy[0].Asc() {
			w.WriteString(" DESC")
		}
		for i := 1; i < n; i++ {
			w.WriteString(",")
			w.WriteString(s.OrderBy[i].ColumnName())
			if !s.OrderBy[i].Asc() {
				w.WriteString(" DESC")
			}
		}
	}
	{{ if ne driver "postgres" -}}
	{{- /* postgres */ -}}
	if s.Limit > 0 {
		w.WriteString(" LIMIT ")
		w.WriteString(strconv.Itoa((int)(s.Limit)))
	}
	{{ end -}}
	w.WriteString(";")
	return nil
}

func ExecStmt[T any, Stmt interface {
	UpdateStmt | DeleteStmt
}](ctx context.Context, db sequel.DB, stmtFunc func(T) Stmt) (sql.Result, error) {
	var v T
	stmt := AcquireStmt()
	defer ReleaseStmt(stmt)

	var err error
	switch vi := any(stmtFunc(v)).(type) {
	case UpdateStmt:
		if vt, ok := any(v).(sequel.Tabler); ok {
			vi.Table = DbTable(vt)
		} else {
			vi.Table = dbName(v) + vi.Table
		}
		err = vi.WriteSQL(stmt)
	case DeleteStmt:
		if vt, ok := any(v).(sequel.Tabler); ok {
			vi.FromTable = DbTable(vt)
		} else {
			vi.FromTable = dbName(v) + vi.FromTable
		}
		err = vi.WriteSQL(stmt)
	}
	if err != nil {
		return nil, err
	}
	return db.ExecContext(ctx, stmt.Query(), stmt.Args()...)
}

var (
	pool = sync.Pool{
		New: func() any {
			return new(SqlStmt)
		},
	}
)

func AcquireStmt() *SqlStmt {
	return pool.Get().(*SqlStmt)
}

func ReleaseStmt(stmt *SqlStmt) {
	if stmt != nil {
		stmt.Reset()
		pool.Put(stmt)
	}
}

type SqlStmt struct {
	blr	 strings.Builder
	pos  int
	args []any
}

var (
	_ sequel.Stmt = (*SqlStmt)(nil)
)

func (s *SqlStmt) Var(value any) string {
	s.pos++
	s.args = append(s.args, value)
	{{ if isStaticVar -}}
	return {{ quote varRune }}
	{{ else -}}
	return wrapVar(s.pos)
	{{ end -}}
}

// Vars groups placeholders in parentheses. Empty input returns "()" without
// adding arguments; callers must handle empty IN lists according to their dialect.
func (s *SqlStmt) Vars(values []any) string {
	noOfLen := len(values)
	if noOfLen == 0 {
		return "()"
	}
	s.args = append(s.args, values...)
	{{ if isStaticVar -}}
	return "(" + strings.Repeat(",{{ varRune }}", noOfLen)[1:] + ")"
	{{ else -}}
	buf := new(strings.Builder)
	buf.WriteString("(")
	i := s.pos
	s.pos += noOfLen
	for ; i < s.pos; i++ {
		if i < s.pos-1 {
			buf.WriteString(wrapVar(i+1) + ",")
		} else {
			buf.WriteString(wrapVar(i + 1))
		}
	}
	buf.WriteString(")")
	return buf.String()
	{{ end -}}
}

func (s *SqlStmt) Write(p []byte) (int, error) {
	return s.blr.Write(p)
}

func (s *SqlStmt) WriteString(v string) (int, error) {
	return s.blr.WriteString(v)
}

func (s *SqlStmt) Quote(v string) string {
	{{ if eq driver "postgres" -}}
	return pgutil.Quote(v)
	{{ else -}}
	return strconv.Quote(v)
	{{ end -}}
}

func (s *SqlStmt) Query() string {
	return s.blr.String()
}

func (s SqlStmt) Args() []any {
	return s.args
}

func (s *SqlStmt) Format(f fmt.State, verb rune) {
	str := s.blr.String()
	switch verb {
	case 's':
		f.Write(unsafe.Slice(unsafe.StringData(str), len(str)))
		return
	case 'v':
		if !f.Flag('#') && !f.Flag('+') {
			f.Write(unsafe.Slice(unsafe.StringData(str), len(str)))
			return
		}
	}

	var (
		args = s.args
		idx  int
		i    = 1
	)

	for {
		{{ if isStaticVar -}}
		idx = strings.Index(str, "?")
		{{ else -}}
		placeholder := wrapVar(i)
		idx = strings.Index(str, placeholder)
		{{ end -}}
		if idx < 0 {
			f.Write(unsafe.Slice(unsafe.StringData(str), len(str)))
			break
		}

		f.Write([]byte(str[:idx]))
		v := strf(args[0])
		f.Write(unsafe.Slice(unsafe.StringData(v), len(v)))
		{{ if isStaticVar -}}
		str = str[idx+1:]
		{{ else -}}
		str = str[idx+len(placeholder):]
		{{ end -}}
		args = args[1:]
		i++
	}
}

func (s *SqlStmt) Reset() {
	s.args = nil
	s.pos = 0
	s.blr.Reset()
}

func DbTable[T sequel.Tabler](model T) string {
	if v, ok := any(model).(sequel.Databaser); ok {
		return v.DatabaseName() + "." + model.TableName()
	}
	return model.TableName()
}

func TableColumns[T sequel.Columner](model T) []string {
	if v, ok := any(model).(sequel.SQLColumner); ok {
		return v.SQLColumns()
	}
	return model.Columns()
}

func dbName(model any) string {
	if v, ok := model.(sequel.Databaser); ok {
		return v.DatabaseName() + "."
	}
	return ""
}

{{ if not isStaticVar -}}
func wrapVar(i int) string {
	return {{ quote varRune }}+ strconv.Itoa(i)
}
{{ end }}

func strf(v any) string {
	switch vi := v.(type) {
	case nil:
		return "NULL"
	case string:
	{{ if eq driver "postgres" -}}
		return pgutil.Quote(vi)
	{{ else -}}
		return strconv.Quote(vi)
	{{ end -}}
	case []byte:
	{{ if eq driver "postgres" -}}
		return pgutil.Quote(unsafe.String(unsafe.SliceData(vi), len(vi)))
	{{ else -}}
		return strconv.Quote(unsafe.String(unsafe.SliceData(vi), len(vi)))
	{{ end -}}
	case bool:
		return strconv.FormatBool(vi)
	case int64:
		return strconv.FormatInt(vi, 10)
	case float64:
		return strconv.FormatFloat(vi, 'f', 10, 64)
	case time.Time:
		return strconv.Quote(vi.Format(time.RFC3339))
	case sql.RawBytes:
		return unsafe.String(unsafe.SliceData(vi), len(vi))
	case driver.Valuer:
		val, _ := vi.Value()
		return strf(val)
	default:
		vt := reflect.ValueOf(v)
		switch vt.Kind() {
		case reflect.String:
		{{ if eq driver "postgres" -}}
			return pgutil.Quote(vt.String())
		{{ else -}}
			return strconv.Quote(vt.String())
		{{ end -}}
		case reflect.Bool:
			return strconv.FormatBool(vt.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return strconv.FormatInt(vt.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return strconv.FormatUint(vt.Uint(), 10)
		case reflect.Float32, reflect.Float64:
			return strconv.FormatFloat(vt.Float(), 'f', 10, 64)
		default:
			panic("unreachable")
		}
	}
}

// writej writes comma-separated strings without creating an intermediate string.
func writej(w io.StringWriter, values []string) {
	if len(values) == 0 {
		return
	}
	if stmt, ok := w.(*SqlStmt); ok {
		w = &stmt.blr
	}
	if b, ok := w.(*strings.Builder); ok {
		size := len(values) - 1
		for _, value := range values {
			size += len(value)
		}
		b.Grow(size)
		b.WriteString(values[0])
		for _, value := range values[1:] {
			b.WriteString(",")
			b.WriteString(value)
		}
		return
	}
	w.WriteString(values[0])
	for _, value := range values[1:] {
		w.WriteString(",")
		w.WriteString(value)
	}
}
