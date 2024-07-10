package pk

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (Car) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{
		PK: &sequel.PrimaryKeyDefinition{
			Columns:    []string{"`id`"},
			Definition: "PRIMARY KEY (`id`)",
		},
		Columns: []sequel.ColumnDefinition{
			{Name: "`id`", Definition: "`id` BIGINT NOT NULL"},
			{Name: "`no`", Definition: "`no` VARCHAR(255) NOT NULL DEFAULT ''"},
			{Name: "`color`", Definition: "`color` INTEGER NOT NULL DEFAULT 0"},
			{Name: "`manuc_date`", Definition: "`manuc_date` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP"},
		},
	}
}
func (Car) TableName() string {
	return "`car`"
}
func (Car) HasPK() {}
func (v Car) PK() (string, int, any) {
	return "`id`", 0, (driver.Valuer)(v.ID)
}
func (Car) Columns() []string {
	return []string{"`id`", "`no`", "`color`", "`manuc_date`"}
}
func (v Car) Values() []any {
	return []any{(driver.Valuer)(v.ID), string(v.No), int64(v.Color), time.Time(v.ManucDate)}
}
func (v *Car) Addrs() []any {
	return []any{types.Integer(&v.ID), types.String(&v.No), types.Integer(&v.Color), (*time.Time)(&v.ManucDate)}
}
func (Car) InsertPlaceholders(row int) string {
	return "(?,?,?,?)"
}
func (v Car) InsertOneStmt() (string, []any) {
	return "INSERT INTO `car` (`id`,`no`,`color`,`manuc_date`) VALUES (?,?,?,?);", v.Values()
}
func (v Car) FindOneByPKStmt() (string, []any) {
	return "SELECT `id`,`no`,`color`,`manuc_date` FROM `car` WHERE `id` = ? LIMIT 1;", []any{(driver.Valuer)(v.ID)}
}
func (v Car) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE `car` SET `no` = ?,`color` = ?,`manuc_date` = ? WHERE `id` = ? LIMIT 1;", []any{string(v.No), int64(v.Color), time.Time(v.ManucDate), (driver.Valuer)(v.ID)}
}
func (v Car) GetID() sequel.ColumnValuer[PK] {
	return sequel.Column("`id`", v.ID, func(val PK) driver.Value { return (driver.Valuer)(val) })
}
func (v Car) GetNo() sequel.ColumnValuer[string] {
	return sequel.Column("`no`", v.No, func(val string) driver.Value { return string(val) })
}
func (v Car) GetColor() sequel.ColumnValuer[Color] {
	return sequel.Column("`color`", v.Color, func(val Color) driver.Value { return int64(val) })
}
func (v Car) GetManucDate() sequel.ColumnValuer[time.Time] {
	return sequel.Column("`manuc_date`", v.ManucDate, func(val time.Time) driver.Value { return time.Time(val) })
}

func (User) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{
		PK: &sequel.PrimaryKeyDefinition{
			Columns:    []string{"`id`"},
			Definition: "PRIMARY KEY (`id`)",
		},
		Columns: []sequel.ColumnDefinition{
			{Name: "`id`", Definition: "`id` BIGINT NOT NULL"},
			{Name: "`name`", Definition: "`name` VARCHAR(255) NOT NULL DEFAULT ''"},
			{Name: "`age`", Definition: "`age` TINYINT UNSIGNED NOT NULL DEFAULT 0"},
			{Name: "`email`", Definition: "`email` VARCHAR(255) NOT NULL DEFAULT ''"},
		},
	}
}
func (User) TableName() string {
	return "`user`"
}
func (User) HasPK() {}
func (v User) PK() (string, int, any) {
	return "`id`", 0, int64(v.ID)
}
func (User) Columns() []string {
	return []string{"`id`", "`name`", "`age`", "`email`"}
}
func (v User) Values() []any {
	return []any{int64(v.ID), string(v.Name), int64(v.Age), string(v.Email)}
}
func (v *User) Addrs() []any {
	return []any{types.Integer(&v.ID), types.String(&v.Name), types.Integer(&v.Age), types.String(&v.Email)}
}
func (User) InsertPlaceholders(row int) string {
	return "(?,?,?,?)"
}
func (v User) InsertOneStmt() (string, []any) {
	return "INSERT INTO `user` (`id`,`name`,`age`,`email`) VALUES (?,?,?,?);", v.Values()
}
func (v User) FindOneByPKStmt() (string, []any) {
	return "SELECT `id`,`name`,`age`,`email` FROM `user` WHERE `id` = ? LIMIT 1;", []any{int64(v.ID)}
}
func (v User) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE `user` SET `name` = ?,`age` = ?,`email` = ? WHERE `id` = ? LIMIT 1;", []any{string(v.Name), int64(v.Age), string(v.Email), int64(v.ID)}
}
func (v User) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column("`id`", v.ID, func(val int64) driver.Value { return int64(val) })
}
func (v User) GetName() sequel.ColumnValuer[LongText] {
	return sequel.Column("`name`", v.Name, func(val LongText) driver.Value { return string(val) })
}
func (v User) GetAge() sequel.ColumnValuer[uint8] {
	return sequel.Column("`age`", v.Age, func(val uint8) driver.Value { return int64(val) })
}
func (v User) GetEmail() sequel.ColumnValuer[string] {
	return sequel.Column("`email`", v.Email, func(val string) driver.Value { return string(val) })
}

func (House) Schemas() sequel.TableDefinition {
	return sequel.TableDefinition{
		PK: &sequel.PrimaryKeyDefinition{
			Columns:    []string{"`id`"},
			Definition: "PRIMARY KEY (`id`)",
		},
		Columns: []sequel.ColumnDefinition{
			{Name: "`id`", Definition: "`id` INTEGER UNSIGNED NOT NULL"},
			{Name: "`no`", Definition: "`no` VARCHAR(255) NOT NULL DEFAULT ''"},
		},
	}
}
func (House) TableName() string {
	return "`house`"
}
func (House) HasPK() {}
func (v House) PK() (string, int, any) {
	return "`id`", 0, int64(v.ID)
}
func (House) Columns() []string {
	return []string{"`id`", "`no`"}
}
func (v House) Values() []any {
	return []any{int64(v.ID), string(v.No)}
}
func (v *House) Addrs() []any {
	return []any{types.Integer(&v.ID), types.String(&v.No)}
}
func (House) InsertPlaceholders(row int) string {
	return "(?,?)"
}
func (v House) InsertOneStmt() (string, []any) {
	return "INSERT INTO `house` (`id`,`no`) VALUES (?,?);", v.Values()
}
func (v House) FindOneByPKStmt() (string, []any) {
	return "SELECT `id`,`no` FROM `house` WHERE `id` = ? LIMIT 1;", []any{int64(v.ID)}
}
func (v House) UpdateOneByPKStmt() (string, []any) {
	return "UPDATE `house` SET `no` = ? WHERE `id` = ? LIMIT 1;", []any{string(v.No), int64(v.ID)}
}
func (v House) GetID() sequel.ColumnValuer[uint] {
	return sequel.Column("`id`", v.ID, func(val uint) driver.Value { return int64(val) })
}
func (v House) GetNo() sequel.ColumnValuer[string] {
	return sequel.Column("`no`", v.No, func(val string) driver.Value { return string(val) })
}
