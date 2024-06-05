package postgres

import (
	"github.com/si3nloong/sqlgen/codegen/templates"
	"github.com/si3nloong/sqlgen/sequel/strpool"
)

func (d *postgresDriver) AlterTableStmt(n string, model *templates.Model) string {
	buf := strpool.AcquireString()
	defer strpool.ReleaseString(buf)
	if model.HasTableName {
		buf.WriteString("`ALTER TABLE `+ " + n + ".TableName() +` (")
	} else {
		buf.WriteString("`ALTER TABLE " + d.QuoteIdentifier(model.TableName) + " (")
	}
	for i, f := range model.Fields {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString("MODIFY " + d.QuoteIdentifier(f.ColumnName) + " " + d.dataType(f))
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
