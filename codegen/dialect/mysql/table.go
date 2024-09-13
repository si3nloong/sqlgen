package mysql

import (
	"github.com/si3nloong/sqlgen/codegen/dialect"
)

func (s *mysqlDriver) CreateTableStmt(w dialect.Writer, schema dialect.Schema) error {
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
	w.WriteString("\n) ENGINE=InnoDB;")
	return nil
}

func (s *mysqlDriver) AlterTableStmt(w dialect.Writer, schema dialect.Schema) error {
	w.WriteString("ALTER TABLE " + s.QuoteIdentifier(schema.TableName()) + " (\n")
	columns := schema.Columns()
	for i := range columns {
		column := schema.ColumnGoType(i)
		if i > 0 {
			w.WriteString(",\n\t" + s.QuoteIdentifier(column.Name()) + " AFTER " + columns[i-1])
		} else {
			w.WriteString("\t" + s.QuoteIdentifier(column.Name()) + " FIRST")
		}
	}
	w.WriteString("\n) ENGINE=InnoDB;")
	return nil
}
