package pk

type User struct {
	ID    int64 `sql:",pk"`
	Name  LongText
	Age   uint8
	Email string
}

type House struct {
	ID uint `sql:",pk"`
	No string
}
