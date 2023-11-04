package db

type EmptyResult struct{}

func (EmptyResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (EmptyResult) RowsAffected() (int64, error) {
	return 0, nil
}
