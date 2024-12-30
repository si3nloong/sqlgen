package pk

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/encoding"
)

func (Car) TableName() string {
	return "car"
}
func (Car) HasPK() {}
func (v Car) PK() (string, int, any) {
	return "id", 0, v.ID
}
func (Car) Columns() []string {
	return []string{"id", "no", "color", "manuc_date"} // 4
}
func (v Car) Values() []any {
	return []any{
		v.ID,             // 0 - id
		v.No,             // 1 - no
		(int64)(v.Color), // 2 - color
		v.ManucDate,      // 3 - manuc_date
	}
}
func (v *Car) Addrs() []any {
	return []any{
		encoding.Int64Scanner[PK](&v.ID),     // 0 - id
		&v.No,                                // 1 - no
		encoding.IntScanner[Color](&v.Color), // 2 - color
		&v.ManucDate,                         // 3 - manuc_date
	}
}
func (Car) InsertPlaceholders(row int) string {
	return "(?,?,?,?)" // 4
}
func (v Car) InsertOneStmt() (string, []any) {
	return "INSERT INTO car (id,no,color,manuc_date) VALUES (?,?,?,?);", v.Values()
}
func (v Car) FindOneByPKStmt() (string, []any) {
	return "SELECT id,no,color,manuc_date FROM car WHERE id = ? LIMIT 1;", []any{v.ID}
}
func (v Car) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE car SET no = ?,color = ?,manuc_date = ? WHERE id = ?;", []any{v.No, (int64)(v.Color), v.ManucDate, v.ID}
}
func (v Car) IDValue() driver.Value {
	return v.ID
}
func (v Car) NoValue() driver.Value {
	return v.No
}
func (v Car) ColorValue() driver.Value {
	return (int64)(v.Color)
}
func (v Car) ManucDateValue() driver.Value {
	return v.ManucDate
}
func (v Car) GetID() sequel.ColumnValuer[PK] {
	return sequel.Column("id", v.ID, func(val PK) driver.Value {
		return val
	})
}
func (v Car) GetNo() sequel.ColumnValuer[string] {
	return sequel.Column("no", v.No, func(val string) driver.Value {
		return val
	})
}
func (v Car) GetColor() sequel.ColumnValuer[Color] {
	return sequel.Column("color", v.Color, func(val Color) driver.Value {
		return (int64)(val)
	})
}
func (v Car) GetManucDate() sequel.ColumnValuer[time.Time] {
	return sequel.Column("manuc_date", v.ManucDate, func(val time.Time) driver.Value {
		return val
	})
}

func (House) TableName() string {
	return "house"
}
func (House) HasPK() {}
func (v House) PK() (string, int, any) {
	return "id", 0, (int64)(v.ID)
}
func (House) Columns() []string {
	return []string{"id", "no"} // 2
}
func (v House) Values() []any {
	return []any{
		(int64)(v.ID), // 0 - id
		v.No,          // 1 - no
	}
}
func (v *House) Addrs() []any {
	return []any{
		encoding.UintScanner[uint](&v.ID), // 0 - id
		&v.No,                             // 1 - no
	}
}
func (House) InsertPlaceholders(row int) string {
	return "(?,?)" // 2
}
func (v House) InsertOneStmt() (string, []any) {
	return "INSERT INTO house (id,no) VALUES (?,?);", v.Values()
}
func (v House) FindOneByPKStmt() (string, []any) {
	return "SELECT id,no FROM house WHERE id = ? LIMIT 1;", []any{(int64)(v.ID)}
}
func (v House) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE house SET no = ? WHERE id = ?;", []any{v.No, (int64)(v.ID)}
}
func (v House) IDValue() driver.Value {
	return (int64)(v.ID)
}
func (v House) NoValue() driver.Value {
	return v.No
}
func (v House) GetID() sequel.ColumnValuer[uint] {
	return sequel.Column("id", v.ID, func(val uint) driver.Value {
		return (int64)(val)
	})
}
func (v House) GetNo() sequel.ColumnValuer[string] {
	return sequel.Column("no", v.No, func(val string) driver.Value {
		return val
	})
}

func (User) TableName() string {
	return "user"
}
func (User) HasPK() {}
func (v User) PK() (string, int, any) {
	return "id", 0, v.ID
}
func (User) Columns() []string {
	return []string{"id", "name", "age", "email"} // 4
}
func (v User) Values() []any {
	return []any{
		v.ID,             // 0 - id
		(string)(v.Name), // 1 - name
		(int64)(v.Age),   // 2 - age
		v.Email,          // 3 - email
	}
}
func (v *User) Addrs() []any {
	return []any{
		&v.ID, // 0 - id
		encoding.StringScanner[LongText](&v.Name), // 1 - name
		encoding.Uint8Scanner[uint8](&v.Age),      // 2 - age
		&v.Email,                                  // 3 - email
	}
}
func (User) InsertPlaceholders(row int) string {
	return "(?,?,?,?)" // 4
}
func (v User) InsertOneStmt() (string, []any) {
	return "INSERT INTO user (id,name,age,email) VALUES (?,?,?,?);", v.Values()
}
func (v User) FindOneByPKStmt() (string, []any) {
	return "SELECT id,name,age,email FROM user WHERE id = ? LIMIT 1;", []any{v.ID}
}
func (v User) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE user SET name = ?,age = ?,email = ? WHERE id = ?;", []any{(string)(v.Name), (int64)(v.Age), v.Email, v.ID}
}
func (v User) IDValue() driver.Value {
	return v.ID
}
func (v User) NameValue() driver.Value {
	return (string)(v.Name)
}
func (v User) AgeValue() driver.Value {
	return (int64)(v.Age)
}
func (v User) EmailValue() driver.Value {
	return v.Email
}
func (v User) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column("id", v.ID, func(val int64) driver.Value {
		return val
	})
}
func (v User) GetName() sequel.ColumnValuer[LongText] {
	return sequel.Column("name", v.Name, func(val LongText) driver.Value {
		return (string)(val)
	})
}
func (v User) GetAge() sequel.ColumnValuer[uint8] {
	return sequel.Column("age", v.Age, func(val uint8) driver.Value {
		return (int64)(val)
	})
}
func (v User) GetEmail() sequel.ColumnValuer[string] {
	return sequel.Column("email", v.Email, func(val string) driver.Value {
		return val
	})
}
