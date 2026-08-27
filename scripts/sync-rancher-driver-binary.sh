#!/usr/bin/env bash
#
# Operational workaround for a known Rancher HA bug (rancher/rancher#42128,
# #42302) where a custom node driver binary only gets downloaded onto some
# replicas of a multi-replica Rancher server deployment, leaving
# /usr/share/rancher/ui/assets/ inconsistent across pods. Rancher's own
# provisioning/dynamicschema controllers run on whichever pod holds
# leadership at a given moment, so an inconsistent set of replicas can
# produce exactly the symptom this repo's README documents under "Rancher
# HA: Keeping the Driver Binary in Sync": dynamically-generated node-driver
# CRDs (NutanixConfig, NutanixMachine, NutanixMachineTemplate) that never
# settle, and general Rancher/Fleet management-cluster instability that
# recurs every time the driver is (re-)registered.
#
# This is a bug in rancher/rancher itself, not something fixable from this
# driver's code - this script is the operational workaround: verify (and
# optionally repair) that every Rancher server replica has an identical
# copy of the driver binary.
#
# Requires: kubectl, with a working context against the Rancher management
# cluster. --fix additionally needs curl and sha256sum (or shasum) if no
# replica has a copy to source from, so it falls back to downloading a
# fresh, checksum-verified copy from this repo's GitHub Releases.
set -euo pipefail

NAMESPACE="cattle-system"
SELECTOR="app=rancher"
BINARY="docker-machine-driver-nutanix"
VERSION=""
FIX=false
ASSET_DIR="/usr/share/rancher/ui/assets"
REPO="jenkira/nutanix-docker-machine"

usage() {
  cat <<EOF
Usage: $(basename "$0") [--fix] [options]

Check that every Rancher server replica has an identical copy of the
$BINARY driver binary in $ASSET_DIR, and optionally repair any that don't.

  --fix                Repair pods missing the binary (default: check only,
                        no changes made).
  --namespace NAME      Rancher namespace (default: $NAMESPACE)
  --selector SELECTOR   Pod label selector for Rancher server pods
                        (default: $SELECTOR)
  --binary NAME         Driver binary filename (default: $BINARY)
  --version TAG         Release tag to download from if --fix has to fetch
                        a fresh copy (default: latest release on $REPO)
  -h, --help            Show this help text
EOF
  exit "${1:-0}"
}

while [ $# -gt 0 ]; do
  case "$1" in
    --fix) FIX=true; shift ;;
    --namespace) NAMESPACE="$2"; shift 2 ;;
    --selector) SELECTOR="$2"; shift 2 ;;
    --binary) BINARY="$2"; shift 2 ;;
    --version) VERSION="$2"; shift 2 ;;
    -h|--help) usage 0 ;;
    *) echo "Unknown argument: $1" >&2; usage 1 ;;
  esac
done

command -v kubectl >/dev/null 2>&1 || { echo "kubectl is required" >&2; exit 1; }

mapfile -t PODS < <(kubectl -n "$NAMESPACE" get pods -l "$SELECTOR" -o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}')

if [ "${#PODS[@]}" -eq 0 ]; then
  echo "No pods found in namespace '$NAMESPACE' matching selector '$SELECTOR'." >&2
  echo "If your Rancher deployment uses different labels, pass --namespace/--selector explicitly." >&2
  exit 1
fi

echo "Checking $BINARY across ${#PODS[@]} Rancher pod(s) in $NAMESPACE (selector: $SELECTOR)"
echo

HAVE=()
MISSING=()
ERRORED=()

for pod in "${PODS[@]}"; do
  if OUT=$(kubectl -n "$NAMESPACE" exec "$pod" -- sh -c "test -f '$ASSET_DIR/$BINARY'" 2>&1); then
    echo "  OK       $pod"
    HAVE+=("$pod")
  elif kubectl -n "$NAMESPACE" exec "$pod" -- true >/dev/null 2>&1; then
    echo "  MISSING  $pod"
    MISSING+=("$pod")
  else
    echo "  ERROR    $pod (could not exec into pod: $OUT)"
    ERRORED+=("$pod")
  fi
done

echo

if [ "${#ERRORED[@]}" -gt 0 ]; then
  echo "${#ERRORED[@]} pod(s) could not be checked (see ERROR lines above)." >&2
  echo "Resolve connectivity/RBAC for those pods and re-run - their state is unknown, not confirmed OK." >&2
  exit 1
fi

if [ "${#MISSING[@]}" -eq 0 ]; then
  echo "All ${#PODS[@]} pod(s) have $BINARY. Nothing to do."
  exit 0
fi

if [ "$FIX" != true ]; then
  echo "${#MISSING[@]} pod(s) are missing $BINARY - re-run with --fix to repair."
  exit 1
fi

TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT
LOCAL_BIN="$TMPDIR/$BINARY"

if [ "${#HAVE[@]}" -gt 0 ]; then
  SRC="${HAVE[0]}"
  echo "Copying $BINARY out of $SRC ..."
  kubectl -n "$NAMESPACE" exec "$SRC" -- cat "$ASSET_DIR/$BINARY" > "$LOCAL_BIN"
else
  command -v curl >/dev/null 2>&1 || { echo "curl is required to download a fresh copy" >&2; exit 1; }

  if [ -z "$VERSION" ]; then
    echo "No pod has $BINARY - resolving the latest release of $REPO ..."
    VERSION="$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
      | grep -m1 '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')"
    if [ -z "$VERSION" ]; then
      echo "Could not resolve the latest release tag - pass --version <tag> explicitly." >&2
      exit 1
    fi
  fi

  echo "Downloading $BINARY $VERSION ..."
  curl -fsSL -o "$LOCAL_BIN" \
    "https://github.com/$REPO/releases/download/$VERSION/${BINARY}"
  curl -fsSL -o "$TMPDIR/checksums.txt" \
    "https://github.com/$REPO/releases/download/$VERSION/${BINARY}_${VERSION}_checksums.txt"

  EXPECTED="$(grep "  $BINARY\$" "$TMPDIR/checksums.txt" | awk '{print $1}')"
  if [ -z "$EXPECTED" ]; then
    echo "Could not find a checksum for $BINARY in the release's checksums file - aborting." >&2
    exit 1
  fi

  if command -v sha256sum >/dev/null 2>&1; then
    ACTUAL="$(sha256sum "$LOCAL_BIN" | awk '{print $1}')"
  else
    ACTUAL="$(shasum -a 256 "$LOCAL_BIN" | awk '{print $1}')"
  fi

  if [ "$EXPECTED" != "$ACTUAL" ]; then
    echo "Checksum mismatch for downloaded $BINARY (expected $EXPECTED, got $ACTUAL) - aborting." >&2
    exit 1
  fi
  echo "Checksum verified: $ACTUAL"
fi

chmod +x "$LOCAL_BIN"

for pod in "${MISSING[@]}"; do
  echo "Copying $BINARY into $pod ..."
  kubectl -n "$NAMESPACE" exec -i "$pod" -- sh -c "cat > $ASSET_DIR/$BINARY && chmod +x $ASSET_DIR/$BINARY" < "$LOCAL_BIN"
done

echo
echo "Done. $BINARY is now consistent across all ${#PODS[@]} checked pod(s)."
echo "If the nutanix NodeDriver's dynamically-generated CRDs (NutanixConfig/NutanixMachine/"
echo "NutanixMachineTemplate) were already stuck, toggle it Inactive and back to Active"
echo "(Cluster Management > Drivers > Node Drivers) to force Rancher to re-run registration"
echo "against every replica now that the binary is consistent."
