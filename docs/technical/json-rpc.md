---
title: "JSON-RPC reference"
description: The Ethereum JSON-RPC methods that the public Upside mainnet endpoint serves, with examples.
slug: /technical/json-rpc
---

# JSON-RPC reference

Upside serves the Ethereum JSON-RPC interface, so Ethereum wallets, libraries, and tools can connect to it. This page lists the `eth`, `net`, and `web3` methods. The public mainnet and testnet endpoints serve only these three namespaces. The `txpool` methods are available only on your own node. Most methods have an example.

| Transport | Mainnet endpoint |
| --- | --- |
| HTTP | `https://rpc.mainnet.upside.net` |
| WebSocket | `wss://rpc.mainnet.upside.net/ws` |

For chain settings, fees, and the differences from Ethereum, see [Build on the EVM](build-on-the-evm.md).

Limits and enabled methods can differ between RPC providers.

## Send a request

Send a JSON-RPC 2.0 request with the HTTP `POST` method and the `Content-Type: application/json` header:

```bash
curl -s -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","id":1,"method":"eth_chainId","params":[]}' \
  https://rpc.mainnet.upside.net
```

```json
{"jsonrpc":"2.0","id":1,"result":"0xf1b1"}
```

## Conventions

- **Quantities.** Numbers are hexadecimal strings with a `0x` prefix, such as `0x3f0687`. Amounts of the native token are in `aup`.
- **Data.** Byte values are hexadecimal strings with a `0x` prefix and two digits per byte.
- **Block parameter.** Methods that read state take a block number in hexadecimal or one of the tags `latest`, `earliest`, `pending`, and `finalized`. The node treats `finalized` the same as `latest`, and `earliest` is block 1. The parameter is required unless a method lists it as optional.
- **Examples.** Where a method takes a block, the examples use block 4,130,439, which is `0x3f0687`. The account in the examples, `0x111f75f912cbad1ac95fda5f9d0fd9457560a978`, belongs to `upside-1`, a validator that the Upside Foundation runs. Values in angle brackets, such as `<transaction-hash>`, are placeholders. Long values and long lists are shortened with `...`.

## Behavior to know

- **Log range.** In `eth_getLogs`, the difference between `toBlock` and `fromBlock` can be at most 10,000.
- **Filters.** A filter lives in the memory of the node that created it. The node can remove it once five minutes pass without an `eth_getFilterChanges` call. A later request routed to another backend cannot find that filter. Use `eth_getLogs` or a WebSocket subscription when you need reliable results.
- **Receipts.** The `effectiveGasPrice` field is present only in receipts for type 2 transactions in this release.
- **Gas estimates.** Use `eth_estimateGas`. Don't rely on the gas value from `eth_fillTransaction`. The node computes that value against the state at block 1, so it's too low for a call to a contract that was deployed later.
- **Node keys.** The public endpoint returns an empty `eth_accounts` list. Sign transactions locally and submit signed bytes with `eth_sendRawTransaction`. `eth_sendTransaction`, `eth_sign`, `eth_signTypedData`, and `eth_resend` sign with a key stored on the node. The public endpoint holds no keys, so these methods fail there.
- **Coinbase.** `eth_coinbase` returns the address of the validator that the node runs. The public RPC node isn't a validator, so the method returns an error there.
- **Mining.** Upside is a proof-of-stake chain. The mining methods of the Ethereum interface aren't served.

## Methods

