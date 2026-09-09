# Changelog

All notable changes to this project are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the project uses
[semantic versioning](https://semver.org/spec/v2.0.0.html). On `v0.x` the minor
version is bumped for breaking changes.

## [Unreleased]

### Added

- `Parse` and `ParseTime` accept a date as a string: `2006-01-02`, RFC 3339,
  `2006/01/02`, `02/01/2006` and `02-01-2006`. A year of 2400 or more is read as
  a Buddhist Era year, so `"2569-09-09"` and `"2026-09-09"` mean the same day.
- `Date` marshals to JSON as its Thai string form.

## [0.1.0] - 2026-09-09

First release.

### Added

- `FromTime` and `FromGregorian` convert a Gregorian date to the Thai lunar
  calendar, returning a `Date` with the day, phase (ข้างขึ้น/ข้างแรม), lunar
  month, whether it falls in เดือน 8 หลัง, and the year type.
- `ToGregorian` and `ToGregorianAll` convert back. `ToGregorianAll` returns every
  match, because เดือน 5 and เดือน 6 can each occur twice in one civil year.
- `YearTypeOf` reports whether a year is ปกติมาส ปกติวาร, อธิกวาร or อธิกมาส.
- `Date.IsFullMoon`, `Date.IsNewMoon` and `Date.IsWanPhra`.
- `HolidayDate` for วันมาฆบูชา, วันวิสาขบูชา, วันอาสาฬหบูชา, วันเข้าพรรษา and
  วันออกพรรษา, shifting a month later in อธิกมาส years.
- Thai and English formatting via `Date.String` and `Date.Format`, with Thai
  numerals.

### Notes

- `Date.Year` is the civil พ.ศ., which rolls over on 1 January. `Date.LunarYear`
  rolls over at lunar new year, so the two differ between January and roughly
  April.
- Conversion is verified across 1950-2100: all 54,787 days form an unbroken
  sequence and round-trip back. Holy days and year types are checked against the
  published ปฏิทินหลวง.

[Unreleased]: https://github.com/goichi-dev/thai-lunar/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/goichi-dev/thai-lunar/releases/tag/v0.1.0
