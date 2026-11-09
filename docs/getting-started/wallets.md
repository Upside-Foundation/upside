---
title: Wallets and accounts
description: Create, recover, and use wallets for Upside L1 accounts.
slug: /getting-started/wallets
---

# Wallets and accounts

A wallet manages the keys that authorize transactions. You can use it through a public RPC endpoint without running a full node.

## Before you start

Install the correct [binary](install.md) and set the [network profile](networks.md). A profile is the set of shell variables that selects your binary, network, RPC endpoint, fees, and wallet storage.

Use separate wallets for mainnet and testnet. The `file` keyring encrypts your keys with a password. The `test` keyring stores unencrypted keys and is suitable only for disposable local accounts.

## Create a wallet

```bash
export KEY_NAME="my-wallet"
"$UPSIDE_BIN" keys add "$KEY_NAME" \
  --key-type eth_secp256k1 --coin-type 60 \
  --keyring-backend "$KEYRING_BACKEND" --home "$WALLET_HOME"
```

The command creates a local key and displays a recovery phrase. Record it privately before leaving the terminal. Never paste it into an explorer, support message, chat, source repository, or network configuration. Anyone with the phrase and any associated passphrase can control the account.

After creating a key, fund its address on the selected network. A newly generated address may appear in chain queries only after it receives funds.

## Recover a wallet

On a trusted machine, select the wallet home where you want to restore the account. Use a new local name and enter the phrase at the prompt:

```bash
"$UPSIDE_BIN" keys add "recovered-wallet" --recover \
  --key-type eth_secp256k1 --coin-type 60 \
  --keyring-backend "$KEYRING_BACKEND" --home "$WALLET_HOME"
```

Recovery needs the key algorithm, derivation path, account/index, and any recovery passphrase you used. The derivation path selects an account from the recovery phrase. These examples use Ethereum coin type `60` and the first account at `m/44'/60'/0'/0/0`.

The example assumes no BIP39 passphrase. If the original wallet used one, add `--interactive` to the recovery command and enter that passphrase when prompted. This is separate from the password that encrypts the local keyring.

To test recovery while keeping the original keyring, use a separate, newly created `WALLET_HOME`. This binary rejects adding the same account twice to one keyring, even under different names. Do not delete the original key to work around that error. Use the recovery home for the address check below, then restore your original wallet profile when finished.

Check the restored address before sending funds:

```bash
"$UPSIDE_BIN" keys show recovered-wallet -a \
  --keyring-backend "$KEYRING_BACKEND" --home "$WALLET_HOME"
```

Compare it with your recorded public address. If it differs, check the network build and original recovery settings before using the restored wallet.

## List keys and show an address

```bash
"$UPSIDE_BIN" keys list \
  --keyring-backend "$KEYRING_BACKEND" --home "$WALLET_HOME"

export ACCOUNT_ADDR="$("$UPSIDE_BIN" keys show "$KEY_NAME" -a \
  --keyring-backend "$KEYRING_BACKEND" --home "$WALLET_HOME")"
printf '%s\n' "$ACCOUNT_ADDR"
```

On both public networks, the account address begins with `up1`. Share this public address to receive UP. Never share a private key or recovery phrase to receive funds.

## Address formats


| Form          | Purpose                                                                   |
| ------------- | ------------------------------------------------------------------------- |
| `up1…`        | Cosmos account address on both public networks, used for native transfers and delegations |
| `0x…`         | Ethereum representation used by EVM wallets and applications              |
| `upvaloper1…` | Validator operator address, used as the target of a delegation            |
| `upvalcons1…` | Consensus address derived from the node's signing key                     |


An Ethereum-compatible account has both a Bech32 address (`up1…`) and a hexadecimal address (`0x…`). These represent the same account on the selected chain. Always check the network as well as the address when sending funds.

The consensus address comes from the validator's consensus signing key. The wallet key controls its operator account.

## Check funds

```bash
"$UPSIDE_BIN" query bank balances "$ACCOUNT_ADDR" \
  --node "$RPC_URL" --home "$WALLET_HOME" --output json

"$UPSIDE_BIN" query bank spendable-balances "$ACCOUNT_ADDR" \
  --node "$RPC_URL" --home "$WALLET_HOME" --output json
```

Spendable funds are the funds available for a new transfer or delegation. Staked and unbonding tokens are queried separately. Keep some spendable UP for fees, including future reward withdrawals and unbonding transactions.

## Confirm a transaction

After a command submits a transaction, save its hash and query it until the node returns a result. A successful result has a `code` of `0`. A nonzero code means the chain rejected the transaction. Check the result before treating a transfer, delegation, validator update, or withdrawal as complete:

```bash
export TX_HASH="PASTE_TRANSACTION_HASH"
"$UPSIDE_BIN" query tx "$TX_HASH" \
  --node "$RPC_URL" --home "$WALLET_HOME" --output json
```

Use the transaction result and the relevant account, staking, or validator query to confirm the state change. Keep the network profile and wallet home unchanged while checking the result.

Next: [send a transfer](transactions.md) or [delegate stake](../staking/delegation.md).