| Group | Methods |
| --- | --- |
| Node and network | `web3_clientVersion`, `web3_sha3`, `net_version`, `net_listening`, `net_peerCount`, `eth_chainId`, `eth_protocolVersion`, `eth_syncing`, `eth_blockNumber` |
| Fees | `eth_gasPrice`, `eth_maxPriorityFeePerGas`, `eth_feeHistory` |
| Accounts and state | `eth_getBalance`, `eth_getTransactionCount`, `eth_getCode`, `eth_getStorageAt`, `eth_getProof`, `eth_accounts` |
| Calls and gas estimates | `eth_call`, `eth_estimateGas` |
| Blocks | `eth_getBlockByNumber`, `eth_getBlockByHash`, `eth_getBlockTransactionCountByNumber`, `eth_getBlockTransactionCountByHash`, `eth_getBlockReceipts`, `eth_getUncleCountByBlockNumber`, `eth_getUncleCountByBlockHash`, `eth_getUncleByBlockNumberAndIndex`, `eth_getUncleByBlockHashAndIndex` |
| Transactions | `eth_getTransactionByHash`, `eth_getTransactionByBlockHashAndIndex`, `eth_getTransactionByBlockNumberAndIndex`, `eth_getTransactionReceipt`, `eth_getTransactionLogs`, `eth_sendRawTransaction`, `eth_getPendingTransactions`, `eth_fillTransaction` |
| Node keys | `eth_coinbase`, `eth_sendTransaction`, `eth_sign`, `eth_signTypedData`, `eth_resend` (see [Behavior to know](#behavior-to-know)) |
| Logs and filters | `eth_getLogs`, `eth_newFilter`, `eth_newBlockFilter`, `eth_newPendingTransactionFilter`, `eth_getFilterChanges`, `eth_getFilterLogs`, `eth_uninstallFilter` |
| Subscriptions | `eth_subscribe`, `eth_unsubscribe` |

`eth_getUncleByBlockNumberAndIndex` and `eth_getUncleByBlockHashAndIndex` always return `null`, because Upside blocks have no uncles. `eth_getTransactionLogs` takes a transaction hash and returns the logs of that transaction.

The public endpoints serve only the `eth`, `net`, and `web3` namespaces. Methods in other namespaces, including `txpool`, `personal`, and `debug`, and the `rpc_modules` method return `Method not found` (code `-32601`). A batch request that contains one of those methods is rejected. The `txpool` methods are described in [Methods available only on your own node](#methods-available-only-on-your-own-node).

## Node and network

### web3_clientVersion

Returns the version string of the node software.

**Parameters** None.

**Returns** A string.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"web3_clientVersion","params":[]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":"Version dev ()\nCompiled at  using Go go1.25.7 (amd64)"}
```

### web3_sha3

Returns the Keccak-256 hash of the parameter.

**Parameters**

1. The text to hash.

**Returns** The 32-byte hash.

The node hashes the characters of the parameter as text. It doesn't decode a hexadecimal string first, which differs from Ethereum clients. The example hashes the 12 characters `0x68656c6c6f`.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"web3_sha3","params":["0x68656c6c6f"]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":"0x43e4886c04c4563d8da337f177f98d3b7b8c9f78e70ece7566829146fd4fb54d"}
```

### net_version

Returns the network ID.

**Parameters** None.

**Returns** The network ID as a decimal string. On mainnet it equals the chain ID.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"net_version","params":[]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":"61873"}
```

### net_listening

Reports whether the node is listening for peer connections.

**Parameters** None.

**Returns** A Boolean value.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"net_listening","params":[]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":true}
```

### net_peerCount

Returns the number of peers that the node is connected to.

**Parameters** None.

**Returns** A JSON number. It isn't a hexadecimal string. The value depends on which node answers.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"net_peerCount","params":[]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":8}
```

### eth_chainId

Returns the chain ID.

**Parameters** None.

**Returns** The chain ID. `0xf1b1` is 61873.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_chainId","params":[]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":"0xf1b1"}
```

### eth_protocolVersion

Returns the Ethereum protocol version that the node reports.

**Parameters** None.

**Returns** The version number.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_protocolVersion","params":[]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":"0x41"}
```

### eth_syncing

Reports whether the node is still catching up with the chain.

**Parameters** None.

**Returns** `false` when the node is in sync. Otherwise, an object that describes the sync progress.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_syncing","params":[]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":false}
```

### eth_blockNumber

Returns the number of the latest block.

**Parameters** None.

**Returns** The block number.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_blockNumber","params":[]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":"0x3f07f9"}
```


## Fees

### eth_gasPrice

Returns a gas price for a new transaction.

**Parameters** None.

**Returns** The base fee plus the suggested priority fee, in `aup`.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_gasPrice","params":[]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":"0x7"}
```

### eth_maxPriorityFeePerGas

Returns a suggested priority fee for a dynamic fee transaction.

**Parameters** None.

**Returns** The suggestion in `aup`. It's the base fee divided by eight and rounded down.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_maxPriorityFeePerGas","params":[]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":"0x0"}
```

### eth_feeHistory

Returns the base fee and gas usage of recent blocks.

**Parameters**

1. The number of blocks to return.
2. The newest block of the range. A block number in hexadecimal, or a block tag.
3. A list of percentiles. For each block, the response gives the priority fee at each percentile.

**Returns** An object that describes the range. It holds the oldest block, the base fee of each block and of the block after the range, the share of the gas limit that each block used, and the priority fees.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_feeHistory","params":["0x2","0x3f0687",[50]]}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "oldestBlock": "0x3f0686",
    "reward": [
      [
        "0x0"
      ],
      [
        "0x0"
      ]
    ],
    "baseFeePerGas": [
      "0x7",
      "0x7",
      "0x7"
    ],
    "gasUsedRatio": [
      0,
      0
    ]
  }
}
```


## Accounts and state

### eth_getBalance

Returns the balance of an account.

**Parameters**

1. The account address.
2. A block number in hexadecimal, or a block tag.

**Returns** The balance in `aup`.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_getBalance","params":["0x111f75f912cbad1ac95fda5f9d0fd9457560a978","0x3f0687"]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":"0x14cd477bcb0b40bac"}
```

### eth_getTransactionCount

Returns the nonce of an account, which is the number of transactions that the account has sent.

**Parameters**

1. The account address.
2. A block number in hexadecimal, or a block tag.

**Returns** The nonce.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_getTransactionCount","params":["0x111f75f912cbad1ac95fda5f9d0fd9457560a978","0x3f0687"]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":"0x9"}
```

