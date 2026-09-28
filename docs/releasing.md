# Release guide

## Before the first release

1. Replace github.com/example/your-service-go in go.mod, examples, tests, and
   the release workflow with the public module path.
2. Replace the service package name and SERVICE_* configuration names with
   names appropriate for the library. Move the synthetic root package into
   `pkg/<provider>` and update imports, examples, tests, and compatibility
   package settings before publishing.
3. Replace the example endpoint, resource, and wire response with verified
   provider behavior. Remove any operation the library does not support.
4. Review the included Apache-2.0 license and confirm it is appropriate for the
   new library before publishing.
5. Run the checks in verification.md and review the public README examples.
6. Replace the owner, repository, module, and site URL placeholders in the
   README badges, and update the API title and schema paths in
   `.github/workflows/docs.yml`.
7. Work through every current item in `docs/library-standards.md` as a checked
   and evidenced library-specific migration checklist, including independent
   review by someone who did not implement the migration. In particular, verify
   the `pkg/<provider>` import path, coverage threshold and exclusions,
   option-based initialization, client/session state boundary, every network
   injection seam, and caller-visible token refresh before tagging.
8. Run schema generation and drift, method/path/channel inventory, race tests,
   and the minimum coverage gate on the exact commit to tag. Review release-note
   copy and every link after moving guides into the Pages site.

## Tag and verify

Use semantic version tags of the form vMAJOR.MINOR.PATCH. Go modules use a
new major module path for v2 and later releases. Tag a reviewed commit and push
the tag to start .github/workflows/release.yml.

The release workflow compares the public API with the prior stable tag, applies
the version policy, runs tests and static checks, downloads the tagged module
through the public Go module proxy in a separate temporary module, compiles a
small consumer against the configured public packages, then publishes GitHub release
notes. It must rerun generation/drift, endpoint inventory, and coverage gates on
the tag itself before publication. Set the workflow PUBLIC_MODULE to the same path as go.mod and
PUBLIC_PACKAGES to the exported packages before creating the first tag.

Breaking API changes are allowed when the major version increases. Before v1,
they are also allowed when the minor version increases. A v0 patch release and
a stable minor release must preserve compatibility.

Verify the published package independently with:

```sh
mkdir module-check
cd module-check
go mod init example.com/module-check
go get github.com/example/your-service-go@v0.1.0
go list -m github.com/example/your-service-go
```

Replace the example module path and version with the values for the release.

The documentation workflow generates the API reference from checked-in schemas
with the shared Fumadocs action and publishes it to GitHub Pages. Enable GitHub
Pages with GitHub Actions as the publishing source in repository settings
before the first deployment. See [website publishing](website.md) for setup.
