{{- reserveImport "context" }}
{{- reserveImport "database/sql" }}
{{- reserveImport "database/sql/driver" }}
{{- reserveImport "strings" }}
{{- reserveImport "strconv" }}
{{- reserveImport "fmt" }}
{{- reserveImport "sync" }}
{{- reserveImport "unsafe" }}
{{- reserveImport "github.com/si3nloong/sqlgen/sequel" }}
{{- reserveImport "github.com/si3nloong/sqlgen/sequel/strpool" }}
{{- reserveImport "github.com/si3nloong/sqlgen/sequel/driver/postgres/pgutil" }}

{{ if eq driver "postgres" -}}
{{- /* postgres */ -}}
func InsertOne[T sequel.TableColumnValuer, Ptr interface {
	sequel.TableColumnValuer
	sequel.PtrScanner[T]
}](ctx context.Context, sqlConn sequel.DB, model Ptr) error {
	switch v := any(model).(type) {
	case sequel.SingleInserter:
		query, args := v.InsertOneStmt()
		return sqlConn.QueryRowContext(ctx, query, args...).Scan(model.Addrs()...)
	default:
		columns, values := model.Columns(), model.Values()
		stmt := strpool.AcquireString()
		defer strpool.ReleaseString(stmt)
		cols := strings.Join(columns, ",")
		stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + cols + ") VALUES (")
		for i := range values {
			if i > 0 {
				stmt.WriteString("," + wrapVar(i+1))
			} else {
				// argument always started from 1
				stmt.WriteString(wrapVar(i + 1))
			}
		}
		stmt.WriteString(") RETURNING " + strings.Join(TableColumns(model), ",") + ";")
		return sqlConn.QueryRowContext(ctx, stmt.String(), values...).Scan(model.Addrs()...)
	}
}	
{{ else }}
type autoIncrKeyInserter interface {
	sequel.AutoIncrKeyer
	sequel.SingleInserter
	ScanAutoIncr(int64) error
}

func InsertOne[T sequel.TableColumnValuer, Ptr interface {
	sequel.TableColumnValuer
	sequel.PtrScanner[T]
}](ctx context.Context, sqlConn sequel.DB, model Ptr) (sql.Result, error) {
	anyv := any(model)
	switch v := anyv.(type) {
	case autoIncrKeyInserter:
		query, args := v.InsertOneStmt()
		result, err := sqlConn.ExecContext(ctx, query, args...)
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
		return sqlConn.ExecContext(ctx, query, args...)
	default:
		columns, values := model.Columns(), model.Values()
		s := strings.Repeat("{{ quoteVar 1 }},", len(columns))
		return sqlConn.ExecContext(ctx, "INSERT INTO "+DbTable(model)+" ("+strings.Join(columns, ",")+") VALUES ("+s[:len(s)-1]+");", values...)
	}
}
{{ end }}

{{ if eq driver "postgres" -}}
{{- /* postgres */ -}}
// Insert is a helper function to insert multiple records.
func Insert[T sequel.Inserter, Ptr sequel.PtrScanner[T]](ctx context.Context, sqlConn sequel.DB, data []T) (sql.Result, error) {
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
		noOfCols := len(columns)
		cols := strings.Join(columns, ",")
		args := make([]any, 0, noOfCols*noOfData)
		stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + cols + ") VALUES ")
		for i := range data {
			if i > 0 {
				stmt.WriteString("," + model.InsertPlaceholders(i))
			} else {
				stmt.WriteString(model.InsertPlaceholders(i))
			}
			values := data[i].Values()
			values = append(values[:idx], values[idx+1:]...)
			args = append(args, values...)
		}
		stmt.WriteString(" RETURNING " + strings.Join(TableColumns(model), ",") + ";")
		rows, err := sqlConn.QueryContext(ctx, stmt.String(), args...)
		strpool.ReleaseString(stmt) // Deallocate statement
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var i int64
		for rows.Next() {
			if err := rows.Scan(Ptr(&data[i]).Addrs()...); err != nil {
				return nil, err
			}
			i++
		}
		return sequel.NewRowsAffectedResult(i), rows.Close()
	default:
		noOfCols := len(columns)
		cols := strings.Join(columns, ",")
		args := make([]any, 0, noOfCols*noOfData)
		stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + cols + ") VALUES ")
		for i := range data {
			if i > 0 {
				stmt.WriteString("," + model.InsertPlaceholders(i))
			} else {
				stmt.WriteString(model.InsertPlaceholders(i))
			}
			args = append(args, data[i].Values()...)
		}
		stmt.WriteString(" RETURNING " + strings.Join(TableColumns(model), ",") + ";")
		rows, err := sqlConn.QueryContext(ctx, stmt.String(), args...)
		strpool.ReleaseString(stmt) // Deallocate statement
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var i int64
		for rows.Next() {
			if err := rows.Scan(Ptr(&data[i]).Addrs()...); err != nil {
				return nil, err
			}
			i++
		}
		return sequel.NewRowsAffectedResult(i), rows.Close()
	}
}
{{ else }}
// Insert is a helper function to insert multiple records.
func Insert[T sequel.TableColumnValuer](ctx context.Context, sqlConn sequel.DB, data []T) (sql.Result, error) {
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
		noOfCols := len(columns)
		cols := strings.Join(columns, ",")
		args := make([]any, 0, noOfCols*noOfData)
		stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + cols + ") VALUES ")
		placeholder := "(" + strings.Repeat(",{{ quoteVar 1 }}", noOfCols)[1:] + ")"
		for i := range data {
			if i > 0 {
				stmt.WriteString("," + placeholder)
			} else {
				stmt.WriteString(placeholder)
			}
			values := data[i].Values()
			values = append(values[:idx], values[idx+1:]...)
			args = append(args, values...)
		}
		stmt.WriteByte(';')
		result, err := sqlConn.ExecContext(ctx, stmt.String(), args...)
		strpool.ReleaseString(stmt)
		return result, err
	default:
		noOfCols := len(columns)
		cols := strings.Join(columns, ",")
		args := make([]any, 0, noOfCols*noOfData)
		stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + cols + ") VALUES ")
		placeholder := "(" + strings.Repeat(",{{ quoteVar 1 }}", noOfCols)[1:] + ")"
		for i := range data {
			if i > 0 {
				stmt.WriteString("," + placeholder)
			} else {
				stmt.WriteString(placeholder)
			}
			args = append(args, data[i].Values()...)
		}
		stmt.WriteByte(';')
		result, err := sqlConn.ExecContext(ctx, stmt.String(), args...)
		strpool.ReleaseString(stmt)
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

