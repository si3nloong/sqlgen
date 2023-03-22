package sqlgen

import (
	"context"
	"database/sql"
	"log"
)

type Scanner[T any] interface {
	*T
	Addrs() []any
}

type Valuer[T any] interface {
	Columns() []string
	Values() []any
}

type DB interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type emptyResult struct{}

func (emptyResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (emptyResult) RowsAffected() (int64, error) {
	return 0, nil
}

type sqlConn struct {
	db   DB
	name string
}

func SqlDB(db DB) {

}

func Query[T any, Ptr interface{ *T }](
	ctx context.Context,
	db DB,
	query string,
	args []any,
	getAddrs func(Ptr) []any,
) ([]T, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []T
	for rows.Next() {
		var result T
		if err = rows.Scan(getAddrs(&result)...); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func QueryScan[T any, Ptr Scanner[T]](
	ctx context.Context,
	db DB,
	query string,
	args ...any,
) ([]T, error) {
	log.Println(query)
	return Query(ctx, db, query, args, func(p Ptr) []any {
		return p.Addrs()
	})
}
