# Verification and fixture guidance

## Local checks

Format changed Go files, then run the checks before a release:

```sh
gofmt -w .
make check
```

`make lint` runs `go vet ./...` and the pinned golangci-lint v2 release from
CI. `.golangci.yml` must keep the literal `linters.default: all`; CI runs every
repository package, treats findings as failures, and does not filter to new
issues. `make check` runs that lint target, builds every package and example,
and runs tests with the race detector. CI also checks formatting and module
metadata on each pull request.

Keep all linters enabled. If a finding requires an exception, scope it to the
specific rule and affected path and explain the reason in a config exclusion or
beside the source annotation. Never set the issue exit code to zero or continue
after a lint failure. A reviewer who did not implement the change confirms the
passing blocking CI result for the exact final commit before checklist item 5
is checked.

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

## Model audit

Run `go run ./tools/schemaexamples` (also the blocking `make schema-examples` target) to validate
every declared schema sample and all schema-bearing OpenAPI request, success, failure and AsyncAPI
message groups. Mark synthetic declarations `x-example-evidence: synthetic`; for AsyncAPI, put each `examples[]` sample's evidence on its enclosing Message Object and direct payload-schema sample evidence on that schema; the closed Example Object does not take evidence.
Local and relative references must resolve. See [canonical schema examples](schema-examples.md) for the ownership and
group rules. Keep `schema-examples` in the default check, pull-request CI and release workflow.

Validate each canonical request, response and event example against its owning schema.
Cover every known dispatch variant with complete envelopes and nested payloads, including
resource updates, relation changes, empty acknowledgements and representative errors.
Test identifier/payload correlation: a known command must reject another command's
parameters and malformed known input must not escape through a future-value branch.
Inspect the rendered reference and generated request snippets as well as schema output;
required nested fields and useful example values must be discoverable to a customer.

For interactive authorization, replay browser consent, callback state validation and PKCE
binding, token exchange and provider-required activation calls as one complete flow.
Keep browser, loopback listener and transport injectable; test denial, invalid or duplicate
callback parameters, wrong redirects, cancellation, repeated completion and cleanup.
Assert ordinary output excludes credentials and explicit export is required.

Review changed serialization code and enum conversions using the canonical schemas,
generated definitions and paired request/result tests. Include representative invalid
payloads and known future-value branches where they affect public behavior. Keep the
supported-operation list current and check generation for drift.

Do not build or require compiler/AST provenance engines, exhaustive ownership or alias
proofs, per-value lineage registries, proof receipts or proof caches. Unresolved symbolic
analysis is not a release blocker. Use established compiler, lint and schema tools and
targeted behavioral tests. A separate verification-engine project needs an explicit user
request; it must not emerge as a prerequisite for an otherwise working library.

## Schema-supplied links and rendered references

Keep the documentation workflow's schema inputs explicit and aligned with the
checker inputs. Set discover to false when the published set is explicit, and
include only the canonical assembled OpenAPI, AsyncAPI, GraphQL, and binding
sources that the workflow renders.

After rendering, run the schema-link checker against that exact output and
input set:

~~~sh
go run ./tools/schemalinks -site site -schemas api/openapi.yaml -base-path /<repository-name>
~~~

The checker resolves local YAML references and records URLs from externalDocs,
server URL fields, descriptions, and supported GraphQL schema bindings. Keep
source schema path, JSON pointer, line, and owning operation in the manifest.
Every published externalDocs destination must be reviewed against primary
provider content for the linked operation and path. Record the evidence class
and any limitations in tools/schemalinks/reviews.json. A general provider
guide does not establish an undocumented operation contract; label
implementation-derived routes as such and point to local implementation
evidence.

CI checks same-site Pages destinations against the rendered files and fails
when a schema-supplied internal target is missing. It records external and URL
template destinations for review without making remote requests, so transient
provider availability does not change the build result. Upload the full
rendered site, including schema-links.json, as a review artifact from the
same commit.

## Fixture sources

Each replay case must carry both sides of the exchange. For HTTP, keep the
expected outbound method, origin, escaped path, complete query multimap,
relevant headers, and request body with the response status, headers, and body.
For streams and sockets, keep ordered client and server frames with payloads.
The replay transport must validate a request before returning its paired
response, fail on unexpected or repeated calls, and assert that all expected
exchanges were consumed. Do not use a sequential or operation-name fallback
that returns a response after a request mismatch. Explicitly describe how
volatile or redacted values are matched; compare their structure or decoded
meaning rather than skipping them. Check this for every transport during the
independent review (library standard 15).
For HTTP, validate the effective authority, including a `Request.Host` override,
before returning the paired response. Reject unexpected URL user information,
opaque URLs, and malformed query strings. Include a request-identity mismatch
negative control, and keep secret values out of mismatch diagnostics.
Match full authentication forms and headers, including CSRF, OTP, and token exchanges.
For OAuth, verify state/callback and PKCE challenge/verifier bindings. Validate hardware
identity reuse and volatile formats before replacing values with placeholders. Match full
signaling handshakes and envelopes. Finish lifecycle assertions after the expected close or
teardown; neither a startup notification nor a terminal read timeout proves cleanup.

