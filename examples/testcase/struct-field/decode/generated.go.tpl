package decode

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/examples/testencoding"
	"github.com/si3nloong/sqlgen/sequel"
)

func (v Model) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `model` (`id` TINYINT UNSIGNED NOT NULL,`text` VARCHAR(255) NOT NULL,`t` DATETIME NOT NULL);"
}
func (Model) TableName() string {
	return "model"
}
func (Model) InsertOneStmt() string {
	return "INSERT INTO model (id,text,t) VALUES (?,?,?);"
}
func (Model) InsertVarQuery() string {
	return "(?,?,?)"
}
func (Model) Columns() []string {
	return []string{"id", "text", "t"}
}
func (v Model) Values() []any {
	return []any{int64(v.ID), string(v.Text), time.Time(v.T)}
}
func (v *Model) Addrs() []any {
	return []any{testencoding.UnmarshalAny(&v.ID), testencoding.UnmarshalString(&v.Text), (*time.Time)(&v.T)}
}
func (v Model) GetID() sequel.ColumnValuer[uint8] {
	return sequel.Column("id", v.ID, func(val uint8) driver.Value { return int64(val) })
}
func (v Model) GetText() sequel.ColumnValuer[LongText] {
	return sequel.Column("text", v.Text, func(val LongText) driver.Value { return string(val) })
}
func (v Model) GetT() sequel.ColumnValuer[time.Time] {
	return sequel.Column("t", v.T, func(val time.Time) driver.Value { return time.Time(val) })
}
