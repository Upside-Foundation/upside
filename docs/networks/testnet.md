---
title: Connect to the public testnet
description: Connect a wallet or node to the Upside L1 public testnet.
slug: /networks/testnet
---

# Connect to the public testnet

Testnet balances have no real-world economic value and are separate from mainnet. The network may change or reset due to software upgrades and testing.

## Network settings

| Setting | Value |
| --- | --- |
| Cosmos chain ID | `upside_61872-1` |
| EVM chain ID | `61872` |
| Account prefix | `up` |
| Base denomination | `aup` |
| Decimal places | 18 |

The testnet uses the same `up` address prefix and `aup` denomination as mainnet. The chain ID is what separates the two networks, so check it before signing.

## Public endpoints

| Service | Address |
| --- | --- |
| EVM JSON-RPC | `https://rpc.testnet.upside.net` |
| EVM JSON-RPC over WebSocket | `wss://rpc.testnet.upside.net/ws` |
| CometBFT RPC | `https://rpc.testnet.upside.net/tendermint/` |
| Block explorer | [explorer.testnet.upside.net](https://explorer.testnet.upside.net) |
| Faucet | [faucet.testnet.upside.net](https://faucet.testnet.upside.net) |

The CometBFT RPC address needs its trailing slash when you pass it to `upsided` with `--node`. The public CometBFT RPC and EVM JSON-RPC serve the same routes and namespaces as [mainnet's public endpoints](mainnet.md#public-endpoints).

## Connect a wallet

1. Install the `upsided` binary. The testnet uses the same address format as mainnet, so both use the mainnet address-prefix build. Select the software version separately for the operation you intend to run.
2. Set the variables from [networks](../getting-started/networks.md) with `CHAIN_ID="upside_61872-1"`, `RPC_URL="https://rpc.testnet.upside.net/tendermint/"`, and `WALLET_HOME="$HOME/.upside-wallet-testnet"`.
3. Create a separate test wallet and verify its address prefix.
4. Obtain test tokens from the [faucet](https://faucet.testnet.upside.net).
5. Practice [transfers](../getting-started/transactions.md) and [delegation](../staking/delegation.md). Use the wallet guide to check transaction results.

## Join with a full node or validator

The testnet hasn't applied the `v2.0.0` upgrade, so the `v2.0.0` release doesn't match its state. A `v2.0.0` or `v2.0.1` binary would compute minting differently from the testnet's own nodes. The binary that the testnet runs isn't published yet. Its nodes report no application version, and its EVM JSON-RPC reports a build with Go 1.25.6.

From the source checkout, set the profile:

```bash
export UPSIDE_BIN="$HOME/.local/bin/upsided-mainnet"
export CHAIN_ID="upside_61872-1"
export NODE_HOME="$HOME/.upside-testnet"
export MONIKER="my-testnet-node"
export GENESIS_FILE="$PWD/assets/testnet/genesis.json"
export GENESIS_SHA256="e2d9ed9255de41e6fdc2de34bdd2ac54d0d3f39092530363682524798ceaedc3"
export NODE_MIN_GAS_PRICES="0aup"
export SEEDS=""
export PEERS="9c43ac9d3f4edf87caee48d102d2099d9369f672@162.19.232.154:26656"
```

`upsided-mainnet` is the name that [Install the node and account CLI](../getting-started/install.md) gives a build with the `up` address prefix. For the testnet, it must be a build of the testnet's release, not `v2.0.0`.

The public testnet peer is:

```text
9c43ac9d3f4edf87caee48d102d2099d9369f672@162.19.232.154:26656
```

Set `persistent_peers` to this testnet P2P endpoint and leave `seeds` empty.

No public testnet snapshot source is published yet. Once you have a binary and synchronization state compatible with the testnet, follow [public-node setup](../nodes/public-node.md), then [validator setup](../validators/create.md). Query the testnet's current staking and slashing parameters before using it.

If a mainnet and testnet node run on one host, each needs distinct P2P, RPC, API, gRPC, EVM RPC, and metrics ports for enabled services, as well as separate homes and service units.
