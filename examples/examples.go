package examples

import (
	"database/sql"
	"errors"
	"os"
)

var (
	dbConn *sql.DB
)

func openSqlConn(driver string) (*sql.DB, error) {
	switch driver {
	case "mysql":
		return sql.Open("mysql", "root:abcd1234@/sqlbench?parseTime=true")
	case "sqlite":
		os.Remove("./sqlite.db")
		return sql.Open("sqlite3", "./sqlite.db")
	default:
		return nil, errors.New("unsupported sql driver")
	}
}

func mustValue[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func mustNoError(err error) {
	if err != nil {
		panic(err)
	}
}
