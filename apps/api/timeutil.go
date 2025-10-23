package main

import (
	"errors"
	"time"
)

// 兼容 YYYY-MM-DD / RFC3339 / "2006-01-02 15:04:05"
func parsePublishedAt(s string) (time.Time, error) {
	if s == "" {
		return time.Now(), nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local); err == nil {
		return t, nil
	}
	return time.Time{}, errors.New("bad published_at, use RFC3339 or YYYY-MM-DD")
}
