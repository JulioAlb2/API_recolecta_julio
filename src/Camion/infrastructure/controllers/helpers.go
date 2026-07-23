package controllers

import (
	"fmt"
	"time"
)

func parseFlexTime(s string) (time.Time, error) {
	formatos := []string{
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, f := range formatos {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("formato de fecha no reconocido: %s", s)
}
