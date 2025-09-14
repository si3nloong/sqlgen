package mysql

import (
	"context"
	"database/sql"

	"github.com/si3nloong/sqlgen/cmd/internal/sqltype"
)

type column struct {
	ColumnName         string
	DataType           string
	ColumnType         string
	Default            sql.NullString
	IsNullable         sqltype.Bool
	CharacterMaxLength sql.NullInt64
	NumericPrecision   sql.NullInt64
	DatetimePrecision  sql.NullInt64
	Comment            string
}

func (s *mysqlDriver) tableColumns(ctx context.Context, sqlConn *sql.DB, dbName, tableName string) ([]column, error) {
	rows, err := sqlConn.QueryContext(ctx, `SELECT 
	COLUMN_NAME,
	COLUMN_DEFAULT,
	IS_NULLABLE,
	DATA_TYPE,
	COLUMN_TYPE,
	CHARACTER_MAXIMUM_LENGTH,
	NUMERIC_PRECISION,
	DATETIME_PRECISION,
	COLUMN_COMMENT
FROM 
	information_schema.columns
WHERE
	TABLE_SCHEMA = ? AND 
	TABLE_NAME = ?
ORDER BY
	ORDINAL_POSITION;`, dbName, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols := make([]column, 0)
	for rows.Next() {
		var col column
		if err := rows.Scan(
			&col.ColumnName,
			&col.Default,
			&col.IsNullable,
			&col.DataType,
			&col.ColumnType,
			&col.CharacterMaxLength,
			&col.NumericPrecision,
			&col.DatetimePrecision,
			&col.Comment,
		); err != nil {
			return nil, err
		}
		cols = append(cols, col)
	}
	return cols, rows.Err()
}
