# thai-lunar

[![CI](https://github.com/goichi-dev/thai-lunar/actions/workflows/ci.yml/badge.svg)](https://github.com/goichi-dev/thai-lunar/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/goichi-dev/thai-lunar.svg)](https://pkg.go.dev/github.com/goichi-dev/thai-lunar)
[![Go Report Card](https://goreportcard.com/badge/github.com/goichi-dev/thai-lunar)](https://goreportcard.com/report/github.com/goichi-dev/thai-lunar)

แปลงวันที่สุริยคติ (Gregorian) เป็นจันทรคติไทย — ข้างขึ้นข้างแรม เดือน ปีนักษัตร
วันพระ และวันสำคัญทางพุทธศาสนา — และแปลงกลับได้ด้วย

A Go library for the Thai lunar calendar (ปฏิทินจันทรคติไทย). It converts
Gregorian dates to ข้างขึ้น/ข้างแรม and back, using the Suriyayatra arithmetic
reckoning that the printed Thai royal calendar (ปฏิทินหลวง) is built on. No
dependencies, no data files, no network — just integer arithmetic.

## Install

```sh
go get github.com/goichi-dev/thai-lunar
```

## Usage

```go
package main

import (
	"fmt"
	"time"

	thailunar "github.com/goichi-dev/thai-lunar"
)

func main() {
	d := thailunar.FromGregorian(2026, time.May, 31)

	fmt.Println(d)           // ขึ้น ๑๕ ค่ำ เดือน ๗ ปี ๒๕๖๙
	fmt.Println(d.Format())  // waxing 15, month 7, BE 2569
	fmt.Println(d.IsFullMoon(), d.IsWanPhra()) // true true
	fmt.Println(d.YearType.Thai())             // อธิกมาส
}
```

รับวันที่เป็น string ได้ด้วย — ปี ค.ศ. หรือ พ.ศ. ก็ได้:

```go
d, err := thailunar.Parse("2026-09-09")   // "2569-09-09" ก็ได้ผลเดียวกัน
// แรม ๑๒ ค่ำ เดือน ๙ ปี ๒๕๖๙
```

`Parse` รับ `2006-01-02`, RFC 3339, `2006/01/02`, `02/01/2006` และ `02-01-2006`
ส่วนปีตั้งแต่ 2400 ขึ้นไปจะถูกอ่านเป็น พ.ศ. ให้อัตโนมัติ ถ้ามี `time.Time` อยู่แล้วใช้
`FromTime` ได้เลย

Convert a lunar date back to Gregorian:

```go
t, err := thailunar.ToGregorian(2569, 7, thailunar.Waxing, 15, false)
// 2026-05-31
```

And the Buddhist holy days, which fall on full moons and shift by a month in an
อธิกมาส year:

```go
for _, h := range []thailunar.Holiday{
	thailunar.Makha, thailunar.Visakha, thailunar.Asalha,
	thailunar.KhaoPhansa, thailunar.OkPhansa,
} {
	t, _ := thailunar.HolidayDate(h, 2569)
	fmt.Printf("%-14s %s\n", h.Thai(), t.Format("2006-01-02"))
}
// วันมาฆบูชา     2026-03-03
// วันวิสาขบูชา    2026-05-31
// วันอาสาฬหบูชา   2026-07-29
// วันเข้าพรรษา    2026-07-30
// วันออกพรรษา    2026-10-26
```

## How the calendar works

A Thai lunar year is one of three kinds, and the library reports which via
`YearTypeOf`:

| ชนิดปี | Type | Days | |
|---|---|---|---|
| ปกติมาส ปกติวาร | `Normal` | 354 | 12 months, alternating 29 and 30 days |
| อธิกวาร | `Athikawan` | 355 | เดือน 7 gains a day |
| อธิกมาส | `Athikamas` | 384 | เดือน 8 repeats (เดือน 8 หลัง) |

In an อธิกมาส year the holy days shift a month later — วันวิสาขบูชา moves from
เดือน 6 to เดือน 7, and วันอาสาฬหบูชา lands in เดือน 8 หลัง. `HolidayDate`
handles that for you.

## Two things worth knowing about the API

**`Year` is the civil พ.ศ., not the lunar year.** `Date.Year` rolls over on
1 January like the Gregorian year, which is what Thai civil usage expects.
`Date.LunarYear` is the lunar year, which rolls over at lunar new year in
roughly April, so the two differ between January and then.

**A lunar month can occur twice in one civil year.** เดือน 5 and เดือน 6 can
appear at both ends of a civil year, so `ToGregorian` returns the first match;
use `ToGregorianAll` when you want every one.

## Accuracy

Verified across 1950–2100: every one of the 54,787 days converts to a date that
follows its predecessor with no gaps or repeats, and round-trips back through
`ToGregorianAll`. The holy days and year types are checked against the published
ปฏิทินหลวง.

## Contributing

Bug reports about a wrong date are especially welcome — please include the source
you checked against. See [CONTRIBUTING.md](CONTRIBUTING.md).

## Credits

The Suriyayatra reckoning follows [pythaidate](https://github.com/markhollow/pythaidate)
by Mark Hollow (MIT). See [NOTICE.md](NOTICE.md).

## License

MIT — see [LICENSE](LICENSE).
