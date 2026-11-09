# Upside Network

`upsided` is the node and account CLI for the Upside Layer 1. It uses the Cosmos SDK, CometBFT consensus, and an Ethermint-based EVM.

This repository contains the running L1 protocol. Application components such as prediction markets are developed separately.

## Networks

| Network | Cosmos chain ID | EVM chain ID | EVM JSON-RPC |
| --- | --- | --- | --- |
| Mainnet | `upside_61873-1` | `61873` | `https://rpc.mainnet.upside.net` |
| Testnet | `upside_61872-1` | `61872` | `https://rpc.testnet.upside.net` |

Both public networks use `up` account addresses and `aup` as the native base denomination. One UP is 10^18 `aup`. The public testnet uses mainnet address prefixes. The `NETWORK=testnet` build option selects the separate `tup` prefix and is incompatible with these public network addresses.

See [network settings](docs/getting-started/networks.md) for wallet connections and [mainnet](docs/networks/mainnet.md) or [testnet](docs/networks/testnet.md) for node setup.

## Build

To run a mainnet node or validator, use the `v2.0.0` release binary described in [mainnet](docs/networks/mainnet.md#release-and-synchronization). Building from source is for development and for reproducing a release.

Install Go 1.23.12 or later and Make. With `GOTOOLCHAIN=auto`, Go selects the `go1.25.7` toolchain specified in `go.mod`.

From a checkout of this repository:

```bash
NETWORK=mainnet LEDGER_ENABLED=false make build
./build/upsided version --long
```

The output is `build/upsided`. Build and install instructions, including the pinned Docker release recipe, are in [Install the node and account CLI](docs/getting-started/install.md).

For an isolated practice chain, follow [Run a localnet](docs/nodes/localnet.md). That guide creates a fresh home and disposable keys. The legacy `start-local-node.sh` helper runs `rm -rf ~/.upside*`. That deletes `~/.upside`, `~/.upside-mainnet`, `~/.upside-testnet`, `~/.upside-wallet-mainnet`, `~/.upside-wallet-testnet`, the localnet homes, and any other path in your home directory whose name begins with `.upside`. Don't use it on a machine containing accounts or node state you need to keep.

## Documentation

| Guide | Contents |
| --- | --- |
| [Documentation index](docs/README.md) | Account, node, validator, staking, and developer guides |
| [Accounts](docs/getting-started/wallets.md) | Create and recover keys, inspect balances, confirm transactions |
| [Validators](docs/validators/create.md) | Register a validator and check its state |
| [Full node and snapshots](docs/nodes/public-node.md#4-install-the-supported-sync-state) | Start a mainnet node with state sync or a signed snapshot archive |
| [Sentries](docs/nodes/sentries.md) | Configure a validator behind sentry nodes |
| [Operations](docs/nodes/operations.md) | Services, monitoring, migration, and upgrades |
| [EVM development](docs/technical/build-on-the-evm.md) | Contract deployment, fees, and compatibility limits |
| [JSON-RPC](docs/technical/json-rpc.md) | Ethereum-compatible RPC methods and examples |
| [Supply and minting](docs/technical/supply-and-minting.md) | The 2.1-billion-UP minting cap and half-life calculation |
| [Upside module](x/upside/spec/README.md) | State, messages, events, parameters, and token conversion |
| [Integration tests](integration_tests/README.md) | Local test tooling |

The source includes IBC and token-conversion support. Deployment settings determine whether those features are active. Mainnet governance disabled IBC token transfers in proposal 6. Native `x/mint` issuance follows the capped half-life curve; the supply reference also describes separate minting paths.

## Source layout

| Directory | Contents |
| --- | --- |
| `app/` | Module wiring, transaction processing, and upgrade handlers |
| `cmd/upsided/` | Daemon entrypoint and configuration |
| `x/upside/` | Upside parameters, minting logic, token mappings, and EVM hooks |
| `x/e2ee/` | Encryption public-key registry |
| `proto/` | Protocol message definitions |
| `contracts/` | Solidity contracts and generated bindings |
| `memiavl/`, `store/`, `versiondb/` | Storage implementations |
| `client/docs/` | API documentation served by the node |
| `assets/` | Public-network genesis files |
| `docs/` | Operator and developer documentation |

## Development

```bash
make test
make lint
```

Some integration and storage tests require additional tooling. See [Contributing](CONTRIBUTING.md), the relevant module README, and [the changelog](CHANGELOG.md).

Report vulnerabilities using the [security policy](SECURITY.md).

## Community

- Join the [Upside Discord server](https://discord.gg/wc7C5AZhh3) for chat and questions.
- Use GitHub Discussions for longer questions and proposals.
- Report bugs as GitHub issues.
- Report vulnerabilities privately as described in the [security policy](SECURITY.md). Don't post them publicly.

## License and attribution

Copyright 2026 Upside Foundation. Derived from Cronos and built with Cosmos SDK, CometBFT, Ethermint, and go-ethereum components.

See [LICENSE](LICENSE) and [NOTICE](NOTICE) for licensing and attribution.
