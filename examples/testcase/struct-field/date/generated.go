package date

import (
	"database/sql/driver"

	"cloud.google.com/go/civil"
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
	return "INSERT INTO user (id,birth_date) VALUES (?,?);", v.Values()
}
func (v User) FindOneByPKStmt() (string, []any) {
	return "SELECT id,birth_date FROM user WHERE id = ? LIMIT 1;", []any{v.ID}
}
func (v User) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE user SET birth_date = ? WHERE id = ?;", []any{encoding.TextValue(v.BirthDate), v.ID}
}
func (v User) GetID() driver.Value {
	return v.ID
}
func (v User) GetBirthDate() driver.Value {
	return encoding.TextValue(v.BirthDate)
}