### eth_getCode

Returns the code of a contract.

**Parameters**

1. The contract address.
2. A block number in hexadecimal, or a block tag.

**Returns** The code. An account without code returns `0x`.

The example reads an account that has no code.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_getCode","params":["0x111f75f912cbad1ac95fda5f9d0fd9457560a978","0x3f0687"]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":"0x"}
```

### eth_getStorageAt

Returns the value of one storage slot of a contract.

**Parameters**

1. The contract address.
2. The slot position.
3. A block number in hexadecimal, or a block tag.

**Returns** The 32-byte value. An unused slot returns zero.

The example reads an account that has no contract storage.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_getStorageAt","params":["0x111f75f912cbad1ac95fda5f9d0fd9457560a978","0x0","0x3f0687"]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":"0x0000000000000000000000000000000000000000000000000000000000000000"}
```

### eth_getProof

Returns an account with proofs of its state.

**Parameters**

1. The account address.
2. A list of storage keys to prove.
3. A block number in hexadecimal, or a block tag.

**Returns** An object with the balance, nonce, code hash, storage hash, and proofs.

The proofs come from the chain's own state tree. They aren't Ethereum Merkle-Patricia proofs, so an Ethereum proof verifier can't check them.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_getProof","params":["0x111f75f912cbad1ac95fda5f9d0fd9457560a978",[],"0x3f0687"]}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "address": "0x111f75f912cbad1ac95fda5f9d0fd9457560a978",
    "accountProof": [
      "0x0ae1040a1501111f75f912cbad1ac95fda5f9d0fd9457560a97812ea010a1e2f...",
      "0x0aff010a0361636312208e392b1e45ebdc19240f357c0422d4a6d2d82d61d84d..."
    ],
    "balance": "0x14cd477bcb0b40bac",
    "codeHash": "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
    "nonce": "0x9",
    "storageHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "storageProof": []
  }
}
```

### eth_accounts

Returns the account addresses exposed by the node's account backend.

**Parameters** None.

**Returns** A list of account addresses exposed by the node's account backend.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_accounts","params":[]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":[]}
```


