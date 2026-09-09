package thailunar

import "time"

// Holiday is a Buddhist holy day that falls on a full moon.
type Holiday int

const (
	// Makha is วันมาฆบูชา, the full moon of month 3 (month 4 in an athikamas year).
	Makha Holiday = iota
	// Visakha is วันวิสาขบูชา, the full moon of month 6 (month 7 in an athikamas year).
	Visakha
	// Asalha is วันอาสาฬหบูชา, the full moon of month 8 (เดือน 8 หลัง in an athikamas year).
	Asalha
	// KhaoPhansa is วันเข้าพรรษา, the day after Asalha.
	KhaoPhansa
	// OkPhansa is วันออกพรรษา, the full moon of month 11.
	OkPhansa
)

func (h Holiday) String() string {
	switch h {
	case Makha:
		return "Makha Bucha"
	case Visakha:
		return "Visakha Bucha"
	case Asalha:
		return "Asalha Bucha"
	case KhaoPhansa:
		return "Khao Phansa"
	default:
		return "Ok Phansa"
	}
}

// Thai returns the Thai name of the holiday.
func (h Holiday) Thai() string {
	switch h {
	case Makha:
		return "วันมาฆบูชา"
	case Visakha:
		return "วันวิสาขบูชา"
	case Asalha:
		return "วันอาสาฬหบูชา"
	case KhaoPhansa:
		return "วันเข้าพรรษา"
	default:
		return "วันออกพรรษา"
	}
}

// HolidayDate returns the Gregorian date of h in the civil Buddhist Era year be.
//
// In an athikamas year the holy days shift a month later, and Asalha falls in
// the repeated month 8 (เดือน 8 หลัง); this is handled for you.
func HolidayDate(h Holiday, be int) (time.Time, error) {
	if h == KhaoPhansa {
		t, err := HolidayDate(Asalha, be)
		if err != nil {
			return time.Time{}, err
		}
		return t.AddDate(0, 0, 1), nil
	}

	leap := YearTypeOf(be) == Athikamas

	var month int
	var second bool
	switch h {
	case Makha:
		month = 3
		if leap {
			month = 4
		}
	case Visakha:
		month = 6
		if leap {
			month = 7
		}
	case Asalha:
		month = 8
		second = leap
	case OkPhansa:
		month = 11
	default:
		return time.Time{}, ErrNoSuchDate
	}

	all := ToGregorianAll(be, month, Waxing, 15, second)
	if len(all) == 0 {
		return time.Time{}, ErrNoSuchDate
	}
	return all[len(all)-1], nil
}
