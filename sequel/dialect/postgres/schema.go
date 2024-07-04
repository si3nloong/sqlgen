package postgres

import (
	"github.com/si3nloong/sqlgen/sequel"
)

func (s *postgresDriver) TableSchemas(table sequel.GoTableSchema) sequel.TableSchema {
	schema := new(tableDefinition)
	schema.keys = append(schema.keys, table.Keys()...)
	for i := range table.Columns() {
		schema.cols = append(schema.cols, dataType(table.Column(i)))
	}
	return schema
}
