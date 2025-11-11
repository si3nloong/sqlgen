package mysql

import (
	"fmt"
	"io"

	"github.com/si3nloong/sqlgen/cmd/sqlgen/codegen/dialect"
	"github.com/si3nloong/sqlgen/cmd/sqlgen/compiler"
	"github.com/si3nloong/sqlgen/cmd/sqlgen/internal/goutil"
	"github.com/si3nloong/sqlgen/cmd/sqlgen/internal/strfmt"
)

func (s *mysqlDriver) Migrate(t *compiler.Table) (dialect.UpFunc, dialect.DownFunc) {
	return func(w io.Writer) error {
			fmt.Fprint(w, "CREATE TABLE "+s.QuoteIdentifier(t.Name)+" (")
			pk, hasPK := t.PK()
			if n := len(t.Columns); n > 0 {
				pk, _ := pk.(*compiler.AutoIncrPrimaryKey)
				column := t.Columns[0]
				strfmt.Fwprintfln(w, strfmt.Tab, s.QuoteIdentifier(column.Name())+" "+getDataType(pk, column))
				for i := 1; i < n; i++ {
					column = t.Columns[i]
					fmt.Fprint(w, ",")
					strfmt.Fwprintfln(w, strfmt.Tab, s.QuoteIdentifier(column.Name())+" "+getDataType(pk, column))
				}
			}
			if hasPK {
				fmt.Fprintf(w, ",")
				strfmt.Fwprintfln(w, strfmt.Tab, "PRIMARY KEY (")
				switch v := pk.(type) {
				case *compiler.AutoIncrPrimaryKey:
					fmt.Fprint(w, s.QuoteIdentifier(v.Name())+")")
				case *compiler.PrimaryKey:
					fmt.Fprint(w, s.QuoteIdentifier(v.Name())+")")
				case *compiler.CompositePrimaryKey:
					columns := v.Columns
					for i := 0; i < len(columns); i++ {
						fmt.Fprint(w, ","+s.QuoteIdentifier(columns[i].Name()))
					}
					fmt.Fprint(w, ")")
				default:
					return fmt.Errorf(`mysql: invalid primary key`)
				}
			}
			strfmt.Fwprintfln(w, strfmt.NoSpace, ");")
			return nil
		}, func(w io.Writer) error {
			fmt.Fprint(w, "DROP TABLE "+s.QuoteIdentifier(t.Name)+";")
			return nil
		}
}

func getDataType(pk *compiler.AutoIncrPrimaryKey, column compiler.Column) string {
	str := ""
	t := goutil.PointerUnderlyingType(column.GoType())
	if v, ok := typeMap[compiler.GoType{Type: t}.GoString()]; ok {
		str += v
	} else {
		str += "JSON"
	}
	switch v := column.(type) {
	case *compiler.BasicColumn:
		if !v.IsNullable() {
			str += " NOT NULL"
		}
	case *compiler.GeneratedColumn:
	default:
		panic("unreachable")
	}
	if pk != nil && pk.Name() == column.Name() {
		str += " AUTO_INCREMENT"
	}
	return str
}
