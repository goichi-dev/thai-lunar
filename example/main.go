package main

import (
	"fmt"
	"time"

	thailunar "github.com/goichi-dev/thai-lunar"
)

func main() {
	d := thailunar.FromGregorian(2026, time.May, 31)
	fmt.Println(d)
	fmt.Println(d.Format())
	fmt.Println(d.IsFullMoon(), d.IsWanPhra())
	fmt.Println(d.YearType.Thai())

	t, err := thailunar.ToGregorian(2569, 7, thailunar.Waxing, 15, false)
	fmt.Println(t.Format("2006-01-02"), err)

	for _, h := range []thailunar.Holiday{
		thailunar.Makha, thailunar.Visakha, thailunar.Asalha,
		thailunar.KhaoPhansa, thailunar.OkPhansa,
	} {
		t, _ := thailunar.HolidayDate(h, 2569)
		fmt.Printf("%-14s %s\n", h.Thai(), t.Format("2006-01-02"))
	}
}
