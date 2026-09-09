# Security Policy

## Supported versions

thai-lunar is on `v0.x`. Only the most recent release receives fixes.

## Reporting a vulnerability

**Please do not open a public issue for a security problem.**

Report it privately through GitHub's
[private vulnerability reporting](https://github.com/goichi-dev/thai-lunar/security/advisories/new)
on this repository. Include the affected version or commit, a description of the
impact, and the smallest program that reproduces it.

You can expect an acknowledgement within seven days, and an assessment with a
planned fix or an explanation within thirty days.

## Scope

This is a pure computation library. It has no dependencies, opens no files, makes
no network calls and starts no goroutines, so the usual categories do not apply.
What is in scope:

- a panic reachable from any exported function, for any input, including
  nonsensical years and out-of-range months or days,
- an unbounded loop or allocation triggered by an extreme input, such as a year
  far outside the useful range.

Both are ordinary bugs as much as security ones. If you find either, a normal
issue is fine unless you believe it is being exploited somewhere.

## A wrong date is not a vulnerability

An incorrect conversion is a correctness bug — please
[open an issue](https://github.com/goichi-dev/thai-lunar/issues) with the source
you checked against. It is treated seriously, but it does not need private
disclosure.

That said, do not use this library as the sole basis for a decision with legal or
financial consequences without checking the official ปฏิทินหลวง. The Thai lunar
calendar is set by arithmetic reckoning, and the published calendar is the
authority.
