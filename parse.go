package thailunar

import (
	"fmt"
	"time"
)

// ErrInvalidDate is returned when a date string cannot be parsed.
var ErrInvalidDate = fmt.Errorf("thailunar: invalid date")

// Parse converts a date string to its Thai lunar date. It accepts a plain
// date ("2026-09-09"), an RFC 3339 timestamp ("2026-09-09T13:45:00+07:00"), and
// the common Thai written forms "09/09/2026" and "09-09-2026".
//
// A year of 2400 or more is read as a Buddhist Era year, so "2569-09-09" and
// "2026-09-09" mean the same day.
func Parse(s string) (Date, error) {
	t, err := ParseTime(s)
	if err != nil {
		return Date{}, err
	}
	return FromTime(t), nil
}

// ParseTime parses the same layouts as Parse and returns the Gregorian date,
// with a Buddhist Era year converted to the Common Era.
func ParseTime(s string) (time.Time, error) {
	layouts := []string{
		"2006-01-02",
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006/01/02",
		"02/01/2006",
		"02-01-2006",
	}

	for _, l := range layouts {
		t, err := time.Parse(l, s)
		if err != nil {
			continue
		}
		if y := t.Year(); y >= 2400 {
			t = t.AddDate(-543, 0, 0)
		}
		return t, nil
	}
	return time.Time{}, fmt.Errorf("%w: %q", ErrInvalidDate, s)
}
