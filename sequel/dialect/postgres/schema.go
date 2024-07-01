package postgres

import (
	"github.com/si3nloong/sqlgen/sequel"
)

func (s *postgresDriver) TableSchemas(table sequel.TableSchema) sequel.TableDefinition {
	// schemas := make([]string, 0)
	// for _, col := range table.Columns() {
	// 	schemas = append(schemas, dataType(col))
	// }
	// for _, idx := range table.Indexes() {
	// 	log.Println(idx)
	// }
	// if len(table.Keys()) > 0 {
	// 	log.Println("PRIMARY KEY (" + strings.Join(lo.Map(table.Keys(), func(f sequel.ColumnSchema, _ int) string {
	// 		return f.ColumnName()
	// 	}), ",") + ")")
	// }
	// return schemas
	return sequel.TableDefinition{}
}
