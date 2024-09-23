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

func (s *postgresDriver) Migrate(ctx context.Context, dsn string, w dialect.Writer, schema dialect.Schema) error {
	sqlConn, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer sqlConn.Close()

	existedColumns, err := s.tableColumns(ctx, sqlConn, "public", schema.TableName())
	if err != nil {
		return err
	}

	// Check existing columns and new columns is not matching
	// If it's not matching then we need to do alter table
	columns := schema.Columns()
	if reflect.DeepEqual(existedColumns, lo.Map(columns, func(_ string, i int) column {
		field := schema.ColumnGoType(i)
		return column{
			Name:               field.ColumnName(),
			DataType:           field.DataType(),
			IsNullable:         sqltype.Bool(field.GoNullable()),
			CharacterMaxLength: sql.NullInt64{},
			NumericPrecision:   sql.NullInt64{},
			DatetimePrecision:  sql.NullInt64{},
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
		column := schema.ColumnGoType(i)
		w.WriteString("\t" + s.QuoteIdentifier(column.ColumnName()) + " " + column.DataType())
		if !column.GoNullable() {
			w.WriteString(" NOT NULL")
		}
		if val, ok := column.Default(); ok {
			w.WriteString(" DEFAULT " + format(val))
		}
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
	// if indexes := schema.Indexes(); len(indexes) > 0 {
	// 	for i := range indexes {
	// 		log.Println(indexes[i])
	// 	}
	// }
	w.WriteString("\n);")

	var up, down strings.Builder
	up.WriteString("\nALTER TABLE " + s.QuoteIdentifier(schema.TableName()) + "\n")
	down.WriteString("\nALTER TABLE " + s.QuoteIdentifier(schema.TableName()) + "\n")

	for i := range columns {
		if i > 0 {
			up.WriteString(",\n")
			down.WriteString(",\n")
		}

		col := schema.ColumnGoType(i)
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
	keys := schema.Keys()
	if existedIndex, idx, ok := lo.FindIndexOf(existedIndexes, func(v index) bool {
		return v.IsPK
	}); !ok && len(keys) > 0 {
		up.WriteString(",\n\tADD PRIMARY KEY (")
		for i := range keys {
			if i > 0 {
				up.WriteString("," + s.QuoteIdentifier(keys[i]))
			} else {
				up.WriteString(s.QuoteIdentifier(keys[i]))
			}
		}
		up.WriteByte(')')
		down.WriteString("DROP PRIMARY KEY")
	} else {
		existedIndexes = append(existedIndexes[:idx], existedIndexes[idx+1:]...)
		up.WriteString(",\n\tDROP CONSTRAINT " + s.QuoteIdentifier(existedIndex.Name))
		if len(keys) > 0 {
			up.WriteString(",\n\tADD PRIMARY KEY (")
			for i := range keys {
				if i > 0 {
					up.WriteString("," + s.QuoteIdentifier(keys[i]))
				} else {
					up.WriteString(s.QuoteIdentifier(keys[i]))
				}
			}
			up.WriteByte(')')
		}
		down.WriteString(",\n\tADD PRIMARY KEY ()")
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
	if len(existedIndexes) > 0 {
		for i := range existedIndexes {
			idx := existedIndexes[i]
			up.WriteString(",\n\tDROP CONSTRAINT " + s.QuoteIdentifier(idx.Name))
			down.WriteString(",\n\tADD CONSTRAINT " + s.QuoteIdentifier(idx.Name) + "(" + strings.Join(idx.Columns, ",") + ")")
		}
		clear(existedIndexes)
	}

	up.WriteString(";")
	down.WriteString(";")

	w.WriteString(up.String())
	w.WriteString("\n\n-- +migrate Down")
	w.WriteString("\n-- SQL section 'Down' is executed when this migration is rolled back")
	w.WriteString(down.String())
	return nil
}
