package db

import (
	"context"

	"github.com/si3nloong/sqlgen/sequel"
)

// FindOne is to find single record using primary key.
func FindOne[T sequel.KeyValuer[T], Ptr sequel.KeyValueScanner[T]](ctx context.Context, db sequel.DB, v Ptr) error {
	pkName, _, pk := v.PK()
	columns := v.Columns()
	stmt := acquireString()
	defer releaseString(stmt)

	stmt.WriteString("SELECT ")
	for i := range columns {
		if i > 0 {
			stmt.WriteByte(',')
		}
		stmt.WriteString(columns[i])
	}
	stmt.WriteString(" FROM " + v.TableName() + " WHERE " + pkName + " = ? LIMIT 1;")
	return db.QueryRowContext(ctx, stmt.String(), pk).Scan(v.Addrs()...)
}
