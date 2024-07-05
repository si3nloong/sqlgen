package composite

import "github.com/gofrs/uuid/v5"

type Composite struct {
	Flag bool
	Col1 string `sql:",pk"`
	Col2 bool
	Col3 uuid.UUID `sql:",pk"`
}
