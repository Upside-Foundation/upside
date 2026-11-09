---
title: "Query the chain"
description: Read balances, staking, governance, supply, blocks, and transactions from Upside mainnet.
slug: /technical/query-the-chain
---

# Query the chain

This page shows how to read data from Upside mainnet. You don't need a key or an account to run the commands.

The examples use `upside-1`, a validator that the Upside Foundation runs. Its account address is `up1zy0ht7gjewk34j2lmf0e6r7eg46kp2tcelxegc`, and its validator address is `upvaloper1zy0ht7gjewk34j2lmf0e6r7eg46kp2tcvrkyvv`. The example output is from block 4,130,439. Amounts change with every block.

## Endpoints

| Interface | URL | Use it for |
| --- | --- | --- |
| CometBFT RPC | `https://rpc.mainnet.upside.net/tendermint/` | Module queries with `upsided`, blocks, and transactions |
| EVM JSON-RPC | `https://rpc.mainnet.upside.net` | Ethereum-style calls. See the [JSON-RPC reference](json-rpc.md). |

> **Important** Keep the trailing slash on the CometBFT RPC URL. Without it, the server answers with a redirect, and `upsided` stops with the error `invalid character '<' looking for beginning of value`.

The Cosmos SDK gRPC and REST interfaces aren't available on these endpoints. To use those interfaces, run your own node.

## Before you begin

1. Install `upsided`. See [Install the node and wallet CLI](../getting-started/install.md).
2. Select the installed binary and the CometBFT RPC URL in your shell:

   ```bash
   export UPSIDE_BIN="$HOME/.local/bin/upsided-mainnet"
   export RPC_URL="https://rpc.mainnet.upside.net/tendermint/"
   ```

## Check the connection

Run the following command:

```bash
"$UPSIDE_BIN" status --node "$RPC_URL"
```

The command returns JSON. Check that `node_info.network` is `upside_61873-1` and that `sync_info.catching_up` is `false`. The `sync_info` part of the response looks like this:

```json
{
  "latest_block_hash": "5967E4791E6580A2C2F37C9B3A0257400153F08084935E8611A2BA65550B1DEE",
  "latest_app_hash": "D221AB0D714A4BD9118EF3FDB12BD6EB6A5B9DB129C87EA68D460D7311C4259D",
  "latest_block_height": "4130862",
  "latest_block_time": "2026-11-05T22:33:07.561971272Z",
  "earliest_block_hash": "829AB86B94B6B627FA8F32F40D5A151A197B4AB4DD8F9166F510D9A513095D10",
  "earliest_app_hash": "E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855",
  "earliest_block_height": "1",
  "earliest_block_time": "2026-02-18T07:24:13.801872946Z",
  "catching_up": false
}
```

## Read the output

- **Format.** Query commands print YAML by default. Add `--output json` for JSON.
- **Units.** Amounts are in `aup`. One UP is 10^18 `aup`, so `1000000000000000000` `aup` is 1 UP.
- **Scaled decimals.** Some commands print decimal values as integers multiplied by 10^18. For example, a validator commission rate of `50000000000000000` is 0.05, and a `community_tax` of `930000000000000000` is 0.93.
- **Past blocks.** Add `--height` to read the state at an earlier block. For example, `--height 4130439` selects the sample height. The endpoint must retain state at that height. Pruned nodes cannot answer every historical query.
- **Paging.** List commands accept `--page-limit`, `--page-key`, `--page-offset`, `--page-count-total`, and `--page-reverse`.

## Balances

To list every balance of an account, run the following command:

```bash
"$UPSIDE_BIN" query bank balances up1zy0ht7gjewk34j2lmf0e6r7eg46kp2tcelxegc --node "$RPC_URL"
```

```yaml
balances:
- amount: "23982925567956880300"
  denom: aup
pagination:
  total: "1"
```

To read one denomination, name it after the address:

```bash
"$UPSIDE_BIN" query bank balance up1zy0ht7gjewk34j2lmf0e6r7eg46kp2tcelxegc aup --node "$RPC_URL"
```

