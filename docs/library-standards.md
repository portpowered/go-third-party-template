# Library standards

Apply these requirements when creating a library from this template:

1. Keep the public client, examples, README, and site independent of any consuming application. Put application adapters and rollout plans in the consuming repository.
2. Document supported operations, authentication, errors, and transport injection with examples that match the exported API. Distinguish verified behavior from synthetic examples and historical references.
3. Show Go version, CI, coverage, release, Go Reference, license, and documentation badges in the README. Replace every example repository value and point badges to live reports.
4. Build documentation in CI and publish the generated site to GitHub Pages. When API schemas exist, generate the API reference from checked-in schemas and verify it in CI.
5. Run offline build, lint, race, and replay checks before release. Keep real captures sanitized and separate from synthetic fixtures.

See [verification](verification.md), [client design](client-design.md), and
[website publishing](website.md) for the implementation details.