## Calls and gas estimates

### eth_call

Runs a call against the state of a block without creating a transaction.

**Parameters**

1. A call object. It can contain `from`, `to`, `gas`, `gasPrice`, `value`, and `data`.
2. A block number in hexadecimal, or a block tag.

**Returns** The data that the call returns.

The example calls the SHA-256 precompiled contract at address `0x02` with the five bytes of the word `hello`.

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "eth_call",
  "params": [
    {
      "to": "0x0000000000000000000000000000000000000002",
      "data": "0x68656c6c6f"
    },
    "0x3f0687"
  ]
}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":"0x2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"}
```

### eth_estimateGas

Estimates the gas that a transaction needs.

**Parameters**

1. A call object, as for `eth_call`.
2. Optional. A block number in hexadecimal, or a block tag.

**Returns** The estimate.

The estimate doesn't include the minimum gas limit for call data. See [Build on the EVM](build-on-the-evm.md#minimum-gas-limit-for-call-data). The example estimates a transfer of 1 `aup`. The result, `0x5208`, is 21,000.

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "eth_estimateGas",
  "params": [
    {
      "from": "0x111f75f912cbad1ac95fda5f9d0fd9457560a978",
      "to": "0xa54670dfcdf22a68a59dc6cf52f71371c995711e",
      "value": "0x1"
    }
  ]
}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":"0x5208"}
```


## Blocks

### eth_getBlockByNumber

Returns a block by its number.

**Parameters**

1. A block number in hexadecimal, or a block tag.
2. `true` to return full transaction objects, or `false` to return only transaction hashes.

**Returns** A block object, or `null` if the block doesn't exist.

In the block object, `hash` is the CometBFT block hash, and `stateRoot` is the application hash. `gasLimit` is `0x8f0d180`, which is 150,000,000. The `miner` field is the operator address of the validator that proposed the block. The `difficulty`, `nonce`, and `mixHash` fields are always zero, and `uncles` is always empty.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_getBlockByNumber","params":["0x3f0687",false]}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "baseFeePerGas": "0x7",
    "difficulty": "0x0",
    "extraData": "0x",
    "gasLimit": "0x8f0d180",
    "gasUsed": "0x0",
    "hash": "0x04a77ab49b22ca94e8e0e4e020996982cedbd2e6843e99bb2a7294ec9b54f7a2",
    "logsBloom": "0x0000000000000000000000000000000000000000000000000000000000000000...",
    "miner": "0xa54670dfcdf22a68a59dc6cf52f71371c995711e",
    "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "nonce": "0x0000000000000000",
    "number": "0x3f0687",
    "parentHash": "0xfbf69eebc7c2e277f1801eb5fe4a204300a97c19ed21dc230112849e323c163e",
    "receiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "size": "0x46c",
    "stateRoot": "0xb1b252bb6668e8ced24c94b1c470cf0f2613b9c01b90bbc8c7961a2dc6a48741",
    "timestamp": "0x6aecfa99",
    "transactions": [],
    "transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "uncles": []
  }
}
```

### eth_getBlockByHash

Returns a block by its hash.

**Parameters**

1. The block hash.
2. `true` to return full transaction objects, or `false` to return only transaction hashes.

**Returns** A block object, or `null` if the block doesn't exist.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_getBlockByHash","params":["0x04a77ab49b22ca94e8e0e4e020996982cedbd2e6843e99bb2a7294ec9b54f7a2",false]}
```

