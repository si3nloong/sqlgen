package mysql

import (
	"context"
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/si3nloong/sqlgen/cmd/codegen/dialect"
)

func (s *mysqlDriver) Migrate(ctx context.Context, dsn string, w dialect.Writer, schema dialect.TableMigrator) error {
	sqlConn, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer sqlConn.Close()

	if err := sqlConn.PingContext(ctx); err != nil {
		return err
	}

	var dbName string
	if err := sqlConn.QueryRowContext(ctx, "SELECT DATABASE();").Scan(&dbName); err != nil {
		return err
	}

	existingCols, err := s.tableColumns(ctx, sqlConn, dbName, schema.TableName())
	if err != nil {
		return err
	}

	columns := schema.Columns()
	newCols := make([]dialect.GoColumn, 0)
	// updatedCols := make([]*columnInfo, 0)

	// // Check existing columns and new columns is not matching
	// // If it's not matching then we need to do alter table
	// for i := range columns {
	// 	col := schema.ColumnByIndex(i)
	// 	if c, idx, ok := lo.FindIndexOf(existingCols, func(v column) bool {
	// 		return v.ColumnName == col.ColumnName()
	// 	}); ok {
	// 		existingCols = slices.Delete(existingCols, idx, idx+1)
	// 		// If the data type is different then we need to update
	// 		if col.DataType() == c.ColumnDataType() &&
	// 			col.GoNullable() == bool(c.IsNullable) {
	// 			continue
	// 		}
	// 		updatedCols = append(updatedCols, &columnInfo{
	// 			oldColumn: &c,
	// 			newColumn: col,
	// 		})
	// 	} else {
	// 		newCols = append(newCols, col)
	// 	}
	// }

	// if len(newCols) == 0 && len(updatedCols) == 0 {
	// 	return dialect.ErrNoNewMigration
	// }

	w.WriteString("-- +migrate Up")
	w.WriteString("\n-- SQL in section 'Up' is executed when this migration is applied")
	w.WriteString("\nCREATE TABLE IF NOT EXISTS " + s.QuoteIdentifier(schema.TableName()) + " (\n")
	for i := range columns {
		if i > 0 {
			w.WriteString(",\n")
		}
		column := schema.ColumnByIndex(i)
		w.WriteString("\t" + s.QuoteIdentifier(column.ColumnName()) + " " + column.DataType())
		if !column.GoNullable() {
			w.WriteString(" NOT NULL")
		}
		if val, ok := column.Default(); ok {
			w.WriteString(" DEFAULT " + format(val))
		}
	}

	pks := schema.PK()
	if len(pks) > 0 {
		w.WriteString(",\n\tCONSTRAINT " + s.QuoteIdentifier(indexName(pks, pk)) + " PRIMARY KEY (")
		for i := range pks {
			if i > 0 {
				w.WriteString("," + s.QuoteIdentifier(pks[i]))
			} else {
				w.WriteString(s.QuoteIdentifier(pks[i]))
			}
		}
		w.WriteString(")")
	}

	// schema.RangeIndex(func(idx dialect.Index, _ int) {
	// 	if idx.Unique() {
	// 		w.WriteString("ADD CONSTRAINT " + indexName(idx.Columns(), unique) + " UNIQUE (")
	// 	} else {
	// 		w.WriteString("ADD CONSTRAINT " + indexName(idx.Columns(), bTree) + " INDEX (")
	// 	}
	// 	for i, col := range idx.Columns() {
	// 		if i > 0 {
	// 			w.WriteString("," + s.QuoteIdentifier(col))
	// 		} else {
	// 			w.WriteString(s.QuoteIdentifier(col))
	// 		}
	// 	}
	// 	w.WriteString(")")
	// })
	w.WriteString("\n);")

	var up, down strings.Builder
	up.WriteString("\nALTER TABLE " + s.QuoteIdentifier(schema.TableName()) + "\n")
	down.WriteString("\nALTER TABLE " + s.QuoteIdentifier(schema.TableName()) + "\n")
	// for _, col := range updatedCols {
	// 	up.WriteString("\tALTER " + s.QuoteIdentifier(col.newColumn.ColumnName()) + " SET DATA TYPE " + col.newColumn.DataType() + ",\n")
	// 	down.WriteString("\tALTER " + s.QuoteIdentifier(col.oldColumn.ColumnName) + " SET DATA TYPE " + col.oldColumn.ColumnDataType() + ",\n")
	// 	if !col.oldColumn.IsNullable {
	// 		down.WriteString("\tALTER " + s.QuoteIdentifier(col.oldColumn.ColumnName) + " SET NOT NULL,\n")
	// 	}
	// 	if col.oldColumn.Default.Valid {
	// 		down.WriteString("\tALTER " + s.QuoteIdentifier(col.oldColumn.ColumnName) + " SET DEFAULT " + col.oldColumn.Default.String + ",\n")
	// 	}
	// }
	for _, col := range newCols {
		up.WriteString("\tADD COLUMN " + s.QuoteIdentifier(col.ColumnName()) + " " + col.DataType())
		if !col.GoNullable() {
			up.WriteString(" NOT NULL")
		}
		if val, ok := col.Default(); ok {
			up.WriteString(" DEFAULT " + format(val))
		}
		up.WriteString(",\n")
		down.WriteString("\tDROP COLUMN " + s.QuoteIdentifier(col.ColumnName()) + ",\n")
	}
	for _, col := range existingCols {
		up.WriteString("\tDROP COLUMN " + s.QuoteIdentifier(col.ColumnName) + ",\n")
		down.WriteString("\tADD COLUMN " + s.QuoteIdentifier(col.ColumnName))
		if !col.IsNullable {
			down.WriteString(" NOT NULL")
		}
		if col.Default.Valid {
			down.WriteString(" DEFAULT " + col.Default.String)
		}
		down.WriteString(",\n")
	}

	// existingIdxs, err := s.tableIndexes(ctx, sqlConn, schema.TableName())
	// if err != nil {
	// 	return err
	// }

	// Migrate up
	query := up.String()
	if query[len(query)-2] == ',' {
		query = query[:len(query)-2]
	}
	w.WriteString(query + ";")

	// Migrate down
	w.WriteString("\n\n-- +migrate Down")
	w.WriteString("\n-- SQL section 'Down' is executed when this migration is rolled back")
	query = down.String()
	if query[len(query)-2] == ',' {
		query = query[:len(query)-2]
	}
	w.WriteString(query + ";")
	return nil
}
