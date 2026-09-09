// Package thailunar converts Gregorian dates to and from the Thai lunar
// calendar (ปฏิทินจันทรคติไทย) using the Suriyayatra arithmetic reckoning that
// underlies the printed Thai royal calendar.
package thailunar

import (
	"fmt"
	"time"
)

// Phase is the half of the lunar month a date falls in.
type Phase int

const (
	// Waxing is ข้างขึ้น, days 1-15 of the lunar month.
	Waxing Phase = iota
	// Waning is ข้างแรม, days 1-14 or 1-15 after the full moon.
	Waning
)

func (p Phase) String() string {
	if p == Waxing {
		return "waxing"
	}
	return "waning"
}

// Thai returns the Thai name of the phase, ขึ้น or แรม.
func (p Phase) Thai() string {
	if p == Waxing {
		return "ขึ้น"
	}
	return "แรม"
}

// YearType classifies a Thai lunar year by its length.
type YearType int

const (
	// Normal is ปกติมาส ปกติวาร: 12 months, 354 days.
	Normal YearType = iota
	// Athikawan is อธิกวาร: an extra day in month 7, 355 days.
	Athikawan
	// Athikamas is อธิกมาส: a repeated month 8, 384 days.
	Athikamas
)

func (t YearType) String() string {
	switch t {
	case Athikawan:
		return "athikawan"
	case Athikamas:
		return "athikamas"
	default:
		return "normal"
	}
}

// Thai returns the Thai name of the year type.
func (t YearType) Thai() string {
	switch t {
	case Athikawan:
		return "อธิกวาร"
	case Athikamas:
		return "อธิกมาส"
	default:
		return "ปกติมาส ปกติวาร"
	}
}

// Days is the number of days in a year of this type.
func (t YearType) Days() int {
	switch t {
	case Athikawan:
		return 355
	case Athikamas:
		return 384
	default:
		return 354
	}
}

// Date is a date in the Thai lunar calendar.
type Date struct {
	// Day is the day within the phase, 1-15.
	Day int
	// Phase is ขึ้น (waxing) or แรม (waning).
	Phase Phase
	// Month is the lunar month, 1-12.
	Month int
	// IsSecondEighth reports whether this is the repeated month 8
	// (เดือน 8 หลัง) of an athikamas year.
	IsSecondEighth bool
	// Year is the Buddhist Era year (พ.ศ.) by civil reckoning, which rolls
	// over on 1 January like the Gregorian year.
	Year int
	// LunarYear is the Buddhist Era year of the lunar year this date belongs
	// to. It rolls over at lunar new year, so it trails Year for dates between
	// 1 January and lunar new year.
	LunarYear int
	// YearType classifies the length of the lunar year.
	YearType YearType

	jd int
}

// FromTime converts t to its Thai lunar date. Only the calendar date of t is
// used; the clock time and location are ignored.
func FromTime(t time.Time) Date {
	return fromJD(timeToJD(t))
}

// FromGregorian converts a Gregorian year, month and day to a Thai lunar date.
func FromGregorian(year int, month time.Month, day int) Date {
	return fromJD(gregorianToJD(year, int(month), day))
}

func fromJD(jd int) Date {
	c := csFromJulianDay(jd)

	day, phase := c.day, Waxing
	if day > 15 {
		day -= 15
		phase = Waning
	}

	yt := Normal
	switch c.y0.calType {
	case typeAthikawan:
		yt = Athikawan
	case typeAthikamas:
		yt = Athikamas
	}

	gy, _, _ := jdToGregorian(jd)

	return Date{
		Day:            day,
		Phase:          phase,
		Month:          c.month(),
		IsSecondEighth: c.isSecondEighth(),
		Year:           gy + 543,
		LunarYear:      c.year + 1181,
		YearType:       yt,
		jd:             jd,
	}
}

// Time returns the Gregorian date corresponding to d, at midnight UTC.
func (d Date) Time() time.Time {
	y, m, day := jdToGregorian(d.jd)
	return time.Date(y, time.Month(m), day, 0, 0, 0, 0, time.UTC)
}

// Gregorian returns the Gregorian year, month and day corresponding to d.
func (d Date) Gregorian() (int, time.Month, int) {
	y, m, day := jdToGregorian(d.jd)
	return y, time.Month(m), day
}

// IsFullMoon reports whether d is a full moon day (ขึ้น 15 ค่ำ).
func (d Date) IsFullMoon() bool {
	return d.Phase == Waxing && d.Day == 15
}

// IsNewMoon reports whether d is a new moon day, the last day of the waning
// phase. That is แรม 15 ค่ำ in a 30-day month and แรม 14 ค่ำ in a 29-day one.
func (d Date) IsNewMoon() bool {
	if d.Phase != Waning {
		return false
	}
	return fromJD(d.jd+1).Day == 1
}

// IsWanPhra reports whether d is a Buddhist sabbath day (วันพระ): the 8th and
// 15th of the waxing phase, and the 8th and last day of the waning phase.
func (d Date) IsWanPhra() bool {
	if d.Day == 8 {
		return true
	}
	if d.Phase == Waxing {
		return d.Day == 15
	}
	return d.IsNewMoon()
}

// MonthName returns the Thai name of the lunar month, e.g. "เดือน ๘ หลัง".
func (d Date) MonthName() string {
	n := thaiDigits(d.Month)
	if d.IsSecondEighth {
		return "เดือน " + n + " หลัง"
	}
	return "เดือน " + n
}

// String formats d in the conventional Thai form, e.g. "ขึ้น ๑๕ ค่ำ เดือน ๘ ปี ๒๕๖๙".
func (d Date) String() string {
	return fmt.Sprintf("%s %s ค่ำ %s ปี %s",
		d.Phase.Thai(), thaiDigits(d.Day), d.MonthName(), thaiDigits(d.Year))
}

// Format formats d in English, e.g. "waxing 15, month 8, BE 2569".
func (d Date) Format() string {
	m := fmt.Sprintf("month %d", d.Month)
	if d.IsSecondEighth {
		m += " (second)"
	}
	return fmt.Sprintf("%s %d, %s, BE %d", d.Phase, d.Day, m, d.Year)
}

var thaiNumerals = [...]rune{'๐', '๑', '๒', '๓', '๔', '๕', '๖', '๗', '๘', '๙'}

func thaiDigits(n int) string {
	if n == 0 {
		return string(thaiNumerals[0])
	}
	var out []rune
	for n > 0 {
		out = append([]rune{thaiNumerals[n%10]}, out...)
		n /= 10
	}
	return string(out)
}

// YearTypeOf returns the type of the Thai lunar year for the given Buddhist Era
// year (พ.ศ.).
func YearTypeOf(be int) YearType {
	switch calculateYear0(be - 1181).calType {
	case typeAthikawan:
		return Athikawan
	case typeAthikamas:
		return Athikamas
	default:
		return Normal
	}
}
