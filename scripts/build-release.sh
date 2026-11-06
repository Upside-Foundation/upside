#!/usr/bin/env bash
# Builds the mainnet upsided release binary for linux/amd64 in a pinned
# container, so anyone can rebuild it from the same commit and get the same
# SHA-256.
#
#   scripts/build-release.sh <version>      e.g. scripts/build-release.sh v2.0.0
#
# Output: build/release/upsided-<version>-linux-amd64 and a .sha256 file.
#
# Inputs that change the output, all fixed here or taken from the commit:
#   - Go toolchain and C toolchain: the image below, pinned by digest.
#   - Target: linux/amd64, CGO enabled (glibc from the image).
#   - Build tags and ldflags: listed below; no git describe, no timestamps.
#   - Source: a git archive of HEAD, so untracked and ignored files are
#     never built. The tree must be clean.
#   - -trimpath and -buildvcs=false, so the build path and .git state are not
#     embedded. The commit is embedded through ldflags instead.
set -euo pipefail

VERSION="${1:?usage: $0 <version, e.g. v2.0.0>}"

# golang:1.25.7-bookworm. Mainnet's current binary reports go1.25.7.
IMAGE="golang@sha256:564e366a28ad1d70f460a2b97d1d299a562f08707eb0ecb24b659e5bd6c108e1"
PLATFORM="linux/amd64"
BUILD_TAGS="netgo objstore pebbledb mainnet"

cd "$(git rev-parse --show-toplevel)"
if [[ -n "$(git status --porcelain --untracked-files=no)" ]]; then
    echo "tracked files have uncommitted changes; commit or stash them first" >&2
    exit 1
fi
COMMIT="$(git rev-parse HEAD)"

OUT="build/release"
mkdir -p "$OUT"
# Keep the archive inside the repository: Docker on macOS (Colima, Docker
# Desktop) only shares the home directory with containers by default.
SRC_TAR="$PWD/$OUT/src-$COMMIT.tar"
trap 'rm -f "$SRC_TAR"' EXIT
git archive --format=tar HEAD > "$SRC_TAR"

BIN="upsided-${VERSION}-linux-amd64"
TAGS_COMMA="${BUILD_TAGS// /,}"
LDFLAGS="-w -s \
 -X github.com/cosmos/cosmos-sdk/version.Name=upside \
 -X github.com/cosmos/cosmos-sdk/version.AppName=upsided \
 -X github.com/cosmos/cosmos-sdk/version.Version=${VERSION#v} \
 -X github.com/cosmos/cosmos-sdk/version.Commit=${COMMIT} \
 -X github.com/cosmos/cosmos-sdk/version.BuildTags=${TAGS_COMMA}"

docker run --rm --platform "$PLATFORM" \
    -e GOTOOLCHAIN=local -e CGO_ENABLED=1 -e GOFLAGS=-mod=readonly \
    -v "$SRC_TAR:/src.tar:ro" -v "$PWD/$OUT:/out" \
    "$IMAGE" bash -euc "
        mkdir /build
        tar -xf /src.tar -C /build
        cd /build
        go build -trimpath -buildvcs=false -tags '$BUILD_TAGS' -ldflags '$LDFLAGS' \
            -o /out/$BIN ./cmd/upsided
        go version -m /out/$BIN | grep -E '^/out|go1|CGO_ENABLED|GOARCH|-tags'
    "

(cd "$OUT" && sha256sum "$BIN" 2>/dev/null || shasum -a 256 "$BIN") | tee "$OUT/$BIN.sha256"
echo "commit $COMMIT"
