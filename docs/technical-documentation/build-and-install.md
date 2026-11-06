# Build and Install `upsided`

This page describes how to build the Upside node binary (`upsided`) from this repository.

## Prerequisites

- **Go**: 1.23.12 or later. `go.mod` sets `toolchain go1.25.7`, so with the
  default `GOTOOLCHAIN=auto` an older Go downloads and uses go1.25.7, the
  version mainnet binaries are built with.
- **Make**
- (Optional) **gcc**: required only if you build with ledger support (`LEDGER_ENABLED=true`, the default).

## Build (local binary in `./build`)

```bash
# Build a mainnet binary (default): produces ./build/upsided
make build

# If you don’t have gcc/ledger dependencies installed:
LEDGER_ENABLED=false make build
```

## Install (into `GOBIN` / `GOPATH/bin`)

```bash
make install
# or
LEDGER_ENABLED=false make install
```

After install, confirm it’s on your PATH:

```bash
upsided version
upsided version --long
```

## Release build (reproducible)

Mainnet release binaries are built with `scripts/build-release.sh`, so that
validators can rebuild a release from its commit and check that the SHA-256
matches the one in the upgrade proposal.

```bash
git checkout <release commit>
scripts/build-release.sh v2.0.0
cat build/release/upsided-v2.0.0-linux-amd64.sha256
```

The script needs Docker (on macOS, Colima or Docker Desktop). It builds
linux/amd64 in `golang:1.25.7-bookworm`, pinned by image digest, from a
`git archive` of the commit, with these settings:

| Setting      | Value                                                   |
| ------------ | ------------------------------------------------------- |
| Go           | go1.25.7 (`GOTOOLCHAIN=local` inside the image)         |
| CGO          | enabled, glibc and gcc from the pinned image            |
| Build tags   | `netgo objstore pebbledb mainnet`                       |
| Flags        | `-trimpath -buildvcs=false -mod=readonly`               |
| ldflags      | `-w -s`, and the version, commit, and build tags (`-X`) |

`make build` gives a binary with the same behavior but not the same bytes:
it uses whatever Go and gcc are installed, adds the `ledger` tag when gcc is
present, and takes the version from `git describe`. Use it for development,
not for release checksums.

The Nix build (`scripts/release.sh`, `flake.nix`) uses Go 1.23 from the
pinned nixpkgs, so its output also differs from the release binary.
