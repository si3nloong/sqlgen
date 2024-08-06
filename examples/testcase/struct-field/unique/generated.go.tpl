package unique

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (User) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{
		Columns: []sequel.ColumnDefinition{
			{Name: "email", Definition: "email VARCHAR(255) NOT NULL DEFAULT ''"},
			{Name: "age", Definition: "age TINYINT UNSIGNED NOT NULL DEFAULT 0"},
			{Name: "first_name", Definition: "first_name VARCHAR(255) NOT NULL DEFAULT ''"},
			{Name: "last_name", Definition: "last_name VARCHAR(255) NOT NULL DEFAULT ''"},
		},
		Indexes: []sequel.IndexDefinition{
			{Name: "0c83f57c786a0b4a39efab23731c7ebc", Definition: "CONSTRAINT 0c83f57c786a0b4a39efab23731c7ebc UNIQUE (email)"},
			{Name: "a108c4c3e4e0326dab6e4bdd4b13e054", Definition: "CONSTRAINT a108c4c3e4e0326dab6e4bdd4b13e054 UNIQUE (first_name,last_name)"},
		},
	}
}
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
