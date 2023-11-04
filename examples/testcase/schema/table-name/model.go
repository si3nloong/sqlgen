package tablename

import "github.com/si3nloong/sqlgen/sequel"

type CustomTableName1 struct {
	// TODO: support embedded type
	sequel.Name `sql:"CustomTableName_1"`
	Text        string
}

type CustomTableName2 struct {
	priv sequel.Name `sql:"custom_table_2"`
	Text string
}

type CustomTableName3 struct {
	Public sequel.Name `sql:"custom_table3"`
	Text   string
}
