package awsctr

import "time"

func ToMilliseconds(t time.Time) int64 {
	return t.Unix() / 1000
}
