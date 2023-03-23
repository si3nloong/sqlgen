package sql

type emptyResult struct{}

func (emptyResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (emptyResult) RowsAffected() (int64, error) {
	return 0, nil
}
