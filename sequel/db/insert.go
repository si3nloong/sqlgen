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
