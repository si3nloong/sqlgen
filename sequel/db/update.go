package db

import (
	"context"
	"database/sql"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/strpool"
)

// UpdateOne is to update single record using primary key.
func UpdateOne[T sequel.KeyValuer[T]](ctx context.Context, db sequel.DB, v T) (sql.Result, error) {
	var (
		pkName, _, pk   = v.PK()
		dialect         = sequel.DefaultDialect()
		columns, values = v.Columns(), v.Values()
		noOfCols        = len(columns)
		i               = 1
		stmt            = strpool.AcquireString()
	)
	defer strpool.ReleaseString(stmt)
	stmt.WriteString("UPDATE " + dialect.Wrap(v.TableName()) + " SET ")
	for ; i <= noOfCols; i++ {
		if i > 1 {
			stmt.WriteByte(',')
		}
		stmt.WriteString(dialect.Wrap(columns[i]) + " = " + dialect.Var(i))
	}
	stmt.WriteString(" WHERE " + dialect.Wrap(pkName) + " = " + dialect.Var(i) + ";")
	return db.ExecContext(ctx, stmt.String(), append(values, pk)...)
}
