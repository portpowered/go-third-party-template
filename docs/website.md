# API documentation website

The `Documentation` workflow uses the shared
[Fumadocs API documentation action](https://github.com/portpowered/api-docs-website-github-action)
to generate the Pages site from checked-in API schemas. The template includes
`api/openapi.yaml` as a clearly synthetic widget example; it is not a contract
for any real service. Replace it with reviewed schemas for a derived library.
The generator builds the schema reference and publishes customer-facing
operation guides from `docs/guides/` alongside it. Replace the template's
guide starter with MDX pages for the derived library's supported operations.
Keep the route inventory complete: every outbound HTTP endpoint belongs in a
checked-in OpenAPI file, and non-HTTP wire exchanges need an appropriate
checked-in protocol schema. Generate wire definitions from those schemas and
check for drift in CI. Label implementation-derived contracts as such; a route
used by the client is not automatically a provider-verified specification.

## Customize the site

Update the workflow's `title` and schema paths when adapting the template. The
action derives the GitHub Pages project base path from the repository name, so
the workflow does not need a repository-specific path. Replace the example
repository and module values in the README badges before publishing.

## Add operation guides

Write customer-facing MDX pages under `docs/guides/`. The template
workflow passes `guides-directory: docs/guides` to action `v0.3.0`.
The action places them in a **Guides** section at `/docs/guides`, alongside the
schema-generated API reference. Add a `meta.json` file when the guide section
needs a custom title or page order. See the
[Python-to-Go migration playbook](python-provider-migration-playbook.md) for the
required guide content and review steps.

Keep contributor and release process notes in repository Markdown. Before
release, review every rendered guide and reference page for short, direct copy,
working navigation, and links to the matching endpoint. Remove duplicate
repository guide pages and repeated caveats while retaining evidence labels.

## Enable GitHub Pages

In the repository settings, open **Pages** and set the publishing source to
**GitHub Actions**. The `Documentation` workflow deploys the generated site on
pushes to `main`, or when manually started from `main`. Pull requests build the
Fumadocs site without deploying it. The workflow exposes the published URL as
its deployment environment URL.

## Coverage report and badge

On each `main` push, the workflow runs race-enabled Go tests with coverage,
writes the HTML report to `site/coverage.html`, and creates `site/coverage.json`
for the Shields endpoint badge. Both files are added after the API site is
built and included in the Pages deployment. Coverage generation is skipped for
pull requests. The README badge links to the generated report.

## Build locally

The shared action repository contains the same generator used in CI. Check out
the `v0.3.0` tag and run `npm ci` in that repository. From its root, run the
build with the source repository and output directory set explicitly:

```sh
API_DOCS_SOURCE=/absolute/path/to/your-service-go \
API_DOCS_OPENAPI=api/openapi.yaml \
API_DOCS_DISCOVER=false \
API_DOCS_GUIDES=docs/guides \
API_DOCS_OUTPUT=/absolute/path/to/your-service-go/site \
npm run build
```

Replace the paths for your local checkouts. The action's `base-path` input is
optional; on GitHub Actions it defaults to the repository name for project
Pages sites.
