package valuerproperty

type anyType struct{}

func (anyType) Value() (any, error) {
	return nil, nil
}

type B struct {
	ID    int64
	Value anyType
	N     string
}
