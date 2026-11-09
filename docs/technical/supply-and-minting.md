---
title: "Supply and minting reference"
description: How the chain limits minting to 2.1 billion UP and computes the amount minted in each block.
slug: /technical/supply-and-minting
---

# Supply and minting reference

This page describes the code that limits the supply of UP and sets how much is minted in each block. For reward allocation and governance queries, see [Rewards and governance](../staking/governance.md).

## Summary

- The mint module stops minting when the supply of `aup` reaches 2,100,000,000 UP. The cap is a fixed value in the code.
- Each block mints a fixed fraction of the tokens that remain below the cap. The fraction comes from a half-life, which is two years on mainnet.
- Governance can change the half-life and can pause minting. A parameter change can't raise the cap. Changing the cap requires a software upgrade.

## Where the code lives

| File | Contents |
| --- | --- |
| `x/upside/types/supply.go` | The cap, the constants, and `MintInflationRate`, which computes the rate. |
| `x/upside/keeper/mint.go` | `MintInflationFn`, which returns the function that the mint module calls at the start of every block. |
| `x/upside/types/params.go` | The `mint_half_life_seconds` parameter, its default, and its lower bound. |
| `x/upside/keeper/msg_server.go` | The handlers for `MsgUpdateParams` and `MsgSetMintPaused`. |
| `app/upgrades.go` | The v2.0.0 upgrade handler that moved mainnet to this curve. |

The Cosmos SDK mint module, `x/mint`, still does the minting. The Upside code replaces only the function that the mint module uses to compute the inflation rate.

## The curve

The target curve is the following, where `S` is the supply, `cap` is 2,100,000,000 UP, and `t` is time in years:

```text
S(t) = cap - (cap - S0) × e^(-λt)
λ    = ln 2 / half-life
```

The code works from the slope of that curve:

```text
dS/dt = λ × (cap - S)
```

At every moment, a fixed fraction `λ` of the tokens that remain below the cap is minted per year. This form has no start time and no block height. Its only input is the current supply, and the curve continues from whatever the supply is.

## Fixed values

The following values are fixed in `x/upside/types/supply.go` and `x/upside/types/params.go`. No message or parameter can change them.

| Name | Value | Meaning |
| --- | --- | --- |
| `SupplyCap` | `2100000000000000000000000000` | The cap in `aup`. It equals 2,100,000,000 UP. |
| `ln2` | `0.693147180559945309` | The natural logarithm of 2 to 18 decimal places, rounded down. |
| `SecondsPerYear` | `31557600` | A Julian year of 365.25 days. |
| `MinBlocksPerYear` | `2891500` | The lowest value of the mint module's `blocks_per_year` that allows minting. |
| `MaxBlocksPerYear` | `11566000` | The highest value of `blocks_per_year` that allows minting. |
| `MintHalfLifeSecondsDefaultValue` | `63115200` | The default half-life, which is two Julian years. |
| `MinMintHalfLifeSeconds` | `31557600` | The shortest half-life that parameter validation accepts, which is one Julian year. |

## MintInflationRate

```go
func MintInflationRate(supply sdkmath.Int, halfLifeSeconds, blocksPerYear uint64) sdkmath.LegacyDec
```

`MintInflationRate` returns the annual inflation rate that keeps the supply on the curve. It's a pure function. It doesn't read chain state and it doesn't mint.

The function does the following:

1. Returns zero if `halfLifeSeconds` is zero, `blocksPerYear` is zero, `supply` isn't positive, or `supply` is at or above `SupplyCap`.
2. Computes the tokens that remain below the cap as `SupplyCap - supply`.
3. Computes `λ` per year as `ln2 × SecondsPerYear / halfLifeSeconds`.
4. Limits `λ` to `blocksPerYear`. With this limit, one block mints at most all of the remaining tokens, whatever the half-life is.
5. Returns `λ × remaining / supply`.

Every division rounds down to 18 decimal places. With the mainnet half-life of 63,115,200 seconds, `λ` is `0.346573590279972654`.

## MintInflationFn

```go
func (k Keeper) MintInflationFn(stakingKeeper types.StakingKeeper) minttypes.InflationCalculationFn
```

`MintInflationFn` returns the function that the mint module calls at the start of every block. The app passes it to the mint module in place of the Cosmos SDK default.

The returned function does the following:

1. Reads the `x/upside` parameters. If `mint_paused` is `true`, it returns zero.
2. Reads the staking bond denomination. If the lookup fails, or if the bond denomination differs from the mint module's `mint_denom`, it logs an error and returns zero.
3. Checks the mint module's `blocks_per_year`. If the value is below `MinBlocksPerYear` or above `MaxBlocksPerYear`, it logs an error and returns zero.
4. Reads the total supply of `mint_denom` from the bank module and returns `MintInflationRate(supply, mint_half_life_seconds, blocks_per_year)`.

