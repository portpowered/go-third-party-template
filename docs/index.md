# Go Service Client

This site documents a Go client library for one service. Replace the example
module, repository, and service details before publishing a derived library.

Use the guides to shape the public client API, migrate an existing Python
provider client, verify protocol behavior, and prepare a Go module release.

## Guides

- [Client design](client-design.md) describes the provider-specific API,
  optional sessions, and transport ownership.
- [Python-to-Go migration playbook](python-provider-migration-playbook.md)
  provides a task-by-task conversion checklist.
- [Verification](verification.md) defines fixture provenance and CI checks.
- [Releasing](releasing.md) covers module versioning and release verification.
- [Website publishing](website.md) explains how this site is built and deployed.

The generated coverage report is published to the repository wiki after pushes
to `main`. Find it through the coverage badge in the
[repository README](https://github.com/example/your-service-go).
