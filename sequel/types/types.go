// The `types` package is reference from :
//
// https://jackieli.dev/posts/pointers-in-go-used-in-sql-scanner/
//
// This package is a helper library to prevent the value being fallback using reflection in `database/sql`.
package types

import (
	"database/sql"
	"database/sql/driver"
	"unsafe"
)

type addrOrPtr[T any] interface {
	*T | **T
}

type ValueScanner[T any] interface {
	driver.Valuer
	sql.Scanner
}

const nullStr = "null"

var nullBytes = unsafe.Slice(unsafe.StringData(nullStr), len(nullStr))
