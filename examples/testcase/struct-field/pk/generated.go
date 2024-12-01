package pk

import (
	"database/sql/driver"

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
func (v Car) GetID() driver.Value {
	return v.ID
}
func (v Car) GetNo() driver.Value {
	return v.No
}
func (v Car) GetColor() driver.Value {
	return (int64)(v.Color)
}
func (v Car) GetManucDate() driver.Value {
	return v.ManucDate
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
	return "SELECT id,no FROM house WHERE id = ? LIMIT 1;", []any{v.ID}
}
func (v House) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE house SET no = ? WHERE id = ?;", []any{v.No, (int64)(v.ID)}
}
func (v House) GetID() driver.Value {
	return (int64)(v.ID)
}
func (v House) GetNo() driver.Value {
	return v.No
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
func (v User) GetID() driver.Value {
	return v.ID
}
func (v User) GetName() driver.Value {
	return (string)(v.Name)
}
func (v User) GetAge() driver.Value {
	return (int64)(v.Age)
}
func (v User) GetEmail() driver.Value {
	return v.Email
}
