# Library standards

Apply these requirements when creating a library from this template:

1. Keep the public client, examples, README, and site independent of any consuming application. Put application adapters and rollout plans in the consuming repository.
2. Document supported operations, authentication, errors, and transport injection with examples that match the exported API. Add customer-facing operation guides for important workflows, and distinguish verified behavior from synthetic examples and historical references.
3. Show Go version, CI, coverage, release, Go Reference, license, and documentation badges in the README. Replace every example repository value and point badges to live reports.
4. Generate the API reference in CI with the shared Fumadocs action and publish it to GitHub Pages. Keep schemas checked in and reviewed, generate schema-defined Go models from them, and check for stale generated output in CI. Publish authored Markdown or MDX operation guides alongside the reference. Replace the template's clearly synthetic widget schema before presenting a provider API as supported behavior.
5. Run offline build, lint, race, and replay checks before release. Keep real captures sanitized and separate from synthetic fixtures.

See [verification](verification.md), [client design](client-design.md), and
[website publishing](website.md) for the implementation details.
