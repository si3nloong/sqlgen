package unique

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (User) TableName() string {
	return "user"
}
func (User) Columns() []string {
	return []string{"email", "age", "first_name", "last_name"}
}
func (v User) Values() []any {
	return []any{string(v.Email), int64(v.Age), string(v.FirstName), string(v.LastName)}
}
func (v *User) Addrs() []any {
	return []any{types.String(&v.Email), types.Integer(&v.Age), types.String(&v.FirstName), types.String(&v.LastName)}
}
func (User) InsertPlaceholders(row int) string {
	return "(?,?,?,?)"
}
func (v User) InsertOneStmt() (string, []any) {
	return "INSERT INTO user (email,age,first_name,last_name) VALUES (?,?,?,?);", v.Values()
}
func (v User) GetEmail() sequel.ColumnValuer[string] {
	return sequel.Column("email", v.Email, func(val string) driver.Value { return string(val) })
}
func (v User) GetAge() sequel.ColumnValuer[uint8] {
	return sequel.Column("age", v.Age, func(val uint8) driver.Value { return int64(val) })
}
func (v User) GetFirstName() sequel.ColumnValuer[string] {
	return sequel.Column("first_name", v.FirstName, func(val string) driver.Value { return string(val) })
}
func (v User) GetLastName() sequel.ColumnValuer[string] {
	return sequel.Column("last_name", v.LastName, func(val string) driver.Value { return string(val) })
}