**Response**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "baseFeePerGas": "0x7",
    "difficulty": "0x0",
    "extraData": "0x",
    "gasLimit": "0x8f0d180",
    "gasUsed": "0x0",
    "hash": "0x04a77ab49b22ca94e8e0e4e020996982cedbd2e6843e99bb2a7294ec9b54f7a2",
    "logsBloom": "0x0000000000000000000000000000000000000000000000000000000000000000...",
    "miner": "0xa54670dfcdf22a68a59dc6cf52f71371c995711e",
    "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "nonce": "0x0000000000000000",
    "number": "0x3f0687",
    "parentHash": "0xfbf69eebc7c2e277f1801eb5fe4a204300a97c19ed21dc230112849e323c163e",
    "receiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "size": "0x46c",
    "stateRoot": "0xb1b252bb6668e8ced24c94b1c470cf0f2613b9c01b90bbc8c7961a2dc6a48741",
    "timestamp": "0x6aecfa99",
    "transactions": [],
    "transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "uncles": []
  }
}
```

### eth_getBlockTransactionCountByNumber

Returns the number of EVM transactions in a block.

**Parameters**

1. A block number in hexadecimal, or a block tag.

**Returns** The count.

The example block holds no EVM transactions.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_getBlockTransactionCountByNumber","params":["0x3f0687"]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":"0x0"}
```

### eth_getBlockTransactionCountByHash

Returns the number of EVM transactions in a block.

**Parameters**

1. The block hash.

**Returns** The count.

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "eth_getBlockTransactionCountByHash",
  "params": [
    "0x04a77ab49b22ca94e8e0e4e020996982cedbd2e6843e99bb2a7294ec9b54f7a2"
  ]
}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":"0x0"}
```

### eth_getBlockReceipts

Returns the receipts of every EVM transaction in a block.

**Parameters**

1. A block number in hexadecimal, or a block tag.

**Returns** A list of [receipt objects](#receipt-object). The list is empty for a block without EVM transactions.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_getBlockReceipts","params":["0x3f0687"]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":[]}
```

### eth_getUncleCountByBlockNumber

Returns the number of uncle blocks. Upside has no uncle blocks.

**Parameters**

1. A block number in hexadecimal, or a block tag.

**Returns** Always `0x0`.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_getUncleCountByBlockNumber","params":["0x3f0687"]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":"0x0"}
```

### eth_getUncleCountByBlockHash

Returns the number of uncle blocks. Upside has no uncle blocks.

**Parameters**

1. The block hash.

**Returns** Always `0x0`.

**Request**

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "eth_getUncleCountByBlockHash",
  "params": [
    "0x04a77ab49b22ca94e8e0e4e020996982cedbd2e6843e99bb2a7294ec9b54f7a2"
  ]
}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":"0x0"}
```


## Transactions

The methods in this group take or return the following objects.

### Transaction object

| Field | Description |
| --- | --- |
| `hash` | The transaction hash. |
| `blockHash`, `blockNumber`, `transactionIndex` | The block that holds the transaction, and the position of the transaction in that block. |
| `from`, `to` | The sender and the recipient. |
| `nonce` | The nonce of the sender. |
| `value` | The amount transferred, in `aup`. |
| `input` | The call data. |
| `gas` | The gas limit. |
| `gasPrice` | The gas price in `aup`. |
| `maxFeePerGas`, `maxPriorityFeePerGas` | The fee limits of a dynamic fee transaction. |
| `type` | The transaction type. `0x0` is legacy, `0x1` is access list, `0x2` is dynamic fee, and `0x4` is EIP-7702 set-code. |
| `accessList` | The access list of the transaction. |
| `chainId` | The chain ID. |
| `v`, `r`, `s`, `yParity` | The signature. |

### Receipt object

