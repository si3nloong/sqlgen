package tablename

import "github.com/si3nloong/sqlgen/sequel"

type CustomTableName1 struct {
	sequel.Table `sql:"CustomTableName_1"`
	Text         string
}

type CustomTableName2 struct {
	priv sequel.Table `sql:"table_2"`
	Text string
}

type CustomTableName3 struct {
	Public sequel.Table `sql:"table_3"`
	Text   string
}
