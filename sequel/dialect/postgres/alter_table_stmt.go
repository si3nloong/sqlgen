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
		if model.PK != nil && model.PK.Field == f {
			buf.WriteString(" PRIMARY KEY")
		}
	}
	buf.WriteString(");`")
	return buf.String()
}
