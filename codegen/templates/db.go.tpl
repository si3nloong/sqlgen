{{- reserveImport "context" }}
{{- reserveImport "database/sql" }}
{{- reserveImport "strings" }}
{{- reserveImport "sync" }}
{{- reserveImport "github.com/si3nloong/sqlgen/sequel" }}
{{- reserveImport "github.com/si3nloong/sqlgen/sequel/strpool" }}
func InsertOne[T sequel.KeyValuer[T], Ptr interface {
	sequel.KeyValuer[T]
	sequel.Scanner[T]
}](ctx context.Context, db sequel.DB, v Ptr) (sql.Result, error) {
	columns, args := v.Columns(), v.Values()
	switch vi := any(v).(type) {
	case sequel.Keyer:
		if vi.IsAutoIncr() {
            // If it's a auto increment primary key
            // We don't need to pass the value
			_, idx, _ := vi.PK()
			columns = append(columns[:idx], columns[idx+1:]...)
			args = append(args[:idx], args[idx+1:]...)
		}
	}
	var (
		noOfCols = len(columns)
		stmt     = strpool.AcquireString()
	)
	defer strpool.ReleaseString(stmt)
	stmt.WriteString("INSERT INTO " + v.TableName() + " (")
	for i := 0; i < noOfCols; i++ {
		if i > 0 {
			stmt.WriteString("," + columns[i])
		} else {
			stmt.WriteString(columns[i])
		}
	}
	stmt.WriteString(") VALUES (")
	for i := range args {
		if i > 0 {
			stmt.WriteByte(',')
		}
		stmt.WriteByte('{{ var 1 }}')
	}
	stmt.WriteString(");")
	return db.ExecContext(ctx, stmt.String(), args...)
}

// FindByID is to find single record using primary key.
func FindByID[T sequel.KeyValuer[T], Ptr sequel.KeyValueScanner[T]](ctx context.Context, db sequel.DB, v Ptr) error {
	var (
		pkName, _, pk = v.PK()
		columns       = v.Columns()
		stmt          = strpool.AcquireString()
	)
	defer strpool.ReleaseString(stmt)
	stmt.WriteString("SELECT ")
	for i := range columns {
		if i > 0 {
			stmt.WriteByte(',')
		}
		stmt.WriteString(columns[i])
	}
	stmt.WriteString(" FROM " + v.TableName() + " WHERE " + pkName + " = {{ var 1 }} LIMIT 1;")
	return db.QueryRowContext(ctx, stmt.String(), pk).Scan(v.Addrs()...)
}

// UpdateByID is to update single record using primary key.
func UpdateByID[T sequel.KeyValuer[T]](ctx context.Context, db sequel.DB, v T) (sql.Result, error) {
	var (
		pkName, idx, pk = v.PK()
		columns, values = v.Columns(), v.Values()
		stmt            = strpool.AcquireString()
	)
    columns = append(columns[:idx], columns[idx+1:]...)
    values = append(values[:idx], values[idx+1:]...)
    var noOfCols = len(columns)
	defer strpool.ReleaseString(stmt)
	stmt.WriteString("UPDATE " + v.TableName() + " SET ")
	for i := 0; i < noOfCols; i++ {
		if i > 1 {
			stmt.WriteByte(',')
		}
		stmt.WriteString(columns[i] + " = {{ var 1 }}")
	}
	stmt.WriteString(" WHERE " + pkName + " = {{ var 1 }};")
	return db.ExecContext(ctx, stmt.String(), append(values, pk)...)
}

// DeleteByID is to update single record using primary key.
func DeleteByID[T sequel.KeyValuer[T]](ctx context.Context, db sequel.DB, v T) (sql.Result, error) {
	var (
		pkName, _, pk = v.PK()
		stmt          = strpool.AcquireString()
	)
	defer strpool.ReleaseString(stmt)
	stmt.WriteString("DELETE FROM " + v.TableName() + " WHERE " + pkName + " = {{ var 1 }};")

	return db.ExecContext(ctx, stmt.String(), pk)
}