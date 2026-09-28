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

For extracted libraries, report coverage for each public and transport package
and for non-generated production code in combination. Exclude generated files
explicitly rather than counting them as uncovered or using generated code to
inflate the percentage. Reach at least 80% combined coverage and target 90%.
Enforce the 80% floor in the extracted library's CI. Use deterministic
synthetic inputs to cover success, provider errors, invalid responses,
cancellation, token rotation, and session lifecycle where applicable.
Review uncovered behavior before adding tests; a percentage alone is not a
behavioral sign-off. Keep live integration results separate from synthetic
unit and replay coverage.

After moving the synthetic scaffold into `pkg/<provider>`, run the reusable
coverage gate from the module root:

```sh
go test -coverpkg=./pkg/... -coverprofile=coverage.out ./pkg/...
go run ./tools/coverage -profile coverage.out -min 80
```

The gate reports each package and a statement-weighted combined percentage.
It excludes generated files marked `Code generated ... DO NOT EDIT` and files
named `*.gen.go`; review the exclusions alongside the result.
Use `-percent-only -min 0` when producing a coverage badge from the same
non-generated measurement. Pass `-filtered-profile coverage.filtered.out` to
write a profile for an HTML report that uses the same exclusions.

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

## Independent standards review

Assign a reviewer who did not implement the migration to inspect the library
at a specific commit. The reviewer checks each numbered library standard against
the exported API, schemas and generated files, client call sites, synthetic and
captured evidence, CI workflows and results, package layout, rendered Pages
site, and release/history record. Record the reviewer, commit, evidence, and
each finding in the library checklist. Keep the independent-review item open
until every finding and other checklist item is resolved. Tracking an open
finding is not sign-off. Rerun the affected checks and have the reviewer
verify the fixes at the final commit before marking the item complete. Keep
provider behavior that lacks documented account evidence labeled
implementation-derived.

## Continuous integration

Before signing off schema generation, list every outbound HTTP method and path,
GraphQL operation, event channel, stream frame, signaling route, and other wire
exchange the client can initiate or consume. Match each inventory row to a
checked-in schema entry, its generated endpoint definition and wire types, and
the client call site that uses them. Include query/header parameter names,
nested event and protocol-specific payload properties, and JSON or other
structures embedded in strings or encrypted wrappers, not just outer envelopes. Mark any
implementation-derived entry as such. An endpoint with only a schema entry or
only a generated reference page is incomplete. Add a CI check that detects
new method-and-path pairs, channels, and call sites without schema entries and
fails when regeneration changes checked-in output or creates an untracked
generated file. A plain `git diff` does not detect new untracked model files;
check tracked drift and untracked generated paths. Add negative tests that
introduce an unschematized route, change a method while retaining its path,
and add an unschematized channel; each must fail. Repeat this inventory for
each library; partial provider catalogues are not a sign-off.

Run the same schema generation, drift, route/channel inventory, and minimum
coverage gates on the exact release tag commit. Passing a prior `main` run is
not evidence that a different tag commit meets the release standard.

The CI workflow builds, tests with the race detector, runs go vet, checks
gofmt, and confirms go mod tidy leaves module files unchanged. Keep checks
offline and deterministic. Live endpoint tests should be separate, opt-in, and
must not require credentials in pull request CI. Pull requests also generate
the API reference from `api/openapi.yaml` with the shared Fumadocs action. The
checked-in widget schema is synthetic example data, not evidence of a real
provider contract. Every push to `main` runs the race-enabled tests with
coverage instrumentation, adds the HTML coverage report and Shields endpoint
JSON to the generated site, and deploys it to GitHub Pages; see
[website publishing](website.md).
