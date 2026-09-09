package thailunar

import (
	"testing"
	"time"
)

func TestFromGregorian(t *testing.T) {
	tests := []struct {
		y, m, d int
		want    string
	}{
		{2026, 5, 31, "ขึ้น ๑๕ ค่ำ เดือน ๗ ปี ๒๕๖๙"},
		{2025, 5, 11, "ขึ้น ๑๕ ค่ำ เดือน ๖ ปี ๒๕๖๘"},
		{2024, 5, 22, "ขึ้น ๑๕ ค่ำ เดือน ๖ ปี ๒๕๖๗"},
		{2026, 7, 29, "ขึ้น ๑๕ ค่ำ เดือน ๘ หลัง ปี ๒๕๖๙"},
	}
	for _, tt := range tests {
		got := FromGregorian(tt.y, time.Month(tt.m), tt.d).String()
		if got != tt.want {
			t.Errorf("FromGregorian(%d, %d, %d) = %q, want %q", tt.y, tt.m, tt.d, got, tt.want)
		}
	}
}

func TestYearType(t *testing.T) {
	tests := []struct {
		be   int
		want YearType
	}{
		{2566, Athikamas},
		{2567, Normal},
		{2568, Athikawan},
		{2569, Athikamas},
		{2570, Normal},
		{2573, Athikawan},
	}
	for _, tt := range tests {
		if got := YearTypeOf(tt.be); got != tt.want {
			t.Errorf("YearTypeOf(%d) = %v, want %v", tt.be, got, tt.want)
		}
	}
}

func TestHolidays(t *testing.T) {
	tests := []struct {
		h    Holiday
		be   int
		want string
	}{
		{Makha, 2567, "2024-02-24"},
		{Makha, 2568, "2025-02-12"},
		{Makha, 2569, "2026-03-03"},
		{Visakha, 2566, "2023-06-03"},
		{Visakha, 2567, "2024-05-22"},
		{Visakha, 2568, "2025-05-11"},
		{Visakha, 2569, "2026-05-31"},
		{Asalha, 2566, "2023-08-01"},
		{Asalha, 2568, "2025-07-10"},
		{KhaoPhansa, 2568, "2025-07-11"},
		{OkPhansa, 2568, "2025-10-07"},
	}
	for _, tt := range tests {
		got, err := HolidayDate(tt.h, tt.be)
		if err != nil {
			t.Errorf("HolidayDate(%v, %d): %v", tt.h, tt.be, err)
			continue
		}
		if g := got.Format("2006-01-02"); g != tt.want {
			t.Errorf("HolidayDate(%v, %d) = %s, want %s", tt.h, tt.be, g, tt.want)
		}
	}
}

func TestFullMoonIsWanPhra(t *testing.T) {
	d := FromGregorian(2025, time.May, 11)
	if !d.IsFullMoon() {
		t.Error("Visakha 2568 should be a full moon")
	}
	if !d.IsWanPhra() {
		t.Error("a full moon should be wan phra")
	}
}

func TestRoundTrip(t *testing.T) {
	start := gregorianToJD(1950, 1, 1)
	end := gregorianToJD(2100, 1, 1)

	prev := fromJD(start)
	for jd := start + 1; jd < end; jd++ {
		d := fromJD(jd)

		if d.Day < 1 || d.Day > 15 || d.Month < 1 || d.Month > 12 {
			t.Fatalf("jd %d out of range: %+v", jd, d)
		}

		ok := d.Day == prev.Day+1 ||
			(prev.Phase == Waxing && prev.Day == 15 && d.Phase == Waning && d.Day == 1) ||
			(prev.Phase == Waning && d.Phase == Waxing && d.Day == 1)
		if !ok {
			y, m, dd := jdToGregorian(jd)
			t.Fatalf("%04d-%02d-%02d: %s does not follow %s", y, m, dd, d.Format(), prev.Format())
		}
		prev = d
	}
}

func TestToGregorianRoundTrip(t *testing.T) {
	for _, be := range []int{2566, 2567, 2568, 2569, 2570} {
		start := gregorianToJD(be-543, 1, 1)
		end := gregorianToJD(be-543+1, 1, 1)
		for jd := start; jd < end; jd++ {
			d := fromJD(jd)
			all := ToGregorianAll(d.Year, d.Month, d.Phase, d.Day, d.IsSecondEighth)
			found := false
			for _, b := range all {
				if timeToJD(b) == jd {
					found = true
					break
				}
			}
			if !found {
				y, m, dd := jdToGregorian(jd)
				t.Fatalf("%04d-%02d-%02d %s not found in reverse conversion", y, m, dd, d.Format())
			}
		}
	}
}

func TestErrNoSuchDate(t *testing.T) {
	if _, err := ToGregorian(2567, 8, Waxing, 15, true); err == nil {
		t.Error("second month 8 in a normal year should fail")
	}
	if _, err := ToGregorian(2568, 13, Waxing, 1, false); err == nil {
		t.Error("month 13 should fail")
	}
}

func TestParse(t *testing.T) {
	want := FromGregorian(2026, time.September, 9)

	for _, s := range []string{
		"2026-09-09",
		"2569-09-09",
		"2026-09-09T13:45:00+07:00",
		"2026-09-09T13:45:00Z",
		"2026/09/09",
		"09/09/2026",
		"09-09-2026",
		"09/09/2569",
	} {
		got, err := Parse(s)
		if err != nil {
			t.Errorf("Parse(%q): %v", s, err)
			continue
		}
		if got.String() != want.String() {
			t.Errorf("Parse(%q) = %s, want %s", s, got, want)
		}
	}
}

func TestParseDayMonthOrder(t *testing.T) {
	got, err := Parse("05/03/2026")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	y, m, d := got.Gregorian()
	if y != 2026 || m != time.March || d != 5 {
		t.Errorf("05/03/2026 = %04d-%02d-%02d, want 2026-03-05", y, m, d)
	}
}

func TestParseInvalid(t *testing.T) {
	for _, s := range []string{"", "hello", "2026-13-45", "2026-3-5"} {
		if _, err := Parse(s); err == nil {
			t.Errorf("Parse(%q) should fail", s)
		}
	}
}
