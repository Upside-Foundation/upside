---
title: Create a validator
description: Register and verify an Upside L1 validator after synchronizing a full node.
slug: /validators/create
---

# Create a validator

A validator signs consensus messages when it belongs to the active set, which is the group of validators selected to produce and finalize blocks. Its operator account controls staking and validator settings. Delegators bond tokens to its operator address and share rewards after commission.

Use this guide to register a new validator. To move an existing validator to another machine, follow the [migration guidance](../nodes/operations.md).

## Registration

- Meet the [compute requirements](../nodes/requirements.md). The approximate mainnet baseline is 8 CPU cores, 16 GB RAM, 1 TB storage, and 1 Gbps networking.
- Complete [public-node setup](../nodes/public-node.md) for the chosen network and verify recent, advancing blocks.
- Have a funded operator [wallet](../getting-started/wallets.md).


| Key or address                                                      | Purpose                                                                      |
| ------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| Operator wallet key                                                 | Authorizes validator registration, edits, staking, and commission withdrawal |
| Operator address (`upvaloper1…` on both public networks)            | Identifies the validator in staking transactions                             |
| Consensus key (`config/priv_validator_key.json` for a local signer) | Signs consensus messages                                                     |
| Signing state (`data/priv_validator_state.json` for a local signer) | Records the last consensus signature to help prevent double signing          |
| P2P key (`config/node_key.json`)                                    | Identifies the node to peers                                                 |


Keep the operator wallet on a trusted machine. The node signs consensus messages with its consensus key, and registration uses the corresponding public key.

## Protect the validator with sentry nodes

We strongly recommend running every mainnet validator behind sentry nodes.

A sentry is an ordinary full node that you operate. It connects to the public network on the validator's behalf, so the validator's own address is not exposed to other peers. Operators commonly run more than one sentry, with different hosting providers or regions.

Two principles apply to any such setup:

- Publish only your sentries' addresses. Keep the validator's address private.
- The consensus key belongs on the validator only. Never copy the validator's consensus private key onto a sentry. A sentry has its own independently generated key file.

Follow [Configure validator sentries](../nodes/sentries.md) for the P2P settings and topology checks.

## Collect the public key

On the synchronized node, with its `UPSIDE_BIN` and `NODE_HOME`:

```bash
"$UPSIDE_BIN" comet show-validator --home "$NODE_HOME"
```

Save the JSON public key output. Never substitute the wallet public key or copy `priv_validator_key.json` into the transaction.

Older examples use `tendermint show-validator`. This release accepts `tendermint` as an alias, but its primary command is `comet`. To check your installed binary, run `"$UPSIDE_BIN" comet show-validator --help`.

## Prepare the validator definition

On the wallet machine, set its network profile, `KEY_NAME`, and `ACCOUNT_ADDR`. Create `validator.json` with the following structure, replacing the public key and operator information:

```json
{
  "pubkey": {
    "@type": "/cosmos.crypto.ed25519.PubKey",
    "key": "REPLACE_WITH_NODE_CONSENSUS_PUBLIC_KEY_BASE64"
  },
  "amount": "1000000000000000000000aup",
  "moniker": "my-validator",
  "identity": "",
  "website": "",
  "security": "",
  "details": "",
  "commission-rate": "0.10",
  "commission-max-rate": "0.20",
  "commission-max-change-rate": "0.01",
  "min-self-delegation": "100000000000000000000"
}
```

The example uses these settings:


| Setting                         | Example value                            |
| ------------------------------- | ---------------------------------------- |
| Initial self-delegation         | 1,000 UP                                 |
| Minimum self-delegation         | 100 UP                                   |
| Commission                      | 10% of the validator's allocated rewards |
| Maximum commission              | 20%                                      |
| Maximum daily commission change | One percentage point                     |


Choose the self-delegation and commission settings for your operation within the network's current limits. These amounts are examples. Commission applies to the validator's staking rewards. See [rewards and the governance treasury](../staking/governance.md).

Commission caps constrain future edits. You can increase minimum self-delegation through an edit transaction. Falling below that commitment can jail the validator. Active-set membership depends on the validator's eligible stake and the network's validator limit.

## Submit registration

```bash
"$UPSIDE_BIN" tx staking create-validator validator.json \
  --from "$KEY_NAME" --chain-id "$CHAIN_ID" --node "$RPC_URL" \
  --home "$WALLET_HOME" --keyring-backend "$KEYRING_BACKEND" \
  --gas auto --gas-adjustment 1.5 --gas-prices "$GAS_PRICES"
```

Inspect the transaction before approving it. Use [wallet transaction confirmation](../getting-started/wallets.md#confirm-a-transaction) after submission.

## Verify registration and active status

Derive the operator address from the wallet and query its record:

```bash
export VALIDATOR="$("$UPSIDE_BIN" keys show "$KEY_NAME" -a --bech val \
  --keyring-backend "$KEYRING_BACKEND" --home "$WALLET_HOME")"

"$UPSIDE_BIN" query staking validator "$VALIDATOR" \
  --node "$RPC_URL" --home "$WALLET_HOME" --output json
```

Check the returned record:

1. Match its consensus public key to the node's public key.
2. Verify the stake and commission.
3. Check that `jailed` is `false`. This binary can omit the field when it is false.
4. Check for `BOND_STATUS_BONDED`, which indicates active-set membership.

A registered validator can remain unbonded if it does not have enough stake to enter the active set.

On the node, check status and signing activity. Its consensus public key must match the registered validator and its consensus voting power should become positive when active. Check block signatures and missed-block counters to confirm that the validator is participating.

## Commission and ongoing changes

Query accrued commission:

```bash
"$UPSIDE_BIN" query distribution commission "$VALIDATOR" \
  --node "$RPC_URL" --home "$WALLET_HOME" --output json
```

Withdraw self-delegation rewards and commission:

```bash
"$UPSIDE_BIN" tx distribution withdraw-rewards "$VALIDATOR" --commission \
  --from "$KEY_NAME" --chain-id "$CHAIN_ID" --node "$RPC_URL" \
  --home "$WALLET_HOME" --keyring-backend "$KEYRING_BACKEND" \
  --gas auto --gas-adjustment 1.5 --gas-prices "$GAS_PRICES"
```

To change the displayed validator name:

```bash
"$UPSIDE_BIN" tx staking edit-validator --new-moniker "my-updated-validator" \
  --from "$KEY_NAME" --chain-id "$CHAIN_ID" --node "$RPC_URL" \
  --home "$WALLET_HOME" --keyring-backend "$KEYRING_BACKEND" \
  --gas auto --gas-adjustment 1.5 --gas-prices "$GAS_PRICES"
```

Query the validator record and check `description.moniker`. Review the installed command's `--help` before changing commission or minimum self-delegation. Use the [staking guide](../staking/delegation.md) to change self-delegation.

Continue with [operations and monitoring](../nodes/operations.md).
