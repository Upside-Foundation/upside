<div align="center">

# Upside Network

**The settlement layer for a new kind of prediction network.**

An EVM-compatible, IBC-connected Layer 1 built on the Cosmos SDK — engineered for high-throughput, verifiable state and designed from the ground up to anchor Upside's decentralized prediction and resolution network.

[![Go Version](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go&logoColor=white)](go.mod)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Cosmos SDK](https://img.shields.io/badge/Cosmos%20SDK-v0.50-5064FB)](https://github.com/crypto-org-chain/cosmos-sdk)
[![CometBFT](https://img.shields.io/badge/CometBFT-v0.38-1D1D1D)](https://github.com/cometbft/cometbft)
[![EVM Compatible](https://img.shields.io/badge/EVM-Compatible-627EEA?logo=ethereum&logoColor=white)](#architecture-highlights)
[![IBC Enabled](https://img.shields.io/badge/IBC-v10-6F4CFF)](https://github.com/cosmos/ibc-go)

[Networks](#networks) · [Quick Start](#quick-start) · [Architecture](#architecture-highlights) · [Docs](#documentation) · [Roadmap](#roadmap) · [Contributing](#contributing)

</div>

---

## What is Upside?

Upside is a sovereign Layer 1 blockchain that pairs the **full Ethereum Virtual Machine** with the **Cosmos SDK** and **CometBFT** consensus. You get Solidity smart contracts, MetaMask/ethers.js/Hardhat compatibility, and Ethereum JSON-RPC out of the box — plus native interoperability with the wider Cosmos ecosystem through IBC.

This repository contains **`upsided`**, the reference node implementation and the core protocol of the Upside chain:

- **EVM execution** via Ethermint, with EIP-1559 fee markets, EIP-712 signing, and EIP-7702 support
- **Cross-chain connectivity** via IBC v10 (transfers and interchain accounts)
- **Token bridging between Cosmos and the EVM**: native and IBC tokens can be mapped to CRC20 contracts so they work with web3 tooling (automatic wrapping of inbound IBC tokens is off on mainnet; `enable_auto_deployment` is `false`)
- **A high-performance storage engine** (`memiavl` + `versiondb`) that dramatically reduces disk footprint and speeds up state access
- **On-chain encryption keys** (`x/e2ee`) for end-to-end encrypted messaging between network participants
- **A capped native supply**: `x/mint` mints UP along a half-life curve toward a cap of 2.1 billion UP and never mints past it (see the [supply curve spec](x/upside/spec/08_supply.md), which also covers the IBC voucher conversion, a separate path that the cap does not limit)

> **Scope of this repo.** This is the *base network* — the consensus, execution, and storage layer. The prediction-network components (markets, resolution, oracle and judging infrastructure) are being developed as separate modules and will ship in dedicated repositories under the [Upside Foundation](https://github.com/Upside-Foundation) organization. See [Roadmap](#roadmap).

---

## Networks

| Network | Chain ID (EVM) | Cosmos Chain ID | RPC | Explorer | Faucet |
| ------- | -------------- | --------------- | --- | -------- | ------ |
| **Mainnet** | `61873` | `upside_61873-1` | [rpc.mainnet.upside.net](https://rpc.mainnet.upside.net) | [explorer.mainnet.upside.net](https://explorer.mainnet.upside.net) | — |
| **Testnet** | `61872` | `upside_61872-1` | [rpc.testnet.upside.net](https://rpc.testnet.upside.net) | [explorer.testnet.upside.net](https://explorer.testnet.upside.net) | [faucet.testnet.upside.net](https://faucet.testnet.upside.net) |

**Token details**

| | |
| --- | --- |
| Native token | **UP** |
| Base denom | `aup` (1 UP = 10¹⁸ `aup`) |
| Bech32 prefix | `up` |
| Key algorithm | `eth_secp256k1` (Ethereum-compatible addresses) |
| Default node home | `~/.upside` |

Genesis files for both networks are checked into this repo: [`assets/mainnet/genesis.json`](assets/mainnet/genesis.json) and [`assets/testnet/genesis.json`](assets/testnet/genesis.json).

---

## Quick Start

### Prerequisites

- **Go 1.23.x** (see [`go.mod`](go.mod))
- **Make**
- **gcc** (only if building with Ledger support, which is the default) — or set `LEDGER_ENABLED=false`
- **jq** (used by the local-node helper script)

### Build & install `upsided`

```bash
git clone https://github.com/Upside-Foundation/upside-network.git
cd upside-network

# Build a mainnet binary into ./build/upsided
make build

# ...or install into $GOPATH/bin
make install

# Without Ledger hardware-wallet support (no gcc required)
LEDGER_ENABLED=false make install

# Build a testnet binary instead
NETWORK=testnet make install

upsided version --long
```

### Run a local single-node devnet

The fastest way to get a chain running on your machine:

```bash
./start-local-node.sh
```

This builds `upsided`, initializes a fresh chain (`upside_61872-1`), funds a few deterministic dev accounts, creates a genesis validator, and starts the node with the EVM JSON-RPC (`:8545`), CometBFT RPC (`:26657`), and REST API (`:1317`) enabled.

> ⚠️ The script wipes `~/.upside*` before starting. Don't run it on a machine with chain state you care about.

Point MetaMask or any Ethereum tool at `http://localhost:8545` with chain ID `61872` and you're live.

### Join a public network

Step-by-step guides for running full nodes and validators:

| Guide | Description |
| ----- | ----------- |
| [Build & Install](docs/technical-documentation/build-and-install.md) | Build `upsided` for mainnet or testnet |
| [Local Node Setup](docs/technical-documentation/local-node-setup.md) | Single- or multi-node local development network |
| [Join Local Network](docs/technical-documentation/join-local-network.md) | Add a validator to an existing local network |
| [Join Testnet](docs/technical-documentation/join-testnet.md) | Sync a full node or validator on testnet |
| [Join Mainnet](docs/technical-documentation/join-mainnet.md) | Sync a full node or validator on mainnet |

---

## Architecture Highlights

Upside is built from battle-tested components and extends them where it matters.

### ⚡ Dual execution & API surface

- **EVM** (via the Crypto.org forks of [Ethermint](https://github.com/crypto-org-chain/ethermint) and [go-ethereum](https://github.com/crypto-org-chain/go-ethereum)): full Ethereum JSON-RPC (`eth`, `net`, `web3`, `txpool`, `debug`), tracing, filters and subscriptions, EIP-1559 dynamic fees, EIP-712 typed data, EIP-7702 account delegation.
- **Cosmos**: gRPC, gRPC-gateway REST, and CometBFT RPC. Every module — including the EVM — is queryable through Cosmos APIs, and Swagger docs ship with the node ([`client/docs`](client/docs)).
- **Rosetta** API support for exchange and custody integrations.

### 🌉 Native interoperability

- **IBC v10**: fungible token transfers, interchain accounts (ICA), and light-client updates.
- **Automatic CRC20 wrapping**: the [`x/upside`](x/upside) module maintains a denom ↔ contract mapping so inbound IBC tokens are minted as CRC20 tokens (an ERC20 superset) and can be moved back out through EVM calls.
- **Conversion middleware** handles the 8 ↔ 18 decimal translation between IBC assets and the EVM transparently.
- **EVM precompiles** expose native Cosmos functionality — bank transfers, interchain accounts, and relayer operations — directly to Solidity contracts.

### 🗄️ High-performance state

- **[`memiavl`](memiavl)**: a re-engineered IAVL tree that keeps the working set in memory-mapped snapshots and uses change sets as a write-ahead log. Better write, read, *and* space amplification than stock IAVL — archived change sets can be ~100× smaller than the equivalent IAVL nodes.
- **[`versiondb`](versiondb)**: a RocksDB-backed, timestamp-indexed key-value store for historical state queries. Lets archive and RPC nodes serve full history without carrying a full IAVL tree, and powers the `--grpc-only` standalone query mode.
- **Pluggable backends**: RocksDB, PebbleDB, LevelDB, and more via build tags.
- **Optimistic execution** and mempool fee-bump / tx-replacement support for lower latency under load.

### 🔐 Protocol-level safety & privacy

- **[`x/e2ee`](x/e2ee)**: participants register encryption public keys on-chain and exchange end-to-end encrypted payloads on or off chain — foundational plumbing for private coordination in the prediction network.
- **Address block lists**: a per-node list checked in the ante handler, and an encrypted on-chain list that validators apply when building blocks. Both cover EIP-7702 authorization lists.
- **Governance-controlled parameters** and coordinated upgrades via `x/upgrade` + Cosmovisor.
- **Permissioned admin actions** (token mapping, block list) gated by an on-chain `upside_admin` account with permission bits.

---

## Repository Layout

```
.
├── app/                    # Application wiring: module manager, keepers, ante handlers, upgrades
├── cmd/upsided/            # CLI entrypoint and node configuration
├── x/upside/               # Core Upside module: CRC20 token mapping, IBC hooks, EVM handlers, precompiles
│   └── spec/               # Module specification (concepts, state, messages, events, params)
├── x/e2ee/                 # End-to-end encryption key registry module
├── proto/                  # Protobuf definitions for x/upside and x/e2ee
├── contracts/              # Solidity sources for module-owned CRC20/CRC21 contracts
├── memiavl/                # High-performance IAVL implementation (separate Go module)
├── store/                  # Custom rootmulti store integrating memiavl and object stores
├── versiondb/              # RocksDB-based historical state store (separate Go module)
├── client/docs/            # Swagger / OpenAPI docs served by the node
├── integration_tests/      # Python end-to-end test suite (pytest, run via Nix)
├── testground/             # Load and network testing harness
├── nix/                    # Reproducible build and dev-shell definitions
├── scripts/                # Devnet, protobuf, release, and contract-binding helpers
├── docs/                   # Operator and developer documentation
└── assets/                 # Mainnet and testnet genesis files
```

---

## Documentation

| Resource | What you'll find |
| -------- | ---------------- |
| [Technical Documentation](docs/technical-documentation/README.md) | Build, run, and operate `upsided` nodes and validators |
| [`x/upside` Module Spec](x/upside/spec/README.md) | Concepts, state, state transitions, messages, ABCI hooks, events, parameters, and the supply curve |
| [`x/upside` Events](x/upside/events/README.md) | EVM event decoding for cross-module hooks |
| [`memiavl` Design](memiavl/README.md) | Snapshot & WAL file formats, design rationale |
| [`versiondb` Guide](versiondb/README.md) | Enabling versiondb, migration, and catch-up procedures |
| [Integration Tests](integration_tests/README.md) | Running the end-to-end test suite |
| [Changelog](CHANGELOG.md) | Release history and upgrade notes |

---

## Roadmap

This repository is **layer one of a multi-repository system**. It establishes the chain that everything else settles to.

| Layer | Status | Where |
| ----- | ------ | ----- |
| **Core network** — consensus, EVM, IBC, storage, e2ee | ✅ This repo | [`Upside-Foundation/upside-network`](https://github.com/Upside-Foundation/upside-network) |
| **Prediction network modules** — markets, resolution, oracle & judging infrastructure | 🔜 Coming next | Dedicated repos under [Upside Foundation](https://github.com/Upside-Foundation) |
| **Client SDKs, indexers, and operator tooling** | 🔜 Planned | Dedicated repos under [Upside Foundation](https://github.com/Upside-Foundation) |

The prediction-network layer will be documented in its own repositories as it lands. Watch the [Upside Foundation](https://github.com/Upside-Foundation) organization for releases.

---

## Contributing

Contributions are welcome — bug reports, documentation, tests, and protocol improvements alike.

### Development workflow

```bash
# Run the Go unit tests (includes memiavl and store sub-modules)
make test

# Run versiondb tests (requires RocksDB)
make test-versiondb

# Lint Go code
make lint

# Regenerate protobuf code / Swagger
make proto-gen
make proto-swagger-gen

# Vulnerability scan
make vulncheck
```

### Integration tests

The end-to-end suite lives in [`integration_tests/`](integration_tests) and runs against a real multi-node devnet inside a reproducible [Nix](https://nixos.org) environment:

```bash
make run-integration-tests
```

A `flake.nix` and `.envrc.example` are provided for a fully pinned dev shell (`direnv allow` or `nix develop`).

### Guidelines

- Open an issue before starting large changes so we can align on design.
- Add an entry to [`CHANGELOG.md`](CHANGELOG.md) under **UNRELEASED** for user-facing changes.
- Keep commits focused; include tests for new behavior.
- Run `make lint` before opening a pull request.

---

## Security

If you discover a vulnerability, **please do not open a public issue**. Follow the coordinated disclosure process in [`SECURITY.md`](SECURITY.md).

---

## Acknowledgements

Upside builds on the shoulders of the Cosmos and Ethereum communities — in particular the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk), [CometBFT](https://github.com/cometbft/cometbft), [IBC-Go](https://github.com/cosmos/ibc-go), [Ethermint](https://github.com/evmos/ethermint), [go-ethereum](https://github.com/ethereum/go-ethereum), and the [Cronos](https://github.com/crypto-org-chain/cronos) project whose architecture this chain extends.

---

## License

Copyright 2026 Upside Foundation. Derived from [Cronos](https://github.com/crypto-org-chain/cronos), Copyright 2021–present Cronos Labs.

Licensed under the [Apache License, Version 2.0](LICENSE). See [`NOTICE`](NOTICE) for attribution and for components under other licenses.
