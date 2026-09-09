package thailunar

import "time"

func gregorianToJD(y, m, d int) int {
	a := (14 - m) / 12
	yy := y + 4800 - a
	mm := m + 12*a - 3
	return d + (153*mm+2)/5 + 365*yy + yy/4 - yy/100 + yy/400 - 32045
}

func jdToGregorian(jd int) (int, int, int) {
	a := jd + 32044
	b := (4*a + 3) / 146097
	c := a - 146097*b/4
	d := (4*c + 3) / 1461
	e := c - 1461*d/4
	m := (5*e + 2) / 153

	day := e - (153*m+2)/5 + 1
	month := m + 3 - 12*(m/10)
	year := 100*b + d - 4800 + m/10
	return year, month, day
}

func timeToJD(t time.Time) int {
	y, m, d := t.Date()
	return gregorianToJD(y, int(m), d)
}
