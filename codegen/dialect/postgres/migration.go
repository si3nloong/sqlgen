package postgres

import (
	"context"
	"database/sql"
	"reflect"
	"strings"

	"github.com/samber/lo"
	"github.com/si3nloong/sqlgen/codegen/dialect"
	"github.com/si3nloong/sqlgen/internal/sqltype"
)

func (s *postgresDriver) Migrate(ctx context.Context, dsn string, w dialect.Writer, schema dialect.TableMigrator) error {
	sqlConn, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer sqlConn.Close()

	if err := sqlConn.PingContext(ctx); err != nil {
		return err
	}

	var dbName string
	if err := sqlConn.QueryRowContext(ctx, `SELECT CURRENT_DATABASE();`).Scan(&dbName); err != nil {
		return err
	}

	existedColumns, err := s.tableColumns(ctx, sqlConn, dbName, schema.TableName())
	if err != nil {
		return err
	}

	// Check existing columns and new columns is not matching
	// If it's not matching then we need to do alter table
	columns := schema.Columns()
	// log.Println(existedColumns)
	// log.Println(reflect.DeepEqual(existedColumns, lo.Map(columns, func(_ string, i int) column {
	// 	field := schema.ColumnByIndex(i)
	// 	return column{
	// 		Name:       field.ColumnName(),
	// 		DataType:   field.DataType(),
	// 		IsNullable: sqltype.Bool(field.GoNullable()),
	// 		// CharacterMaxLength: toNullInt64(field.CharacterMaxLength()),
	// 		// NumericPrecision:   toNullInt64(field.NumericPrecision()),
	// 		// DatetimePrecision:  toNullInt64(field.DatetimePrecision()),
	// 	}
	// })))
	if reflect.DeepEqual(existedColumns, lo.Map(columns, func(_ string, i int) column {
		field := schema.ColumnByIndex(i)
		return column{
			Name:       field.ColumnName(),
			DataType:   field.DataType(),
			IsNullable: sqltype.Bool(field.GoNullable()),
			// CharacterMaxLength: toNullInt64(field.CharacterMaxLength()),
			// NumericPrecision:   toNullInt64(field.NumericPrecision()),
			// DatetimePrecision:  toNullInt64(field.DatetimePrecision()),
		}
	})) {
		return dialect.ErrNoNewMigration
	}

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

	schema.RangeIndex(func(idx dialect.Index, _ int) {
		if idx.Unique() {
			w.WriteString("ADD CONSTRAINT " + indexName(idx.Columns(), unique) + " UNIQUE (")
		} else {
			w.WriteString("ADD CONSTRAINT " + indexName(idx.Columns(), bTree) + " INDEX (")
		}
		for i, col := range idx.Columns() {
			if i > 0 {
				w.WriteString("," + s.QuoteIdentifier(col))
			} else {
				w.WriteString(s.QuoteIdentifier(col))
			}
		}
		w.WriteString(")")
	})
	w.WriteString("\n);")

	var up, down strings.Builder
	up.WriteString("\nALTER TABLE " + s.QuoteIdentifier(schema.TableName()) + "\n")
	down.WriteString("\nALTER TABLE " + s.QuoteIdentifier(schema.TableName()) + "\n")

	for i := range columns {
		if i > 0 {
			up.WriteString(",\n")
			down.WriteString(",\n")
		}

		col := schema.ColumnByIndex(i)
		if prevColumn, idx, _ := lo.FindIndexOf(existedColumns, func(v column) bool {
			return v.Name == col.ColumnName()
		}); idx < 0 {
			up.WriteString("\tADD COLUMN " + s.QuoteIdentifier(col.ColumnName()))
			if !col.GoNullable() {
				up.WriteString(" NOT NULL")
			}
			if val, ok := col.Default(); ok {
				up.WriteString(" DEFAULT " + format(val))
			}
			down.WriteString("\tDROP COLUMN " + s.QuoteIdentifier(col.ColumnName()))
		} else {
			// Remove existing columns
			existedColumns = append(existedColumns[:idx], existedColumns[idx+1:]...)
			up.WriteString("\tALTER " + s.QuoteIdentifier(col.ColumnName()) + " SET DATA TYPE " + col.DataType())
			down.WriteString("\tALTER " + s.QuoteIdentifier(col.ColumnName()) + " SET DATA TYPE " + prevColumn.ColumnType())
			// If prev setup is different with the current one
			if bool(prevColumn.IsNullable) != col.GoNullable() {
				if !col.GoNullable() {
					up.WriteString(",\n\tALTER " + s.QuoteIdentifier(col.ColumnName()) + " SET NOT NULL")
					down.WriteString(",\n\tALTER " + s.QuoteIdentifier(col.ColumnName()) + " DROP NOT NULL")
				} else {
					up.WriteString(",\n\tALTER " + s.QuoteIdentifier(col.ColumnName()) + " DROP NOT NULL")
					down.WriteString(",\n\tALTER " + s.QuoteIdentifier(col.ColumnName()) + " SET NOT NULL")
				}
			}
			// if val, ok := column.Default(); ok {
			// 	up.WriteString(",\n\tALTER " + s.QuoteIdentifier(column.ColumnName()) + " SET DEFAULT " + format(val))
			// } else {
			// 	up.WriteString(",\n\tALTER " + s.QuoteIdentifier(column.ColumnName()) + " DROP DEFAULT")
			// }
		}
	}

	existedIndexes, err := s.tableIndexes(ctx, sqlConn, schema.TableName())
	if err != nil {
		return err
	}

	if pkIndex, _, ok := lo.FindIndexOf(existedIndexes, func(v index) bool {
		return v.IsPK
	}); !ok && len(pks) > 0 {
		up.WriteString(",\n\tADD CONSTRAINT " + s.QuoteIdentifier(indexName(pks, pk)) + " PRIMARY KEY (")
		for i := range pks {
			if i > 0 {
				up.WriteString("," + s.QuoteIdentifier(pks[i]))
			} else {
				up.WriteString(s.QuoteIdentifier(pks[i]))
			}
		}
		up.WriteByte(')')
		down.WriteString(",\n\tDROP PRIMARY KEY")
	} else {
		// log.Println(existedIndexes, len(existedIndexes))

		// existedIndexes = append(existedIndexes[:idx], existedIndexes[idx+1:]...)
		up.WriteString(",\n\tDROP CONSTRAINT " + s.QuoteIdentifier(pkIndex.Name))
		// If the current primary key index name is not similar to the previous one
		// we definitely need to replace it
		if len(pks) > 0 && indexName(pks, pk) != pkIndex.Name {
			up.WriteString(",\n\tADD " + s.QuoteIdentifier(indexName(pks, pk)) + " PRIMARY KEY (")
			for i := range pks {
				if i > 0 {
					up.WriteString("," + s.QuoteIdentifier(pks[i]))
				} else {
					up.WriteString(s.QuoteIdentifier(pks[i]))
				}
			}
			up.WriteByte(')')
		}
		// down.WriteString(",\n\tADD PRIMARY KEY ()")
	}
	if len(existedColumns) > 0 {
		for i := range existedColumns {
			column := existedColumns[i]
			up.WriteString(",\n\tDROP COLUMN " + s.QuoteIdentifier(column.Name))
			down.WriteString(",\n\tADD COLUMN " + s.QuoteIdentifier(column.Name))
			if !column.IsNullable {
				down.WriteString(" NOT NULL")
			}
		}
		clear(existedColumns)
	}
	// if len(existedIndexes) > 0 {
	// 	for i := range existedIndexes {
	// 		idx := existedIndexes[i]
	// 		up.WriteString(",\n\tDROP CONSTRAINT " + s.QuoteIdentifier(idx.Name))
	// 		down.WriteString(",\n\tADD CONSTRAINT " + s.QuoteIdentifier(idx.Name) + " (")
	// 		for j := range idx.Columns {
	// 			if j > 0 {
	// 				down.WriteString("," + s.QuoteIdentifier(idx.Columns[j]))
	// 			} else {
	// 				down.WriteString(s.QuoteIdentifier(idx.Columns[j]))
	// 			}
	// 		}
	// 		down.WriteByte(')')
	// 	}
	// 	clear(existedIndexes)
	// }
	up.WriteString(";")
	down.WriteString(";")

	w.WriteString(up.String())
	w.WriteString("\n\n-- +migrate Down")
	w.WriteString("\n-- SQL section 'Down' is executed when this migration is rolled back")
	w.WriteString(down.String())
	return nil
}

func toNullInt64(n int64, ok bool) sql.NullInt64 {
	return sql.NullInt64{Int64: n, Valid: ok}
}
