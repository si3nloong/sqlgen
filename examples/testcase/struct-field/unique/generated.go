package unique

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
)

func (User) TableName() string {
	return "user"
}
func (User) Columns() []string {
	return []string{"email", "age", "first_name", "last_name"} // 4
}
func (v User) Values() []any {
	return []any{
		v.Email,        // 0 - email
		(int64)(v.Age), // 1 - age
		v.FirstName,    // 2 - first_name
		v.LastName,     // 3 - last_name
	}
}
func (v *User) Addrs() []any {
	return []any{
		&v.Email,                             // 0 - email
		encoding.Uint8Scanner[uint8](&v.Age), // 1 - age
		&v.FirstName,                         // 2 - first_name
		&v.LastName,                          // 3 - last_name
	}
}
func (User) InsertPlaceholders(row int) string {
	return "(?,?,?,?)" // 4
}
func (v User) InsertOneStmt() (string, []any) {
	return "INSERT INTO user (email,age,first_name,last_name) VALUES (?,?,?,?);", v.Values()
}
func (v User) EmailValue() driver.Value {
	return v.Email
}
func (v User) AgeValue() driver.Value {
	return (int64)(v.Age)
}
func (v User) FirstNameValue() driver.Value {
	return v.FirstName
}
func (v User) LastNameValue() driver.Value {
	return v.LastName
}
func (v User) GetEmail() sequel.ColumnValuer[string] {
	return sequel.Column("email", v.Email, func(val string) driver.Value {
		return val
	})
}
func (v User) GetAge() sequel.ColumnValuer[uint8] {
	return sequel.Column("age", v.Age, func(val uint8) driver.Value {
		return (int64)(val)
	})
}
func (v User) GetFirstName() sequel.ColumnValuer[string] {
	return sequel.Column("first_name", v.FirstName, func(val string) driver.Value {
		return val
	})
}
func (v User) GetLastName() sequel.ColumnValuer[string] {
	return sequel.Column("last_name", v.LastName, func(val string) driver.Value {
		return val
	})
}
