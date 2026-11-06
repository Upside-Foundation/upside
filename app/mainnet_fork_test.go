//go:build mainnetfork

package app

// Runs the v2.0.0 upgrade against a copy of real mainnet application state.
//
//	MAINNET_DB=<dir containing application.db> MAINNET_APP_HASH=<base64> \
//	  go test -tags "objstore pebbledb mainnet mainnetfork" ./app -run TestMainnetFork -v
//
// MAINNET_APP_HASH is the app_hash in the header of the block after the one
// the copy was taken at. Nothing is committed: Commit is never called, so the
// changes are discarded. Run it on a clone of the database anyway.

import (
	"encoding/base64"
	"math"
	"os"
	"sync"
	"testing"
	"time"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	upsidetypes "github.com/0xUp5ide/upside-foundation/x/upside/types"
)

var setConfigOnce sync.Once

func TestMainnetFork(t *testing.T) {
	// A real node sets the "up" address prefix at startup; without it,
	// module addresses don't round-trip and bank sends fail.
	setConfigOnce.Do(SetConfig)
	dir := os.Getenv("MAINNET_DB")
	if dir == "" {
		t.Skip("MAINNET_DB not set")
	}
	db, err := dbm.NewGoLevelDB("application", dir, nil)
	require.NoError(t, err)
	defer db.Close()

	app := New(log.NewNopLogger(), db, nil, true,
		simtestutil.AppOptionsMap{flags.FlagHome: t.TempDir()},
		baseapp.SetChainID("upside_61873-1"))

	height := app.LastBlockHeight()
	hash := app.LastCommitID().Hash
	t.Logf("loaded mainnet state at height %d, app hash %X", height, hash)

	// 1. The copy is mainnet state: its stored app hash matches the next
	// block's header. LastCommitID is read from the database, not recomputed,
	// so this does not show that this code computes the same state.
	if want := os.Getenv("MAINNET_APP_HASH"); want != "" {
		wantBytes, err := base64.StdEncoding.DecodeString(want)
		require.NoError(t, err)
		require.Equal(t, wantBytes, hash, "app hash differs from mainnet")
		t.Logf("app hash matches mainnet block %d header", height+1)
	}

	ctx := app.NewUncachedContext(false, cmtproto.Header{
		ChainID: "upside_61873-1", Height: height + 1, Time: time.Now().UTC(),
	})
	// Changes stay in the working store; Commit is never called and the
	// database is a throwaway clone.

	denom := "aup"
	supply := app.BankKeeper.GetSupply(ctx, denom).Amount
	mintParams, err := app.MintKeeper.Params.Get(ctx)
	require.NoError(t, err)
	t.Logf("before: supply %s UP, mint %+v", sdkmath.LegacyNewDecFromIntWithPrec(supply, 18).String(), mintParams)
	require.Zero(t, app.UpsideKeeper.GetParams(ctx).MintHalfLifeSeconds, "field 6 should be absent on mainnet")

	// 2. Apply v2.0.0 to real state.
	require.NoError(t, app.UpgradeKeeper.ApplyUpgrade(ctx, upgradetypes.Plan{Name: UpgradeName, Height: ctx.BlockHeight()}))
	up := app.UpsideKeeper.GetParams(ctx)
	require.Equal(t, upsidetypes.MintHalfLifeSecondsDefaultValue, up.MintHalfLifeSeconds)
	mintParams, err = app.MintKeeper.Params.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, uint64(MainnetBlocksPerYear), mintParams.BlocksPerYear)
	t.Logf("after upgrade: x/upside params %+v", up)

	// 3. Mint a few blocks with the module as app.go wires it.
	feeCollector := app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	c, _ := sdkmath.LegacyNewDecFromInt(upsidetypes.SupplyCap).Float64()
	for i := 0; i < 5; i++ {
		s := app.BankKeeper.GetSupply(ctx, denom).Amount
		fees := app.BankKeeper.GetBalance(ctx, feeCollector, denom).Amount
		mintBeginBlock(t, app, ctx)
		minted := app.BankKeeper.GetSupply(ctx, denom).Amount.Sub(s)
		require.Equal(t, minted, app.BankKeeper.GetBalance(ctx, feeCollector, denom).Amount.Sub(fees))
		sf, _ := sdkmath.LegacyNewDecFromInt(s).Float64()
		got, _ := sdkmath.LegacyNewDecFromInt(minted).Float64()
		want := math.Ln2 / 2 * (c - sf) / MainnetBlocksPerYear
		require.InEpsilon(t, want, got, 1e-9)
		t.Logf("block +%d: minted %.6f UP (expected %.6f)", i+1, got/1e18, want/1e18)
	}
	minter, err := app.MintKeeper.Minter.Get(ctx)
	require.NoError(t, err)
	t.Logf("inflation after upgrade: %s", minter.Inflation)
}
