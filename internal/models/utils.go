package models

import (
	"time"
)

func TimeToString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func TimeFromString(t string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", t)
}
