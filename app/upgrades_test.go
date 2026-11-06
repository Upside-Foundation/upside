package app

import (
	"math"
	"testing"

	"cosmossdk.io/core/appmodule"
	sdkmath "cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	upsidekeeper "github.com/0xUp5ide/upside-foundation/x/upside/keeper"
	upsidetypes "github.com/0xUp5ide/upside-foundation/x/upside/types"
)

// setupSupplyCapTest returns an app whose state looks like mainnet before
// the upgrade: total supply of 615.26M UP and x/upside params stored without
// mint_half_life_seconds.
func setupSupplyCapTest(t *testing.T) (*App, sdk.Context, string) {
	t.Helper()
	app := Setup(t, "")
	ctx := app.BaseApp.NewContext(false)

	denom, err := app.StakingKeeper.BondDenom(ctx)
	require.NoError(t, err)

	target, ok := sdkmath.NewIntFromString("615260010397900599274361153")
	require.True(t, ok)
	extra := target.Sub(app.BankKeeper.GetSupply(ctx, denom).Amount)
	coins := sdk.NewCoins(sdk.NewCoin(denom, extra))
	require.NoError(t, app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins))
	holder := sdk.AccAddress([]byte("supply-cap-test-holder"))
	require.NoError(t, app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, holder, coins))
	require.Equal(t, target, app.BankKeeper.GetSupply(ctx, denom).Amount)

	// Mainnet's stored params predate field 6. Write them raw, skipping
	// Validate, the way they exist on chain today.
	p := app.UpsideKeeper.GetParams(ctx)
	p.MintHalfLifeSeconds = 0
	ctx.KVStore(app.keys[upsidetypes.StoreKey]).Set(upsidetypes.ParamsKey, app.appCodec.MustMarshal(&p))
	require.Zero(t, app.UpsideKeeper.GetParams(ctx).MintHalfLifeSeconds)

	return app, ctx, denom
}

func mintBeginBlock(t *testing.T, app *App, ctx sdk.Context) {
	t.Helper()
	// The mint module exactly as app.go registered it, with its inflation function.
	mod, ok := app.ModuleManager.Modules[minttypes.ModuleName].(appmodule.HasBeginBlocker)
	require.True(t, ok)
	require.NoError(t, mod.BeginBlock(ctx))
}

func TestUpgradeAppliesSupplyCap(t *testing.T) {
	app, ctx, denom := setupSupplyCapTest(t)
	feeCollector := app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)

	// Before the upgrade the stored half-life is zero, so the curve mints nothing.
	before := app.BankKeeper.GetSupply(ctx, denom).Amount
	mintBeginBlock(t, app, ctx)
	require.Equal(t, before, app.BankKeeper.GetSupply(ctx, denom).Amount)

	plan := upgradetypes.Plan{Name: UpgradeName, Height: ctx.BlockHeight()}
	require.NoError(t, app.UpgradeKeeper.ApplyUpgrade(ctx, plan))

	require.Equal(t, upsidetypes.MintHalfLifeSecondsDefaultValue, app.UpsideKeeper.GetParams(ctx).MintHalfLifeSeconds)
	mintParams, err := app.MintKeeper.Params.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(MainnetBlocksPerYear), mintParams.BlocksPerYear)

	feesBefore := app.BankKeeper.GetBalance(ctx, feeCollector, denom).Amount
	supply := app.BankKeeper.GetSupply(ctx, denom).Amount
	mintBeginBlock(t, app, ctx)
	minted := app.BankKeeper.GetSupply(ctx, denom).Amount.Sub(supply)

	// One block mints ln2/2 / 5,783,000 of the tokens below the cap: about 88.98 UP.
	s, _ := sdkmath.LegacyNewDecFromInt(supply).Float64()
	c, _ := sdkmath.LegacyNewDecFromInt(upsidetypes.SupplyCap).Float64()
	got, _ := sdkmath.LegacyNewDecFromInt(minted).Float64()
	require.InEpsilon(t, math.Ln2/2*(c-s)/MainnetBlocksPerYear, got, 1e-9)
	require.InDelta(t, 88.98, got/1e18, 0.01)

	// All of it goes to the fee collector; x/distribution splits it from there.
	require.Equal(t, minted, app.BankKeeper.GetBalance(ctx, feeCollector, denom).Amount.Sub(feesBefore))

	minter, err := app.MintKeeper.Minter.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, upsidetypes.MintInflationRate(supply, upsidetypes.MintHalfLifeSecondsDefaultValue, MainnetBlocksPerYear), minter.Inflation)
}

