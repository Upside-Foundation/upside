---
title: Run a localnet
description: Run a private Upside L1 practice chain with disposable accounts.
slug: /nodes/localnet
---

# Run a localnet

A localnet is a private chain with its own genesis and disposable funds. Use it for development, and to practice wallet commands, transfers, delegation, and validator operations.

These instructions create a new home, fresh development keys, and a separate chain ID. Accounts use `up1…` addresses and hold local test funds in `aup`. The chain uses the binary's genesis defaults with the changes below, including the default local issuance settings.

## 1. Prepare a separate home

Install the [`upsided` binary](../getting-started/install.md). You also need Bash, `jq`, `curl`, and available ports `27656`, `27657`, `18545`, `18546`, and `19090`. Do not use the repository's existing `start-local-node.sh`. It runs `rm -rf ~/.upside*`, which deletes every file and directory in your home directory whose name begins with `.upside`. That includes the homes that these guides use: `~/.upside` (the default home of `upsided`), `~/.upside-mainnet`, `~/.upside-testnet`, `~/.upside-wallet-mainnet`, `~/.upside-wallet-testnet`, and the `~/.upside-localnet.*` homes that this guide creates.

Open a dedicated Bash shell, so that a failed command closes only that shell:

```bash
bash
```

In that shell, run the setup in order. `set -euo pipefail` stops the setup at the first failed command:

```bash
set -euo pipefail
export UPSIDE_BIN="$HOME/.local/bin/upsided-mainnet"
export CHAIN_ID="upside_61874-1"
export DENOM="aup"
export NODE_HOME="$(mktemp -d "$HOME/.upside-localnet.XXXXXX")"
export WALLET_HOME="$NODE_HOME"
export KEYRING_BACKEND="test"
export RPC_URL="http://127.0.0.1:27657"
export GAS_PRICES="1000000000aup"
export KEY_NAME="local-validator"
printf 'Localnet home: %s\n' "$NODE_HOME"

"$UPSIDE_BIN" init "local-validator" --chain-id "$CHAIN_ID" \
  --default-denom "$DENOM" --home "$NODE_HOME"
```

Record the printed home path. The `test` keyring stores unencrypted keys.

## 2. Create local accounts

```bash
"$UPSIDE_BIN" keys add "$KEY_NAME" --key-type eth_secp256k1 --coin-type 60 \
  --keyring-backend test --home "$NODE_HOME"
"$UPSIDE_BIN" keys add "local-user" --key-type eth_secp256k1 --coin-type 60 \
  --keyring-backend test --home "$NODE_HOME"

export ACCOUNT_ADDR="$("$UPSIDE_BIN" keys show "$KEY_NAME" -a \
  --keyring-backend test --home "$NODE_HOME")"
export LOCAL_USER_ADDR="$("$UPSIDE_BIN" keys show local-user -a \
  --keyring-backend test --home "$NODE_HOME")"
```

## 3. Configure local genesis and balances

Use the same denomination for EVM balances and native staking. Set a block gas limit and an initial base fee for the practice chain:

```bash
GENESIS_FILE="$NODE_HOME/config/genesis.json"
jq --arg denom "$DENOM" '
  .app_state.evm.params.evm_denom = $denom |
  .app_state.feemarket.params.base_fee = "1000000000" |
  .consensus.params.block.max_gas = "10000000"
' "$GENESIS_FILE" > "$NODE_HOME/config/genesis.local.json"
mv "$NODE_HOME/config/genesis.local.json" "$GENESIS_FILE"

"$UPSIDE_BIN" genesis add-genesis-account "$ACCOUNT_ADDR" "10000000000000000000000${DENOM}" \
  --home "$NODE_HOME"
"$UPSIDE_BIN" genesis add-genesis-account "$LOCAL_USER_ADDR" "10000000000000000000000${DENOM}" \
  --home "$NODE_HOME"
```

Each account receives 10,000 local tokens when this new chain starts.

Create the initial validator with 1,000 local tokens:

```bash
"$UPSIDE_BIN" genesis gentx "$KEY_NAME" "1000000000000000000000${DENOM}" \
  --chain-id "$CHAIN_ID" --keyring-backend test --home "$NODE_HOME"
"$UPSIDE_BIN" genesis collect-gentxs --home "$NODE_HOME"
"$UPSIDE_BIN" genesis validate --home "$NODE_HOME"
```

If genesis validation fails, check that staking, mint, governance deposits, and EVM denominations are consistent. The `--default-denom` initialization flag supplies native module defaults.

## 4. Start on loopback

The profiling listener is disabled to keep it from occupying its default port, `6060`. gRPC runs on loopback to enable the transaction service used by automatic gas simulation through RPC.

```bash
"$UPSIDE_BIN" start --home "$NODE_HOME" \
  --p2p.laddr tcp://127.0.0.1:27656 \
  --rpc.laddr tcp://127.0.0.1:27657 \
  --rpc.pprof_laddr "" \
  --minimum-gas-prices "$GAS_PRICES" \
  --api.enable=false --grpc.enable=true \
  --grpc.address 127.0.0.1:19090 --grpc-web.enable=false \
  --json-rpc.enable=true --json-rpc.api eth,net,web3 \
  --json-rpc.address 127.0.0.1:18545 \
  --json-rpc.ws-address 127.0.0.1:18546
```

Keep this terminal running. In another terminal:

```bash
curl --fail --silent --show-error http://127.0.0.1:27657/status \
  | jq '.result | {chain_id: .node_info.network, sync: .sync_info}'
```

Check for chain ID `upside_61874-1`, a recent block time, advancing height, and eventually `catching_up: false`. This local chain runs with one validator.

## 5. Local wallets

In a separate terminal, set the same `UPSIDE_BIN`, `CHAIN_ID`, `DENOM`, `RPC_URL`, `GAS_PRICES`, and `KEY_NAME` as above. Set `WALLET_HOME` to the exact home path printed by step 1 and `KEYRING_BACKEND=test`. Retrieve `ACCOUNT_ADDR` with `keys show` as in step 2.

Follow [transfers](../getting-started/transactions.md), using the `local-user` address as recipient. For staking, obtain the initial validator operator address:

```bash
export VALIDATOR="$("$UPSIDE_BIN" keys show local-validator -a --bech val \
  --keyring-backend test --home "$WALLET_HOME")"
```

Set `KEY_NAME=local-user` and retrieve its `ACCOUNT_ADDR` before using the [staking guide](../staking/delegation.md) as a delegator. The local chain uses the default unbonding period, so unbonding does not complete immediately.
