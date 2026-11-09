package app

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	upsidetypes "github.com/Upside-Foundation/upside/x/upside/types"
)

// UpgradeName must match the name in the governance proposal's
// MsgSoftwareUpgrade plan exactly.
const UpgradeName = "v2.0.0"

// MainnetBlocksPerYear is a Julian year divided by mainnet's measured average
// block time of 5.457 s (genesis 2026-02-18 to height 3,568,196 on
// 2026-10-01). x/mint spreads a year's minting over this many blocks, so it
// sets how closely the real half-life tracks the mint_half_life_seconds param.
const MainnetBlocksPerYear = 5_783_000

// RegisterUpgradeHandlers returns if store loader is overridden
func (app *App) RegisterUpgradeHandlers(cdc codec.BinaryCodec, maxVersion int64) bool {
	app.UpgradeKeeper.SetUpgradeHandler(UpgradeName, func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		vm, err := app.ModuleManager.RunMigrations(ctx, app.configurator, fromVM)
		if err != nil {
			return vm, err
		}
		if err := app.applySupplyCap(ctx); err != nil {
			return vm, fmt.Errorf("%s: %w", UpgradeName, err)
		}
		return vm, nil
	})
	return false
}

// applySupplyCap switches x/mint to the capped half-life curve. From the next
// mint onward x/upside's MintInflationFn sets the rate; this sets the two
// params it reads.
func (app *App) applySupplyCap(ctx context.Context) error {
	mintParams, err := app.MintKeeper.Params.Get(ctx)
	if err != nil {
		return err
	}
	// The curve reads the supply of the mint denom, and x/mint scales the rate
	// by the supply of the bond denom. They must be the same token.
	bondDenom, err := app.StakingKeeper.BondDenom(ctx)
	if err != nil {
		return err
	}
	if mintParams.MintDenom != bondDenom {
		return fmt.Errorf("mint denom %q differs from bond denom %q", mintParams.MintDenom, bondDenom)
	}

	mintParams.BlocksPerYear = MainnetBlocksPerYear
	if err := mintParams.Validate(); err != nil {
		return err
	}
	if err := app.MintKeeper.Params.Set(ctx, mintParams); err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	upsideParams := app.UpsideKeeper.GetParams(sdkCtx)
	upsideParams.MintHalfLifeSeconds = upsidetypes.MintHalfLifeSecondsDefaultValue
	return app.UpsideKeeper.SetParams(sdkCtx, upsideParams) // validates all x/upside params
}
