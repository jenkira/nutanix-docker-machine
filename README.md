# Nutanix Rancher Node Driver

This repository contains the Rancher Node Driver for Nutanix. Nutanix Node driver are used to provision hosts on Nutanix Enterprise Cloud, which Rancher uses to launch and manage Kubernetes clusters.


---

[![Go Report Card](https://goreportcard.com/badge/github.com/jenkira/nutanix-docker-machine)](https://goreportcard.com/report/github.com/jenkira/nutanix-docker-machine)
![CI](https://github.com/jenkira/nutanix-docker-machine/actions/workflows/integration.yml/badge.svg)
![Release](https://github.com/jenkira/nutanix-docker-machine/actions/workflows/release.yml/badge.svg)

[![release](https://img.shields.io/github/release-pre/jenkira/nutanix-docker-machine.svg)](https://github.com/jenkira/nutanix-docker-machine/releases)
[![License](https://img.shields.io/badge/License-MPL%202.0-blue.svg)](https://github.com/jenkira/nutanix-docker-machine/blob/main/LICENSE)
![Proudly written in Golang](https://img.shields.io/badge/written%20in-Golang-92d1e7.svg)
[![Releases](https://img.shields.io/github/downloads/jenkira/nutanix-docker-machine/total.svg)](https://github.com/jenkira/nutanix-docker-machine/releases)

---

## Features


- Configure Prism Central and corresponding user to talk to Nutanix platform
- Define target cluster to deploy VM
- Ability to set a custom name for the newly created VM
- Ability to select VM's Main Memory in Megabytes
- Ability to select VM's vCPU count
- Ability to set the number of cores per vCPU
- Ability to specify the network(s) of the VM (Classic or VPC)
- Ability to specify the template disk in the VM by image name and modify his size (increase only)
- Ability to specify categories to applied to the VM ( flow, leap, ...)
- Ability to add one additional disk by specifying disk-size and storage-container
- Enable passthrough the host's CPU features to the newly created VM
- Define a Cloud-init user-data to send to the newly created VM
- Project support
- Serial Port support
- Boot type selection : Legacy or UEFI 
- GPU support
- Prism Central Service Accounts support
- Windows guest OS support (via cloudbase-init)


## Installation


If you want to use Nutanix Node Driver, you need add it in order to start using them to create node templates and eventually node pools for your Kubernetes cluster.

1. From the Home view, choose *Cluster Management* > *Drivers* in the navigation bar. From the Drivers page, select the *Node Drivers* tab.
2. Click *Add Node Driver*.
3. Complete the Add Node Driver form. Then click Create.

    - *Download URL*: `https://github.com/jenkira/nutanix-docker-machine/releases/download/v3.10.0/docker-machine-driver-nutanix`  
    - *Checksum*: `99d9cd74870b059721acb0b83bc45c381a831196a37a99340769a68c96c04392`  

    This fork does not host the legacy `component.js` Custom UI - it's RKE1-only
    tooling whose official scaffold Rancher archived in 2024 (see
    [Windows Node Support](#windows-node-support)). Leave *Custom UI URL* and
    *Whitelist Domains* blank for a driver-only install using Rancher's generic
    form.

    For RKE2/K3s node pools, install the companion [Rancher UI Extension](./ui)
    instead - it gives the driver proper Cloud Credential and Machine Config
    forms (OS picker, static IP, etc.) instead of one plain input per flag. It's
    published as a Helm chart repository via GitHub Pages:

    - In Rancher, go to *☰ > Extensions*, open the *⋮* menu (top right) >
      *Manage Repositories* > *Create*.
    - Set *Index URL* to `https://jenkira.github.io/nutanix-docker-machine/`.
    - Once the repository syncs, install **Nutanix Node Driver UI** from the
      *Available* tab on the Extensions page.

    See [`ui/`](./ui) for the extension's source, or
    [`ui/pkg/nutanix`](./ui/pkg/nutanix) for its own README (build/dev
    instructions, current scaffold status).

<img width="948" height="474" alt="image" src="https://github.com/user-attachments/assets/e2383b37-bb55-4242-ae6e-fec392da9577" />







4. Wait for the driver to become "Active"
5. Go to *RKE1 Configuration > Node Templates*, you can create a Nutanix Template using Rancher's generic driver form (all `nutanix-*` args below become fields).

![image](https://github.com/nutanix/docker-machine/assets/180613/8c56a022-ad6b-406b-80e6-10c5673c0d9e)



## Driver Args

| Arg                          | Description                                                                                      | Required | Default                                   |
|------------------------------|:-------------------------------------------------------------------------------------------------|:---------|-------------------------------------------|
| `nutanix-endpoint`           | The hostname/ip-address of the Prism Central                                                     | yes      |                                           |
| `nutanix-port`               | The port to connect to Prism Central                                                             | no       | 9440                                      |
| `nutanix-username`           | The username of the nutanix management account                                                   | yes      |                                           |
| `nutanix-password`           | The password of the nutanix management account                                                   | yes      |                                           |
| `nutanix-insecure`           | Set to true to force SSL insecure connection                                                     | no       | false                                     |
| `nutanix-cluster`            | The name (case sensitive) or UUID of the cluster to deploy the VM on                             | yes      |                                           |
| `nutanix-boot-type`          | The boot type of the VM (legacy or uefi)                                                         | no       | legacy                                    |
| `nutanix-vm-mem`             | The amount of RAM of the newly created VM (MB)                                                   | no       | 2 GB                                      |
| `nutanix-vm-cpus`            | The number of cpus in the newly created VM (core)                                                | no       | 2                                         |
| `nutanix-vm-cores`           | The number of cores per vCPU                                                                     | no       | 1                                         |
| `nutanix-vm-network`         | The network(s) to which the VM is attached to ( name or UUID )                                   | yes      |                                           |
| `nutanix-vm-ip`              | Request a specific static IPv4 for the VM's first NIC (needs a Nutanix-managed/IPAM subnet)      | no       |                                           |
| `nutanix-vm-image`           | The name of the Disk Image template we use for the newly created VM (must support cloud-init)    | yes      |                                           |
| `nutanix-vm-image-size`      | The new size of the Image we use as a template (in GiB)                                          | no       |                                           |
| `nutanix-vm-categories`      | The name of the categories who will be applied to the newly created VM                           | no       |                                           |
| `nutanix-vm-gpu`             | The list of GPU device names to attach to the newly created VM (can be specified multiple times) | no       |                                           |
| `nutanix-project`            | The name of the project where deploy the VM (default if empty)                                   | no       | default                                   |
| `nutanix-disk-size`          | The size of the additional disk to add to the VM (in GiB)                                        | no       |                                           |
| `nutanix-storage-container`  | The storage container UUID of the additional disk to add to the VM                               | no       |                                           |
| `nutanix-cloud-init`         | Cloud-init to provide to the VM (will be patched with rancher root user)                         | no       |                                           |
| `nutanix-vm-os`              | The guest OS family of the VM (linux or windows); Windows uses cloudbase-init for cloud-init     | no       | linux                                     |
| `nutanix-vm-ssh-user`        | SSH username to provision via cloud-init on the newly created VM                                 | no       | root (linux) / Administrator (windows)    |
| `nutanix-vm-cpu-passthrough` | Enable passthrough the host's CPU features to the newly created VM                               | no       | false                                     |
| `nutanix-vm-serial-port`     | Attach a serial port to the newly created VM                                                     | no       | false                                     |
| `nutanix-vm-description`     | The description of the newly created VM                                                          | no       | VM created by Nutanix Rancher Node Driver |



## Project support

Starting `v3.3.0` the Rancher Node driver implements Nutanix Project support. The  prerequisite needed to be able to use the Rancher Node Driver is the following:
- Target cluster and network available in the Project
- Role with the following recommended permission:
  - VM Full Access
  - Cluster View Access
  - Image View Only

## Service Accounts support

Starting `v3.9.0` the Rancher Node Driver support Prism Central Service Accounts. 
To use a Service Account, you need to provide `X-ntnx-api-key` as the user name and the corresponding API Key as the password.

## Rancher HA: Keeping the Driver Binary in Sync

If your Rancher server runs multiple replicas, be aware of a known Rancher bug ([rancher/rancher#42128](https://github.com/rancher/rancher/issues/42128), [#42302](https://github.com/rancher/rancher/issues/42302)): a custom node driver binary like this one can end up downloaded onto only *some* replicas' `/usr/share/rancher/ui/assets/`, not all of them. Because Rancher's node-driver registration and dynamic CRD generation (the `NutanixConfig`/`NutanixMachine`/`NutanixMachineTemplate` resources backing RKE2/K3s node pools) can run on whichever replica currently holds leadership, an inconsistent set of replicas can leave those resources unable to settle - showing up as recurring `failed to sync cache` errors and general Rancher/Fleet management-cluster instability, especially right after adding or updating this driver.

This is a bug in Rancher itself, not something this driver's code can fix. [`scripts/sync-rancher-driver-binary.sh`](./scripts/sync-rancher-driver-binary.sh) is the operational workaround: it checks every Rancher server pod for a consistent copy of the driver binary, and with `--fix`, copies it from whichever pod has it (or downloads a fresh, checksum-verified copy from this repo's releases if no pod has it) to every pod that doesn't.

```bash
# Check only - reports which pods are missing the binary, changes nothing
./scripts/sync-rancher-driver-binary.sh

# Repair - syncs the binary across every Rancher pod
./scripts/sync-rancher-driver-binary.sh --fix
```

Requires `kubectl` with a working context against the Rancher management cluster. Run `./scripts/sync-rancher-driver-binary.sh --help` for the full option list - namespace, pod selector, and binary name are all overridable if your deployment doesn't use Rancher's defaults.

## Cluster Name Ambiguity

`nutanix-cluster` accepts either a name or a UUID. Nutanix does not enforce unique cluster names within a single Prism Central, so if two clusters share the configured name, the driver can't guess which one you meant - it fails with an error listing the conflicting UUIDs rather than silently picking one. If you hit this, set `nutanix-cluster` to one of the UUIDs from that error message instead of the name.

This matters more than a normal validation error: without it, RKE2/K3s machine provisioning retries the same unresolvable failure indefinitely, which can degrade the whole Rancher/Fleet management cluster over time (see [rancher/rancher#47493](https://github.com/rancher/rancher/issues/47493)). Switching to the UUID is the only real fix - there's no way to disambiguate identically-named clusters by name alone.

## GPU support

The Rancher Node Driver supports attaching GPU devices to VMs. To use GPUs:
- Specify GPU devices by their name using the `nutanix-vm-gpu` parameter
- Multiple GPUs can be attached by specifying the parameter multiple times
- Only UNUSED GPUs from the target Prism Element cluster will be selected
- GPU names must match exactly with the GPU names available in the cluster
- The driver will search for available GPUs across all hosts in the specified cluster

## Static IP Assignment

By default the driver waits for whatever IP the target subnet's DHCP hands out (see `nutanix-vm-network`). If you need a specific, predictable address instead - the same problem VMware's vSphere integrations usually solve by pushing a static IP into the guest via `guestinfo` - Nutanix solves it differently: **AHV's own IPAM, not guest-side configuration.**

To use it:
- The subnet named/UUID'd by `nutanix-vm-network` (the *first* one, if you pass it multiple times for a multi-NIC VM) must be a **Managed** subnet in Prism Central - i.e. IP Address Management (IPAM) enabled, with an IP pool configured. Unmanaged subnets (external/customer-run DHCP) have no pool for the driver's request to land in.
- Set `nutanix-vm-ip` to the address you want, e.g. `--nutanix-vm-ip 10.0.0.50`. It must fall inside that subnet's configured pool.
- If the subnet isn't Managed, or the address isn't available (outside the pool, already leased, etc.), VM creation fails with the error Nutanix itself reports - the driver doesn't pre-validate this ahead of time, it's a straight passthrough to the API.

**No cloudbase-init or cloud-init changes are needed, and there are no variables to reference.** This is not the same mechanism as VMware's guestinfo-delivered static networking (where cloud-init/cloudbase-init writes a static IP/netmask/gateway into the guest's own network config, replacing DHCP). Nutanix's `IPEndpointList`/`ASSIGNED` request just tells AHV's *internal* DHCP server which address to hand out to that NIC's MAC address - the guest still runs an ordinary DHCP client and receives the reserved address exactly like it would any other DHCP lease. Any image already prepared for this driver (cloud-init on Linux, cloudbase-init on Windows per below) works as-is; nothing in `nutanix-cloud-init` needs to change to pick this up.

## Windows Node Support

The Rancher Node Driver can provision Windows Server VMs as Rancher-managed nodes, using the same `GuestCustomization` cloud-init mechanism as Linux VMs. This mirrors how Rancher's own vSphere node driver supports Windows: rather than using Sysprep, the Windows VM template runs [cloudbase-init](https://cloudbase-init.readthedocs.io/) to consume the same `#cloud-config` payload cloud-init consumes on Linux, so existing cloudbase-init-based templates and cloud-init scripts can be reused unchanged.

To use it, set:
- `nutanix-vm-os=windows`
- `nutanix-vm-ssh-user` if the target account isn't the built-in `Administrator` (the default when unset)

Your existing `nutanix-cloud-init` payload can be passed through as-is; the driver only changes how it injects the SSH-enabled user into the `users` list (Windows accounts are enabled via `inactive: false`, since cloudbase-init leaves local accounts, including the built-in Administrator, disabled by default).

The Windows image referenced by `nutanix-vm-image` must be prepared before use, the same way a Linux image must already support cloud-init:
- [cloudbase-init](https://cloudbase-init.readthedocs.io/) installed, configured with a datasource compatible with AHV's cloud-init guest customization (the ConfigDrive datasource is the most likely match — verify against your template before relying on it in production)
- Win32-OpenSSH Server feature installed (the driver does not install it for you)
- Nutanix VirtIO drivers installed
- Generalized with `sysprep /generalize /oobe /shutdown` (or the modern in-box equivalent) so the template boots to OOBE, ready for per-VM customization

Windows worker nodes are a Rancher/RKE2 capability, not something this driver enforces: they're only supported on RKE2/K3s clusters (not the legacy RKE1 flow), the control plane must remain Linux-only with Windows nodes as workers, and Calico is required as the CNI. See [Rancher's Windows cluster documentation](https://ranchermanager.docs.rancher.com/) for details.

## Development

### Build Instructions

build linux/amd64 binary => `make` 
build local binary => `make local`

## History

* v1 is the original Nutanix docker machine driver that connect to Prism Element
* v2.x add Rancher 2.0 support
* v3.x is a rewrite of the driver that connect to Prism Central

## Support
This code is developed in the open with input from the community through issues and PRs. A Nutanix engineering team serves as the maintainer. Documentation is available in the project repository. Issues and enhancement requests can be submitted in the [Issues tab of this repository](../../issues). Please search for and review the existing open issues before submitting a new issue.

## License
Copyright 2022-2023 Nutanix, Inc.

The project is released under Mozilla Public License Version 2.0.
