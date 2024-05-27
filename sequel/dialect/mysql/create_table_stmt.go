package mysql

import (
	"github.com/si3nloong/sqlgen/codegen/templates"
	"github.com/si3nloong/sqlgen/sequel/strpool"
)

func (d *mysqlDriver) CreateTableStmt(n string, model *templates.Model) string {
	buf := strpool.AcquireString()
	defer strpool.ReleaseString(buf)

	if model.HasTableName {
		buf.WriteString(`"CREATE TABLE IF NOT EXISTS "+ ` + n + `.TableName() +" (`)
	} else {
		buf.WriteString(`"CREATE TABLE IF NOT EXISTS ` + d.QuoteIdentifier(model.TableName) + ` (`)
	}
	for i, f := range model.Fields {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(d.QuoteIdentifier(f.ColumnName) + " " + dataType(f))
		if model.PK != nil && model.PK.Field == f && model.PK.IsAutoIncr {
			buf.WriteString(" AUTO_INCREMENT")
		}
	}
	if model.PK != nil {
		buf.WriteString(",PRIMARY KEY (" + d.QuoteIdentifier(model.PK.Field.ColumnName) + ")")
	}
	buf.WriteString(`);"`)
	return buf.String()
}