| Field | Description |
| --- | --- |
| `transactionHash`, `transactionIndex`, `blockHash`, `blockNumber` | The transaction, and where it was included. |
| `from`, `to` | The sender and the recipient. |
| `status` | `0x1` for success, or `0x0` for failure. |
| `gasUsed` | The gas that the transaction was charged for. |
| `cumulativeGasUsed` | The gas used in the block up to and including this transaction. |
| `effectiveGasPrice` | The price paid per unit of gas. The field is present only for a type 2 transaction in this release. |
| `contractAddress` | The address of the contract that the transaction created, or `null`. |
| `logs` | A list of [log objects](#log-object). |
| `logsBloom` | The bloom filter of the logs. |
| `type` | The transaction type. |

### eth_getTransactionByHash

Returns a transaction by its Ethereum hash.

**Parameters**

1. The transaction hash.

**Returns** A [transaction object](#transaction-object), or `null` if the node doesn't know the transaction.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_getTransactionByHash","params":["<transaction-hash>"]}
```

### eth_getTransactionByBlockHashAndIndex

Returns a transaction by its position in a block.

**Parameters**

1. The block hash.
2. The index of the transaction in the block.

**Returns** A [transaction object](#transaction-object), or `null`.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_getTransactionByBlockHashAndIndex","params":["<block-hash>","0x0"]}
```

### eth_getTransactionByBlockNumberAndIndex

Returns a transaction by its position in a block.

**Parameters**

1. A block number in hexadecimal, or a block tag.
2. The index of the transaction in the block.

**Returns** A [transaction object](#transaction-object), or `null`.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_getTransactionByBlockNumberAndIndex","params":["<block-number>","0x0"]}
```

### eth_getTransactionReceipt

Returns the receipt of a transaction.

**Parameters**

1. The transaction hash.

**Returns** A [receipt object](#receipt-object), or `null` when this endpoint cannot return a receipt.

A missing receipt does not prove that the transaction was never included. Index availability and some Cosmos-level execution failures can prevent a receipt from being returned. Check the original Cosmos transaction result and account state before submitting a replacement. See [transaction lookup](query-the-chain.md#blocks-and-transactions).

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_getTransactionReceipt","params":["<transaction-hash>"]}
```

### eth_sendRawTransaction

Submits a signed transaction.

**Parameters**

1. The signed transaction, encoded as hexadecimal data.

**Returns** The transaction hash.

The transaction must be signed with the chain ID and must use the next nonce of the account. To learn whether it succeeded, read its receipt.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_sendRawTransaction","params":["<signed-transaction>"]}
```

### eth_getPendingTransactions

Returns pending transactions that the node knows about.

**Parameters** None.

**Returns** A list of up to 30 [transaction objects](#transaction-object).

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_getPendingTransactions","params":[]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":[]}
```


## Logs and filters

### Log object

| Field | Description |
| --- | --- |
| `address` | The contract that emitted the log. |
| `topics` | The indexed values of the event. The first topic is usually the hash of the event signature. |
| `data` | The values of the event that aren't indexed. |
| `blockNumber`, `blockHash`, `transactionHash`, `transactionIndex`, `logIndex` | Where the log was emitted. |
| `removed` | `false`. Upside blocks are final, so a log isn't removed later. |

### eth_getLogs

Returns the logs that match a filter.

**Parameters**

1. A filter object. It can contain `fromBlock`, `toBlock`, `address`, `topics`, and `blockHash`.

**Returns** A list of [log objects](#log-object).

The difference between `toBlock` and `fromBlock` can be at most 10,000. The example block has no logs.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_getLogs","params":[{"fromBlock":"0x3f0687","toBlock":"0x3f0687"}]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":[]}
```

### eth_newFilter

Creates a log filter on the node.

**Parameters**

1. A filter object, as for `eth_getLogs`.

**Returns** The filter ID.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_newFilter","params":[{"fromBlock":"<block-number>","toBlock":"latest","address":"<contract-address>"}]}
```

### eth_newBlockFilter

Creates a filter that reports new blocks.

**Parameters** None.

**Returns** The filter ID.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_newBlockFilter","params":[]}
```

### eth_newPendingTransactionFilter

Creates a filter that reports new pending transactions.

**Parameters** None.

**Returns** The filter ID.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_newPendingTransactionFilter","params":[]}
```

### eth_getFilterChanges

Returns what a filter has collected since the last poll.

**Parameters**

1. The filter ID.

**Returns** A list of block hashes, transaction hashes, or [log objects](#log-object), depending on the kind of filter.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_getFilterChanges","params":["<filter-id>"]}
```

### eth_getFilterLogs

Returns the logs that match a log filter. The block range limit of `eth_getLogs` applies.

**Parameters**

1. The ID of a log filter.

**Returns** A list of [log objects](#log-object).

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_getFilterLogs","params":["<filter-id>"]}
```

### eth_uninstallFilter

Removes a filter from the node.

**Parameters**

1. The filter ID.

**Returns** `true` if the filter was removed.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_uninstallFilter","params":["<filter-id>"]}
```


## Subscriptions

### eth_subscribe and eth_unsubscribe

These methods work only over the WebSocket endpoint. `eth_subscribe` opens a subscription and returns its ID. The node then sends an `eth_subscription` message for each event. `eth_unsubscribe` closes the subscription.

Over WebSocket, the request `id` must be a number. The node rejects a text ID, such as a UUID, with the error `invalid type for connection ID`.

**Parameters of `eth_subscribe`**

1. The kind of subscription, either `newHeads` or `logs`.
2. For `logs`, a filter object with `address` and `topics`.

The following exchange subscribes to new block headers, receives one header, and unsubscribes. In the header message, `gasLimit` and `gasUsed` are `0x0`, and `logsBloom` is zero. To get those values, read the block with `eth_getBlockByNumber`.

**Sent**

```json
{"jsonrpc":"2.0","id":1,"method":"eth_subscribe","params":["newHeads"]}
```

**Received**

```json
{"jsonrpc":"2.0","result":"0x38699130f6005ea5e7691f737c4cdfa","id":1}
```

**Received**

```json
{
  "jsonrpc": "2.0",
  "method": "eth_subscription",
  "params": {
    "subscription": "0x38699130f6005ea5e7691f737c4cdfa",
    "result": {
      "parentHash": "0xd36fb87cf7638f04dd048502f890cf58cca3e39266e2b28e8f7237e966806834",
      "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
      "miner": "0xf5bf5DcADF41445d8365D08Bb61603C954843633",
      "stateRoot": "0xa228920a933ca4e4eb03c1d9df0d25686772c0ac3b0938da9af2c75806aaaae2",
      "transactionsRoot": "0xe3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
      "receiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
      "logsBloom": "0x0000000000000000000000000000000000000000000000000000000000000000...",
      "difficulty": "0x0",
      "number": "0x3f0928",
      "gasLimit": "0x0",
      "gasUsed": "0x0",
      "timestamp": "0x6aed09c9",
      "extraData": "0x",
      "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
      "nonce": "0x0000000000000000",
      "baseFeePerGas": "0x7",
      "withdrawalsRoot": null,
      "hash": "0x6f103df021ecf57ea1bb080d94974991749a8393ab033c4fe702e8bb488869d3"
    }
  }
}
```

**Sent**

```json
{"jsonrpc":"2.0","id":2,"method":"eth_unsubscribe","params":["0x38699130f6005ea5e7691f737c4cdfa"]}
```

**Received**

```json
{"jsonrpc":"2.0","result":true,"id":2}
```


## Methods available only on your own node

Use these methods on a node that you run. The public endpoints don't serve them. A node serves them only if you add `txpool` to `json-rpc.api` in its `app.toml`; the default is `eth,net,web3`.

### txpool_status

Returns fixed zero counts for pending and queued transactions.

**Parameters** None.

**Returns** An object in which both counts are always zero.

To see pending EVM transactions, use `eth_getPendingTransactions`.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"txpool_status","params":[]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":{"pending":"0x0","queued":"0x0"}}
```

### txpool_content

Returns empty `pending` and `queued` objects regardless of the transaction pool's contents.

**Parameters** None.

**Returns** An object in which `pending` and `queued` are always empty.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"txpool_content","params":[]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":{"pending":{},"queued":{}}}
```

### txpool_inspect

Returns empty `pending` and `queued` objects regardless of the transaction pool's contents.

**Parameters** None.

**Returns** An object in which `pending` and `queued` are always empty.

**Request**

```json
{"jsonrpc":"2.0","id":1,"method":"txpool_inspect","params":[]}
```

**Response**

```json
{"jsonrpc":"2.0","id":1,"result":{"pending":{},"queued":{}}}
```
