package value

import "time"

type Period struct {
	// From is the start of the period.
	From time.Time
	// To is the end of the period.
	To time.Time
}
