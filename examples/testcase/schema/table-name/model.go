package tablename

import (
	"github.com/si3nloong/sqlgen/sequel"
)

type CustomTableName1 struct {
	_    sequel.TableName `sql:"CustomTableName_1"`
	Text string
}

type CustomTableName2 struct {
	priv sequel.TableName `sql:"table_2"`
	Text string
}

type CustomTableName3 struct {
	Public sequel.TableName `sql:"table_3"`
	Text   string
}
