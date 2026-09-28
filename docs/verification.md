# Verification and fixture guidance

## Local checks

Format changed Go files, then run the checks before a release:

```sh
gofmt -w .
make check
```

make check builds every package and example, runs go vet, and runs tests with
the race detector. CI also checks formatting and module metadata on each pull
request.

CI reports public API changes against the pull request base. The report is
non-blocking so reviewers can decide whether an intentional break is appropriate.
The package list defaults to the module root and httpclient; pass -module and
-packages to tools/compatibility when the module path or public packages change.
Run make api-compatibility to compare against the most recent stable release tag.

## Fixture provenance

Keep real captures, synthetic examples, and historical references separate:

- Put observed and sanitized payloads under
  tests/replay/fixtures/captured/.
- Put hand-authored examples and generated test payloads under
  tests/replay/fixtures/synthetic/.
- Put old or superseded documentation under docs/reference/.

Do not call a fixture captured unless its source and collection date are known.
For each captured file, record the operation, capture date in UTC, source
category, redactions, and the behavior the fixture supports in a neighboring
provenance note. Remove access tokens, cookies, personal data, and account or
device identifiers before adding it. Keep credentials out of fixtures and CI.

Mark synthetic fixtures as synthetic in their filename or metadata. A synthetic
fixture can test an intended contract, but it is not evidence that a real service
behaves that way. Historical references are context and do not establish current
behavior.

The template includes guidance directories only; it does not contain captures
or claim any provider behavior has been verified.

## Continuous integration

The CI workflow builds, tests with the race detector, runs go vet, checks
gofmt, and confirms go mod tidy leaves module files unchanged. Keep checks
offline and deterministic. Live endpoint tests should be separate, opt-in, and
must not require credentials in pull request CI. Pull requests also build the
documentation site with strict checks. Pushes to `main` publish coverage
reports to the repository wiki and documentation to GitHub Pages; see
[website publishing](website.md).
