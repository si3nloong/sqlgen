package postgres

import (
	"database/sql"
	"slices"
	"strings"

	"github.com/samber/lo"
	"github.com/si3nloong/sqlgen/codegen/dialect"
	"github.com/si3nloong/sqlgen/sequel/strpool"
)

const _getColumns = `SELECT "column_name" FROM "information_schema"."columns" WHERE "table_name" = $1 ORDER BY "ordinal_position";`

func (s *postgresDriver) CreateTableStmt2(sqlDSN string, w dialect.Writer, schema dialect.Schema) error {
	sqlConn, err := sql.Open("pgx", sqlDSN)
	if err != nil {
		return err
	}
	defer sqlConn.Close()

	rows, err := sqlConn.Query(_getColumns, schema.TableName())
	if err != nil {
		return err
	}
	defer rows.Close()

	columns := make([]string, 0)
	for rows.Next() {
		var column string
		if err := rows.Scan(&column); err != nil {
			return err
		}
		columns = append(columns, column)
	}

	up := strpool.AcquireString()
	defer strpool.ReleaseString(up)

	down := strpool.AcquireString()
	defer strpool.ReleaseString(down)

	up.WriteString("\nALTER TABLE " + s.QuoteIdentifier(schema.TableName()) + " (\n")
	down.WriteString("\nALTER TABLE " + s.QuoteIdentifier(schema.TableName()) + " (\n")
	cols := schema.Columns()
	for i := range cols {
		if i > 0 {
			up.WriteString(",\n")
			down.WriteString(",\n")
		}
		c := schema.ColumnGoType(i)
		if idx := lo.IndexOf(columns, cols[i]); idx >= 0 {
			columns = slices.Delete(columns, idx, idx+1)
			up.WriteString("\tALTER " + s.QuoteIdentifier(cols[i]) + " " + c.DataType())
			down.WriteString("\tALTER " + s.QuoteIdentifier(cols[i]))
		} else {
			if i == 0 {
				up.WriteString("\tADD " + s.QuoteIdentifier(cols[i]) + " " + c.DataType())
			} else {
				up.WriteString("\tADD " + s.QuoteIdentifier(cols[i]) + " " + c.DataType())
			}
			down.WriteString("\tDROP COLUMN " + s.QuoteIdentifier(cols[i]))
		}
	}
	if keys := schema.Keys(); len(keys) > 0 {
		up.WriteString(",\n\tADD PRIMARY KEY (" + s.QuoteIdentifier(strings.Join(keys, s.QuoteIdentifier(","))) + ")")
		down.WriteString(",\n\tDROP PRIMARY KEY")
	}
	if len(columns) > 0 {
		for i := range columns {
			up.WriteString(",\n\tDROP COLUMN " + s.QuoteIdentifier(columns[i]))
		}
	}
	clear(columns)
	up.WriteString(");\n")

	w.WriteString(`-- +migrate Up
	-- SQL in section 'Up' is executed when this migration is applied`)
	w.WriteString(up.String())
	w.WriteString(`-- +migrate Down
	-- SQL section 'Down' is executed when this migration is rolled back`)
	w.WriteString(down.String())
	return nil
}

func (s *postgresDriver) CreateTableStmt(w dialect.Writer, schema dialect.Schema) error {
	w.WriteString("CREATE TABLE " + s.QuoteIdentifier(schema.TableName()) + " IF NOT EXISTS (\n")
	columns := schema.Columns()
	for i := range columns {
		if i > 0 {
			w.WriteString(",\n")
		}
		column := schema.ColumnGoType(i)
		w.WriteString("\t" + s.QuoteIdentifier(column.Name()) + " " + column.DataType())
	}
	if keys := schema.Keys(); len(keys) > 0 {
		w.WriteString(",\n\tPRIMARY KEY (")
		for i := range keys {
			if i > 0 {
				w.WriteString("," + s.QuoteIdentifier(keys[i]))
			} else {
				w.WriteString(s.QuoteIdentifier(keys[i]))
			}
		}
		w.WriteString(")")
	}
	w.WriteString("\n);")
	return nil
}

func (s *postgresDriver) AlterTableStmt(w dialect.Writer, schema dialect.Schema) error {
	w.WriteString("ALTER TABLE " + s.QuoteIdentifier(schema.TableName()) + " (\n")
	columns := schema.Columns()
	for i := range columns {
		column := schema.ColumnGoType(i)
		if i > 0 {
			w.WriteString(",\n\t" + s.QuoteIdentifier(column.Name()))
		} else {
			w.WriteString("\t" + s.QuoteIdentifier(column.Name()))
		}
	}
	w.WriteString("\n);")
	return nil
}
