# Library standards

Apply these requirements when creating a library from this template:

1. Keep the public client, examples, README, and site independent of any consuming application. Put application adapters and rollout plans in the consuming repository.
2. Document supported operations, authentication, errors, and transport injection with examples that match the exported API. Add customer-facing operation guides for important workflows, and distinguish verified behavior from synthetic examples and historical references.
3. Show Go version, CI, coverage, release, Go Reference, license, and documentation badges in the README. Replace every example repository value and point badges to live reports.
4. Generate the API reference in CI with the shared Fumadocs action and publish it to GitHub Pages. Keep schemas checked in and reviewed, generate schema-defined Go models from them, and check for stale generated output in CI. Publish authored Markdown or MDX operation guides alongside the reference. Replace the template's clearly synthetic widget schema before presenting a provider API as supported behavior.
5. Run offline build, lint, race, and replay checks before release. Keep real captures sanitized and separate from synthetic fixtures.
6. Add deterministic synthetic request, response, error, and session fixtures for supported behavior. Measure coverage of non-generated production code by public and transport package and in combination; reach at least 80% combined coverage and target 90%. Report generated-code exclusions and remaining uncovered behavior rather than adding tests solely to raise a number.
7. Put the reusable public provider package under `pkg/<provider>` and its private wire types under an appropriate `internal` package. Keep examples, generated models, and transport packages in clear, separate locations. Verify public import paths from a separate consumer module.
8. Initialize clients from explicit options or a configuration value with sensible defaults and validation. Keep account credentials out of reusable client configuration when the client serves multiple accounts.
9. Keep the reusable client stateless with respect to accounts and connections. Return explicit session objects for login, event streams, sockets, RTC, or other stateful lifecycles; make ownership, close, errors, and token state visible to callers.
10. Allow callers to inject the transport at every network edge the library uses, including HTTP, HTTP/2, WebSocket, MQTT, and RTC signaling as applicable. Test request and response behavior through those seams without real credentials or network access.
11. Expose token exchange and refresh as explicit operations that return the current credentials to the caller. Do not silently refresh or retain updated tokens inside a reusable client; document caller storage and renewal responsibilities.

See [verification](verification.md), [client design](client-design.md), and
[website publishing](website.md) for the implementation details.
