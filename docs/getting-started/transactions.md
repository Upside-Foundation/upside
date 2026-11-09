---
title: Transfers
description: Send UP and check native transaction fees on Upside L1.
slug: /getting-started/transactions
---

# Transfers

Select a network profile before using these commands. Choose transfer amounts appropriate to your task.

## Receive funds

Give the sender the public account address and the exact network. On mainnet, receive native UP at an `up1…` account address. On the public testnet, obtain test tokens from the [faucet](https://faucet.testnet.upside.net). On localnet, use the genesis-funded practice account.

Check your [wallet balance](wallets.md) after the transfer is confirmed. The [wallet guide](wallets.md#confirm-a-transaction) contains the transaction confirmation procedure.

## Send a native transfer

Start with the profile and `KEY_NAME` from the [wallet guide](wallets.md). Obtain the destination account address from the recipient. The following example sends one whole token on a network using 18 decimals:

```bash
RECIPIENT="REPLACE_WITH_RECIPIENT_ACCOUNT_ADDRESS"
AMOUNT="1000000000000000000${DENOM}"

"$UPSIDE_BIN" tx bank send "$KEY_NAME" "$RECIPIENT" "$AMOUNT" \
  --chain-id "$CHAIN_ID" --node "$RPC_URL" \
  --home "$WALLET_HOME" --keyring-backend "$KEYRING_BACKEND" \
  --gas auto --gas-adjustment 1.5 --gas-prices "$GAS_PRICES"
```

`--gas auto` simulates the transaction and `--gas-adjustment 1.5` adds a margin to its gas limit. Set the gas price separately using the network's current fee guidance.

This build applies the dynamic base fee to native Cosmos transactions, including transfers and staking, when the fee market is active. The chain can also enforce a global minimum gas price. A node's local `minimum-gas-prices` setting also applies where required by its fee-checking path. A simulation can succeed even when the submitted transaction's fee is too low.

## Transaction result

Use the [wallet transaction confirmation procedure](wallets.md#confirm-a-transaction) after submission. Then check the destination balance. An included transaction with a nonzero code failed and can still have consumed a fee.

## Common failures


| Symptom                   | Check                                                                                |
| ------------------------- | ------------------------------------------------------------------------------------ |
| Wrong address prefix      | Network build and recipient address type                                             |
| Insufficient funds        | Spendable balance covers both amount and fee                                         |
| Insufficient fee          | Current fee guidance and EVM base-fee conditions                                     |
| Account sequence mismatch | Query previous transaction hashes to check whether the account sequence has advanced |
| Account not found         | Correct network and whether the address has received funds                           |
| Out of gas                | Execution error and whether a new estimate/margin is needed                          |
| Transaction not found     | Inclusion delay, endpoint synchronization, and indexing                              |


If the terminal request times out, check the original transaction hash before sending another transfer.
