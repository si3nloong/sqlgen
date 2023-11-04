package binary

import (
	"time"

	"github.com/google/uuid"
)

type Binary struct {
	ID   uuid.UUID `sql:",binary,pk"`
	Str  string
	Time time.Time
}
