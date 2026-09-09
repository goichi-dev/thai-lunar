package thailunar

const (
	daysIn800Years     = 292207
	timeUnitsIn1Day    = 800
	epochOffset        = 373
	uccaponConstant    = 2611
	apogeeRotationDays = 3232

	csJulianDayOffset = 1954167
)

const (
	typeNormal    = 'A'
	typeAthikawan = 'B'
	typeAthikamas = 'C'
	typeCoincide  = 'c'
)

var calTypeDays = map[byte]int{
	typeNormal:    354,
	typeAthikawan: 355,
	typeAthikamas: 384,
	typeCoincide:  384,
}

var monthCumulativeDays = map[byte][]int{
	typeNormal:    {0, 29, 59, 88, 118, 147, 177, 206, 236, 265, 295, 324, 354, 383},
	typeAthikawan: {0, 29, 59, 89, 119, 148, 178, 207, 237, 266, 296, 325, 355, 384},
	typeAthikamas: {0, 29, 59, 88, 118, 148, 177, 207, 236, 266, 295, 325, 354, 384},
}

var lunarMonths = []int{0, 5, 6, 7, 8, 9, 10, 11, 12, 1, 2, 3, 4, 8, 88, 15, 16}

var monthPositionAB = []int{0, 5, 6, 7, 8, 9, 10, 11, 12, 1, 2, 3, 4, 15, 16}

var monthPositionC = []int{0, 5, 6, 7, 8, 88, 9, 10, 11, 12, 1, 2, 3, 4, 15, 16}
