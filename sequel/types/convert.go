package types

import (
	"database/sql/driver"

	"golang.org/x/exp/constraints"
)

func ConvertInt[T constraints.Integer](v T) driver.Value {
	return (int64)(v)
}

func ConvertFloat[T constraints.Float](v T) driver.Value {
	return (float64)(v)
}

func ConvertBool[T ~bool](v T) driver.Value {
	return (bool)(v)
}

func ConvertString[T ~string](v T) driver.Value {
	return (string)(v)
}

func ConvertBytes[T ~[]byte](v T) driver.Value {
	return ([]byte)(v)
}
