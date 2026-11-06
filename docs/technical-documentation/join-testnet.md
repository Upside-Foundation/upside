# Join Testnet (Full Node / Validator)

This guide covers joining an existing Upside **testnet** as:

- a **full node** (sync + serve RPC)
- optionally, a **validator** (create-validator transaction)

Because testnet parameters (chain-id, genesis URL, seeds/peers, minimum gas prices, snapshots) are network-ops controlled, you must obtain those values from the testnet launch/ops source of truth.

## 0) Get the correct binary

To join testnet, you need the same binary version the testnet runs. There
is no published download at the moment: the previous Backblaze link returns
HTTP 401. Ask the testnet operators for the binary and its SHA-256, or build
from source with the testnet build tag:

```bash
NETWORK=testnet make build   # produces ./build/upsided with the testnet tag
mv ./build/upsided /usr/local/bin/
upsided version --long        # build_tags must include "testnet"
```

A binary built from source only matches the testnet's if it is built from
the same commit the testnet runs.

## 1) Initialize home directory

Default home is `~/.upside`.

```bash
MONIKER="my-testnet-node"
upsided init "$MONIKER" --home "$HOME/.upside"
```

## 2) Fetch and install `genesis.json`

Upside testnet genesis file is available at `docs/technical-documentation/genesis-testnet.json` ([genesis-testnet.json](./genesis-testnet.json)).

```bash
cp ./docs/technical-documentation/genesis-testnet.json "$HOME/.upside/config/genesis.json"

# Verify chain-id recorded in genesis
jq -r .chain_id "$HOME/.upside/config/genesis.json"
```

## 3) Configure P2P (persistent peers)

Edit `~/.upside/config/config.toml` and set `persistent_peers` in the `[p2p]` section using the appropriate list below.

**For Validator**

```toml
persistent_peers = "9c43ac9d3f4edf87caee48d102d2099d9369f672@162.19.232.154:26656,a02b9006343bb7b5d5be096c77fcf4dda69127ce@162.19.232.149:26656,4c03728a407e7183726cfe0fa309174baf8b0ad4@162.19.232.151:26656"
```

**For RPC (public access)**

```toml
persistent_peers = "9c43ac9d3f4edf87caee48d102d2099d9369f672@162.19.232.154:26656"
```

## 4) Configure minimum gas prices + APIs

Minimum gas prices are chain-specific. Set them to match testnet requirements.

```bash
APP_TOML="$HOME/.upside/config/app.toml"

# Example placeholder:
# MIN_GAS_PRICES="0.025taup"
sed -i.bak -E "s/^minimum-gas-prices *=.*/minimum-gas-prices = \"${MIN_GAS_PRICES}\"/" "$APP_TOML"
```

If you need EVM JSON-RPC / gRPC / REST for dev integration, enable the relevant endpoints in `app.toml` and start with the needed flags. A common dev setup is:

```bash
upsided start \
  --home "$HOME/.upside" \
  --rpc.unsafe \
  --json-rpc.api "eth,net,web3,debug,txpool" \
  --api.enable \
  --minimum-gas-prices "${MIN_GAS_PRICES}"
```

## 5) Start the node and wait for sync

```bash
upsided start --home "$HOME/.upside"
```

Check sync status:

```bash
curl -s http://localhost:26657/status | jq -r .result.sync_info.catching_up
```

Proceed only when it returns `false`.

## 6) (Optional) Become a validator on testnet

### 6.1 Create or recover a key

```bash
KEY_NAME="my-testnet-val"
upsided keys add "$KEY_NAME" --keyring-backend os --algo eth_secp256k1

VAL_ADDR=$(upsided keys show "$KEY_NAME" -a --keyring-backend os)
echo "$VAL_ADDR"
```

Fund `VAL_ADDR` with testnet tokens using the faucet / token distribution process.

### 6.2 Create validator tx

You must use the **exact testnet chain-id** (from genesis):

```bash
CHAIN_ID=$(jq -r .chain_id "$HOME/.upside/config/genesis.json")
PUB_KEY_JSON=$(upsided tendermint show-validator --home "$HOME/.upside")

AMOUNT="1000000taup" # replace (must include denom)

# Create the validator definition file required by the CLI.
cat > validator.json <<EOF
{
  "pubkey": $PUB_KEY_JSON,
  "amount": "$AMOUNT",
  "moniker": "$MONIKER",
  "identity": "",
  "website": "",
  "security": "",
  "details": "",
  "commission-rate": "0.10",
  "commission-max-rate": "0.20",
  "commission-max-change-rate": "0.01",
  "min-self-delegation": "1"
}
EOF

upsided tx staking create-validator validator.json \
  --chain-id "$CHAIN_ID" \
  --gas auto \
  --gas-adjustment 1.5 \
  --gas-prices "${MIN_GAS_PRICES}" \
  --from "$KEY_NAME" \
  --keyring-backend os \
  --home "$HOME/.upside" \
  -y
```

## Operational checklist (testnet)

- The running binary was built with `NETWORK=testnet` (Bech32 prefix `tup`).
- `genesis.json` matches the testnet you’re joining.
- `seeds`/`persistent_peers` are set and reachable.
- `minimum-gas-prices` matches network requirements.
- If you’re validating: keys are backed up and your node is fully synced before `create-validator`.
