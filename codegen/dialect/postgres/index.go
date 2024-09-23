package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type index struct {
	Name    string
	IsPK    bool
	Columns pq.StringArray
}

func (s *postgresDriver) tableIndexes(ctx context.Context, sqlConn *sql.DB, tableName string) ([]index, error) {
	// c.relnamespace::regnamespace as schema_name,
	// c.relname as table_name,
	rows, err := sqlConn.QueryContext(ctx, `SELECT
    i.relname as index_name,
	ix.indisprimary as is_pk,
    ARRAY_AGG(a.attname) AS column_names
FROM
    pg_class t,
    pg_class i,
    pg_index ix,
    pg_attribute a
WHERE
    t.oid = ix.indrelid
    and i.oid = ix.indexrelid
    and a.attrelid = t.oid
    and a.attnum = ANY(ix.indkey)
    and t.relname = $1
GROUP BY index_name, is_pk;`, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	idxs := make([]index, 0)
	for rows.Next() {
		var idx index
		if err := rows.Scan(&idx.Name, &idx.IsPK, &idx.Columns); err != nil {
			return nil, err
		}
		idxs = append(idxs, idx)
	}
	return idxs, rows.Err()
}
