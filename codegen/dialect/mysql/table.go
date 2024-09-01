package mysql

import (
	"github.com/si3nloong/sqlgen/codegen/dialect"
)

func (s *mysqlDriver) CreateTableStmt(w dialect.Writer, schema dialect.Schema) error {
	w.WriteString("CREATE TABLE " + s.QuoteIdentifier(schema.TableName()) + " IF NOT EXISTS (\n")
	columns := schema.Columns()
	for i, column := range columns {
		if i > 0 {
			w.WriteString(",")
		}
		w.WriteString("\t" + s.QuoteIdentifier(column) + " NOT NULL")
	}
	w.WriteString(");\n")
	return nil
}

func (s *mysqlDriver) AlterTableStmt(w dialect.Writer, schema dialect.Schema) error {
	w.WriteString("ALTER TABLE " + s.QuoteIdentifier(schema.TableName()) + " (\n")
	columns := schema.Columns()
	for i, column := range columns {
		if i > 0 {
			w.WriteString(",")
		}
		w.WriteString("\t" + s.QuoteIdentifier(column) + " AFTER ")
	}
	w.WriteString(");\n\n")
	return nil
}

// func (s *mysqlDriver) TableSchemas(table sequel.GoTableSchema) sequel.TableSchema {
// 	schema := new(tableDefinition)
// 	autoIncrKey, hasAutoIncr := table.AutoIncrKey()
// 	schema.keys = append(schema.keys, table.Keys()...)
// 	for i := range table.Columns() {
// 		col := table.Column(i)
// 		colDef := dataType(col)
// 		// key shouldn't have default value
// 		if lo.IndexOf(schema.keys, col.ColumnName()) > -1 {
// 			colDef.defaultValue = nil
// 		}
// 		if hasAutoIncr && autoIncrKey == col {
// 			colDef.extra = "AUTO_INCREMENT"
// 		}
// 		schema.cols = append(schema.cols, colDef)
// 	}
// 	for i := range table.Indexes() {
// 		idx := table.Index(i)
// 		idxDef := indexDefinition{}
// 		idxDef.cols = idx.Columns()
// 		switch idx.Type() {
// 		case "unique":
// 			idxDef.indexType = unique
// 		}
// 		schema.idxs = append(schema.idxs, &idxDef)
// 	}
// 	return schema
// }
