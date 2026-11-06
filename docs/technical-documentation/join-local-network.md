# Joining a Running Local Node as a Validator

This guide provides step-by-step instructions on how to join a pre-existing local Upside network as a new validator.

## Prerequisites

- **Upside Daemon (`upsided`)**: Ensure it is compiled and accessible in your system's `PATH`.
- **Network Access**: You need the IP address and the P2P connection details of a running node (the "seed" or "RPC" node) to sync the blockchain.
- **Genesis File**: You must obtain the exact `genesis.json` file from the existing network.

## Step 1: Initialize Your Node

First, initialize your node. This creates your data directories and default configuration files. Replace `my-new-validator` with your preferred moniker.

```bash
MONIKER="my-new-validator"
CHAINID="<CHAIN_ID_FROM_GENESIS>" # Must match the existing network's chain ID
NODE_HOME="$HOME/.upside_new"

upsided init "$MONIKER" --chain-id "$CHAINID" --home "$NODE_HOME"
```

## Step 2: Replace the Genesis File

Your node needs the identical genesis state to sync with the network. Copy the `genesis.json` from the existing node to your new node's configuration directory.

```bash
# Assuming the existing node genesis is in the current directory or a known path
cp path/to/existing/genesis.json "$NODE_HOME/config/genesis.json"
```

## Step 3: Configure Network Peers

To discover and sync with the network, your node needs to know where to connect. We configure this in `config.toml`.

Get the Node ID of the existing running node you wish to connect to. On the existing node, run:

```bash
upsided tendermint show-node-id --home <path_to_existing_node_home>
```

Construct the peer string combining the ID, IP address, and P2P port. For example, if the ID is `1234abcd` and the IP is `192.168.1.100` on port `26656`.

```bash
CONFIG_FILE="$NODE_HOME/config/config.toml"
PEER_ID="1234abcd@192.168.1.100:26656"

# Add the peer to persistent_peers
sed -i.bak "s/^persistent_peers *=.*/persistent_peers = \"$PEER_ID\"/" "$CONFIG_FILE"
```

## Step 4: Configure Application Settings

Set your minimum gas prices in `app.toml` to match the network requirements.

```bash
APP_FILE="$NODE_HOME/config/app.toml"

sed -i.bak 's/^minimum-gas-prices *=.*/minimum-gas-prices = "<MIN_GAS_PRICES>"/' "$APP_FILE"
```

## Step 5: Start Your Node and Sync

Start your node. It will connect to the peer you specified and begin downloading the blockchain data.

```bash
upsided start --home "$NODE_HOME"
```

**Important:** Wait for your node to fully sync to the latest block before proceeding to the next step. You can check your sync status by querying the RPC endpoint (e.g., `curl -s http://localhost:26657/status | jq .result.sync_info.catching_up`). It must return `false` before you can become a validator.

## Step 6: Create or Recover Your Validator Key

You need an account with funds (AUP tokens) to stake. You can either create a new key and fund it, or recover an existing one.

```bash
KEY_NAME="my-new-val-key"

# To create a new key:
upsided keys add "$KEY_NAME" --keyring-backend "file" --algo "eth_secp256k1" --home "$NODE_HOME"

# Remember your mnemonic!
VAL_ADDRESS=$(upsided keys show "$KEY_NAME" -a --keyring-backend "file" --home "$NODE_HOME")
```

If you created a new account, you **must fund it** by sending AUP tokens from an existing account on the network to `$VAL_ADDRESS`.

## Step 7: Submit the Create Validator Transaction

Once your node is fully synced and your account has sufficient funds, broadcast a transaction to the network to declare yourself as a validator.

```bash
AMOUNT="<STAKE_AMOUNT><DENOM>" # e.g. "1000000aup" (must match the chain denom)
PUB_KEY_JSON=$(upsided tendermint show-validator --home "$NODE_HOME")

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
  --chain-id "$CHAINID" \
  --gas "auto" \
  --gas-adjustment "1.5" \
  --gas-prices "<MIN_GAS_PRICES>" \
  --from "$KEY_NAME" \
  --keyring-backend "file" \
  --home "$NODE_HOME" \
  -y
```

### Transaction Parameters Explained

- `validator.json.amount`: The amount of AUP you are bonding/staking (must include the chain denom, e.g. `1000000aup`).
- `validator.json.pubkey`: Your node's consensus public key (from `upsided tendermint show-validator`).
- `validator.json.moniker`: The validator name shown on explorers.
- `validator.json.commission-rate`: The percentage of rewards you keep from delegators (e.g., `0.10` = 10%).
- `validator.json.commission-max-rate`: The maximum commission you can ever charge.
- `validator.json.commission-max-change-rate`: The maximum your commission can change per day.
- `validator.json.min-self-delegation`: The minimum stake you must maintain to remain an active validator.
- `--from`: The local key name (or address) that signs and pays for the transaction.

Once the transaction is included in a block, your node will become an active validator and participate in consensus!
