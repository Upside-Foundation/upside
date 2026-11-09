---
title: Install the node and account CLI
description: Build upsided from source or install a verified release artifact.
slug: /getting-started/install
---

# Install the node and account CLI

`upsided` provides the node process and account commands. Account commands can use a public RPC endpoint without running a node.

## Build from source

Use a checkout of the release tag selected for your network. You need Git, Make, and Go 1.23.12 or later. With the default `GOTOOLCHAIN=auto`, Go selects the `go1.25.7` toolchain requested in `go.mod`. Ledger support also requires a C compiler.

Clone the repository and check out the release tag that the [mainnet](../networks/mainnet.md) or [testnet](../networks/testnet.md) page names. Replace `RELEASE_TAG` with that tag:

```bash
git clone https://github.com/Upside-Foundation/upside.git
cd upside
git checkout RELEASE_TAG
```

From the repository root, build without Ledger support:

```bash
NETWORK=mainnet LEDGER_ENABLED=false make build
./build/upsided version --long
mkdir -p "$HOME/.local/bin"
install -m 0755 ./build/upsided "$HOME/.local/bin/upsided-mainnet"
export UPSIDE_BIN="$HOME/.local/bin/upsided-mainnet"
```

For Ledger support, install a C compiler first, and then omit `LEDGER_ENABLED=false`.

Both public networks use the `up` address prefix. Use `NETWORK=mainnet` for that address format. `NETWORK=testnet` selects `tup` and is incompatible with the public testnet's addresses. Select the software version and synchronization state using the [mainnet](../networks/mainnet.md) or [testnet](../networks/testnet.md) instructions.

The repository and the root Go module are both `github.com/Upside-Foundation/upside`. Local replacement directives select the `memiavl`, `store`, and `versiondb` modules from the same checkout. Build from the checkout. `go install ...@latest` does not support this module's replacement directives.

## Install a release binary

Mainnet nodes and validators run the `v2.0.0` release binary for Linux on amd64. It is published at `https://github.com/Upside-Foundation/upside/releases/download/v2.0.0/upsided-v2.0.0-linux-amd64`, with its checksum in a file of the same name plus `.sha256`. The SHA-256 checksum is `1279d570c0a031bcd0b7c7ce4aa379b3de9f6bae984feffe894bc357de6ca40b`, which is also in [`docs/network-manifest.json`](../network-manifest.json). On other systems, use the source build above.

Download the binary as `upsided`, and check its checksum before you run it:

```bash
curl -fL -o upsided https://github.com/Upside-Foundation/upside/releases/download/v2.0.0/upsided-v2.0.0-linux-amd64
EXPECTED_SHA256="1279d570c0a031bcd0b7c7ce4aa379b3de9f6bae984feffe894bc357de6ca40b"
ACTUAL_SHA256="$(sha256sum ./upsided | awk '{print $1}')" || ACTUAL_SHA256=""
```

On macOS, replace `sha256sum` with `shasum -a 256`. The next command installs the binary only if the two values match, and otherwise prints `Stop:` and installs nothing:

```bash
if [ -n "$ACTUAL_SHA256" ] && [ "$ACTUAL_SHA256" = "$EXPECTED_SHA256" ]
then
  mkdir -p "$HOME/.local/bin"
  install -m 0755 ./upsided "$HOME/.local/bin/upsided-mainnet"
  export UPSIDE_BIN="$HOME/.local/bin/upsided-mainnet"
  "$UPSIDE_BIN" version --long
else
  printf 'Stop: upsided has SHA-256 %s, not %s\n' "${ACTUAL_SHA256:-empty}" "$EXPECTED_SHA256" >&2
fi
```

Check the embedded commit, Go version, and build tags against the release instructions.

## Reproduce a release with Docker

The repository's release script builds Linux/amd64 in a pinned Docker image, so it needs Docker. On Apple silicon, Docker must also be able to run amd64 images through emulation. From a clean checkout of the intended release commit:

```bash
RELEASE_VERSION="v2.0.1"
bash scripts/build-release.sh "$RELEASE_VERSION"
cat "build/release/upsided-${RELEASE_VERSION}-linux-amd64.sha256"
```

Set `RELEASE_VERSION` to the release you intend to reproduce. The script builds `HEAD` and embeds that commit. To embed a different commit, set `UPSIDE_BUILD_COMMIT` to its full 40-character hash. [Mainnet](../networks/mainnet.md#release-and-synchronization) uses this to reproduce the `v2.0.0` artifact. It requires a clean tracked working tree, creates a source archive, and writes the binary and checksum under `build/release/`.

| Setting | Release recipe |
| --- | --- |
| Target | Linux/amd64 |
| Go | 1.25.7 |
| Image | `golang@sha256:564e366a28ad1d70f460a2b97d1d299a562f08707eb0ecb24b659e5bd6c108e1` |
| CGO | Enabled, using the image's C toolchain |
| Build tags | `netgo,objstore,pebbledb,mainnet` |
| Build flags | `-trimpath -buildvcs=false -mod=readonly` |

Compare checksums only for the same source commit, release label, target, and build recipe. `make build` also includes `nativebyteorder`, normally includes `ledger`, and obtains the version from `git describe`. Its checksum can differ. The Nix release definitions use a separate toolchain and build configuration.

## Versions and upgrades

The release label, source commit, and on-chain upgrade name are separate identifiers. The `v2.0.1` release registers only the `v2.0.0` upgrade handler. An on-chain upgrade named `v2.0.1` requires its own handler.

To upgrade a running node, change its service binary according to the release's compatibility and migration instructions.

Next: choose a [network profile](networks.md), then [create an account](wallets.md) or [run a localnet](../nodes/localnet.md).
