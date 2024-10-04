package pkautoincr

import "github.com/si3nloong/sqlgen/sequel"

type Flag bool

type LongText string

type Model struct {
	_    sequel.TableName `sql:"AutoIncrPK"`
	Name LongText
	F    Flag
	ID   uint `sql:",pk,auto_increment"`
	N    int64
}
