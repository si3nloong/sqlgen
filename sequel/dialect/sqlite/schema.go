package sqlite

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/si3nloong/sqlgen/sequel"
)

func (s *sqliteDriver) TableSchemas(table sequel.TableSchema) sequel.TableDefinition {
	def := sequel.TableDefinition{}
	for _, col := range table.Columns() {
		if k, ok := table.AutoIncrKey(); ok && k == col {
			def.Columns = append(def.Columns, sequel.ColumnDefinition{
				Definition: dataType(col) + " AUTO_INCREMENT",
			})
		} else {
			def.Columns = append(def.Columns, sequel.ColumnDefinition{
				Definition: dataType(col),
			})
		}
	}
	if keys := table.Keys(); len(keys) > 0 {
		keyCols := lo.Map(keys, func(v sequel.ColumnSchema, _ int) string {
			return v.ColumnName()
		})
		def.PK.Columns = append(def.PK.Columns, keyCols...)
		def.PK.Definition = fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(keyCols, ","))
	}
	if idxs := table.Indexes(); len(idxs) > 0 {
		def.Indexes = append(def.Indexes, sequel.IndexDefinition{
			Definition: "",
		})
	}
	return def
}
