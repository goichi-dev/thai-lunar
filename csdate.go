package thailunar

type monthDay struct {
	month int
	day   int
}

var findDateTable = map[byte][][2]int{
	typeNormal: {
		{383, 16}, {354, 15}, {324, 12}, {295, 11}, {265, 10}, {236, 9},
		{206, 8}, {177, 7}, {147, 6}, {118, 5}, {88, 4}, {59, 3}, {29, 2},
	},
	typeAthikawan: {
		{384, 16}, {355, 15}, {325, 12}, {296, 11}, {266, 10}, {237, 9},
		{207, 8}, {178, 7}, {148, 6}, {119, 5}, {89, 4}, {59, 3}, {29, 2},
	},
	typeAthikamas: {
		{384, 15}, {354, 12}, {325, 11}, {295, 10}, {266, 9}, {236, 8},
		{207, 7}, {177, 6}, {148, 5}, {118, 14}, {88, 13}, {59, 3}, {29, 2},
	},
}

func findDate(cal byte, days int) monthDay {
	d := days
	month := lunarMonths[1]
	for _, row := range findDateTable[cal] {
		if d > row[0] {
			d -= row[0]
			month = lunarMonths[row[1]]
			break
		}
		month = lunarMonths[1]
	}
	return monthDay{month: month, day: d}
}

type csDate struct {
	year     int
	monthRaw int
	day      int
	days     int
	y0       *lunarYear

	horakhun   int
	kammacapon int
	uccapon    int
	avoman     int
	masaken    int
	tithi      int
}

func indexOf(s []int, v int) int {
	for i, x := range s {
		if x == v {
			return i
		}
	}
	return -1
}

func newCsDate(year, month, day int) *csDate {
	c := &csDate{year: year, monthRaw: month, day: day}
	c.y0 = calculateYear0(year)

	dateOffset := -1
	switch month {
	case 5:
		dateOffset = day
	case 6:
		dateOffset = 29 + day
	}

	mp := monthPositionAB
	if c.y0.calType == typeAthikamas {
		mp = monthPositionC
	}
	tmonth := indexOf(mp, month)
	if tmonth < 1 {
		return nil
	}
	if dateOffset >= 0 && dateOffset < c.y0.offsetDays {
		if c.y0.calType == typeAthikamas {
			tmonth += 13
		} else {
			tmonth += 12
		}
		c.monthRaw += 10
	}
	cum := monthCumulativeDays[c.y0.calType]
	if tmonth-1 >= len(cum) {
		return nil
	}
	c.days = cum[tmonth-1] + day - c.y0.offsetDays

	c.calculate()
	return c
}

func (c *csDate) calculate() {
	c.horakhun = (c.year*daysIn800Years+epochOffset)/timeUnitsIn1Day + 1 + c.days
	c.kammacapon = timeUnitsIn1Day - (c.year*daysIn800Years+epochOffset)%timeUnitsIn1Day
	c.uccapon = (c.horakhun + uccaponConstant) % apogeeRotationDays
	c.avoman = (c.horakhun*11 + 650) % 692
	if c.avoman == 0 {
		c.avoman = 692
	}
	c.masaken = (((c.horakhun+c.days)*11+650)/692 + c.horakhun) / 30
	c.tithi = ((c.horakhun*11+650)/692 + c.horakhun) % 30
}

func csFromYearDays(year, days int) *csDate {
	y0 := calculateYear0(year)
	daysInYear := 365
	if y0.leapday {
		daysInYear++
	}
	yr, d := year, days
	for d > daysInYear {
		yr++
		d -= daysInYear
		y0 = calculateYear0(yr)
		daysInYear = 365
		if y0.leapday {
			daysInYear++
		}
	}
	md := findDate(y0.calType, y0.offsetDays+d)
	return newCsDate(yr, md.month, md.day)
}

func csFromJulianDay(jd int) *csDate {
	hk := jd - csJulianDayOffset
	year := (hk*800 - 373) / 292207
	var days int
	if hk%292207 == 95333 {
		year--
		days = 365
	} else {
		days = hk - yearHorakhun(year)
	}
	return csFromYearDays(year, days)
}

func (c *csDate) month() int {
	if c.monthRaw == 15 || c.monthRaw == 16 {
		return c.monthRaw - 10
	}
	if c.monthRaw == 88 {
		return 8
	}
	return c.monthRaw
}

func (c *csDate) isSecondEighth() bool { return c.monthRaw == 88 }
