package load

import "time"

type Profile interface {
	At(elapsed time.Duration) int
	TotalDuration() time.Duration
}
