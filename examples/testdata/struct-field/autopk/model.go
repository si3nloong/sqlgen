package autopk

type Flag bool

type LongText string

type Model struct {
	Name LongText
	F    Flag
	ID   uint `sql:",pk,auto_increment"`
	N    int64
}
