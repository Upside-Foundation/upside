package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/0xUp5ide/upside-foundation/x/upside/types"
)

// MintInflationFn returns the inflation function x/mint calls at the start of
// every block. It replaces x/mint's default (a rate between inflation_min and
// inflation_max that follows the bonded ratio) with the capped half-life
// curve in types.MintInflationRate. x/mint still mints and pays the result
// to the fee collector; the bonded ratio and the min/max params are unused.
//
// The curve reads the supply of the mint denom, but x/mint multiplies the
// rate by the supply of the staking bond denom. If governance ever makes the
// two denoms differ, the result no longer follows the curve and can pass the
// cap, so the function returns zero and minting stops until they match again.
//
// While the mint_paused param is true the function also returns zero. The
// curve depends only on the current supply, so after the pause minting
// continues at the rate it had before the pause.
//
// It also returns zero while x/mint's blocks_per_year is outside
// [types.MinBlocksPerYear, types.MaxBlocksPerYear]. blocks_per_year scales
// the real half-life, so without the bound one x/mint proposal could mint
// most of the tokens below the cap within a few blocks.
func (k Keeper) MintInflationFn(stakingKeeper types.StakingKeeper) minttypes.InflationCalculationFn {
	return func(ctx context.Context, _ minttypes.Minter, params minttypes.Params, _ sdkmath.LegacyDec) sdkmath.LegacyDec {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		upsideParams := k.GetParams(sdkCtx)
		if upsideParams.MintPaused {
			return sdkmath.LegacyZeroDec()
		}
		bondDenom, err := stakingKeeper.BondDenom(ctx)
		if err != nil || bondDenom != params.MintDenom {
			sdkCtx.Logger().Error("capped mint skipped: mint denom must equal bond denom",
				"mint_denom", params.MintDenom, "bond_denom", bondDenom, "err", err)
			return sdkmath.LegacyZeroDec()
		}
		if params.BlocksPerYear < types.MinBlocksPerYear || params.BlocksPerYear > types.MaxBlocksPerYear {
			sdkCtx.Logger().Error("capped mint skipped: blocks_per_year out of range",
				"blocks_per_year", params.BlocksPerYear, "min", types.MinBlocksPerYear, "max", types.MaxBlocksPerYear)
			return sdkmath.LegacyZeroDec()
		}
		supply := k.bankKeeper.GetSupply(ctx, params.MintDenom).Amount
		return types.MintInflationRate(supply, upsideParams.MintHalfLifeSeconds, params.BlocksPerYear)
	}
}
