package testencoding

import "database/sql/driver"

func MarshalAny(v any) driver.Value {
	switch v.(type) {
	case int64:
	case string:
	}
	return nil
}

func MarshalGenericString[T ~string](v T) driver.Value {
	return string(v)
}

func UnmarshalAny(v any) driver.Value {
	switch v.(type) {
	case *int64:
	case *string:
	}
	return nil
}

func UnmarshalString[T ~string](v *T) driver.Value {
	return nil
}
