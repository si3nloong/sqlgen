package db

import (
	"context"
	"database/sql"

	"github.com/si3nloong/sqlgen/sequel"
)

func InsertOne[T sequel.KeyValuer[T], Ptr interface {
	sequel.KeyValuer[T]
	sequel.Scanner[T]
}](ctx context.Context, db sequel.DB, v Ptr) (sql.Result, error) {
	columns, args := v.Columns(), v.Values()
	switch vi := any(v).(type) {
	case sequel.Keyer:
		if vi.IsAutoIncr() {
			_, idx, _ := vi.PK()
			columns = append(columns[:idx], columns[idx+1:]...)
			args = append(args[:idx], args[idx+1:]...)
		}
	}

	noOfCols := len(columns)
	stmt := acquireString()
	defer releaseString(stmt)
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
		// stmt.WriteString(dialect.Var(i + 1))
	}
	stmt.WriteByte(')')
	// stmt.WriteString(`) RETURNING `)
	// for i, col := range v.Columns() {
	// 	if i > 0 {
	// 		stmt.WriteString("," + dialect.Wrap(col))
	// 	} else {
	// 		stmt.WriteString(dialect.Wrap(col))
	// 	}
	// }
	// return db.QueryRowContext(ctx, stmt.String(), args...).Scan(v.Addrs()...)
	stmt.WriteByte(';')
	return db.ExecContext(ctx, stmt.String(), args...)
}

// InsertInto is a helper function to insert your records.
func InsertInto[T sequel.KeyValuer[T]](ctx context.Context, db sequel.DB, data []T) (sql.Result, error) {
	n := len(data)
	if n == 0 {
		return new(EmptyResult), nil
	}

	var (
		model   T
		columns = model.Columns()
		idx     = -1
	)
	switch vi := any(model).(type) {
	case sequel.Keyer:
		if vi.IsAutoIncr() {
			_, idx, _ = vi.PK()
			columns = append(columns[:idx], columns[idx+1:]...)
		}
	}
	var (
		noOfCols = len(columns)
		args     = make([]any, 0)
		stmt     = acquireString()
	)
	defer releaseString(stmt)
	stmt.WriteString("INSERT INTO " + model.TableName() + " (")
	for i := 0; i < noOfCols; i++ {
		if i > 0 {
			stmt.WriteString("," + columns[i])
		} else {
			stmt.WriteString(columns[i])
		}
	}
	stmt.WriteString(") VALUES ")
	var (
		pos int
	)
	for i := range data {
		if i > 0 {
			stmt.WriteString(",(")
		} else {
			stmt.WriteByte('(')
		}
		for j := 0; j < noOfCols; j++ {
			if j > 0 {
				stmt.WriteByte(',')
			}
			stmt.WriteString("?")
			pos++
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