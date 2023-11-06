package date

import (
	"cloud.google.com/go/civil"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `sql:",pk"`
	BirthDate civil.Date
}
