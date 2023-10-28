package db

import (
	"context"
	"database/sql"

	"github.com/si3nloong/sqlgen/sequel"
)

// UpdateOne is to update single record using primary key.
func UpdateOne[T sequel.KeyValuer[T]](ctx context.Context, db sequel.DB, v T) (sql.Result, error) {
	pkName, _, pk := v.PK()
	stmt := acquireString()
	defer releaseString(stmt)
	columns, values := v.Columns(), v.Values()
	stmt.WriteString("UPDATE " + v.TableName() + " SET ")
	noOfCols := len(columns)
	for i := 0; i < noOfCols; i++ {
		if i > 0 {
			stmt.WriteByte(',')
		}
		stmt.WriteString(columns[i] + " = ?")
	}
	stmt.WriteString(" WHERE " + pkName + " = ?;")
	return db.ExecContext(ctx, stmt.String(), append(values, pk)...)
}
