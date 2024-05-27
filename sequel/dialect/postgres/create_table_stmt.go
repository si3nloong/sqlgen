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
		if model.PK != nil && model.PK.Field == f && model.PK.IsAutoIncr {
			buf.WriteString(" AUTO_INCREMENT")
		}
	}
	if model.PK != nil {
		buf.WriteString(",PRIMARY KEY (" + d.QuoteIdentifier(model.PK.Field.ColumnName) + ")")
	}
	buf.WriteString(");`")
	return buf.String()
}
