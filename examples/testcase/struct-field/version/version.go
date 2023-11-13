package version

import "github.com/gofrs/uuid/v5"

type Version struct {
	ID uuid.UUID `sql:",pk"`
}
