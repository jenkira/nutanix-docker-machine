# gh-pages

This branch hosts the published Helm chart repository for this project's
Rancher UI Extensions (`ui/pkg/nutanix`), served via GitHub Pages.

It is populated automatically by the "Publish UI Extension" GitHub Actions
workflow (`.github/workflows/publish-ui-extension.yml` on `main`), which
wraps Rancher's own `create-pr-build-extension-charts.yml` reusable
workflow. Each run opens a pull request against this branch with the built
`assets/`, `charts/`, `extensions/`, and `index.yaml` - nothing is pushed
here directly.

Do not edit this branch by hand; changes will be overwritten by the next
publish run.