func TestUpgradeRejectsMintDenomMismatch(t *testing.T) {
	app, ctx, _ := setupSupplyCapTest(t)
	mintParams, err := app.MintKeeper.Params.Get(ctx)
	require.NoError(t, err)
	mintParams.MintDenom = "otherdenom"
	require.NoError(t, app.MintKeeper.Params.Set(ctx, mintParams))

	err = app.UpgradeKeeper.ApplyUpgrade(ctx, upgradetypes.Plan{Name: UpgradeName, Height: ctx.BlockHeight()})
	require.ErrorContains(t, err, "differs from bond denom")
}

// A later change that makes the mint denom differ from the bond denom stops
// minting instead of applying the rate to the wrong supply.
func TestMintDenomMismatchAfterUpgradeMintsNothing(t *testing.T) {
	app, ctx, denom := setupSupplyCapTest(t)
	require.NoError(t, app.UpgradeKeeper.ApplyUpgrade(ctx, upgradetypes.Plan{Name: UpgradeName, Height: ctx.BlockHeight()}))

	mintParams, err := app.MintKeeper.Params.Get(ctx)
	require.NoError(t, err)
	mintParams.MintDenom = "otherdenom"
	require.NoError(t, app.MintKeeper.Params.Set(ctx, mintParams))

	before := app.BankKeeper.GetSupply(ctx, denom).Amount
	mintBeginBlock(t, app, ctx)
	require.Equal(t, before, app.BankKeeper.GetSupply(ctx, denom).Amount)
	require.True(t, app.BankKeeper.GetSupply(ctx, "otherdenom").Amount.IsZero())

	minter, err := app.MintKeeper.Minter.Get(ctx)
	require.NoError(t, err)
	require.True(t, minter.Inflation.IsZero())

	// Setting the denom back resumes minting.
	mintParams.MintDenom = denom
	require.NoError(t, app.MintKeeper.Params.Set(ctx, mintParams))
	mintBeginBlock(t, app, ctx)
	require.True(t, app.BankKeeper.GetSupply(ctx, denom).Amount.GT(before))
}

func TestUpgradeAtCapMintsNothing(t *testing.T) {
	app, ctx, denom := setupSupplyCapTest(t)
	require.NoError(t, app.UpgradeKeeper.ApplyUpgrade(ctx, upgradetypes.Plan{Name: UpgradeName, Height: ctx.BlockHeight()}))

	extra := upsidetypes.SupplyCap.Sub(app.BankKeeper.GetSupply(ctx, denom).Amount)
	coins := sdk.NewCoins(sdk.NewCoin(denom, extra))
	require.NoError(t, app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins))
	require.NoError(t, app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, sdk.AccAddress([]byte("supply-cap-test-holder")), coins))
	require.Equal(t, upsidetypes.SupplyCap, app.BankKeeper.GetSupply(ctx, denom).Amount)

	mintBeginBlock(t, app, ctx)
	require.Equal(t, upsidetypes.SupplyCap, app.BankKeeper.GetSupply(ctx, denom).Amount)
}

