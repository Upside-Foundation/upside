---
title: Staking and delegation
description: Delegate UP, withdraw rewards, redelegate, and unbond stake.
slug: /staking/delegation
---

# Staking and delegation

Delegation locks tokens in support of a validator. You keep authority to move the stake to another validator or request its return. The validator cannot spend your delegated tokens, but its performance affects your rewards and its penalties can reduce your stake.

Mainnet currently directs 93% of block rewards to the governance treasury and 7% to staking rewards. Governance can change this split, which is the `community_tax` parameter. Validators take commission from their share of the staking allocation. Read [rewards and the governance treasury](governance.md) before comparing your rewards with the network's issuance rate.

## Choose a validator and inspect parameters

```bash
"$UPSIDE_BIN" query staking validators \
  --node "$RPC_URL" --home "$WALLET_HOME" --output json
"$UPSIDE_BIN" query staking params \
  --node "$RPC_URL" --home "$WALLET_HOME" --output json
"$UPSIDE_BIN" query slashing params \
  --node "$RPC_URL" --home "$WALLET_HOME" --output json
```

Queries can be paginated. Review status, commission, self-delegation, operating history, and voting-power concentration. Check whether the validator is bonded and jailed. Choose a bonded validator to participate in active-validator rewards.

The current unbonding period and maximum entries come from `staking params`.

## Delegate

Set `VALIDATOR` to the chosen operator address, which begins with `upvaloper1` on both public networks. This example delegates 100 tokens:

```bash
export VALIDATOR="REPLACE_WITH_VALIDATOR_OPERATOR_ADDRESS"
"$UPSIDE_BIN" tx staking delegate "$VALIDATOR" "100000000000000000000${DENOM}" \
  --from "$KEY_NAME" --chain-id "$CHAIN_ID" --node "$RPC_URL" \
  --home "$WALLET_HOME" --keyring-backend "$KEYRING_BACKEND" \
  --gas auto --gas-adjustment 1.5 --gas-prices "$GAS_PRICES"
```

After the transaction succeeds, inspect your positions. Use [wallet transaction confirmation](../getting-started/wallets.md#confirm-a-transaction) after submission:

```bash
"$UPSIDE_BIN" query staking delegations "$ACCOUNT_ADDR" \
  --node "$RPC_URL" --home "$WALLET_HOME" --output json
"$UPSIDE_BIN" query distribution rewards "$ACCOUNT_ADDR" \
  --node "$RPC_URL" --home "$WALLET_HOME" --output json
```

## Withdraw rewards

Check the configured withdrawal address before claiming. It defaults to your account. If you previously changed it, rewards will go to that destination.

```bash
"$UPSIDE_BIN" query distribution delegator-withdraw-address \
  --delegator-address "$ACCOUNT_ADDR" \
  --node "$RPC_URL" --home "$WALLET_HOME" --output json
```

Claim the rewards from one validator:

```bash
"$UPSIDE_BIN" tx distribution withdraw-rewards "$VALIDATOR" \
  --from "$KEY_NAME" --chain-id "$CHAIN_ID" --node "$RPC_URL" \
  --home "$WALLET_HOME" --keyring-backend "$KEYRING_BACKEND" \
  --gas auto --gas-adjustment 1.5 --gas-prices "$GAS_PRICES"
```

To claim from every validator you delegate to, use `tx distribution withdraw-all-rewards` with the same flags.

## Redelegate to another validator

This moves an existing staking position without first making it liquid:

```bash
DEST_VALIDATOR="REPLACE_WITH_DESTINATION_OPERATOR_ADDRESS"
"$UPSIDE_BIN" tx staking redelegate "$VALIDATOR" "$DEST_VALIDATOR" "100000000000000000000${DENOM}" \
  --from "$KEY_NAME" --chain-id "$CHAIN_ID" --node "$RPC_URL" \
  --home "$WALLET_HOME" --keyring-backend "$KEYRING_BACKEND" \
  --gas auto --gas-adjustment 1.5 --gas-prices "$GAS_PRICES"
```

The chain limits the number and timing of redelegations. You might have to wait before moving the same stake again. Penalties for earlier offenses by the source validator can still apply during the relevant period.

## Unbond and return tokens to the wallet

```bash
"$UPSIDE_BIN" tx staking unbond "$VALIDATOR" "100000000000000000000${DENOM}" \
  --from "$KEY_NAME" --chain-id "$CHAIN_ID" --node "$RPC_URL" \
  --home "$WALLET_HOME" --keyring-backend "$KEYRING_BACKEND" \
  --gas auto --gas-adjustment 1.5 --gas-prices "$GAS_PRICES"
```

To inspect the unbonding entries:

```bash
"$UPSIDE_BIN" query staking unbonding-delegations "$ACCOUNT_ADDR" \
  --node "$RPC_URL" --home "$WALLET_HOME" --output json
```

The tokens stop earning delegation rewards and remain unavailable to spend until unbonding completes. The standard staking process returns them to your wallet automatically, without a separate withdrawal transaction. Penalties for certain past validator offenses can still apply during unbonding.

## Managing existing positions

The CLI also includes `cancel-unbond`, which identifies an unbonding entry by its creation height, and `distribution set-withdraw-addr`, which changes the reward destination.
