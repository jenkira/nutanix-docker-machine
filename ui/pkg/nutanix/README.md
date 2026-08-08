# Nutanix Node Driver UI

A Rancher UI Extension providing the **Cloud Credential** and **Machine Config**
components for the `nutanix` node driver, used when creating RKE2/K3s node
pools from Cluster Management. This is what makes the driver's fields -
including `nutanix-vm-os` (Windows guest support) - a proper form instead of
Rancher's generic fallback (a plain input per flag).

This replaces the legacy `component.js`/"Custom UI URL" mechanism the rest of
this repo's README documents for RKE1 Node Templates. That mechanism is
Ember-based, its official scaffold (`rancher/ui-driver-skel`) was archived by
Rancher in 2024, and RKE1 never supported Windows worker nodes in the first
place - so it isn't a usable path for what this package is for. RKE2/K3s node
pools use this newer, Vue-based Extensions framework instead.

## Status: scaffold, not yet verified against a live Rancher instance

This was built directly off Rancher's own current reference example
(`rancher/ui-plugin-examples`, `pkg/node-driver`, an OpenStack Cloud
Credential + Machine Config pair) - itself labeled "an in progress example"
by Rancher. The package/plugin structure (`index.ts`, `package.json` shape,
`machine-config/<driver>.vue` and `cloud-credential/<driver>.vue` naming,
`importTypes`/`plugin.metadata` registration) is copied from that verified
source.

What is **not** independently verified here, because `extensions.rancher.io`
was unreachable from this environment when this was written:

- **Field name mapping.** Each form field assumes a driver flag
  `nutanix-vm-os` becomes the camelCased `vmOs` key on the machine config
  resource (`nutanix-endpoint` -> `endpoint` on the cloud credential's
  `decodedData`, etc.). This matches Rancher's documented general convention
  for turning docker-machine flags into API schema fields, but confirm it
  against your target Rancher version before shipping.
- **Exact component prop/event contract.** `value`/`credentialId`/`mode`
  props and direct mutation of `value.*` follow Rancher Dashboard's
  established shared-component conventions, but the OpenStack reference
  component has more validation/event wiring (e.g. a `validationChanged`
  emit on the cloud credential form) that this scaffold only stubs out
  (`validate()` is unused by anything yet - wire it up per whatever your
  target Rancher version's node-driver docs specify).
- **Live data.** `Network(s)`, `Categories`, and `GPU device names` use
  taggable multi-selects instead of plain text (a real improvement for
  those array-typed flags), but nothing here calls the Prism Central API to
  populate them, a cluster/image/network picker, the way the OpenStack
  example calls out to OpenStack's API through Rancher's `/meta/proxy/`
  passthrough. That's the natural next enhancement, not included here.

## Developing

From `ui/`:

```
yarn install
API=https://<your-rancher-server> yarn dev
```

Then add this driver in Rancher (Cluster Management > Drivers > Node
Drivers) with the dev server's URL as the Custom UI URL, per Rancher's
Extensions docs.

## Building

```
yarn build-pkg nutanix
```

Output lands in `dist-pkg/nutanix/`.

## Published

Tagged builds are packaged into a Helm chart and published to this repo's
`gh-pages` branch by `.github/workflows/publish-ui-extension.yml`, served via
GitHub Pages at `https://jenkira.github.io/nutanix-docker-machine/`. Add that
URL as a repository in Rancher (*☰ > Extensions > ⋮ > Manage Repositories*)
to install without a local dev server - see the main [README](../../../README.md#installation)
for the full steps.
