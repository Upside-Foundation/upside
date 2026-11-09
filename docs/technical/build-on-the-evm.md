---
title: "Build on the EVM"
description: Network settings, fees, differences from Ethereum, and tools for building on the Upside EVM.
slug: /technical/build-on-the-evm
---

# Build on the EVM

Upside L1 runs the Ethereum Virtual Machine (EVM). Build contracts with Solidity compilers and connect Ethereum libraries through JSON-RPC. Review the compatibility limits below before deployment.

## Network settings

| Setting | Mainnet value |
| --- | --- |
| Chain ID | `61873` (`0xf1b1`) |
| HTTP endpoint | `https://rpc.mainnet.upside.net` |
| WebSocket endpoint | `wss://rpc.mainnet.upside.net/ws` |
| Currency symbol | UP |
| Decimals | 18 |
| Block explorer | `https://explorer.mainnet.upside.net` |

Use these values to add the network to a wallet or a library that supports custom EVM networks. For the methods that the endpoint serves, see the [JSON-RPC reference](json-rpc.md).

## Accounts and units

- **Keys.** By default, accounts use Ethereum secp256k1 keys and the derivation path `m/44'/60'/0'/0/0`. With these defaults, a recovery phrase produces the same account in `upsided` and in an Ethereum wallet.
- **Addresses.** An account has a `0x` address for Ethereum tools and an `up1` address for `upsided`. Both are the same 20 bytes. To convert between them, see [Convert an address](query-the-chain.md#convert-an-address).
- **Balances.** The native `aup` balance is shared between Cosmos bank queries and `eth_getBalance` for the same account. Bank accounts can also hold other denominations, and ERC-20 balances are separate contract state.
- **Changes outside the EVM.** A balance can change without an EVM transaction, for example through a transfer sent with `upsided`. That change has no EVM transaction hash or EVM logs. Cosmos transactions can still emit Cosmos events.
- **Nonce.** The EVM nonce is the account's sequence number. A transaction sent with `upsided` from the same account also advances it.
- **Units.** The smallest unit is `aup`. One UP is 10^18 `aup`, so `aup` takes the place of wei. Gas prices are in `aup` per unit of gas.

## EVM version

The EVM dependency includes Prague opcodes and precompiled contracts, with the compatibility limits below. Compile contracts for the `prague` or `cancun` EVM version. Fee accounting differs from Ethereum, as described in [Fees](#fees).

- `PUSH0`, transient storage (`TSTORE` and `TLOAD`), and `MCOPY` are available.
- The precompiled contracts at addresses `0x01` through `0x11` are available. This range includes the BLS12-381 precompiles that Prague added.
- The app also registers two precompiled contracts that are specific to Upside and inherited from Cronos ([`app/app.go`](https://github.com/Upside-Foundation/upside/blob/main/app/app.go)). The relayer precompile at `0x0000000000000000000000000000000000000065` accepts IBC client, connection, channel, and packet messages from EVM transactions ([`Relayer.sol`](https://github.com/Upside-Foundation/upside/blob/main/x/upside/events/bindings/src/Relayer.sol)). The interchain accounts precompile at `0x0000000000000000000000000000000000000066` lets a contract register an interchain account on another chain and submit transactions through it ([`ICA.sol`](https://github.com/Upside-Foundation/upside/blob/main/x/upside/events/bindings/src/ICA.sol)). Both work only over IBC connections, and mainnet has no IBC channels.
- The `CLZ` opcode from the Osaka upgrade isn't defined, so code that uses it fails.

Set the EVM target explicitly in `foundry.toml` so a toolchain update does not select unsupported opcodes:

```toml
[profile.default]
evm_version = "prague"
```

> **Important** Don't read `block.blobbasefee`. A call that runs the `BLOBBASEFEE` opcode fails with an internal error. This can also cause a submitted transaction to fail without an Ethereum receipt. A missing receipt does not establish that the nonce and fee were untouched. Reconcile the original transaction through the Cosmos transaction result and account state before retrying.

## Fees

A transaction pays its gas price multiplied by the gas it's charged for. Fees are paid in `aup`.

### Base fee

Upside uses the EIP-1559 fee model. The base fee moves with how full the previous block was, measured against a target of half the block gas limit.

| Parameter | Mainnet value |
| --- | --- |
| Block gas limit | 150,000,000 |
| Gas target per block | 75,000,000 |
| Largest change in the base fee per block | One eighth. An increase is at least 1 `aup`. |
| Base fee at block 4,130,439 | 7 `aup` |

With these parameters the base fee stops falling at 7 `aup`, because a smaller decrease rounds down to zero. At that base fee, a transfer that uses 21,000 gas costs 147,000 `aup`, which is 0.000000000000147 UP.

To read the current values, call `eth_gasPrice` or read `baseFeePerGas` from the latest block.

### Priority fee

A dynamic fee transaction can add a priority fee. The priority fee goes to the fee collector with the rest of the fee.

`eth_maxPriorityFeePerGas` returns a suggested value, which is the base fee divided by eight and rounded down. At a base fee of 7 `aup`, the suggestion is zero. `eth_gasPrice` returns the base fee plus that suggestion.

A legacy transaction pays its full `gasPrice` for every unit of gas.

### Minimum gas charge

A transaction is charged for at least half of its gas limit. The gas charged is the larger of the gas that the transaction used and half of the gas limit. The fraction is the governance parameter `min_gas_multiplier`, which is 0.5 on mainnet.

For example, a plain transfer uses 21,000 gas. With a gas limit of 21,000, it's charged for 21,000 gas. With a gas limit of 100,000, it's charged for 50,000 gas.

Set the gas limit close to what the transaction needs. Unused gas above the charged amount is refunded.

Gas refunds for clearing storage aren't applied. A transaction that clears a storage slot is charged for the gas that it used before the refund.

### Minimum gas limit for call data

The node rejects a transaction whose gas limit is below 21,000, plus 10 for each zero byte of call data and 40 for each non-zero byte. The error is `insufficient gas for floor data gas cost`.

`eth_estimateGas` doesn't include this minimum. For a transaction that carries much call data and does little work, the estimate can be below it. For example, a transfer with 100 non-zero bytes of call data has an estimate of 22,600 and a minimum of 25,000. In that case, set the gas limit to the minimum.

### Where fees go

The base fee and the priority fee both go to the fee collector account. The distribution module then splits that balance between the community pool and stakers.

## Transactions

- **Types.** Upside supports legacy (0), access list (1), dynamic fee (2), and EIP-7702 set-code (4) transactions. Type 4 changes account delegation and requires compatible signing tools. The examples below use ordinary transfers and contract calls.
- **Nonces.** A transaction must use the account's next nonce. The node rejects a transaction with a higher nonce and returns the error `invalid nonce`. It doesn't queue the transaction.
- **Replay protection.** A transaction must be signed with the chain ID. The node rejects a transaction without it and returns the error `only replay-protected (EIP-155) transactions allowed over RPC`.
- **Receipts.** The `effectiveGasPrice` field is present only in receipts for type 2 transactions in this release.
- **Finality.** CometBFT commits provide finality under its consensus assumptions. Block intervals vary. Track transaction inclusion to determine when execution completes.

## Block values in contracts

| Solidity value | What it returns on Upside |
| --- | --- |
| `block.chainid` | `61873` on mainnet |
| `block.number` | The block height |
| `block.timestamp` | The block time in seconds |
| `block.basefee` | The base fee in `aup` |
| `block.coinbase` | The operator address of the validator that proposed the block |
| `block.prevrandao` | Always zero. Don't use it as a source of randomness. |
| `block.gaslimit` | The configured block gas limit in a transaction (150,000,000 in the mainnet sample). Zero in `eth_call` and `eth_estimateGas`. |

Because `block.gaslimit` is zero in a simulation, a gas estimate can be too low for code whose cost depends on that value.

## Logs and events

- **Range limit.** With the documented RPC configuration, in `eth_getLogs` the difference between `toBlock` and `fromBlock` can be at most 10,000. A larger range fails with the error `maximum [from, to] blocks distance: 10000`.
- **Filters.** A node keeps a filter in memory and can remove it once five minutes pass without an `eth_getFilterChanges` call. Requests routed to another backend cannot find that filter. For reliable results, use `eth_getLogs` with a block range, or use a WebSocket subscription.
- **Subscriptions.** `eth_subscribe` works over the WebSocket endpoint for new block headers and for logs. Over WebSocket, the request `id` must be a number.

## Tools

The examples use Foundry 1.9.0, ethers 6.17.0, viem 2.57.16, and web3.py 8.0.0. Pin your application's tool versions and test custom-chain support, signing, and transaction handling on localnet or testnet before mainnet deployment.

### Use viem for contract writes

This version of viem asks the node to fill in a transaction with the `eth_fillTransaction` method. The node computes that gas estimate against the state at block 1. For a contract deployed after block 1, the estimate covers only the base cost of the transaction, and the node then rejects the transaction or the transaction runs out of gas.

To avoid the problem, estimate the gas with an explicit block tag and pass the result:

```javascript
const gas = await publicClient.estimateContractGas({
  address,
  abi,
  functionName: "setNumber",
  args: [11n],
  account,
  blockTag: "latest",
});

const hash = await walletClient.writeContract({
  address,
  abi,
  functionName: "setNumber",
  args: [11n],
  gas,
});
```

## Connect a library

Each of the following examples reads the chain ID, the latest block number, and the balance of the `upside-1` validator account from mainnet.

### ethers

```javascript
import { ethers } from "ethers";

const provider = new ethers.JsonRpcProvider("https://rpc.mainnet.upside.net");
const address = "0x111f75f912cbad1ac95fda5f9d0fd9457560a978";

console.log((await provider.getNetwork()).chainId);
console.log(await provider.getBlockNumber());
console.log(ethers.formatEther(await provider.getBalance(address)));
```

### viem

```javascript
import { createPublicClient, defineChain, formatEther, http } from "viem";

const upside = defineChain({
  id: 61873,
  name: "Upside",
  nativeCurrency: { name: "Upside", symbol: "UP", decimals: 18 },
  rpcUrls: {
    default: {
      http: ["https://rpc.mainnet.upside.net"],
      webSocket: ["wss://rpc.mainnet.upside.net/ws"],
    },
  },
  blockExplorers: {
    default: { name: "Upside Explorer", url: "https://explorer.mainnet.upside.net" },
  },
});

const client = createPublicClient({ chain: upside, transport: http() });
const address = "0x111f75f912cbad1ac95fda5f9d0fd9457560a978";

console.log(await client.getChainId());
console.log(await client.getBlockNumber());
console.log(formatEther(await client.getBalance({ address })));
```

### web3.py

```python
from web3 import Web3

w3 = Web3(Web3.HTTPProvider("https://rpc.mainnet.upside.net"))
address = Web3.to_checksum_address("0x111f75f912cbad1ac95fda5f9d0fd9457560a978")

print(w3.eth.chain_id)
print(w3.eth.block_number)
print(w3.from_wei(w3.eth.get_balance(address), "ether"))
```

## Deploy a contract with Foundry

This walkthrough deploys a small contract on public testnet and calls it with Foundry. Use a separate account funded with test tokens. Each write sends a transaction to the selected network.

1. Create a project folder with a `foundry.toml` file:

   ```toml
   [profile.default]
   src = "src"
   out = "out"
   evm_version = "prague"
   ```

2. Add the contract as `src/Counter.sol`:

   ```solidity
   // SPDX-License-Identifier: MIT
   pragma solidity ^0.8.20;

   contract Counter {
       uint256 public number;

       event NumberSet(address indexed by, uint256 newNumber);

       function setNumber(uint256 newNumber) public {
           number = newNumber;
           emit NumberSet(msg.sender, newNumber);
       }

       function increment() public {
           number++;
           emit NumberSet(msg.sender, number);
       }
   }
   ```

3. Compile the contract:

   ```bash
   forge build
   ```

4. Import your private key into an encrypted keystore. The command prompts for the key and for a password.

   ```bash
   cast wallet import deployer --interactive
   ```

5. Set the endpoint and check the chain ID. The second command must print `61872` for public testnet.

   ```bash
   export ETH_RPC_URL="https://rpc.testnet.upside.net"
   cast chain-id
   ```

6. Deploy the contract:

   ```bash
   forge create src/Counter.sol:Counter --rpc-url "$ETH_RPC_URL" --account deployer --broadcast
   ```

   The command prints the address of the new contract:

   ```text
   Deployer: <your-address>
   Deployed to: <contract-address>
   Transaction hash: <transaction-hash>
   ```

7. Set `CONTRACT_ADDRESS` to the deployed contract address, then read the stored number. The command prints `0`.

   ```bash
   cast call "$CONTRACT_ADDRESS" "number()(uint256)"
   ```

8. Send a transaction that sets the number:

   ```bash
   cast send "$CONTRACT_ADDRESS" "setNumber(uint256)" 42 --account deployer
   ```

   The output includes a `status` line that reads `1 (success)`.

9. Read the number again. The command prints `42`.

   ```bash
   cast call "$CONTRACT_ADDRESS" "number()(uint256)"
   ```