The function ignores the bonded ratio and the stored minter. As a result, the mint module's `inflation_min`, `inflation_max`, `inflation_rate_change`, and `goal_bonded` parameters have no effect on mainnet.

## What happens in each block

The mint module runs at the start of every block and does the following:

1. Sets the inflation rate to the value that the function from `MintInflationFn` returns.
2. Sets the annual provisions to the inflation rate multiplied by the total supply of the bond denomination.
3. Divides the annual provisions by `blocks_per_year` and drops any fraction of an `aup`. The result is the amount for the block.
4. Mints that amount and sends it to the fee collector account.
5. Emits a `mint` event whose attributes include `bonded_ratio`, `inflation`, `annual_provisions`, and `amount`.

The net effect is that each block mints `λ / blocks_per_year` of the tokens that remain below the cap.

The distribution module runs next in the same block. It splits the fee collector balance, which holds the newly minted tokens and the transaction fees from the previous block. The `community_tax` share goes to the community pool, and the rest goes to validators and their delegators. On mainnet, `community_tax` is 0.93.

## Parameters

| Parameter | Module | Mainnet value | Effect |
| --- | --- | --- | --- |
| `mint_half_life_seconds` | `x/upside` | `63115200` | The half-life of the curve. Validation rejects values below `31557600`. The value must fit in an unsigned 64-bit integer. |
| `mint_paused` | `x/upside` | `false` | When `true`, the inflation rate is zero. Only `MsgSetMintPaused` changes it. |
| `blocks_per_year` | `x/mint` | `5783000` | The number of blocks over which the mint module spreads a year of minting. Minting stops if the value is outside 2,891,500 to 11,566,000. |
| `mint_denom` | `x/mint` | `aup` | The token that the mint module mints. Minting stops if it differs from the staking `bond_denom`. |

Governance can change every value in this table.

The curve advances with the block count, so the half-life in real time depends on the real block time:

```text
real half-life = mint_half_life_seconds × blocks_per_year / (blocks produced in a real year)
```

The mainnet value of `blocks_per_year` assumes an average block time of 5.457 seconds. If blocks are slower on average, the real half-life is longer than the parameter. If blocks are faster, it's shorter.

The bounds limit how far governance can shorten the real half-life. With the shortest half-life and the lowest `blocks_per_year` that allows minting, the real half-life is about six months at a block time of 5.457 seconds.

### Read the current values

Set `RPC_URL` as described in [Query the chain](query-the-chain.md), and then run the following commands:

```bash
upsided query upside params --node "$RPC_URL"
upsided query mint params --node "$RPC_URL"
upsided query mint inflation --node "$RPC_URL"
upsided query mint annual-provisions --node "$RPC_URL"
upsided query bank total-supply-of aup --node "$RPC_URL"
```

The first command returns output like this:

```yaml
enable_auto_deployment: false
ibc_timeout: "86400000000000"
ibc_upside_denom: ibc/6B5A664BF0AF4F71B2F0BAA33141E2F1321242FBD5D19762F541EC971ACB0865
max_callback_gas: "50000"
mint_half_life_seconds: "63115200"
mint_paused: false
upside_admin: ""
```

## Messages

Both messages can be sent only by the governance module account, so they take effect through a governance proposal. On mainnet, that account is `up10d07y265gmmuvt4z0w9aw880jnsr700jwh8mz8`.

### MsgSetMintPaused

`MsgSetMintPaused` sets the `mint_paused` parameter. Its type URL is `/upside.MsgSetMintPaused`.

| Field | Type | Description |
| --- | --- | --- |
| `authority` | string | The address of the governance module account. |
| `paused` | bool | `true` to pause minting, `false` to resume it. |

A governance proposal runs its messages at the end of the block in which voting ends. From the next block, the inflation function returns zero and the mint module mints nothing. The inflation rate and the annual provisions read zero until minting resumes. During a pause, the supply changes only if another module mints or burns tokens. With the mainnet governance parameters, voting lasts 48 hours, or 24 hours for an expedited proposal. The handler emits a `set_mint_paused` event with a `paused` attribute.

A second proposal with `paused` set to `false` resumes minting. The rate depends only on the current supply, so minting continues along the curve from the supply at that time. The tokens that the paused blocks would have minted aren't minted later. If supply and the other parameters remain unchanged during the pause, subsequent issuance resumes from the same point on the curve.

