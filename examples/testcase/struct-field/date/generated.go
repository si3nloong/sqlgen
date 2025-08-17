package date

import (
	"database/sql/driver"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
)

func (User) TableName() string {
	return "user"
}
func (User) HasPK() {}
func (v User) PK() (string, int, any) {
	return "id", 0, v.ID
}
func (User) Columns() []string {
	return []string{"id", "birth_date"} // 2
}
func (v User) Values() []any {
	return []any{
		v.ID,                            // 0 - id
		encoding.TextValue(v.BirthDate), // 1 - birth_date
	}
}
func (v *User) Addrs() []any {
	return []any{
		&v.ID, // 0 - id
		encoding.TextScanner[civil.Date](&v.BirthDate), // 1 - birth_date
	}
}
func (User) InsertPlaceholders(row int) string {
	return "(?,?)" // 2
}
func (v User) InsertOneStmt() (string, []any) {
	return "INSERT INTO `user` (`id`,`birth_date`) VALUES (?,?);", v.Values()
}
func (v User) FindOneByPKStmt() (string, []any) {
	return "SELECT `id`,`birth_date` FROM `user` WHERE `id` = ? LIMIT 1;", []any{v.ID}
}
func (v User) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE `user` SET `birth_date` = ? WHERE `id` = ?;", []any{encoding.TextValue(v.BirthDate), v.ID}
}
func (v User) IDValue() any {
	return v.ID
}
func (v User) BirthDateValue() any {
	return encoding.TextValue(v.BirthDate)
}
func (v User) ColumnID() sequel.ColumnValuer[uuid.UUID] {
	return sequel.Column("id", v.ID, func(val uuid.UUID) driver.Value {
		return val
	})
}
func (v User) ColumnBirthDate() sequel.ColumnValuer[civil.Date] {
	return sequel.Column("birth_date", v.BirthDate, func(val civil.Date) driver.Value {
		return encoding.TextValue(val)
	})
}
