package sql

import "context"

// Query
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
	return Query(ctx, db, query, args, func(p Ptr) []any {
		return p.Addrs()
	})
}
