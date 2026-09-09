# Contributing to thai-lunar

Thanks for taking the time to contribute.

## Getting started

```bash
git clone https://github.com/goichi-dev/thai-lunar.git
cd thai-lunar
go build ./...
go run ./example
```

The library requires Go 1.26 or later and has no dependencies. Please keep it
that way: the whole point is that a calendar conversion should not pull in a
module graph.

## Quality gate

Run all four before opening a pull request. CI runs the same checks.

```bash
gofmt -w .
go vet ./...
staticcheck ./...
go test ./...
```

## Correctness comes from the calendar, not from the code

A date conversion that compiles, vets and passes staticcheck can still be wrong
on every single day of the year. The tests are the real gate here, and two of
them matter more than the rest:

- `TestRoundTrip` walks every day from 1950 to 2100 and asserts that each date
  follows its predecessor with no gap, no repeat and no jump. It catches almost
  any arithmetic mistake immediately.
- `TestHolidays` checks the Buddhist holy days against dates published in the
  ปฏิทินหลวง.

**If you change anything in `year.go`, `csdate.go` or `constants.go`, say in the
pull request which published dates you checked against.** "The tests pass" is
necessary but not sufficient — an anchor that is wrong in the test is wrong
everywhere.

When you add a case, prefer a date you can point to a source for: a government
holiday announcement, a printed calendar, or the Royal Institute's tables.

## Reporting a wrong date

This is the most useful kind of issue for this project. Please include:

- the Gregorian date,
- what the library returned,
- what the printed calendar says, **and where you read it**.

A disagreement with an online converter is a starting point, not a bug report —
several of them use the Chinese or Vietnamese lunar calendar, which genuinely
differs from the Thai one and will not agree.

## Scope

The library covers the Thai lunar calendar (ปฏิทินจันทรคติไทย). Requests to add
the Chinese, Vietnamese or Burmese lunar calendars belong in their own package;
they are different systems, not options on this one.

## Commit messages

Short, imperative, capitalised, no prefix:

```
Fix athikamas month 8 numbering
Add Ok Phansa holiday
```

## Pull requests

- Keep changes focused; unrelated refactors make review harder.
- Explain the behaviour before and after, not just the diff.
- Note any breaking change explicitly — the project is on `v0.x`, so breaking
  changes are allowed, but they must be called out for the changelog.
- New exported identifiers need a doc comment. They become the public reference
  on pkg.go.dev.

## Comments and language

Code, comments and commit messages are written in English. Thai is used where it
belongs: calendar terms, user-facing strings and documentation examples, since
that is what the output actually looks like.
