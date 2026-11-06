<!--
order: 8
-->

# Supply Curve

Total supply of `aup` is capped at 2,100,000,000 UP
(`2100000000000000000000000000` aup). The cap is a constant in
[`x/upside/types/supply.go`](../types/supply.go); changing it requires a
software upgrade.

## How minting works

The chain uses the standard Cosmos SDK `x/mint` module. `x/mint` asks an
inflation function for an annual rate at the start of every block, mints
rate × supply / `blocks_per_year`, and sends it to the fee collector. The app
passes `x/mint` the function in [`x/upside/keeper/mint.go`](../keeper/mint.go)
instead of the SDK default. `x/mint`'s `inflation_min`, `inflation_max`, and
`goal_bonded` parameters are therefore not used.

The function returns

```
inflation = λ × (cap − supply) / supply,   λ = ln 2 / half-life (per year)
```

so each block mints λ / `blocks_per_year` of the tokens still below the cap.
Its only input is the current total supply: there is no start time or block
height. Supply follows

```
S(t) = cap − (cap − S₀) × e^(−λt)
```

from whatever supply S₀ exists when the function first runs. Each half-life,
half of the tokens remaining below the cap are minted, and supply approaches
the cap without passing it.

## Guarantees

- λ is limited so that one block mints at most all of the remaining tokens,
  for any half-life setting.
- Every division rounds down.
- When supply is at or above the cap, the rate is zero.
- When `x/mint`'s `blocks_per_year` is below 2,891,500 or above 11,566,000,
  the rate is zero (see [Parameters](#parameters)).
- When the `x/mint` mint denom differs from the `x/staking` bond denom, the
  rate is zero. `x/mint` multiplies the rate by the bond denom's supply, so
  the curve only holds when both are the same token. The `v2.0.0` upgrade
  fails if they differ; after it, a governance change that makes them differ
  stops minting until they match again.

The cap limits `x/mint` only. Other paths that create `aup`, such as the IBC
voucher conversion, are not limited by it; when they raise supply, `x/mint`
mints correspondingly less.

## Pausing minting

Governance can stop minting by passing a proposal with
`MsgSetMintPaused{paused: true}`. The proposal executes at the end of a block;
from the next block, the inflation function returns zero, `x/mint` mints
nothing, and supply stays where it is. The `x/mint` minter's `inflation` and
`annual_provisions` read zero from that next block until minting resumes.

A second proposal with `MsgSetMintPaused{paused: false}` resumes minting. The rate
depends only on the current supply, so minting resumes at the rate it had when
it was paused. The tokens that the paused blocks would have minted are not
minted later; the whole supply schedule moves later by the length of the
pause.

`MsgSetMintPaused` is the only message that changes `mint_paused`.
`MsgUpdateParams` replaces every parameter, so it fails when its `mint_paused`
differs from the stored value. A params proposal written without the field
therefore fails while minting is paused, instead of resuming minting. This
includes params proposals that were already in their voting period when the
pause passed: they carry `mint_paused: false` and fail when they execute, so
they have to be submitted again with `mint_paused: true`.

Setting `x/mint`'s `inflation_min` and `inflation_max` to zero does not stop
minting, because the inflation function does not read them.

With the governance parameters in mainnet genesis, a pause takes effect when
the vote ends: 48 hours after the proposal reaches its minimum deposit for a
normal proposal, 24 hours for an expedited one. An expedited proposal that
fails the expedited threshold becomes a normal proposal and votes for another
48 hours.

## Distribution of minted tokens

Minted tokens go to the fee collector, and `x/distribution` splits them with
the rest of the block's fees: the `community_tax` share goes to the community
pool, and the remainder goes to validators and delegators.

## Parameters

- `mint_half_life_seconds` and `mint_paused` in this module (see
  [Parameters](07_params.md)).
- `blocks_per_year` in `x/mint`. It should match the real average block time;
  if it is too high, the real half-life is longer than the parameter says, and
  if it is too low, the real half-life is shorter. x/mint accepts any value
  from 1 to 2^63 − 1; at 1, each block would mint λ of the remaining tokens,
  about 35% at a 2-year half-life. The inflation function therefore returns
  zero while `blocks_per_year` is outside 2,891,500 to 11,566,000, half and
  double the 5,783,000 measured on mainnet (average blocks of 10.9 s and
  2.7 s). The bounds limit what a governance change to the parameter can do;
  they do not follow the real block time, which is set by each node's
  `timeout_commit` and by network conditions. At 5.457 s blocks, the shortest
  real half-life within the bounds is half of `mint_half_life_seconds`, at
  `blocks_per_year` = 2,891,500. If blocks become faster, the real half-life
  becomes shorter in the same proportion. Governance proposals that change
  `blocks_per_year` should be checked against the real block time.
