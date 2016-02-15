package util

import (
	"time"
)

func ParseTime(s string) (time.Time, error) {
	layout := "2006-01-02 15:04 Uhr"
	return time.Parse(layout, s)
}
