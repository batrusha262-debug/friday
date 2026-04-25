package tests

import "time"

// Now returns current UTC time truncated to millisecond for DB comparison.
func Now() time.Time {
	return time.Now().UTC().Truncate(time.Millisecond)
}
