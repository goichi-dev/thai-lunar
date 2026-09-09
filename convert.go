package thailunar

import (
	"errors"
	"time"
)

// ErrNoSuchDate is returned when a lunar date does not occur in the given year,
// such as month 8 หลัง in a year that is not athikamas.
var ErrNoSuchDate = errors.New("thailunar: no such lunar date")

// ToGregorian converts a Thai lunar date to its Gregorian date.
//
// The year is the Buddhist Era year by civil reckoning (as reported by
// Date.Year), day is 1-15 within the phase, and month is 1-12. Set secondEighth
// for เดือน 8 หลัง in an athikamas year.
//
// A lunar month can fall twice within one civil year, once in January and again
// in December; ToGregorian returns the first. Use ToGregorianAll to get both.
//
// It returns ErrNoSuchDate if the date does not occur in that year.
func ToGregorian(be, month int, phase Phase, day int, secondEighth bool) (time.Time, error) {
	all := ToGregorianAll(be, month, phase, day, secondEighth)
	if len(all) == 0 {
		return time.Time{}, ErrNoSuchDate
	}
	return all[0], nil
}

// ToGregorianAll returns every Gregorian date in the civil Buddhist Era year be
// that matches the given lunar date, in chronological order. The result is
// usually one date, but a lunar month falling at both ends of a civil year
// yields two. It returns nil if the date does not occur in that year.
func ToGregorianAll(be, month int, phase Phase, day int, secondEighth bool) []time.Time {
	if month < 1 || month > 12 || day < 1 || day > 15 {
		return nil
	}

	gy := be - 543
	start := gregorianToJD(gy, 1, 1)
	end := gregorianToJD(gy+1, 1, 1)

	var out []time.Time
	for jd := start; jd < end; jd++ {
		d := fromJD(jd)
		if d.Day == day && d.Phase == phase && d.Month == month &&
			d.IsSecondEighth == secondEighth {
			out = append(out, d.Time())
			jd += 20
		}
	}
	return out
}
