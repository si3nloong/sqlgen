// Code generated by sqlgen, version 1.0.0. DO NOT EDIT.

package pk

import (
	"database/sql/driver"

	"github.com/si3nloong/sqlgen/sql/types"
)

// Implements `sql.Valuer` interface.
func (Car) Table() string {
	return "car"
}

func (Car) Columns() []string {
	return []string{"id", "no", "color", "manuc_date"}
}

func (Car) PKName() string {
	return "id"
}

func (v Car) PK() (driver.Value, error) {
	return ((driver.Valuer)(v.ID)).Value()
}

func (v Car) Values() []any {
	return []any{(driver.Valuer)(v.ID), v.No, int64(v.Color), v.ManucDate}
}

func (v *Car) Addrs() []any {
	return []any{types.Integer(&v.ID), &v.No, types.Integer(&v.Color), &v.ManucDate}
}

// Implements `sql.Valuer` interface.
func (House) Table() string {
	return "house"
}

func (House) Columns() []string {
	return []string{"id", "no"}
}

func (House) PKName() string {
	return "id"
}

func (v House) PK() (driver.Value, error) {
	return int64(v.ID), nil
}

func (v House) Values() []any {
	return []any{int64(v.ID), v.No}
}

func (v *House) Addrs() []any {
	return []any{types.Integer(&v.ID), &v.No}
}

// Implements `sql.Valuer` interface.
func (User) Table() string {
	return "user"
}

func (User) Columns() []string {
	return []string{"id", "name", "age", "email"}
}

func (User) PKName() string {
	return "id"
}

func (v User) PK() (driver.Value, error) {
	return v.ID, nil
}

func (v User) Values() []any {
	return []any{v.ID, string(v.Name), int64(v.Age), v.Email}
}

func (v *User) Addrs() []any {
	return []any{&v.ID, types.String(&v.Name), types.Integer(&v.Age), &v.Email}
}
