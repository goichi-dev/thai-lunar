package thailunar

type lunarYear struct {
	year       int
	horakhun   int
	kammacapon int
	uccapon    int
	avoman     int
	masaken    int
	tithi      int
	weekday    int
	langsak    int
	nyd        int
	nextNyd    int
	leapday    bool
	calType    byte
	caldays    int
	offset     bool

	firstMonth byte
	firstDay   int
	offsetDays int
}

func yearHorakhun(year int) int {
	return (year*daysIn800Years+epochOffset)/timeUnitsIn1Day + 1
}

func newLunarYear(year int) *lunarYear {
	y := &lunarYear{year: year}

	y.horakhun = yearHorakhun(year)
	y.kammacapon = timeUnitsIn1Day - (year*daysIn800Years+epochOffset)%timeUnitsIn1Day
	y.uccapon = (uccaponConstant + y.horakhun) % apogeeRotationDays

	avoQuot := (y.horakhun*11 + 650) / 692
	y.avoman = (y.horakhun*11 + 650) % 692
	if y.avoman == 0 {
		y.avoman = 692
	}
	y.masaken = (avoQuot + y.horakhun) / 30
	y.tithi = (avoQuot + y.horakhun) % 30
	if y.avoman == 692 {
		y.tithi--
	}
	y.weekday = y.horakhun % 7

	hk1 := yearHorakhun(year + 1)
	tithi1 := ((hk1*11+650)/692 + hk1) % 30

	y.langsak = y.tithi
	if y.langsak < 1 {
		y.langsak = 1
	}
	n := y.langsak
	if n < 6 {
		n += 29
	}
	y.nyd = (y.weekday - n + 1 + 35) % 7

	y.leapday = y.kammacapon <= 207

	y.calType = typeNormal
	if y.tithi > 24 || y.tithi < 6 {
		y.calType = typeAthikamas
	}
	if y.tithi == 25 && tithi1 == 5 {
		y.calType = typeNormal
	}
	if (y.leapday && y.avoman <= 126) || (!y.leapday && y.avoman <= 137) {
		if y.calType != typeAthikamas {
			y.calType = typeAthikawan
		} else {
			y.calType = typeCoincide
		}
	}

	switch y.calType {
	case typeNormal:
		y.nextNyd = (y.nyd + 4) % 7
	case typeAthikawan:
		y.nextNyd = (y.nyd + 5) % 7
	default:
		y.nextNyd = (y.nyd + 6) % 7
	}
	y.caldays = calTypeDays[y.calType]

	return y
}

func calculateYear0(year int) *lunarYear {
	y := [5]*lunarYear{
		newLunarYear(year - 2),
		newLunarYear(year - 1),
		newLunarYear(year),
		newLunarYear(year + 1),
		newLunarYear(year + 2),
	}

	if y[2].tithi == 24 && y[3].tithi == 6 {
		for i := range y {
			y[i].calType = typeAthikamas
			y[i].nextNyd = (y[i].nextNyd + 2) % 7
		}
	}

	for _, i := range []int{1, 2, 3} {
		if y[i].calType == typeCoincide {
			j := -1
			if y[i].nyd == y[i-1].nextNyd {
				j = 1
			}
			y[i+j].calType = typeAthikawan
			y[i+j].nextNyd = (y[i+j].nextNyd + 1) % 7
		}
	}

	for _, i := range []int{1, 2, 3} {
		if y[i-1].nextNyd != y[i].nyd && y[i].nextNyd != y[i+1].nyd {
			y[i].offset = true
			y[i].langsak++
			y[i].nyd = (y[i].nyd + 6) % 7
			y[i].nextNyd = (y[i].nextNyd + 6) % 7
		}
	}

	for i := range y {
		if y[i].calType == typeCoincide {
			y[i].calType = typeAthikamas
		}
		y[i].caldays = calTypeDays[y[i].calType]
	}

	c := y[2]
	c.firstMonth = 'C'
	c.firstDay = c.langsak
	c.offsetDays = c.langsak
	min := 6
	if c.offset {
		min = 7
	}
	if c.offsetDays < min {
		c.firstMonth = 'V'
		c.firstDay = c.offsetDays
		c.offsetDays += 29
	}
	return c
}
