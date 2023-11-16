package embedded

import (
	"time"

	c "cloud.google.com/go/civil"
)

type B struct {
	// Embedded struct
	c.DateTime
	t time.Time

	// Embedded int64
	time.Duration
}
