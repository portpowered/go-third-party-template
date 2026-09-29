# Independent review: all-linters CI gate

Reviewed commit: `071e51e6370997e9174a86b1d4816c8feb6c849a`

Change under review: `a854780` and its depguard correction in `071e51e`.

## Verdict

**Pass for the new all-linters CI requirement.** The configuration enables all
linters, the blocking workflows pin and run golangci-lint v2 over the full
repository, and the exact commit's CI run passed. The depguard rule is
restrictive in `strict` mode. I found no source behavior regression in the
lint-driven refactors.

## Evidence

- `.golangci.yml` contains the literal `linters.default: all`. Depguard rule
  `Main` uses `list-mode: strict` and allows only `$gostd` and the module's own
  import path. The generated-file mode is `strict`; this repository contains
  no generated Go source files.
- Both CI and release workflows use `golangci/golangci-lint-action@v9`, pin
  golangci-lint to `v2.14.0`, set `only-new-issues: false`, and do not set
  `continue-on-error` or `--issues-exit-code=0`. The CI action has no package
  restriction. `Makefile` makes the full scope explicit with
  `golangci-lint run --timeout=5m ./...`.
- [GitHub Actions run 36514746866](https://github.com/portpowered/go-third-party-template/actions/runs/36514746866)
  is for this exact SHA. All six OS/Go matrix jobs and the API compatibility
  job passed, including every `Lint (all linters)` step.
- Local `make check` passed with Go 1.24.2 and golangci-lint 2.3.0 after setting
  a process-local Git safe-directory value. The local linter reported zero
  issues, with a deprecation warning for `wsl`; the exact CI run verifies the
  pinned 2.14.0 release.
- The changed client and command refactors preserve the existing request,
  response, API compatibility, and coverage behavior on inspection. The source
  `#nosec` annotations name the specific rules and give local reasons.

## Scaffold note

`tests/replay/fixtures/captured/` and `tests/replay/fixtures/synthetic/` contain
guidance README files only. There is no replay test or replay check in
`make check`, CI, or the release workflow. Add provider-specific replay checks
during migration before release; this scaffold gap is separate from the
all-linters CI requirement reviewed here.
