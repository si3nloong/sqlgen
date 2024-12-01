package unique

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel/encoding"
)

func (User) TableName() string {
	return "user"
}
func (User) Columns() []string {
	return []string{"email", "age", "first_name", "last_name"}
}
func (v User) Values() []any {
	return []any{v.Email, (int64)(v.Age), v.FirstName, v.LastName}
}
func (v *User) Addrs() []any {
	return []any{&v.Email, encoding.Uint8Scanner[uint8](&v.Age), &v.FirstName, &v.LastName}
}
func (User) InsertPlaceholders(row int) string {
	return "(?,?,?,?)"
}
func (v User) InsertOneStmt() (string, []any) {
	return "INSERT INTO user (email,age,first_name,last_name) VALUES (?,?,?,?);", v.Values()
}
func (v User) GetEmail() driver.Value {
	return v.Email
}
func (v User) GetAge() driver.Value {
	return (int64)(v.Age)
}
func (v User) GetFirstName() driver.Value {
	return v.FirstName
}
func (v User) GetLastName() driver.Value {
	return v.LastName
}