```yaml
balance:
  amount: "23982925567956880300"
  denom: aup
```

The same account has an Ethereum-style address, and `eth_getBalance` returns its native `aup` balance. Other bank denominations and ERC-20 balances require their own queries. The block number `0x3f0687` is 4,130,439 in hexadecimal, and the result `0x14cd477bcb0b40bac` is 23982925567956880300 in decimal.

```bash
curl -s -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","id":1,"method":"eth_getBalance","params":["0x111f75f912cbad1ac95fda5f9d0fd9457560a978","0x3f0687"]}' \
  https://rpc.mainnet.upside.net
```

```json
{"jsonrpc":"2.0","id":1,"result":"0x14cd477bcb0b40bac"}
```

## Convert an address

An account has one set of 20 bytes and two ways to write them. The `up1` form is used by `upsided`, and the `0x` form is used by Ethereum tools. To convert between them, run `upsided debug addr`. The command works offline.

```bash
"$UPSIDE_BIN" debug addr up1zy0ht7gjewk34j2lmf0e6r7eg46kp2tcelxegc
```

```text
Address: [17 31 117 249 18 203 173 26 201 95 218 95 157 15 217 69 117 96 169 120]
Address (hex): 111F75F912CBAD1AC95FDA5F9D0FD9457560A978
Bech32 Acc: up1zy0ht7gjewk34j2lmf0e6r7eg46kp2tcelxegc
Bech32 Val: upvaloper1zy0ht7gjewk34j2lmf0e6r7eg46kp2tcvrkyvv
Bech32 Con: upvalcons1zy0ht7gjewk34j2lmf0e6r7eg46kp2tccs9cqd
```

`Bech32 Con` above is the same account bytes formatted with a consensus prefix. It is not the validator's actual consensus identity. Obtain that identity from the node's signing public key using `"$UPSIDE_BIN" comet show-address --home "$NODE_HOME"`.

To convert in the other direction, pass the hexadecimal form without the `0x` prefix. The command rejects an address that starts with `0x`.

The command prints the hexadecimal form in uppercase. Some Ethereum libraries accept only lowercase or checksummed addresses, so convert the value to lowercase before you use it there.

## Staking

To list the validators, run the following command:

```bash
"$UPSIDE_BIN" query staking validators --node "$RPC_URL"
```

The entry for `upside-1` looks like this:

```yaml
- commission:
    commission_rates:
      max_change_rate: "1000000000000000000"
      max_rate: "1000000000000000000"
      rate: "1000000000000000000"
    update_time: "2026-08-25T04:27:37.93171827Z"
  consensus_pubkey:
    type: /cosmos.crypto.ed25519.PubKey
    value: qujaQZ9BHjBaUhuaWfnbTmkbtqSuo4h5gnQuai5Mypw=
  delegator_shares: "100000000000000000000000000000000000000000000"
  description:
    moniker: upside-1
  min_self_delegation: "1"
  operator_address: upvaloper1zy0ht7gjewk34j2lmf0e6r7eg46kp2tcvrkyvv
  status: BOND_STATUS_BONDED
  tokens: "100000000000000000000000000"
  unbonding_time: "1970-01-01T00:00:00Z"
```

The `tokens` value is the validator's total stake in `aup`. The commission rates and `delegator_shares` are scaled decimals.

To read the delegations of an account, run the following command:

```bash
"$UPSIDE_BIN" query staking delegations up1zy0ht7gjewk34j2lmf0e6r7eg46kp2tcelxegc --node "$RPC_URL"
```

```yaml
delegation_responses:
- balance:
    amount: "100000000000000000000000000"
    denom: aup
  delegation:
    delegator_address: up1zy0ht7gjewk34j2lmf0e6r7eg46kp2tcelxegc
    shares: "100000000000000000000000000000000000000000000"
    validator_address: upvaloper1zy0ht7gjewk34j2lmf0e6r7eg46kp2tcvrkyvv
pagination:
  total: "1"
```

