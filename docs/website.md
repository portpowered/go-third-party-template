# Documentation website

The Markdown files in `docs/` are built with MkDocs. Pull requests build the
site in strict mode, so broken navigation and links fail CI. Every push to
`main` runs the race-enabled test suite, builds the site, adds a Go coverage
report and badge data, and publishes the result to GitHub Pages.

## Customize the site

Before publishing a derived library, update `site_name`, `site_url`, `repo_url`,
`repo_name`, and `edit_uri` in `mkdocs.yml`. Replace the example repository and
module values in the README badges too. Keep the site name and navigation
specific to the library's supported API and documentation.

## Enable GitHub Pages

In the repository settings, open **Pages** and set the publishing source to
**GitHub Actions**. The `Documentation` workflow deploys the generated site on
every push to `main`, or when manually started from the `main` branch in the
Actions tab. The workflow exposes the published URL as its deployment
environment URL.

## Coverage report and badge

On each `main` push, the workflow runs race-enabled Go tests with coverage,
writes the HTML report to `site/coverage.html`, and creates `site/coverage.json`
for the Shields endpoint badge. Both files are included in the Pages
deployment, so the workflow needs no repository content write permission.
Coverage generation is skipped for pull requests; pull requests still build
the documentation in strict mode. The README badge links to the generated
report.

Run the same documentation build locally with:

```sh
python -m pip install -r requirements-docs.txt
mkdocs build --strict
```

The `Documentation` workflow contains the coverage report and badge generation
steps. Replace the example Pages URL in the README badge after customizing the
site settings.
