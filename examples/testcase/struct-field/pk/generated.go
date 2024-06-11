package pk

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v Car) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `car` (`id` BIGINT NOT NULL,`no` VARCHAR(255) NOT NULL,`color` INTEGER NOT NULL,`manuc_date` DATETIME NOT NULL,PRIMARY KEY (`id`));"
}
func (Car) TableName() string {
	return "car"
}
func (Car) InsertOneStmt() string {
	return "INSERT INTO car (id,no,color,manuc_date) VALUES (?,?,?,?);"
}
func (Car) InsertVarQuery() string {
	return "(?,?,?,?)"
}
func (Car) Columns() []string {
	return []string{"id", "no", "color", "manuc_date"}
}
func (Car) HasPK() {}
func (v Car) PK() ([]string, []int, []any) {
	return []string{"id"}, []int{0}, []any{(driver.Valuer)(v.ID)}
}
func (Car) FindByPKStmt() string {
	return "SELECT id,no,color,manuc_date FROM car WHERE id = ? LIMIT 1;"
}
func (Car) UpdateByPKStmt() string {
	return "UPDATE car SET no = ?,color = ?,manuc_date = ? WHERE id = ? LIMIT 1;"
}
func (v Car) Values() []any {
	return []any{(driver.Valuer)(v.ID), string(v.No), int64(v.Color), time.Time(v.ManucDate)}
}
func (v *Car) Addrs() []any {
	return []any{types.Integer(&v.ID), types.String(&v.No), types.Integer(&v.Color), (*time.Time)(&v.ManucDate)}
}
func (v Car) GetID() sequel.ColumnValuer[PK] {
	return sequel.Column("id", v.ID, func(val PK) driver.Value { return (driver.Valuer)(val) })
}
func (v Car) GetNo() sequel.ColumnValuer[string] {
	return sequel.Column("no", v.No, func(val string) driver.Value { return string(val) })
}
func (v Car) GetColor() sequel.ColumnValuer[Color] {
	return sequel.Column("color", v.Color, func(val Color) driver.Value { return int64(val) })
}
func (v Car) GetManucDate() sequel.ColumnValuer[time.Time] {
	return sequel.Column("manuc_date", v.ManucDate, func(val time.Time) driver.Value { return time.Time(val) })
}

func (v User) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `user` (`id` BIGINT NOT NULL,`name` VARCHAR(255) NOT NULL,`age` TINYINT UNSIGNED NOT NULL,`email` VARCHAR(255) NOT NULL,PRIMARY KEY (`id`));"
}
func (User) TableName() string {
	return "user"
}
func (User) InsertOneStmt() string {
	return "INSERT INTO user (id,name,age,email) VALUES (?,?,?,?);"
}
func (User) InsertVarQuery() string {
	return "(?,?,?,?)"
}
func (User) Columns() []string {
	return []string{"id", "name", "age", "email"}
}
func (User) HasPK() {}
func (v User) PK() ([]string, []int, []any) {
	return []string{"id"}, []int{0}, []any{int64(v.ID)}
}
func (User) FindByPKStmt() string {
	return "SELECT id,name,age,email FROM user WHERE id = ? LIMIT 1;"
}
func (User) UpdateByPKStmt() string {
	return "UPDATE user SET name = ?,age = ?,email = ? WHERE id = ? LIMIT 1;"
}
func (v User) Values() []any {
	return []any{int64(v.ID), string(v.Name), int64(v.Age), string(v.Email)}
}
func (v *User) Addrs() []any {
	return []any{types.Integer(&v.ID), types.String(&v.Name), types.Integer(&v.Age), types.String(&v.Email)}
}
func (v User) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column("id", v.ID, func(val int64) driver.Value { return int64(val) })
}
func (v User) GetName() sequel.ColumnValuer[LongText] {
	return sequel.Column("name", v.Name, func(val LongText) driver.Value { return string(val) })
}
func (v User) GetAge() sequel.ColumnValuer[uint8] {
	return sequel.Column("age", v.Age, func(val uint8) driver.Value { return int64(val) })
}
func (v User) GetEmail() sequel.ColumnValuer[string] {
	return sequel.Column("email", v.Email, func(val string) driver.Value { return string(val) })
}

func (v House) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS `house` (`id` INTEGER NOT NULL,`no` VARCHAR(255) NOT NULL,PRIMARY KEY (`id`));"
}
func (House) TableName() string {
	return "house"
}
func (House) InsertOneStmt() string {
	return "INSERT INTO house (id,no) VALUES (?,?);"
}
func (House) InsertVarQuery() string {
	return "(?,?)"
}
func (House) Columns() []string {
	return []string{"id", "no"}
}
func (House) HasPK() {}
func (v House) PK() ([]string, []int, []any) {
	return []string{"id"}, []int{0}, []any{int64(v.ID)}
}
func (House) FindByPKStmt() string {
	return "SELECT id,no FROM house WHERE id = ? LIMIT 1;"
}
func (House) UpdateByPKStmt() string {
	return "UPDATE house SET no = ? WHERE id = ? LIMIT 1;"
}
func (v House) Values() []any {
	return []any{int64(v.ID), string(v.No)}
}
func (v *House) Addrs() []any {
	return []any{types.Integer(&v.ID), types.String(&v.No)}
}
func (v House) GetID() sequel.ColumnValuer[uint] {
	return sequel.Column("id", v.ID, func(val uint) driver.Value { return int64(val) })
}
func (v House) GetNo() sequel.ColumnValuer[string] {
	return sequel.Column("no", v.No, func(val string) driver.Value { return string(val) })
}
