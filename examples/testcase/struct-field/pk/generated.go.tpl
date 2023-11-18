package pk

import (
	"database/sql/driver"
	"time"

	"github.com/si3nloong/sqlgen/sequel"
	"github.com/si3nloong/sqlgen/sequel/types"
)

func (v Car) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`id` BIGINT NOT NULL,`no` VARCHAR(255) NOT NULL,`color` INTEGER NOT NULL,`manuc_date` DATETIME NOT NULL,PRIMARY KEY (`id`));"
}
func (Car) AlterTableStmt() string {
	return "ALTER TABLE `car` MODIFY `id` BIGINT NOT NULL,MODIFY `no` VARCHAR(255) NOT NULL AFTER `id`,MODIFY `color` INTEGER NOT NULL AFTER `no`,MODIFY `manuc_date` DATETIME NOT NULL AFTER `color`;"
}
func (Car) TableName() string {
	return "`car`"
}
func (v Car) InsertOneStmt() string {
	return "INSERT INTO `car` (`id`,`no`,`color`,`manuc_date`) VALUES (?,?,?,?);"
}
func (Car) InsertVarQuery() string {
	return "(?,?,?,?)"
}
func (Car) Columns() []string {
	return []string{"`id`", "`no`", "`color`", "`manuc_date`"}
}
func (v Car) PK() (columnName string, pos int, value driver.Value) {
	return "`id`", 0, (driver.Valuer)(v.ID)
}
func (v Car) FindByPKStmt() string {
	return "SELECT `id`,`no`,`color`,`manuc_date` FROM `car` WHERE `id` = ? LIMIT 1;"
}
func (v Car) Values() []any {
	return []any{(driver.Valuer)(v.ID), string(v.No), int64(v.Color), time.Time(v.ManucDate)}
}
func (v *Car) Addrs() []any {
	return []any{types.Integer(&v.ID), types.String(&v.No), types.Integer(&v.Color), (*time.Time)(&v.ManucDate)}
}
func (v Car) GetID() sequel.ColumnValuer[PK] {
	return sequel.Column[PK]("`id`", v.ID, func(vi PK) driver.Value { return (driver.Valuer)(vi) })
}
func (v Car) GetNo() sequel.ColumnValuer[string] {
	return sequel.Column[string]("`no`", v.No, func(vi string) driver.Value { return string(vi) })
}
func (v Car) GetColor() sequel.ColumnValuer[Color] {
	return sequel.Column[Color]("`color`", v.Color, func(vi Color) driver.Value { return int64(vi) })
}
func (v Car) GetManucDate() sequel.ColumnValuer[time.Time] {
	return sequel.Column[time.Time]("`manuc_date`", v.ManucDate, func(vi time.Time) driver.Value { return time.Time(vi) })
}

func (v User) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`id` BIGINT NOT NULL,`name` VARCHAR(255) NOT NULL,`age` TINYINT UNSIGNED NOT NULL,`email` VARCHAR(255) NOT NULL,PRIMARY KEY (`id`));"
}
func (User) AlterTableStmt() string {
	return "ALTER TABLE `user` MODIFY `id` BIGINT NOT NULL,MODIFY `name` VARCHAR(255) NOT NULL AFTER `id`,MODIFY `age` TINYINT UNSIGNED NOT NULL AFTER `name`,MODIFY `email` VARCHAR(255) NOT NULL AFTER `age`;"
}
func (User) TableName() string {
	return "`user`"
}
func (v User) InsertOneStmt() string {
	return "INSERT INTO `user` (`id`,`name`,`age`,`email`) VALUES (?,?,?,?);"
}
func (User) InsertVarQuery() string {
	return "(?,?,?,?)"
}
func (User) Columns() []string {
	return []string{"`id`", "`name`", "`age`", "`email`"}
}
func (v User) PK() (columnName string, pos int, value driver.Value) {
	return "`id`", 0, int64(v.ID)
}
func (v User) FindByPKStmt() string {
	return "SELECT `id`,`name`,`age`,`email` FROM `user` WHERE `id` = ? LIMIT 1;"
}
func (v User) Values() []any {
	return []any{int64(v.ID), string(v.Name), int64(v.Age), string(v.Email)}
}
func (v *User) Addrs() []any {
	return []any{types.Integer(&v.ID), types.String(&v.Name), types.Integer(&v.Age), types.String(&v.Email)}
}
func (v User) GetID() sequel.ColumnValuer[int64] {
	return sequel.Column[int64]("`id`", v.ID, func(vi int64) driver.Value { return int64(vi) })
}
func (v User) GetName() sequel.ColumnValuer[LongText] {
	return sequel.Column[LongText]("`name`", v.Name, func(vi LongText) driver.Value { return string(vi) })
}
func (v User) GetAge() sequel.ColumnValuer[uint8] {
	return sequel.Column[uint8]("`age`", v.Age, func(vi uint8) driver.Value { return int64(vi) })
}
func (v User) GetEmail() sequel.ColumnValuer[string] {
	return sequel.Column[string]("`email`", v.Email, func(vi string) driver.Value { return string(vi) })
}

func (v House) CreateTableStmt() string {
	return "CREATE TABLE IF NOT EXISTS " + v.TableName() + " (`id` INTEGER UNSIGNED NOT NULL,`no` VARCHAR(255) NOT NULL,PRIMARY KEY (`id`));"
}
func (House) AlterTableStmt() string {
	return "ALTER TABLE `house` MODIFY `id` INTEGER UNSIGNED NOT NULL,MODIFY `no` VARCHAR(255) NOT NULL AFTER `id`;"
}
func (House) TableName() string {
	return "`house`"
}
func (v House) InsertOneStmt() string {
	return "INSERT INTO `house` (`id`,`no`) VALUES (?,?);"
}
func (House) InsertVarQuery() string {
	return "(?,?)"
}
func (House) Columns() []string {
	return []string{"`id`", "`no`"}
}
func (v House) PK() (columnName string, pos int, value driver.Value) {
	return "`id`", 0, int64(v.ID)
}
func (v House) FindByPKStmt() string {
	return "SELECT `id`,`no` FROM `house` WHERE `id` = ? LIMIT 1;"
}
func (v House) Values() []any {
	return []any{int64(v.ID), string(v.No)}
}
func (v *House) Addrs() []any {
	return []any{types.Integer(&v.ID), types.String(&v.No)}
}
func (v House) GetID() sequel.ColumnValuer[uint] {
	return sequel.Column[uint]("`id`", v.ID, func(vi uint) driver.Value { return int64(vi) })
}
func (v House) GetNo() sequel.ColumnValuer[string] {
	return sequel.Column[string]("`no`", v.No, func(vi string) driver.Value { return string(vi) })
}
