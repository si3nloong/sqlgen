package postgres

import (
	"github.com/si3nloong/sqlgen/codegen/templates"
	"github.com/si3nloong/sqlgen/sequel/strpool"
)

func (d postgresDriver) AlterTableStmt(n string, model *templates.Model) string {
	buf := strpool.AcquireString()
	defer strpool.ReleaseString(buf)
	buf.WriteString("`ALTER TABLE `+ " + n + ".TableName() +` (")
	for i, f := range model.Fields {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString("MODIFY " + d.QuoteIdentifier(f.ColumnName) + " " + d.dataType(f))
		if model.PK.Field == f {
			buf.WriteString(" PRIMARY KEY")
		}
		if i > 0 {
			buf.WriteString(" AFTER " + d.QuoteIdentifier(model.Fields[i-1].ColumnName))
		}
	}
	buf.WriteString(");`")
	return buf.String()
}
