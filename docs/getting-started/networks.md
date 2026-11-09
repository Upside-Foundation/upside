---
title: Networks and connection settings
description: Choose mainnet, testnet, or localnet and configure the matching connection settings.
slug: /getting-started/networks
---

# Networks and connection settings

Choose a network before creating keys or sending transactions. Mainnet, public testnet, and localnet each have their own balances and state. Test funds can be used only on the network that issued them.


| Setting               | Mainnet                    | Public testnet                                                       | Localnet in this handbook              |
| --------------------- | -------------------------- | -------------------------------------------------------------------- | -------------------------------------- |
| Purpose               | Public L1                  | Shared testing                                                       | Private practice                       |
| Cosmos chain ID       | `upside_61873-1`           | `upside_61872-1`                                                     | `upside_61874-1`                       |
| EVM chain ID          | `61873`                    | `61872`                                                              | `61874`                                |
| Account prefix        | `up`                       | `up`                                                                 | `up`                                    |
| Base denomination     | `aup`                      | `aup`                                                                | `aup`                                  |
| Decimal places        | 18                         | 18                                                                   | 18                                     |
| Node home used here   | `~/.upside-mainnet`        | `~/.upside-testnet`                                                  | A newly created localnet directory     |
| Wallet home used here | `~/.upside-wallet-mainnet` | `~/.upside-wallet-testnet`                                           | Localnet directory for disposable keys |


The [localnet guide](../nodes/localnet.md) uses a separate chain ID for local practice.

## Connection types


| Connection         | Used by                                                                   |
| ------------------ | ------------------------------------------------------------------------- |
| CometBFT RPC       | `upsided --node`, chain status, Cosmos transaction submission and queries |
| EVM JSON-RPC       | Ethereum-compatible wallets and applications                              |
| Cosmos REST / gRPC | Account, staking, and module integrations (no public endpoint; run your own node) |
| P2P peer           | Nodes synchronizing blocks and participating in consensus                 |


Use a CometBFT RPC URL with `--node`. EVM wallets use an EVM JSON-RPC URL, and nodes connect to P2P peers using `node-id@host:port`.

## Shell profile used by the wallet guides

For mainnet, select the installed mainnet binary and set the RPC and gas-price values used by your provider or node:

```bash
export UPSIDE_BIN="$HOME/.local/bin/upsided-mainnet"
export CHAIN_ID="upside_61873-1"
export DENOM="aup"
export WALLET_HOME="$HOME/.upside-wallet-mainnet"
export KEYRING_BACKEND="file"
export RPC_URL="https://rpc.mainnet.upside.net/tendermint/"
export GAS_PRICES="1000000000aup"
```

`GAS_PRICES` is an example transaction gas price of one billion `aup` per gas, which is far above the current base fee. The base fee changes with network use. Before signing, query the current value with [`eth_gasPrice`](../technical/json-rpc.md#eth_gasprice) and set `GAS_PRICES` at or above it. A fixed value can become insufficient as the base fee changes. See [transaction fees](transactions.md).

Keep the trailing slash on `RPC_URL`. Without it, the server answers with a redirect and `upsided` commands that use `--node` fail.

The encrypted `file` keyring works without a desktop credential service. You can also use an `os` keyring on a configured desktop. Use the same backend and home for every command that accesses a given key.

For testnet, use the values on the [testnet page](../networks/testnet.md). The [localnet guide](../nodes/localnet.md) supplies a complete profile for its disposable accounts.

Check the endpoint's network before using it:

```bash
curl --fail --silent --show-error "${RPC_URL%/}/status" \
  | jq '.result | {chain_id: .node_info.network, sync: .sync_info}'
```

Use a published provider or your own synchronized node. Check that the returned chain ID matches `CHAIN_ID`, the block time is recent, and the height advances between requests.

## Mainnet token amounts


| UP    | CLI amount                  |
| ----- | --------------------------- |
| 0.01  | `10000000000000000aup`      |
| 1     | `1000000000000000000aup`    |
| 100   | `100000000000000000000aup`  |
| 1,000 | `1000000000000000000000aup` |


Use decimal strings or arbitrary-precision arithmetic for token amounts. These quantities can exceed the safe range of Bash integers and JavaScript floating-point numbers. The CLI accepts base units, regardless of the display unit chosen by a wallet.

## Launch issuance

Native issuance follows a half-life curve toward **2.1 billion UP**. The [supply reference](../technical/supply-and-minting.md) defines the minting cap and its scope. Mainnet currently directs 93% of block rewards to the governance treasury and 7% to staking rewards. The split is the `community_tax` parameter of `x/distribution`, which governance can change. See [rewards and governance](../staking/governance.md) for the allocation, Foundation policy, and what each means for delegators. Use the active on-chain parameters for the current issuance rate.