// Governance pauses minting with MsgSetMintPaused. While paused, x/mint mints
// nothing; after the pause it mints at the rate it had before, with no
// catch-up for the blocks that minted nothing. MsgUpdateParams cannot change
// mint_paused, so a params proposal that leaves it out fails instead of
// resuming minting.
func TestMintPausedByGovernance(t *testing.T) {
	app, ctx, denom := setupSupplyCapTest(t)
	require.NoError(t, app.UpgradeKeeper.ApplyUpgrade(ctx, upgradetypes.Plan{Name: UpgradeName, Height: ctx.BlockHeight()}))
	require.False(t, app.UpsideKeeper.GetParams(ctx).MintPaused)

	msgServer := upsidekeeper.NewMsgServerImpl(app.UpsideKeeper)
	authority := app.UpsideKeeper.GetAuthority()
	setPaused := func(paused bool) {
		_, err := msgServer.SetMintPaused(ctx, upsidetypes.NewMsgSetMintPaused(authority, paused))
		require.NoError(t, err)
		require.Equal(t, paused, app.UpsideKeeper.GetParams(ctx).MintPaused)
	}

	// MsgUpdateParams cannot pause.
	p := app.UpsideKeeper.GetParams(ctx)
	p.MintPaused = true
	_, err := msgServer.UpdateParams(ctx, upsidetypes.NewMsgUpdateParams(authority, p))
	require.ErrorContains(t, err, "use MsgSetMintPaused")
	require.False(t, app.UpsideKeeper.GetParams(ctx).MintPaused)

	// Only the governance authority can pause.
	_, err = msgServer.SetMintPaused(ctx, upsidetypes.NewMsgSetMintPaused(sdk.AccAddress([]byte("not-gov")).String(), true))
	require.ErrorContains(t, err, "invalid authority")

	supply := app.BankKeeper.GetSupply(ctx, denom).Amount
	want := upsidetypes.MintInflationRate(supply, upsidetypes.MintHalfLifeSecondsDefaultValue, MainnetBlocksPerYear)

	ctx = ctx.WithEventManager(sdk.NewEventManager())
	setPaused(true)
	var found bool
	for _, e := range ctx.EventManager().Events() {
		if e.Type == upsidetypes.EventTypeSetMintPaused {
			v, ok := e.GetAttribute(upsidetypes.AttributeKeyPaused)
			found = ok && v.Value == "true"
		}
	}
	require.True(t, found)

	for i := 0; i < 3; i++ {
		mintBeginBlock(t, app, ctx)
	}
	require.Equal(t, supply, app.BankKeeper.GetSupply(ctx, denom).Amount)
	minter, err := app.MintKeeper.Minter.Get(ctx)
	require.NoError(t, err)
	require.True(t, minter.Inflation.IsZero())
	require.True(t, minter.AnnualProvisions.IsZero())

	// A params proposal written without mint_paused (false) fails while
	// paused; one that repeats true succeeds and minting stays paused.
	p = app.UpsideKeeper.GetParams(ctx)
	p.MaxCallbackGas = 60000
	p.MintPaused = false
	_, err = msgServer.UpdateParams(ctx, upsidetypes.NewMsgUpdateParams(authority, p))
	require.ErrorContains(t, err, "mint_paused must stay true")
	p.MintPaused = true
	_, err = msgServer.UpdateParams(ctx, upsidetypes.NewMsgUpdateParams(authority, p))
	require.NoError(t, err)
	require.Equal(t, uint64(60000), app.UpsideKeeper.GetParams(ctx).MaxCallbackGas)
	mintBeginBlock(t, app, ctx)
	require.Equal(t, supply, app.BankKeeper.GetSupply(ctx, denom).Amount)

	setPaused(false)
	mintBeginBlock(t, app, ctx)
	minter, err = app.MintKeeper.Minter.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, want, minter.Inflation)
	minted := app.BankKeeper.GetSupply(ctx, denom).Amount.Sub(supply)
	require.InDelta(t, 88.98, float64(minted.Quo(sdkmath.NewInt(1e15)).Int64())/1e3, 0.01)
}

// The bounds are half and double the measured value the upgrade sets.
func TestBlocksPerYearBounds(t *testing.T) {
	require.Equal(t, MainnetBlocksPerYear, 2*upsidetypes.MinBlocksPerYear)
	require.Equal(t, 2*MainnetBlocksPerYear, upsidetypes.MaxBlocksPerYear)
}

// An x/mint proposal that sets blocks_per_year outside the bounds stops
// minting instead of shortening the half-life. Inside the bounds, minting
// follows the curve.
func TestBlocksPerYearOutOfRangeMintsNothing(t *testing.T) {
	app, ctx, denom := setupSupplyCapTest(t)
	require.NoError(t, app.UpgradeKeeper.ApplyUpgrade(ctx, upgradetypes.Plan{Name: UpgradeName, Height: ctx.BlockHeight()}))

	setBPY := func(bpy uint64) {
		mintParams, err := app.MintKeeper.Params.Get(ctx)
		require.NoError(t, err)
		mintParams.BlocksPerYear = bpy
		require.NoError(t, mintParams.Validate())
		require.NoError(t, app.MintKeeper.Params.Set(ctx, mintParams))
	}

	for _, bpy := range []uint64{1, upsidetypes.MinBlocksPerYear - 1, upsidetypes.MaxBlocksPerYear + 1, math.MaxInt64} {
		setBPY(bpy)
		before := app.BankKeeper.GetSupply(ctx, denom).Amount
		mintBeginBlock(t, app, ctx)
		require.Equal(t, before, app.BankKeeper.GetSupply(ctx, denom).Amount, "blocks_per_year %d", bpy)
	}

	for _, bpy := range []uint64{upsidetypes.MinBlocksPerYear, upsidetypes.MaxBlocksPerYear} {
		setBPY(bpy)
		before := app.BankKeeper.GetSupply(ctx, denom).Amount
		mintBeginBlock(t, app, ctx)
		minter, err := app.MintKeeper.Minter.Get(ctx)
		require.NoError(t, err)
		require.Equal(t, upsidetypes.MintInflationRate(before, upsidetypes.MintHalfLifeSecondsDefaultValue, bpy), minter.Inflation)
		require.True(t, app.BankKeeper.GetSupply(ctx, denom).Amount.GT(before), "blocks_per_year %d", bpy)
	}
}