The following commands read other staking data:

| To read | Command |
| --- | --- |
| One validator | `"$UPSIDE_BIN" query staking validator upvaloper1zy0ht7gjewk34j2lmf0e6r7eg46kp2tcvrkyvv --node "$RPC_URL"` |
| The delegations to a validator | `"$UPSIDE_BIN" query staking delegations-to upvaloper1zy0ht7gjewk34j2lmf0e6r7eg46kp2tcvrkyvv --node "$RPC_URL"` |
| The rewards an account can withdraw | `"$UPSIDE_BIN" query distribution rewards up1zy0ht7gjewk34j2lmf0e6r7eg46kp2tcelxegc --node "$RPC_URL"` |
| The commission a validator can withdraw | `"$UPSIDE_BIN" query distribution commission upvaloper1zy0ht7gjewk34j2lmf0e6r7eg46kp2tcvrkyvv --node "$RPC_URL"` |
| The total bonded and unbonded stake | `"$UPSIDE_BIN" query staking pool --node "$RPC_URL"` |
| The staking parameters | `"$UPSIDE_BIN" query staking params --node "$RPC_URL"` |

A validator address starts with `upvaloper1`. An account address starts with `up1`.

The staking pool and the staking parameters look like this:

```yaml
pool:
  bonded_tokens: "385000030000000000000000000"
  not_bonded_tokens: "8290340000000000000000"
```

```yaml
params:
  bond_denom: aup
  historical_entries: 10000
  max_entries: 7
  max_validators: 100
  min_commission_rate: "0"
  unbonding_time: 504h0m0s
```

An `unbonding_time` of `504h0m0s` is 21 days.

## Governance

To list the proposals, run the following command. Add `--proposal-status passed` to list only the proposals that passed.

```bash
"$UPSIDE_BIN" query gov proposals --node "$RPC_URL"
```

To read one proposal, pass its ID:

```bash
"$UPSIDE_BIN" query gov proposal 6 --node "$RPC_URL"
```

The output includes the following fields:

```yaml
  expedited: true
  id: "6"
  status: PROPOSAL_STATUS_PASSED
  submit_time: "2026-10-22T19:08:50.680867419Z"
  title: Disable IBC token transfers on Upside mainnet
  voting_end_time: "2026-10-23T19:08:50.680867419Z"
  voting_start_time: "2026-10-22T19:08:50.680867419Z"
```

To read the vote count of a proposal, run the following command. Vote counts are in `aup` of voting power.

```bash
"$UPSIDE_BIN" query gov tally 6 --node "$RPC_URL"
```

```yaml
tally:
  abstain_count: "0"
  no_count: "0"
  no_with_veto_count: "0"
  yes_count: "200000000000000000000000000"
```

To read the governance parameters, run the following command:

```bash
"$UPSIDE_BIN" query gov params --node "$RPC_URL"
```

```yaml
params:
  burn_vote_veto: true
  expedited_min_deposit:
  - amount: "100000"
    denom: aup
  expedited_threshold: "0.667000000000000000"
  expedited_voting_period: 24h0m0s
  max_deposit_period: 48h0m0s
  min_deposit:
  - amount: "10000"
    denom: aup
  min_deposit_ratio: "0.010000000000000000"
  min_initial_deposit_ratio: "0.000000000000000000"
  proposal_cancel_ratio: "0.500000000000000000"
  quorum: "0.334000000000000000"
  threshold: "0.500000000000000000"
  veto_threshold: "0.334000000000000000"
  voting_period: 48h0m0s
```

## Supply and minting

| To read | Command | Output at block 4,130,439 |
| --- | --- | --- |
| The total supply of `aup` | `"$UPSIDE_BIN" query bank total-supply-of aup --node "$RPC_URL"` | `668360906886397911462034871` |
| The annual inflation rate | `"$UPSIDE_BIN" query mint inflation --node "$RPC_URL"` | `0.742365971441284236` |
| The annual provisions in `aup` | `"$UPSIDE_BIN" query mint annual-provisions --node "$RPC_URL"` | `496168330220775927410258746.969612305291076872` |
| The community pool in `aup` | `"$UPSIDE_BIN" query distribution community-pool --node "$RPC_URL"` | `138216024116968283840106165.516854950593947283` |

