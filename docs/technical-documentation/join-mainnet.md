# Join Mainnet (Full Node / Validator)

This guide covers joining Upside **mainnet** as a full node, and (if the network is permissionless for validators) as a validator.

Mainnet parameters (chain-id, genesis URL, seeds/peers, minimum gas prices, snapshots, state-sync RPCs) must come from the official mainnet ops documentation.

## 0) Get the binaries

Mainnet has changed binaries at a software upgrade, so a node that syncs from
genesis needs more than one binary:

| Blocks                          | Binary                                       | Download                                                                                  |
| ------------------------------- | -------------------------------------------- | ----------------------------------------------------------------------------------------- |
| 1 to `<V2_UPGRADE_HEIGHT>` − 1  | pre-upgrade binary (mainnet runs a build reporting commit `fd12fc6983c1…`) | `<PRE_V2_BINARY_URL>`, SHA-256 `<PRE_V2_BINARY_SHA256>`, source `<PRE_V2_SOURCE_COMMIT>` |
| `<V2_UPGRADE_HEIGHT>` onward    | `v2.0.0`                                     | `<V2_BINARY_URL>`, SHA-256 `<V2_BINARY_SHA256>`                                           |

Neither binary can replace the other:

- `v2.0.0` cannot start from the mainnet genesis file. Genesis has no
  `mint_half_life_seconds` param, and `v2.0.0` rejects a half-life of zero.
  Before the upgrade it would also mint different amounts than mainnet did,
  so its state would not match.
- The pre-upgrade binary stops at the upgrade height with
  `UPGRADE "v2.0.0" NEEDED at height: <V2_UPGRADE_HEIGHT>`.

Use [Cosmovisor](https://docs.cosmos.network/main/build/tooling/cosmovisor),
which switches binaries at the upgrade height:

```bash
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest

export DAEMON_NAME=upsided
export DAEMON_HOME="$HOME/.upside"

mkdir -p "$DAEMON_HOME/cosmovisor/genesis/bin" "$DAEMON_HOME/cosmovisor/upgrades/v2.0.0/bin"
cp upsided-pre-v2 "$DAEMON_HOME/cosmovisor/genesis/bin/upsided"
cp upsided-v2.0.0 "$DAEMON_HOME/cosmovisor/upgrades/v2.0.0/bin/upsided"
chmod +x "$DAEMON_HOME"/cosmovisor/genesis/bin/upsided "$DAEMON_HOME"/cosmovisor/upgrades/v2.0.0/bin/upsided

# Check each file against the SHA-256 in the table before running it.
shasum -a 256 "$DAEMON_HOME"/cosmovisor/genesis/bin/upsided "$DAEMON_HOME"/cosmovisor/upgrades/v2.0.0/bin/upsided
```

The upgrade directory name must be exactly `v2.0.0`, the name in the upgrade
plan.

If you start from a snapshot or with state sync at a height after the
upgrade, you only need `v2.0.0`.

Use the pre-upgrade binary (from `cosmovisor/genesis/bin`) for the `init` command
below.

## 1) Initialize home directory

Default home is `~/.upside`.

```bash
MONIKER="my-mainnet-node"
upsided init "$MONIKER" --home "$HOME/.upside"
```

## 2) Fetch and install `genesis.json`

Your mainnet genesis file is available at `docs/technical-documentation/genesis-mainnet.json` ([genesis-mainnet.json](./genesis-mainnet.json)).

```bash
cp ./docs/technical-documentation/genesis-mainnet.json "$HOME/.upside/config/genesis.json"

CHAIN_ID=$(jq -r .chain_id "$HOME/.upside/config/genesis.json")
echo "$CHAIN_ID"
```

## 3) Configure P2P (persistent peers)

Edit `~/.upside/config/config.toml` and set `persistent_peers` in the `[p2p]` section using the appropriate list below.

**For Validator**

```toml
persistent_peers = "c0135b9048562977832d28efb4cc1688b35741e0@162.19.139.135:26656,89a31846dedc0cd51af08f8ef28a844e35a5303d@141.95.126.110:26656,1c9d3e7d81d36868fee200b4877eb0ef4bc26563@162.19.139.78:26656"
```

**Note:** These peers require IP whitelisting. Have your validator node's IP whitelisted before connecting.

**For RPC (public access)**

```toml
persistent_peers = "ff1d5d41e1e6c273234b1d7f231d3cb58e2b0d70@162.19.234.18:26656,a2620fcd4ffcc6c1fe7bc145554c3b5180a04f05@162.19.234.17:26656"
```

## 4) Configure minimum gas prices + APIs

```bash
APP_TOML="$HOME/.upside/config/app.toml"

# Example placeholder:
# MIN_GAS_PRICES="0.025aup"
sed -i.bak -E "s/^minimum-gas-prices *=.*/minimum-gas-prices = \"${MIN_GAS_PRICES}\"/" "$APP_TOML"
```

## 5) Start and sync

```bash
cosmovisor run start --home "$HOME/.upside" --minimum-gas-prices "${MIN_GAS_PRICES}"
```

If you sync from a snapshot taken after the upgrade and do not use
Cosmovisor, run `upsided start` with the `v2.0.0` binary instead.

Wait until:

```bash
curl -s http://localhost:26657/status | jq -r .result.sync_info.catching_up
```

returns `false`.

## Upgrading a running node to v2.0.0

The `v2.0.0` upgrade runs at height `<V2_UPGRADE_HEIGHT>`, set by governance
proposal #2.

- With Cosmovisor: put the `v2.0.0` binary in
  `$DAEMON_HOME/cosmovisor/upgrades/v2.0.0/bin/upsided` before that height.
  Cosmovisor stops the node at the height, switches binaries, and restarts.
- Without Cosmovisor: the node stops at that height and logs
  `UPGRADE "v2.0.0" NEEDED at height: <V2_UPGRADE_HEIGHT>`. Replace the
  binary with `v2.0.0` and start the node again.
- Do not start `v2.0.0` before the upgrade height. It stops with
  `BINARY UPDATED BEFORE TRIGGER! UPGRADE "v2.0.0" - in binary but not executed on chain`.

After the upgrade, `upsided query mint inflation` returns the rate from the
capped supply curve, which falls as supply grows (about 0.84 at 615 million
UP), and `upsided query upside params` shows `mint_half_life_seconds` set to
`63115200`.

## 6) (Optional) Become a validator on mainnet

Whether this is allowed depends on mainnet governance/permissioning.

### 6.1 Create/recover a key

```bash
KEY_NAME="my-mainnet-val"
upsided keys add "$KEY_NAME" --keyring-backend os --algo eth_secp256k1

VAL_ADDR=$(upsided keys show "$KEY_NAME" -a --keyring-backend os)
echo "$VAL_ADDR"
```

Fund `VAL_ADDR` with mainnet tokens via your normal treasury/onboarding process.

### 6.2 Create validator tx

```bash
PUB_KEY_JSON=$(upsided tendermint show-validator --home "$HOME/.upside")

AMOUNT="1000000aup" # replace (must include denom)

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

## Operational checklist (mainnet)

- The running binary uses **mainnet prefixes** (`up...`).
- The binary matches the block height: the pre-upgrade binary before `<V2_UPGRADE_HEIGHT>`, `v2.0.0` from that height.
- `genesis.json` is the official mainnet genesis and chain-id matches.
- `seeds`/`persistent_peers` are configured and reachable.
- `minimum-gas-prices` matches network requirements.
- If validating: keys are securely backed up and the node is fully synced.
