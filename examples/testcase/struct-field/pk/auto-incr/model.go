package pkautoincr

type Flag bool

type LongText string

// +sql:table=AutoIncrPK
type Model struct {
	Name LongText
	F    Flag
	ID   uint `sql:",pk,auto_increment"`
	N    int64
}