For how these values are computed, see the [Supply and minting reference](supply-and-minting.md).

## Blocks and transactions

To read a block by height, run the following command:

```bash
"$UPSIDE_BIN" query block --type=height 4130439 --node "$RPC_URL"
```

To read a transaction by its hash, run the following command. The example is the vote that the `upside-1` account cast on proposal 6.

```bash
"$UPSIDE_BIN" query tx A2F8AF557259B669C634CFB76CA8E536B447785CD5BED1C6E73B9274B7A3B478 --node "$RPC_URL"
```

The output includes the following fields:

```yaml
code: 0
gas_used: "69071"
gas_wanted: "200000"
height: "3908130"
timestamp: "2026-10-22T19:19:38Z"
txhash: A2F8AF557259B669C634CFB76CA8E536B447785CD5BED1C6E73B9274B7A3B478
```

For native Cosmos operations, an included transaction with `code: 0` succeeded. For EVM operations, also inspect the Ethereum receipt's `status`. An EVM revert can be recorded within a successful Cosmos transaction.

To search for transactions by event, pass a query. The following command returns the transactions that the `upside-1` account sent, two per page:

```bash
"$UPSIDE_BIN" query txs --query "message.sender='up1zy0ht7gjewk34j2lmf0e6r7eg46kp2tcelxegc'" --limit 2 --node "$RPC_URL"
```

An EVM transaction has an Ethereum hash and a different hash on the Cosmos side. To find one by its Ethereum hash, use `eth_getTransactionByHash` on the EVM JSON-RPC endpoint, or search the event index:

```bash
"$UPSIDE_BIN" query txs --query "ethereum_tx.ethereumTxHash='<ethereum-transaction-hash>'" --node "$RPC_URL"
```

## Query without upsided

The CometBFT RPC endpoint also answers plain HTTP requests. For example, the following command returns the node status:

```bash
curl -s "https://rpc.mainnet.upside.net/tendermint/status"
```

The endpoint serves the following paths:

| Path | Returns |
| --- | --- |
| `/status` | Node and sync status |
| `/health` | An empty result if the node is running |
| `/abci_info` | Application name, version, and latest height |
| `/block?height=4130439` | A block |
| `/validators?height=4130439` | The validator set at a height |
| `/consensus_params?height=4130439` | The consensus parameters at a height |
| `/tx?hash=0xA2F8AF557259B669C634CFB76CA8E536B447785CD5BED1C6E73B9274B7A3B478` | A transaction |

## Query from JavaScript

CosmJS can read from the CometBFT RPC endpoint. The following script uses `@cosmjs/stargate` 0.39.0.

```javascript
import { StargateClient, QueryClient, setupBankExtension, setupStakingExtension } from "@cosmjs/stargate";
import { Comet38Client } from "@cosmjs/tendermint-rpc";

const rpc = "https://rpc.mainnet.upside.net/tendermint/";
const address = "up1zy0ht7gjewk34j2lmf0e6r7eg46kp2tcelxegc";

const client = await StargateClient.connect(rpc);
console.log(await client.getChainId());
console.log(await client.getHeight());
console.log(await client.getBalance(address, "aup"));

const comet = await Comet38Client.connect(rpc);
const query = QueryClient.withExtensions(comet, setupBankExtension, setupStakingExtension);
console.log(await query.bank.supplyOf("aup"));
console.log((await query.staking.validators("BOND_STATUS_BONDED")).validators.length);

client.disconnect();
comet.disconnect();
```

The script prints the chain ID `upside_61873-1`, the current height, the balance of the account, the total supply, and the number of bonded validators.

`client.getAccount` doesn't work with this version of CosmJS. It stops with an `Unsupported type` error, because CosmJS doesn't recognize the account type that the chain uses.
