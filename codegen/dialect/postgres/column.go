package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/si3nloong/sqlgen/internal/sqltype"
)

type column struct {
	Name               string
	DataType           string
	Default            sql.RawBytes
	IsNullable         sqltype.Bool
	CharacterMaxLength sql.NullInt64
	NumericPrecision   sql.NullInt64
	DatetimePrecision  sql.NullInt64
}

func (c column) ColumnType() string {
	switch c.DataType {
	case "varchar":
		if c.CharacterMaxLength.Valid {
			return fmt.Sprintf("%s(%d)", c.DataType, c.CharacterMaxLength.Int64)
		}
	case "timestamptz":
		if c.DatetimePrecision.Valid {
			return fmt.Sprintf("%s(%d)", c.DataType, c.DatetimePrecision.Int64)
		}
	}
	return c.DataType
}

func (s *postgresDriver) tableColumns(ctx context.Context, sqlConn *sql.DB, dbName, tableName string) ([]column, error) {
	rows, err := sqlConn.QueryContext(ctx, `SELECT 
	column_name,
	column_default,
	is_nullable,
	udt_name,
	character_maximum_length,
	numeric_precision,
	datetime_precision
FROM 
	information_schema.columns
WHERE
	table_schema = $1 AND 
	table_name = $2
ORDER BY
	ordinal_position;`, dbName, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols := make([]column, 0)
	for rows.Next() {
		var col column
		if err := rows.Scan(
			&col.Name,
			&col.Default,
			&col.IsNullable,
			&col.DataType,
			&col.CharacterMaxLength,
			&col.NumericPrecision,
			&col.DatetimePrecision,
		); err != nil {
			return nil, err
		}
		cols = append(cols, col)
	}
	return cols, rows.Err()
}
