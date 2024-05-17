package sqlite

import (
	"github.com/si3nloong/sqlgen/codegen/templates"
	"github.com/si3nloong/sqlgen/sequel/strpool"
)

func (d sqliteDriver) AlterTableStmt(n string, model *templates.Model) string {
	buf := strpool.AcquireString()
	defer strpool.ReleaseString(buf)
	buf.WriteString(`"ALTER TABLE "+ ` + n + `.TableName() +" (`)
	for i, f := range model.Fields {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString("MODIFY " + d.QuoteIdentifier(f.ColumnName) + " " + dataType(f))
		if model.PK != nil && model.PK.Field == f && model.PK.IsAutoIncr {
			buf.WriteString(" AUTO_INCREMENT")
		}
		if i > 0 {
			// buf.WriteString(" FIRST")
			buf.WriteString(" AFTER " + d.QuoteIdentifier(model.Fields[i-1].ColumnName))
		}
	}
	buf.WriteString(`);"`)
	return buf.String()
}
