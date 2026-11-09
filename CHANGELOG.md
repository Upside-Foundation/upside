# Changelog

## Unreleased

## v2.0.1

Repository and documentation release, 2026-11-09. No consensus change; the state machine is the same as v2.0.0.

Nodes and validators keep running the `v2.0.0` binary. The v2.0.1 changes ship with the next chain software
upgrade.

- Rename the Go module and import prefix from `github.com/0xUp5ide/upside-foundation` to
  `github.com/Upside-Foundation/upside`, including the three local storage modules and protobuf
  Go-package metadata. The public repository and the module path are both
  `github.com/Upside-Foundation/upside`. Go consumers must update their imports.
  Installation remains a full checkout followed by `make build` or `make install`;
  direct `go install ...@version` and a tagged v2 Go-library release are not supported.
- Consolidate the account, staking, validator, node, and EVM guides under `docs/`,
  with a documentation index and network metadata. Document public mainnet sentries,
  isolated localnet setup, and historical-genesis compatibility. The steps for syncing across the
  `v2.0.0` upgrade, including Cosmovisor, are in `docs/nodes/operations.md`.
- Remove the remaining inherited ADRs and the historical Ethermint audit report from the
  documentation tree. Remove the duplicate genesis copies under `docs/technical-documentation/`;
  the genesis files in `assets/` and the embedded API assets are kept.
- Fix the single-node helper to set its generated EVM denomination to `aup`,
  matching its funded accounts, staking and minting. Existing network genesis
  files and chain state are unchanged.
- `LICENSE` holds only the standard Apache License 2.0 text, so license detection recognizes it.
  The copyright line, the upstream attribution, and the components under other licenses are in
  `NOTICE`.
- `scripts/build-release.sh` accepts `UPSIDE_BUILD_COMMIT` to set the commit embedded in the binary,
  so the `v2.0.0` artifact can be reproduced from the public `v2.0.0` tag. The image, flags, build
  tags, and paths are unchanged.
- Document starting a mainnet node with state sync from the public sentries or from the signed daily
  archive at `https://snapshots.mainnet.upside.net`, with the verification steps. Add the snapshot
  signing key statement, its signature, the signer file, and the public key under `docs/snapshots/`.
- Document that the public EVM JSON-RPC endpoints serve only the `eth`, `net`, and `web3`
  namespaces, and list the CometBFT RPC routes the public endpoints serve.
- Adopt the Contributor Covenant v2.1 as the code of conduct, with reports to contact@upside.net, and
  remove the C4 process references from `CONTRIBUTING.md`.
- Document the mnemonics in `scripts/.env` and `start-local-node.sh` as public test values
  for local networks and integration tests.

These changes do not introduce a consensus upgrade or change the existing
`v2.0.0` upgrade name.

## v2.0.0

Mainnet upgrade at height 3,898,900 through governance proposal 4, upgrade plan `v2.0.0`,
2026-10-22. Source tag `v2.0.0`.

### Features

- `x/upside`: capped native supply. `x/mint` now mints along a half-life curve toward a cap of
  2,100,000,000 UP instead of a fixed inflation rate, and never mints past it. The IBC voucher
  conversion is a separate path that can create `aup` and is not limited by the cap. Adds the
  `mint_half_life_seconds` parameter (default 2 years, minimum 1 year). The `v2.0.0` upgrade handler
  sets it and sets `x/mint`'s `blocks_per_year` to 5,783,000 to match the measured block time. See the
  [supply curve spec](x/upside/spec/08_supply.md).
- `x/upside`: governance can pause minting with `MsgSetMintPaused`, which sets the new `mint_paused`
  parameter (default `false`). While it is `true`, `x/mint` mints nothing; when it is set back to
  `false`, minting resumes along the curve from the current supply, with no catch-up for the paused
  blocks.

### State machine breaking

- `x/upside`: `MsgUpdateParams` must now include `mint_half_life_seconds` of at least 31,557,600
  (1 year). The whole param set is replaced, and the parameters are validated when the proposal is
  submitted, so a proposal that leaves the field out is rejected at submission.
- `x/upside`: `MsgUpdateParams` fails if its `mint_paused` differs from the stored value; only
  `MsgSetMintPaused` changes it.
- `x/upside`: capped minting returns a rate of zero while `x/mint`'s `blocks_per_year` is outside
  2,891,500 to 11,566,000 (half and double the measured mainnet value). At the measured block time,
  one `x/mint` proposal can no longer shorten the real half-life below half of
  `mint_half_life_seconds`.
- `x/upside`: capped minting returns a rate of zero whenever `x/mint`'s mint denom differs from the
  staking bond denom, so the curve is never applied to the wrong supply.
- `x/upside`: the v1 to v2 params migration fills in the default half-life, so chains still at
  consensus version 1 can upgrade.

### Build

- `go.mod` pins `toolchain go1.25.7`, the version the mainnet binary is built with.
- `scripts/build-release.sh` builds the mainnet linux/amd64 binary reproducibly in a pinned container.
- CI (`.github/workflows/test.yml`) builds the mainnet binary, runs `go vet`, and runs the Go tests.

### Documentation

- `LICENSE` now holds the full Apache License 2.0 text; `NOTICE` added with the Cronos attribution and
  the components under other licenses.
- Rewrote `SECURITY.md`, `CODE_OF_CONDUCT.md`, and `CONTRIBUTING.md` for Upside.
- Corrected the `x/upside` spec and removed the gravity bridge documentation and ADRs 001 and 005.
- Updated the mainnet and testnet peer lists in the join guides.
- `join-mainnet.md` explains syncing across the `v2.0.0` upgrade with Cosmovisor; the dead binary
  download links in the join guides were removed.

## v1.6.2

Mainnet launch: chain `upside_61873-1`, genesis 2026-02-18.

Based on [Cronos](https://github.com/crypto-org-chain/cronos) after its v1.5.4 release, including its
unreleased changes at the time of the fork. For history before the fork, see the
[Cronos changelog](https://github.com/crypto-org-chain/cronos/blob/main/CHANGELOG.md).
