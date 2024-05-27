package sequel

import "database/sql"

type result struct {
	lastID       int64
	rowsAffected int64
}

func (r result) LastInsertId() (int64, error) {
	return r.lastID, nil
}

func (r result) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}

func NewRowsAffectedResult(rowsAffected int64) sql.Result {
	return &result{rowsAffected: rowsAffected}
}

type EmptyResult struct{}

func (EmptyResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (EmptyResult) RowsAffected() (int64, error) {
	return 0, nil
}
