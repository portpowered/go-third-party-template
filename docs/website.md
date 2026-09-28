# Documentation website

The Markdown files in `docs/` are built with MkDocs. Pull requests build the
site in strict mode, so broken navigation and links fail CI. Pushes to `main`
that change documentation configuration or content build the site and publish
it to GitHub Pages.

## Customize the site

Before publishing a derived library, update `site_name`, `site_url`, `repo_url`,
`repo_name`, and `edit_uri` in `mkdocs.yml`. Replace the example repository and
module values in the README badges too. Keep the site name and navigation
specific to the library's supported API and documentation.

## Enable GitHub Pages

In the repository settings, open **Pages** and set the publishing source to
**GitHub Actions**. The `Documentation` workflow then deploys the generated
site when documentation changes reach `main`, or when manually started from
the `main` branch in the Actions tab. The workflow exposes the published URL as
its deployment environment URL.

## Publish coverage reports

The CI workflow measures Go test coverage on pushes to `main` and publishes an
HTML report and badge to the repository wiki. Enable the repository wiki and
allow the workflow to write repository contents. The README's coverage badge
links to that generated report. Coverage reporting is skipped for pull
requests, so unmerged changes do not update the public report.

Run the same documentation build locally with:

```sh
python -m pip install -r requirements-docs.txt
mkdocs build --strict
```
