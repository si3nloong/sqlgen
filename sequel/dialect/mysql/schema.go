package mysql

import (
	"github.com/si3nloong/sqlgen/sequel"
)

func (s *mysqlDriver) TableSchemas(table sequel.GoTableSchema) sequel.TableSchema {
	schema := new(tableDefinition)
	schema.keys = append(schema.keys, table.Keys()...)
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