Keep real captures, synthetic examples, and historical references separate:

- Put observed and sanitized payloads under
  tests/replay/fixtures/captured/.
- Put hand-authored examples and generated test payloads under
  tests/replay/fixtures/synthetic/.
- Put old or superseded documentation under docs/reference/.

Do not call a fixture captured unless its source and collection date are known.
For each captured file, record the operation, capture date in UTC, source
category, redactions, and the behavior the fixture supports in a neighboring
source note. Remove access tokens, cookies, personal data, and account or
device identifiers before adding it. Keep credentials out of fixtures and CI.

Mark synthetic fixtures as synthetic in their filename or metadata. A synthetic
fixture can test an intended contract, but it is not evidence that a real service
behaves that way. Historical references are context and do not establish current
behavior.

The template includes guidance directories only; it does not contain captures
or claim any provider behavior has been verified.

## CLI verification

The standalone CLI must exercise the public SDK without a consuming application.
Check help, authentication, read/control commands, event cancellation and cleanup,
JSON output, and failure exit codes with credential-free paired replay tests.
Run lint, build, tests, formatting, and module checks in the CLI's separate module
as blocking CI jobs. Verify the published CLI can be installed from a clean
consumer module; document the commands and explicit credential export in MDX.
Replay the customer guide's ordered workflow: interactive authorization, persisted
account reuse, device enumeration and selection by ID, a useful read, an explicit
control, and logout. The customer must not construct request JSON, copy credentials
between files, or supply default wire parameters. Check invalid authorization,
ambiguous or unknown device selection, profile persistence and removal, secret
redaction, cancellation, and session cleanup. Inspect the rendered guide to confirm
the main instructions match CLI help and the tested command syntax; keep advanced
formats and protocol explanations outside that sequence.

## Independent standards review

Assign two reviewers who did not implement the migration to inspect the library
at the final implementation commit. Each reviewer checks every numbered library standard against
the exported API, schemas and generated files, client call sites, synthetic and
captured evidence, CI workflows and results, package layout, rendered Pages
site, and release/history record. Both reviewers write separate sections in one current review document
with an individual verdict and concrete evidence for each numbered standard,
including supported operations, generated contracts and paired replay results for item 4. Link it
from the library checklist and record each finding there. Keep the
independent-review item open until every finding and other checklist item is
resolved. Tracking an open finding is not sign-off. Rerun the affected checks
and have both reviewers verify the fixes at the final commit before marking the
item complete. Keep provider behavior that lacks documented account evidence
labeled implementation-derived.

Review schemas and generated models by API responsibility, representative serialization
paths, compatibility aliases and public import paths. Use generation drift and functional
tests as evidence; an exhaustive compiler-resolved model or primitive inventory is not
required.

For documentation sign-off, inspect every tracked document, including files excluded from
the site build. Record its audience and purpose, remove duplicate or obsolete internal
material, and check incoming links after deletion. Keep README content useful to callers;
put maintenance details in contributor material. Retain one current checklist and review
record rather than a chain of standalone historical reports.

## Continuous integration

CI runs compilation, pinned lint, formatting, module checks, reproducible schema
generation, race tests and paired replay tests for the SDK and every shipped CLI or
example module. Validate canonical examples with established schema tooling and enforce
the documented coverage floor. Keep a supported-operation list and test each operation
through its injectable transport, including active dependency socket exchanges.

Check tracked and untracked generated output for drift. Verify public consumer installation
and run the same practical checks on the exact release commit. Review concrete unsupported
routes or protocol mismatches as functional defects; do not make exhaustive static
provenance or custom proof-engine completion a CI requirement.

The CI workflow installs the exact golangci-lint v2 version recorded in the
workflow and runs all linters over the full repository as a blocking gate. It
also builds, tests with the race detector, runs go vet, checks gofmt, and
confirms go mod tidy leaves module files unchanged. Keep checks offline and
deterministic. Live endpoint tests should be separate, opt-in, and must not
require credentials in pull request CI. Pull requests also generate
the API reference from `api/openapi.yaml` with the shared Fumadocs action. The
checked-in widget schema is synthetic example data, not evidence of a real
provider contract. Every push to `main` runs the race-enabled tests with
coverage instrumentation, adds the HTML coverage report and Shields endpoint
JSON to the generated site, and deploys it to GitHub Pages; see
[website publishing](website.md).
