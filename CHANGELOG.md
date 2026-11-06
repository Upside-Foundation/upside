# Changelog

## Unreleased

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
  (1 year). The whole param set is replaced, so a proposal that leaves the field out fails when it
  executes.
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
- Corrected the `x/upside` spec and removed the gravity bridge documentation and ADRs.
- Updated the mainnet and testnet peer lists in the join guides.
- `join-mainnet.md` explains syncing across the `v2.0.0` upgrade with Cosmovisor; the dead binary
  download links in the join guides were removed.

## v1.6.2

Mainnet launch: chain `upside_61873-1`, genesis 2026-02-18.

Based on [Cronos](https://github.com/crypto-org-chain/cronos) after its v1.5.4 release, including its
unreleased changes at the time of the fork. For history before the fork, see the
[Cronos changelog](https://github.com/crypto-org-chain/cronos/blob/main/CHANGELOG.md).
