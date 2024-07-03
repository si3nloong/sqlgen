package postgres

import (
	"github.com/si3nloong/sqlgen/sequel"
)

func (s *postgresDriver) TableSchemas(table sequel.GoTableSchema) sequel.TableSchema {
	schema := new(tableDefinition)
	for i := range table.Columns() {
		col := table.Column(i)
		column := columnDefinition{}
		column.name = col.ColumnName()
		column.length = col.Size()
		column.dataType = dataType(col)
		// column.nullable =
		column.comment = ""
		schema.cols = append(schema.cols, &column)
	}
	return schema
}
