package postgres

import (
	"github.com/si3nloong/sqlgen/codegen/templates"
	"github.com/si3nloong/sqlgen/sequel/strpool"
)

func (d *postgresDriver) CreateTableStmt(n string, model *templates.Model) string {
	buf := strpool.AcquireString()
	defer strpool.ReleaseString(buf)

	buf.WriteString("`CREATE TABLE IF NOT EXISTS `+ " + n + ".TableName() +` (")
	for i, f := range model.Fields {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(d.QuoteIdentifier(f.ColumnName) + " " + d.dataType(f))
	}
	if len(model.Keys) > 0 {
		buf.WriteString(",PRIMARY KEY (")
		for i, k := range model.Keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(d.QuoteIdentifier(k.ColumnName))
		}
		buf.WriteByte(')')
	}
	buf.WriteString(");`")
	return buf.String()
}
