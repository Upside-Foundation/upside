---
title: Connect to mainnet
description: Connect a wallet or full node to the Upside L1 mainnet.
slug: /networks/mainnet
---

# Connect to mainnet

Mainnet uses Cosmos chain ID `upside_61873-1`, EVM chain ID `61873`, and the `aup` base denomination.

## Public endpoints

| Service | Address |
| --- | --- |
| EVM JSON-RPC | `https://rpc.mainnet.upside.net` |
| EVM JSON-RPC over WebSocket | `wss://rpc.mainnet.upside.net/ws` |
| CometBFT RPC | `https://rpc.mainnet.upside.net/tendermint/` |
| Block explorer | [explorer.mainnet.upside.net](https://explorer.mainnet.upside.net) |
| Snapshot manifest | `https://snapshots.mainnet.upside.net/manifest.json` |

The CometBFT RPC address needs its trailing slash when you pass it to `upsided` with `--node`.

The public CometBFT RPC serves only these routes: `/health`, `/status`, `/abci_info`, `/abci_query`, `/block`, `/block_by_hash`, `/block_results`, `/blockchain`, `/commit`, `/header`, `/header_by_hash`, `/validators`, `/consensus_params`, `/genesis_chunked`, `/tx`, `/tx_search`, `/block_search`, `/check_tx`, `/broadcast_tx_sync`, `/broadcast_tx_async`, `/num_unconfirmed_txs`. Other routes return `404`, including `/genesis`, `/net_info`, `/subscribe`, and `/broadcast_tx_commit`. Use `/genesis_chunked` or the genesis file in this repository instead of `/genesis`, and run your own node for the other routes. The EVM JSON-RPC endpoints serve only the `eth`, `net`, and `web3` namespaces, as described in [JSON-RPC methods](../technical/json-rpc.md).

Use the seed for peer discovery and the public sentries for persistent connections:

| Role | Node ID | Address |
| --- | --- | --- |
| Seed | `7c1e6a83b2e5308428c1164091c6a5765d79a184` | `3.135.80.240:26656` |
| Public sentry | `1bfb165a3fa4b3a8ccbd55fabd3ecf8117e88bd6` | `3.19.153.80:26656` |
| Public sentry | `9d8723c3fec856b80d5357f480c984b923b2c729` | `18.225.42.62:26656` |

Check any endpoint's chain ID before use, as shown in [networks and connection settings](../getting-started/networks.md).

## Connect a wallet

Install the mainnet binary and select the [mainnet profile](../getting-started/networks.md). The CometBFT RPC reports `upside_61873-1`. Follow [wallets](../getting-started/wallets.md), [transfers](../getting-started/transactions.md), and [staking](../staking/delegation.md).

An EVM wallet uses the EVM JSON-RPC URL and chain ID `61873`. The native currency is UP with 18 decimals. Wallet account and transaction rules are described in the [wallet guide](../getting-started/wallets.md).

## Join with a full node

A full node verifies the chain and maintains its state. Becoming a validator is a separate step after synchronization. Review the [requirements](../nodes/requirements.md) and install the `v2.0.0` binary described in [Release and synchronization](#release-and-synchronization) first.

From the source checkout, set the node profile in a Bash shell:

```bash
export UPSIDE_BIN="$HOME/.local/bin/upsided-mainnet"
export CHAIN_ID="upside_61873-1"
export NODE_HOME="$HOME/.upside-mainnet"
export MONIKER="my-mainnet-node"
export GENESIS_FILE="$PWD/assets/mainnet/genesis.json"
export GENESIS_SHA256="e34bf9628164ba43375186b1118c8aaef82d2be08d4685d9e8cca68ad356aa50"
export SEEDS="7c1e6a83b2e5308428c1164091c6a5765d79a184@3.135.80.240:26656"
export PEERS="1bfb165a3fa4b3a8ccbd55fabd3ecf8117e88bd6@3.19.153.80:26656,9d8723c3fec856b80d5357f480c984b923b2c729@18.225.42.62:26656"
export NODE_MIN_GAS_PRICES="0aup"
```

`SEEDS` is the public mainnet seed node and `PEERS` lists the public mainnet sentries. Each entry has the form `node-id@host:port`, and all of them use the standard P2P port `26656`. `NODE_MIN_GAS_PRICES="0aup"` adds no node-local gas-price floor. The chain's active fee-market checks still apply. Configure wallet transaction gas prices separately from these node settings.

Continue with [public-node setup](../nodes/public-node.md). After synchronization, you can serve your own wallet queries or [create a validator](../validators/create.md).

## Release and synchronization

Mainnet nodes and validators run the `v2.0.0` binary. Proposal 4 activated it at height **3,898,900** through the upgrade plan named `v2.0.0`. Its source is the tag `v2.0.0`. Download the release artifact `upsided-v2.0.0-linux-amd64` from `https://github.com/Upside-Foundation/upside/releases/download/v2.0.0/upsided-v2.0.0-linux-amd64` and verify this SHA-256 checksum before you run it, as described in [Install a release binary](../getting-started/install.md#install-a-release-binary):

```text
1279d570c0a031bcd0b7c7ce4aa379b3de9f6bae984feffe894bc357de6ca40b
```

The artifact embeds commit `8521a42fc1806e06c422e67588dd9dd834a09b76`, which `upsided version --long` reports. That is the commit the artifact was built from in the original private repository, and the tag `v2.0.0` has exactly the same source tree.

To reproduce the artifact, check out the tag `v2.0.0` and build it with the release script from `v2.0.1`. The `v2.0.0` copy of the script has no `UPSIDE_BUILD_COMMIT` override, which sets the embedded commit. The script needs Docker, as described in [Reproduce a release with Docker](../getting-started/install.md#reproduce-a-release-with-docker). From a clean checkout:

```bash
git checkout v2.0.0
git show v2.0.1:scripts/build-release.sh > "$TMPDIR/build-release.sh"
UPSIDE_BUILD_COMMIT=8521a42fc1806e06c422e67588dd9dd834a09b76 bash "$TMPDIR/build-release.sh" v2.0.0
cat build/release/upsided-v2.0.0-linux-amd64.sha256
```

The checksum must be `1279d570c0a031bcd0b7c7ce4aa379b3de9f6bae984feffe894bc357de6ca40b`. If `TMPDIR` isn't set, write the script to any directory outside the checkout.

The `v2.0.1` release changes the repository, the Go module path, the documentation, and the release build script. It doesn't change the state machine, and it keeps the `v2.0.0` upgrade handler, so a node has nothing to build or switch to. A `v2.0.1` build is compatible and computes the same state, but it isn't the binary that governance approved, and its checksum differs from the `v2.0.0` artifact. The `v2.0.1` changes will be included in the next chain software upgrade, which governance schedules.

A new node using the current release starts from state taken after the upgrade. Use state sync from the public sentries, or restore a signed daily archive from `https://snapshots.mainnet.upside.net`. [Install the supported sync state](../nodes/public-node.md#4-install-the-supported-sync-state) gives the steps and the verification checks for both.

Replaying from genesis requires the historical binary sequence: the pre-v2 binary for heights 1 to 3,898,899, then `v2.0.0` from height 3,898,900. The original genesis predates `mint_half_life_seconds`, and the current release's genesis validator requires that parameter. The current minting function also differs from the pre-upgrade function. Keep the original genesis unchanged and use the matching historical software when replaying those blocks. The pre-v2 binary isn't published yet. [Replay mainnet across the v2.0.0 upgrade](../nodes/operations.md#replay-mainnet-across-the-v200-upgrade) describes the sequence.

Native issuance follows the **2.1-billion-UP** minting cap described in [Supply and minting](../technical/supply-and-minting.md).