func UpsertOne[T sequel.KeyValuer, Ptr sequel.KeyValueScanner[T]](ctx context.Context, sqlConn sequel.DB, model Ptr, opts ...UpsertOption) error {
	var opt upsertOpts
	for i := range opts {
		opts[i](&opt)
	}

	columns := model.Columns()
	noOfCols := len(columns)
	values := model.Values()

	stmt := strpool.AcquireString()
	switch v := any(model).(type) {
	case sequel.AutoIncrKeyer:
		pkName, idx, _ := v.PK()
		opt.omitFields = append(opt.omitFields, pkName)
		columns = append(columns[:idx], columns[idx+1:]...)
		noOfCols = len(columns)
		values = append(values[:idx], values[idx+1:]...)
		// Don't include auto increment primary key on INSERT
		stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES (")
	case sequel.PrimaryKeyer:
		pkName, _, _ := v.PK()
		opt.omitFields = append(opt.omitFields, pkName)
		stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES (")
	case sequel.CompositeKeyer:
		keyNames, _, _ := v.CompositeKey()
		opt.omitFields = append(opt.omitFields, keyNames...)
		stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES (")
	default:
		panic("unreachable")
	}
	for i := 1; i <= noOfCols; i++ {
		if i > 1 {
			stmt.WriteString("," + wrapVar(i))
		} else {
			stmt.WriteString(wrapVar(i))
		}
	}
	if len(opt.duplicateKeys) > 0 {
		stmt.WriteString(") ON CONFLICT(" + strings.Join(opt.duplicateKeys, ",") + ")")
	} else {
		switch vi := any(model).(type) {
		case duplicateKeyer:
			keys := vi.DuplicateKeys()
			opt.omitFields = append(opt.omitFields, keys...)
			stmt.WriteString(") ON CONFLICT(" + strings.Join(keys, ",") + ")")
		case sequel.PrimaryKeyer:
			pkName, _, _ := vi.PK()
			stmt.WriteString(") ON CONFLICT(" + pkName + ")")
		case sequel.CompositeKeyer:
			names, _, _ := vi.CompositeKey()
			stmt.WriteString(") ON CONFLICT(" + strings.Join(names, ",") + ")")
		default:
			panic("unreachable")
		}
	}
	if opt.doNothing {
		stmt.WriteString(" DO NOTHING")
	} else {
		stmt.WriteString(" DO UPDATE SET ")
		omitDict := make(map[string]struct{})
		noOfCols = len(columns)
		if omitter, ok := any(model).(setOnUpsertOmitter); ok {
			opt.omitFields = append(opt.omitFields, omitter.OnUpsertOmitColumns()...)
		}
		for i := range opt.omitFields {
			omitDict[opt.omitFields[i]] = struct{}{}
		}
		first := true
		for i := 0; i < noOfCols; i++ {
			if _, ok := omitDict[columns[i]]; ok {
				continue
			}
			if first {
				stmt.WriteString(columns[i] + " = EXCLUDED." + columns[i])
			} else {
				stmt.WriteString(","+ columns[i] + " = EXCLUDED." + columns[i])
			}
			first = false
		}
		clear(omitDict)
	}
	stmt.WriteString(" RETURNING " + strings.Join(TableColumns(model), ",") + ";")
	row := sqlConn.QueryRowContext(ctx, stmt.String(), values...)
	strpool.ReleaseString(stmt)
	return row.Scan(model.Addrs()...)
}

// Upsert is a helper function to upsert multiple records.
func Upsert[T interface {
	sequel.Keyer
	sequel.Inserter
}, Ptr sequel.PtrScanner[T]](ctx context.Context, sqlConn sequel.DB, data []T, opts ...UpsertOption) (sql.Result, error) {
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
		stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES ")
		for i := 0; i < noOfData; i++ {
			if i > 0 {
				stmt.WriteString("," + model.InsertPlaceholders(i))
			} else {
				stmt.WriteString(model.InsertPlaceholders(i))
			}
			values := data[i].Values()
			values = append(values[:idx], values[idx+1:]...)
			args = append(args, values...)
		}
	case sequel.PrimaryKeyer:
		pkName, _, _ := v.PK()
		opt.omitFields = append(opt.omitFields, pkName)
		stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES ")
		for i := 0; i < noOfData; i++ {
			if i > 0 {
				stmt.WriteString("," + model.InsertPlaceholders(i))
			} else {
				stmt.WriteString(model.InsertPlaceholders(i))
			}
			args = append(args, data[i].Values()...)
		}
	case sequel.CompositeKeyer:
		keyNames, _, _ := v.CompositeKey()
		opt.omitFields = append(opt.omitFields, keyNames...)
		stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES ")
		for i := 0; i < noOfData; i++ {
			if i > 0 {
				stmt.WriteString("," + model.InsertPlaceholders(i))
			} else {
				stmt.WriteString(model.InsertPlaceholders(i))
			}
			args = append(args, data[i].Values()...)
		}
	default:
		panic("unreachable")
	}
	if len(opt.duplicateKeys) > 0 {
		stmt.WriteString(" ON CONFLICT(" + strings.Join(opt.duplicateKeys, ",") + ")")
	} else {
		switch vi := any(model).(type) {
		case duplicateKeyer:
			keys := vi.DuplicateKeys()
			opt.omitFields = append(opt.omitFields, keys...)
			stmt.WriteString(" ON CONFLICT(" + strings.Join(keys, ",") + ")")
		case sequel.PrimaryKeyer:
			pkName, _, _ := vi.PK()
			stmt.WriteString(" ON CONFLICT(" + pkName + ")")
		case sequel.CompositeKeyer:
			names, _, _ := vi.CompositeKey()
			// // Exclude primary key, don't update it
			// for i := len(idxs) - 1; i >= 0; i-- {
			// 	columns = append(columns[:idxs[i]], columns[idxs[i]+1:]...)
			// }
			stmt.WriteString(" ON CONFLICT(" + strings.Join(names, ",") + ")")
		default:
			panic("unreachable")
		}
	}
	if opt.doNothing {
		stmt.WriteString(" DO NOTHING")
	} else {
		stmt.WriteString(" DO UPDATE SET ")
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
			if first {
				stmt.WriteString(columns[i] + " = EXCLUDED." + columns[i])
			} else {
				stmt.WriteString(","+ columns[i] + " = EXCLUDED." + columns[i])
			}
			first = false
		}
		clear(omitDict)
	}
	stmt.WriteString(" RETURNING " + strings.Join(TableColumns(model), ",") + ";")
	rows, err := sqlConn.QueryContext(ctx, stmt.String(), args...)
	strpool.ReleaseString(stmt) // Deallocate statement
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var i int64
	for rows.Next() {
		if err := rows.Scan(Ptr(&data[i]).Addrs()...); err != nil {
			return nil, err
		}
		i++
	}
	return sequel.NewRowsAffectedResult(i), rows.Close()
}
{{ else }}
func UpsertOne[T sequel.KeyValuer, Ptr sequel.KeyValueScanner[T]](ctx context.Context, sqlConn sequel.DB, model Ptr, override bool, omittedFields ...string) (sql.Result, error) {
	stmt := strpool.AcquireString()
	defer strpool.ReleaseString(stmt)
	switch v := any(model).(type) {
	case sequel.PrimaryKeyer:
		pkName, idx, _ := v.PK()
		columns := model.Columns()
		if !override {
			stmt.WriteString("INSERT IGNORE INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES (" + strings.Repeat(",?", len(columns))[1:] + ");")
		} else {
			omitDict := map[string]struct{}{pkName: {}}
			for i := range omittedFields {
				omitDict[omittedFields[i]] = struct{}{}
			}
			noOfCols := len(columns)
			stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES (" + strings.Repeat(",?", noOfCols)[1:] + ") ON DUPLICATE KEY UPDATE ")
			first := true
			for i := range columns {
				if _, ok := omitDict[columns[i]]; ok {
					continue
				}
				if first {
					stmt.WriteString(columns[i] + " =VALUES(" + columns[i] + ")")
				} else {
					stmt.WriteString(","+ columns[i] + " =VALUES(" + columns[i] + ")")
				}
				first = false
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
				return nil, fmt.Errorf(`sqlgen: invalid auto increment data type %T`, v)
			}
			return result, nil
		}
		return sqlConn.ExecContext(ctx, stmt.String(), model.Values()...)
	case sequel.CompositeKeyer:
		names, idxs, _ := v.CompositeKey()
		columns := model.Columns()
		if !override {
			stmt.WriteString("INSERT IGNORE INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES (" + strings.Repeat(",?", len(columns))[1:] + ");")
		} else {
			dict := make(map[string]struct{})
			for i := range append(names, omittedFields...) {
				dict[omittedFields[i]] = struct{}{}
			}
			noOfCols := len(columns)
			stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES (" + strings.Repeat(",?", noOfCols)[1:] + ") ON DUPLICATE KEY UPDATE ")
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

// Upsert is a helper function to upsert multiple records.
func Upsert[T interface {
	sequel.Keyer
	sequel.Inserter
}, Ptr sequel.PtrScanner[T]](ctx context.Context, sqlConn sequel.DB, data []T, override bool, omittedFields ...string) (sql.Result, error) {
	noOfData := len(data)
	if noOfData == 0 {
		return new(sequel.EmptyResult), nil
	}

	model := data[0]
	columns := model.Columns()
	noOfCols := len(columns)
	args := make([]any, 0, noOfCols*noOfData)

	stmt := strpool.AcquireString()
	switch v := any(model).(type) {
	case sequel.AutoIncrKeyer:
		pkName, idx, _ := v.PK()
		omittedFields = append(omittedFields, pkName)
		// Don't include auto increment primary key on INSERT
		columns = append(columns[:idx], columns[idx+1:]...)
		if override {
			stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES ")
		} else {
			stmt.WriteString("INSERT IGNORE INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES ")
		}
		placeholder := ",(" + strings.Repeat(",?", noOfCols)[1:] + ")"
		for i := 0; i < noOfData; i++ {
			if i > 0 {
				stmt.WriteString(placeholder)
			} else {
				stmt.WriteString(placeholder[1:])
			}
			values := data[i].Values()
			values = append(values[:idx], values[idx+1:]...)
			args = append(args, values...)
		}
	case sequel.PrimaryKeyer:
		pkName, _, _ := v.PK()
		omittedFields = append(omittedFields, pkName)
		if override {
			stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES ")
		} else {
			stmt.WriteString("INSERT IGNORE INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES ")
		}
		for i := 0; i < noOfData; i++ {
			if i > 0 {
				stmt.WriteString("," + model.InsertPlaceholders(i))
			} else {
				stmt.WriteString(model.InsertPlaceholders(i))
			}
			args = append(args, data[i].Values()...)
		}
	case sequel.CompositeKeyer:
		if override {
			stmt.WriteString("INSERT INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES ")
		} else {
			stmt.WriteString("INSERT IGNORE INTO " + DbTable(model) + " (" + strings.Join(columns, ",") + ") VALUES ")
		}
		for i := 0; i < noOfData; i++ {
			if i > 0 {
				stmt.WriteString("," + model.InsertPlaceholders(i))
			} else {
				stmt.WriteString(model.InsertPlaceholders(i))
			}
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
				stmt.WriteString(columns[i] + "=VALUES(" + columns[i] + ")")
			} else {
				stmt.WriteString(","+ columns[i] + "=VALUES(" + columns[i] + ")")
			}
			first = false
		}
		clear(omitDict)
	}
	stmt.WriteByte(';')
	result, err := sqlConn.ExecContext(ctx, stmt.String(), args...)
	strpool.ReleaseString(stmt)
	return result, err
}
{{ end }}

// FindByPK is to find single record using primary key.
func FindByPK[T sequel.KeyValuer, Ptr sequel.KeyValueScanner[T]](ctx context.Context, sqlConn sequel.DB, model Ptr) error {
	switch v := any(model).(type) {
	case sequel.KeyFinder:
		query, args := v.FindOneByPKStmt()
		return sqlConn.QueryRowContext(ctx, query, args...).Scan(model.Addrs()...)
	case sequel.PrimaryKeyer:
		columns := TableColumns(model)
		pkName, _, pk := v.PK()
		return sqlConn.QueryRowContext(ctx, "SELECT "+strings.Join(columns, ",")+" FROM "+DbTable(model)+" WHERE "+pkName+" = {{ quoteVar 1 }} LIMIT 1;", pk).Scan(model.Addrs()...)
	case sequel.CompositeKeyer:
		columns := TableColumns(model)
		names, _, keys := v.CompositeKey()
		{{ if isStaticVar -}}
		return sqlConn.QueryRowContext(ctx, "SELECT "+strings.Join(columns, ",")+" FROM "+DbTable(model)+" WHERE "+strings.Join(names, " = {{ quoteVar 1 }} AND ")+" = {{ quoteVar 1 }} LIMIT 1;", keys...).Scan(model.Addrs()...)
		{{ else -}}
		stmt := strpool.AcquireString()
		stmt.WriteString("SELECT " + strings.Join(columns, ",") + " FROM " + DbTable(model) + " WHERE ("+ strings.Join(names, ",") +") = (")
		noOfKey := len(names)
		for i := 1; i <= noOfKey; i++ {
			if i > 1 {
				stmt.WriteString(","+ wrapVar(i))
			} else {
				stmt.WriteString(wrapVar(i))
			}
		}
		stmt.WriteString(") LIMIT 1;")
		row := sqlConn.QueryRowContext(ctx, stmt.String(), keys...)
		strpool.ReleaseString(stmt)
		return row.Scan(model.Addrs()...)
		{{ end -}}
	default:
		panic("unreachable")
	}
}

// UpdateByPK is to update single record using primary key.
func UpdateByPK[T sequel.KeyValuer](ctx context.Context, sqlConn sequel.DB, model T) (sql.Result, error) {
	switch v := any(model).(type) {
	case sequel.KeyUpdater:
		query, args := v.UpdateOneByPKStmt()
		return sqlConn.ExecContext(ctx, query, args...)
	case sequel.PrimaryKeyer:
		pkName, pkIdx, pk := v.PK()
		columns := model.Columns()
		columns = append(columns[:pkIdx], columns[pkIdx+1:]...)
		values := model.Values()
		values = append(values[:pkIdx], append(values[pkIdx+1:], pk)...)
		{{ if isStaticVar -}}
		return sqlConn.ExecContext(ctx, "UPDATE "+DbTable(model)+" SET "+strings.Join(columns, " = {{ quoteVar 1 }},")+" = {{ quoteVar 1 }} WHERE "+pkName+" = {{ quoteVar 1 }};", append(values, pk)...)
		{{ else -}}
		stmt := strpool.AcquireString()
		stmt.WriteString("UPDATE "+DbTable(model)+" SET ")
		for idx := range columns {
			if idx > 0 {
				stmt.WriteByte(',')
			}
			stmt.WriteString(columns[idx] +" = "+ wrapVar(idx + 1))
		}
		stmt.WriteString(" WHERE "+ pkName +" = "+ wrapVar(len(columns) + 2)+ ";")
		result, err := sqlConn.ExecContext(ctx, stmt.String(), append(values, pk)...)
		strpool.ReleaseString(stmt)
		return result, err
		{{ end -}}
	default:
		panic("unreachable")
	}
}

// DeleteByPK is to update single record using primary key.
func DeleteByPK[T sequel.KeyValuer](ctx context.Context, sqlConn sequel.DB, model T) (sql.Result, error) {
	switch v := any(model).(type) {
	case sequel.KeyDeleter:
		query, args := v.DeleteOneByPKStmt()
		return sqlConn.ExecContext(ctx, query, args...)
	case sequel.PrimaryKeyer:
		pkName, _, pk := v.PK()
		return sqlConn.ExecContext(ctx, "DELETE FROM "+DbTable(model)+" WHERE "+pkName+" = {{ quoteVar 1 }};", pk)
	case sequel.CompositeKeyer:
		names, _, keys := v.CompositeKey()
		{{ if isStaticVar -}}
		return sqlConn.ExecContext(ctx, "DELETE FROM "+DbTable(model)+" WHERE "+strings.Join(names, " = ? AND ")+" = ?;", keys...)
		{{ else -}}
		stmt := strpool.AcquireString()
		stmt.WriteString("DELETE FROM "+DbTable(model)+" WHERE ")
		noOfKey := len(names)
		for i := 0; i < noOfKey; i++ {
			if i == 1 {
				stmt.WriteString(names[i]+" = "+ wrapVar(i+1))
			} else {
				stmt.WriteString(" AND "+ names[i]+" = "+ wrapVar(i+1))
			}
		}
		stmt.WriteByte(';')
		result, err := sqlConn.ExecContext(ctx, stmt.String(), keys...)
		strpool.ReleaseString(stmt)
		return result, err
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

type Result[T any] struct {
	data []T
	err  error
}

func (r *Result[T]) All(yield func(T) bool) {
	for _, s := range r.data {
		if !yield(s) {
			return
		}
	}
}

func (r *Result[T]) AllWithIndex(yield func(int, T) bool) {
	for i, s := range r.data {
		if !yield(i, s) {
			return
		}
	}
}

func (r *Result[T]) Err() error {
	return r.err
}

type Pager[T sequel.KeyValuer, Ptr sequel.KeyValueScanner[T]] struct {
	stmt *PaginateStmt
}

func (r *Pager[T, Ptr]) Prev(ctx context.Context, sqlConn sequel.DB, cursor ...T) iter.Seq[*Result[T]] {
	return func(yield func(*Result[T]) bool) {
		var (
			v         T
			hasCursor bool
			maxLimit  = r.stmt.Limit
		)
		if len(cursor) > 0 {
			v = cursor[0]
			hasCursor = true
		}

		for {
			if hasCursor {
				if err := FindByPK(ctx, sqlConn, Ptr(&v)); err != nil {
					if !yield(&Result[T]{err: err}) {
						return
					}
					return
				}
			}

			blr := AcquireStmt()
			defer ReleaseStmt(blr)
			switch vi := any(v).(type) {
			case sequel.SQLColumner:
				blr.WriteString("SELECT " + strings.Join(vi.SQLColumns(), ",") + " FROM " + DbTable(v) + " WHERE ")
			default:
				blr.WriteString("SELECT " + strings.Join(v.Columns(), ",") + " FROM " + DbTable(v) + " WHERE ")
			}
			if r.stmt.Where != nil {
				r.stmt.Where(blr)
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
						blr.WriteString(" AND ")
					}
					if orderBy.Asc() {
						// If ascending
						blr.WriteString(colName + " <= " + blr.Var(val))
					} else {
						// If descending
						blr.WriteString(colName + " >= " + blr.Var(val))
					}
				}
				blr.WriteString(" AND (")
				for i, orderBy := range r.stmt.OrderBy {
					colName := orderBy.ColumnName()
					val := values[colDict[colName]]
					if i > 0 {
						blr.WriteString(" OR ")
					}
					if orderBy.Asc() {
						// If ascending
						blr.WriteString(colName + " < " + blr.Var(val))
					} else {
						// If descending
						blr.WriteString(colName + " > " + blr.Var(val))
					}
				}
				clear(colDict)
				blr.WriteString(" OR ")
			} else {
				blr.WriteByte('(')
			}

			// Check the primary key and compare value
			switch vi := any(v).(type) {
			case sequel.CompositeKeyer:
				pkNames, _, vals := vi.CompositeKey()
				blr.WriteString("(" + strings.Join(pkNames, ",") + ") <= " + blr.Vars(vals))
			case sequel.PrimaryKeyer:
				pkName, _, val := vi.PK()
				blr.WriteString(pkName + " <= " + blr.Var(val))
			default:
				panic("unreachable")
			}

			blr.WriteString(") ORDER BY ")
			if len(r.stmt.OrderBy) > 0 {
				for i := range r.stmt.OrderBy {
					if r.stmt.OrderBy[i].Asc() {
						blr.WriteString(r.stmt.OrderBy[i].ColumnName() + " ASC,")
					} else {
						blr.WriteString(r.stmt.OrderBy[i].ColumnName() + " DESC,")
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
					for i := range pkNames {
						if i > 0 {
							blr.WriteString("," + pkNames[i] + suffix)
						} else {
							blr.WriteString(pkNames[i] + suffix)
						}
					}
				case sequel.PrimaryKeyer:
					pkName, _, _ := vi.PK()
					blr.WriteString(pkName + suffix)
				default:
					panic("unreachable")
				}
			} else {
				// If there is no order by clause,
				// ascending is the default order
				switch vi := any(v).(type) {
				case sequel.CompositeKeyer:
					pkNames, _, _ := vi.CompositeKey()
					for i := range pkNames {
						if i > 0 {
							blr.WriteString("," + pkNames[i] + " DESC")
						} else {
							blr.WriteString(pkNames[i] + " DESC")
						}
					}
				case sequel.PrimaryKeyer:
					pkName, _, _ := vi.PK()
					blr.WriteString(pkName + " DESC")
				default:
					panic("unreachable")
				}
			}
			// Add one to limit to find next cursor
			blr.WriteString(" LIMIT " + strconv.FormatUint(uint64(r.stmt.Limit+1), 10) + ";")

			rows, err := sqlConn.QueryContext(ctx, blr.Query(), blr.Args()...)
			if err != nil {
				if !yield(&Result[T]{err: err}) {
					return
				}
				return
			}
			defer rows.Close()

			data := make([]T, 0, r.stmt.Limit+1)
			for rows.Next() {
				var v T
				if err := rows.Scan(Ptr(&v).Addrs()...); err != nil {
					if !yield(&Result[T]{err: err}) {
						return
					}
					return
				}
				data = append(data, v)
			}

			noOfRecord := len(data)
			if uint16(noOfRecord) < maxLimit {
				if !yield(&Result[T]{data: data}) {
					return
				}
				return
			}

			if !yield(&Result[T]{data: data[:noOfRecord-1]}) {
				return
			}

			// Set next cursor
			v = data[noOfRecord-1]
			hasCursor = true
		}
	}
}

func (r *Pager[T, Ptr]) Next(ctx context.Context, sqlConn sequel.DB, cursor ...T) iter.Seq[*Result[T]] {
	return func(yield func(*Result[T]) bool) {
		var (
			v         T
			hasCursor bool
			maxLimit  = r.stmt.Limit
		)
		if len(cursor) > 0 {
			v = cursor[0]
			hasCursor = true
		}

		for {
			if hasCursor {
				if err := FindByPK(ctx, sqlConn, Ptr(&v)); err != nil {
					if !yield(&Result[T]{err: err}) {
						return
					}
					return
				}
			}

			blr := AcquireStmt()
			defer ReleaseStmt(blr)
			switch vi := any(v).(type) {
			case sequel.SQLColumner:
				blr.WriteString("SELECT " + strings.Join(vi.SQLColumns(), ",") + " FROM " + DbTable(v) + " WHERE ")
			default:
				blr.WriteString("SELECT " + strings.Join(v.Columns(), ",") + " FROM " + DbTable(v) + " WHERE ")
			}
			if r.stmt.Where != nil {
				r.stmt.Where(blr)
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
						blr.WriteString(" AND ")
					}
					if orderBy.Asc() {
						// If ascending
						blr.WriteString(colName + " >= " + blr.Var(val))
					} else {
						// If descending
						blr.WriteString(colName + " <= " + blr.Var(val))
					}
				}
				blr.WriteString(" AND (")
				for i, orderBy := range r.stmt.OrderBy {
					colName := orderBy.ColumnName()
					val := values[colDict[colName]]
					if i > 0 {
						blr.WriteString(" OR ")
					}
					if orderBy.Asc() {
						// If ascending
						blr.WriteString(colName + " > " + blr.Var(val))
					} else {
						// If descending
						blr.WriteString(colName + " < " + blr.Var(val))
					}
				}
				clear(colDict)
				blr.WriteString(" OR ")
			} else {
				blr.WriteByte('(')
			}

			// Check the primary key and compare value
			switch vi := any(v).(type) {
			case sequel.CompositeKeyer:
				pkNames, _, vals := vi.CompositeKey()
				blr.WriteString("(" + strings.Join(pkNames, ",") + ") >= " + blr.Vars(vals))
			case sequel.PrimaryKeyer:
				pkName, _, val := vi.PK()
				blr.WriteString(pkName + " >= " + blr.Var(val))
			default:
				panic("unreachable")
			}

			blr.WriteString(") ORDER BY ")
			if len(r.stmt.OrderBy) > 0 {
				for i := range r.stmt.OrderBy {
					if r.stmt.OrderBy[i].Asc() {
						blr.WriteString(r.stmt.OrderBy[i].ColumnName() + " ASC,")
					} else {
						blr.WriteString(r.stmt.OrderBy[i].ColumnName() + " DESC,")
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
					for i := range pkNames {
						if i > 0 {
							blr.WriteString("," + pkNames[i] + suffix)
						} else {
							blr.WriteString(pkNames[i] + suffix)
						}
					}
				case sequel.PrimaryKeyer:
					pkName, _, _ := vi.PK()
					blr.WriteString(pkName + suffix)
				default:
					panic("unreachable")
				}
			} else {
				// If there is no order by clause,
				// ascending is the default order
				switch vi := any(v).(type) {
				case sequel.CompositeKeyer:
					pkNames, _, _ := vi.CompositeKey()
					for i := range pkNames {
						if i > 0 {
							blr.WriteString("," + pkNames[i] + " ASC")
						} else {
							blr.WriteString(pkNames[i] + " ASC")
						}
					}
				case sequel.PrimaryKeyer:
					pkName, _, _ := vi.PK()
					blr.WriteString(pkName + " ASC")
				default:
					panic("unreachable")
				}
			}
			// Add one to limit to find next cursor
			blr.WriteString(" LIMIT " + strconv.FormatUint(uint64(r.stmt.Limit+1), 10) + ";")

			rows, err := sqlConn.QueryContext(ctx, blr.Query(), blr.Args()...)
			if err != nil {
				if !yield(&Result[T]{err: err}) {
					return
				}
				return
			}
			defer rows.Close()

			data := make([]T, 0, r.stmt.Limit+1)
			for rows.Next() {
				var v T
				if err := rows.Scan(Ptr(&v).Addrs()...); err != nil {
					if !yield(&Result[T]{err: err}) {
						return
					}
					return
				}
				data = append(data, v)
			}

			noOfRecord := len(data)
			if uint16(noOfRecord) < maxLimit {
				if !yield(&Result[T]{data: data}) {
					return
				}
				return
			}

			if !yield(&Result[T]{data: data[:noOfRecord-1]}) {
				return
			}

			// Set next cursor
			v = data[noOfRecord-1]
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

func QueryStmt[T any, Ptr sequel.PtrScanner[T], Stmt interface{ 
	SelectStmt | *SqlStmt
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
		blr.WriteString("SELECT ")
		if len(vi.Select) > 0 {
			blr.WriteString(strings.Join(vi.Select, ","))
		} else {
			switch vj := any(v).(type) {
			case sequel.Columner:
				blr.WriteString(strings.Join(TableColumns(vj), ","))
			default:
				blr.WriteByte('*')
			}
		}
		if vi.FromTable != "" {
			blr.WriteString(" FROM " + dbName(v) + vi.FromTable)
		} else {
			switch vj := any(v).(type) {
			case sequel.Tabler:
				blr.WriteString(" FROM " + DbTable(vj))
			default:
				ReleaseStmt(blr)
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
				if vi.OrderBy[i].Asc() {
					blr.WriteString(vi.OrderBy[i].ColumnName() +" ASC")
				} else {
					blr.WriteString(vi.OrderBy[i].ColumnName() +" DESC")
				}
			}
		}
		if vi.Limit > 0 {
			blr.WriteString(" LIMIT " + strconv.FormatUint(uint64(vi.Limit), 10))
		}
		if vi.Offset > 0 {
			blr.WriteString(" OFFSET " + strconv.FormatUint(vi.Offset, 10))
		}
		blr.WriteByte(';')
		rows, err = sqlConn.QueryContext(ctx, blr.Query(), blr.Args()...)
		ReleaseStmt(blr)

	case *SqlStmt:
		rows, err = sqlConn.QueryContext(ctx, vi.Query(), vi.Args()...)
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

type SelectOneStmt struct {
	Select    []string
	FromTable string
	Where     sequel.WhereClause
	OrderBy   []sequel.OrderByClause
	GroupBy   []string
}

func QueryOneStmt[T any, Ptr sequel.PtrScanner[T], Stmt interface {
	SelectOneStmt | *SqlStmt
}](ctx context.Context, sqlConn sequel.DB, stmt Stmt) (Ptr, error) {
	var v T
	switch vi := any(stmt).(type) {
	case SelectOneStmt:
		var blr = AcquireStmt()
		blr.WriteString("SELECT ")
		if len(vi.Select) > 0 {
			blr.WriteString(strings.Join(vi.Select, ","))
		} else {
			switch vj := any(v).(type) {
			case sequel.Columner:
				blr.WriteString(strings.Join(TableColumns(vj), ","))
			default:
				blr.WriteByte('*')
			}
		}
		if vi.FromTable != "" {
			blr.WriteString(" FROM " + dbName(v) + vi.FromTable)
		} else {
			switch vj := any(v).(type) {
			case sequel.Tabler:
				blr.WriteString(" FROM " + DbTable(vj))
			default:
				ReleaseStmt(blr)
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
				if vi.OrderBy[i].Asc() {
					blr.WriteString(vi.OrderBy[i].ColumnName() +" ASC")
				} else {
					blr.WriteString(vi.OrderBy[i].ColumnName() +" DESC")
				}
			}
		}
		blr.WriteString(" LIMIT 1;")
		row := sqlConn.QueryRowContext(ctx, blr.Query(), blr.Args()...)
		ReleaseStmt(blr)
		if err := row.Scan(Ptr(&v).Addrs()...); err != nil {
			return nil, err
		}
		return &v, nil

	case *SqlStmt:
		if err := sqlConn.QueryRowContext(ctx, vi.Query(), vi.Args()...).Scan(Ptr(&v).Addrs()...); err != nil {
			return nil, err
		}
		return &v, nil

	default:
		panic("unreachable")
	}
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

func ExecStmt[T any, Stmt interface {
	UpdateStmt | DeleteStmt
}](ctx context.Context, sqlConn sequel.DB, stmt Stmt) (sql.Result, error) {
	blr := AcquireStmt()
	defer ReleaseStmt(blr)

	var v T
	switch vi := any(stmt).(type) {
	case UpdateStmt:
		if vt, ok := any(v).(sequel.Tabler); ok {
			blr.WriteString("UPDATE " + DbTable(vt))
		} else {
			blr.WriteString("UPDATE " + dbName(v) + vi.Table)
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
				if vi.OrderBy[i].Asc() {
					blr.WriteString(vi.OrderBy[i].ColumnName() +" ASC")
				} else {
					blr.WriteString(vi.OrderBy[i].ColumnName() +" DESC")
				}
			}
		}
		{{ if ne driver "postgres" -}}
		{{- /* postgres */ -}}
		if vi.Limit > 0 {
			blr.WriteString(" LIMIT " + strconv.FormatUint(uint64(vi.Limit), 10))
		}
		{{ end -}}

	case DeleteStmt:
		if vt, ok := any(v).(sequel.Tabler); ok {
			blr.WriteString("DELETE FROM " + DbTable(vt))
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
				if vi.OrderBy[i].Asc() {
					blr.WriteString(vi.OrderBy[i].ColumnName() +" ASC")
				} else {
					blr.WriteString(vi.OrderBy[i].ColumnName() +" DESC")
				}
			}
		}
		{{ if ne driver "postgres" -}}
		{{- /* postgres */ -}}
		if vi.Limit > 0 {
			blr.WriteString(" LIMIT " + strconv.FormatUint(uint64(vi.Limit), 10))
		}
		{{ end -}}
	}
	blr.WriteByte(';')
	return sqlConn.ExecContext(ctx, blr.Query(), blr.Args()...)
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

func (s *SqlStmt) Vars(values []any) string {
	noOfLen := len(values)
	s.args = append(s.args, values...)
	{{ if isStaticVar -}}
	return "(" + strings.Repeat(",{{ varRune }}", noOfLen)[1:] + ")"
	{{ else -}}
	buf := new(strings.Builder)
	buf.WriteByte('(')
	i := s.pos
	s.pos += noOfLen
	for ; i < s.pos; i++ {
		if i < s.pos-1 {
			buf.WriteString(wrapVar(i+1) + ",")
		} else {
			buf.WriteString(wrapVar(i + 1))
		}
	}
	buf.WriteByte(')')
	return buf.String()
	{{ end -}}
}

func (s *SqlStmt) Write(p []byte) (int, error) {
	return s.blr.Write(p)
}

func (s *SqlStmt) WriteString(v string) (int, error) {
	return s.blr.WriteString(v)
}

func (s *SqlStmt) WriteByte(c byte) error {
	return s.blr.WriteByte(c)
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
		args = make([]any, len(s.args))
		idx  int
		i    = 1
	)

	copy(args, s.args)

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
		panic("unreachable")
	}
}