The following proposal file pauses minting:

```json
{
  "messages": [
    {
      "@type": "/upside.MsgSetMintPaused",
      "authority": "up10d07y265gmmuvt4z0w9aw880jnsr700jwh8mz8",
      "paused": true
    }
  ],
  "metadata": "",
  "deposit": "10000aup",
  "title": "Pause minting",
  "summary": "Set mint_paused to true."
}
```

### MsgUpdateParams

`MsgUpdateParams` replaces every `x/upside` parameter. Its type URL is `/upside.MsgUpdateParams`. The message must supply all seven parameters. The governance module checks the parameters when the proposal is submitted, so a proposal that leaves out `mint_half_life_seconds` or `ibc_upside_denom` is rejected at submission. Any other parameter that is left out is written as zero, empty, or `false`.

Two rules apply to the minting parameters:

- **The half-life has a lower bound.** Values of `mint_half_life_seconds` below 31,557,600 fail parameter validation, so the proposal is rejected when it is submitted.
- **The message can't change `mint_paused`.** If its `mint_paused` value differs from the stored value, the message fails when the proposal runs, and the proposal ends with the status `PROPOSAL_STATUS_FAILED`. A parameter proposal that omits `mint_paused` carries `false`, so it fails while minting is paused.

The following example changes the half-life to three Julian years. Before constructing a proposal, query the current `x/upside` parameters and copy every unchanged value, including `mint_paused`. These values are a historical example. Submitting them unchanged can overwrite newer governance decisions:

```json
{
  "messages": [
    {
      "@type": "/upside.MsgUpdateParams",
      "authority": "up10d07y265gmmuvt4z0w9aw880jnsr700jwh8mz8",
      "params": {
        "ibc_upside_denom": "ibc/6B5A664BF0AF4F71B2F0BAA33141E2F1321242FBD5D19762F541EC971ACB0865",
        "ibc_timeout": "86400000000000",
        "upside_admin": "",
        "enable_auto_deployment": false,
        "max_callback_gas": "50000",
        "mint_half_life_seconds": "94672800",
        "mint_paused": false
      }
    }
  ],
  "metadata": "",
  "deposit": "10000aup",
  "title": "Set the mint half-life to three years",
  "summary": "Set mint_half_life_seconds to 94672800."
}
```

## The v2.0.0 upgrade

Mainnet adopted this curve in the v2.0.0 upgrade at block 3,898,900 on October 22, 2026 (UTC). After the module migrations, the upgrade handler did the following:

1. Checked that the mint denomination equals the staking bond denomination. The upgrade fails if they differ.
2. Set the mint module's `blocks_per_year` to 5,783,000.
3. Set `mint_half_life_seconds` to the default of 63,115,200.

The mint module has taken its rate from `MintInflationFn` since block 3,898,900. That block minted 86.996573607358165295 UP.

## Worked example

This example works through the minting in mainnet block 4,130,439.

| Step | Value |
| --- | --- |
| Supply after block 4,130,438 | 668,360,821.088657945911986902 UP |
| Remaining below the cap | 1,431,639,178.911342054088013098 UP |
| `λ` | 0.346573590279972654 |
| Inflation rate, `λ × remaining / supply` | 0.742365971441284236 |
| Annual provisions, rate multiplied by supply | 496,168,330.220775927410258746969612305291076872 UP |
| Minted in the block, annual provisions divided by 5,783,000 and rounded down to a whole `aup` | 85.797739965550047969 UP |
| Supply after block 4,130,439 | 668,360,906.886397911462034871 UP |

To read these values from the chain, run the following commands:

```bash
upsided query bank total-supply-of aup --height 4130438 --node "$RPC_URL"
upsided query mint inflation --height 4130439 --node "$RPC_URL"
upsided query mint annual-provisions --height 4130439 --node "$RPC_URL"
upsided query bank total-supply-of aup --height 4130439 --node "$RPC_URL"
```

## Scope of the minting rules

- **Other minting paths.** The cap applies only to the mint module. Other code paths can create `aup` outside this limit, such as the conversion of an IBC voucher. If such a path raises the supply, the mint module mints less. On mainnet, governance has turned IBC token transfers off, and `aup` is the only token in circulation.
- **Burns.** Other modules handle token burns, including those caused by validator slashing. A burn reduces the supply, and the minting curve continues from the new supply.
- **Reward allocation.** The distribution module allocates newly minted tokens.
- **Block timing.** `blocks_per_year` is a governance parameter with fixed bounds. Adjusting it to reflect changes in block time requires a governance update.
