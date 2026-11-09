package types_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/Upside-Foundation/upside/x/upside/types"
)

const (
	twoYears          = types.MintHalfLifeSecondsDefaultValue
	mainnetBlocksYear = 5_783_000
)

// Mainnet supply at height 3,568,196 (2026-10-01).
var mainnetSupply, _ = sdkmath.NewIntFromString("615260010397900599274361153")

func toFloat(i sdkmath.Int) float64 {
	f, _ := sdkmath.LegacyNewDecFromInt(i).Float64()
	return f
}

// mintBlock runs one block of x/mint's arithmetic with the custom rate and
// returns the amount minted.
func mintBlock(supply sdkmath.Int, halfLife, blocksPerYear uint64) sdkmath.Int {
	params := minttypes.DefaultParams()
	params.MintDenom = "aup"
	params.BlocksPerYear = blocksPerYear
	minter := minttypes.DefaultInitialMinter()
	minter.Inflation = types.MintInflationRate(supply, halfLife, blocksPerYear)
	minter.AnnualProvisions = minter.NextAnnualProvisions(params, supply)
	return minter.BlockProvision(params).Amount
}

func TestSupplyCap(t *testing.T) {
	require.Equal(t, "2100000000000000000000000000", types.SupplyCap.String())
}

func TestDefaultHalfLife(t *testing.T) {
	require.Equal(t, uint64(63_115_200), types.MintHalfLifeSecondsDefaultValue)
	require.NoError(t, types.DefaultParams().Validate())

	for _, tc := range []struct {
		halfLife uint64
		valid    bool
	}{
		{0, false},
		{1, false},
		{types.MinMintHalfLifeSeconds - 1, false},
		{types.MinMintHalfLifeSeconds, true},
		{types.MintHalfLifeSecondsDefaultValue, true},
		{math.MaxUint64, true},
	} {
		p := types.DefaultParams()
		p.MintHalfLifeSeconds = tc.halfLife
		if tc.valid {
			require.NoError(t, p.Validate(), "half-life %d", tc.halfLife)
		} else {
			require.ErrorContains(t, p.Validate(), "at least 31557600 seconds", "half-life %d", tc.halfLife)
		}
	}
}

func TestMintInflationRateZeroCases(t *testing.T) {
	one := sdkmath.OneInt()
	for name, tc := range map[string]struct {
		supply            sdkmath.Int
		halfLife, perYear uint64
	}{
		"zero half-life":   {mainnetSupply, 0, mainnetBlocksYear},
		"zero blocks/year": {mainnetSupply, twoYears, 0},
		"zero supply":      {sdkmath.ZeroInt(), twoYears, mainnetBlocksYear},
		"negative supply":  {one.Neg(), twoYears, mainnetBlocksYear},
		"supply at cap":    {types.SupplyCap, twoYears, mainnetBlocksYear},
		"supply above cap": {types.SupplyCap.Add(one), twoYears, mainnetBlocksYear},
	} {
		t.Run(name, func(t *testing.T) {
			require.True(t, types.MintInflationRate(tc.supply, tc.halfLife, tc.perYear).IsZero())
		})
	}
}

// At today's mainnet supply and a 2-year half-life the rate is
// ln2/2 * (2100 - 615.26) / 615.26 ≈ 0.8364.
func TestMintInflationRateMainnet(t *testing.T) {
	got, err := types.MintInflationRate(mainnetSupply, twoYears, mainnetBlocksYear).Float64()
	require.NoError(t, err)
	s, c := toFloat(mainnetSupply), toFloat(types.SupplyCap)
	want := math.Ln2 / 2 * (c - s) / s
	require.InDelta(t, want, got, 1e-12)
	require.InDelta(t, 0.8364, got, 1e-4)

	// One block mints λ/blocksPerYear of the remaining tokens: about 88.98 UP.
	minted := toFloat(mintBlock(mainnetSupply, twoYears, mainnetBlocksYear)) / 1e18
	require.InDelta(t, math.Ln2/2*(c-s)/1e18/mainnetBlocksYear, minted, 1e-6)
	require.InDelta(t, 88.98, minted, 0.01)
}

// Running x/mint's arithmetic block by block for one half-life must halve
// the tokens remaining below the cap, and supply must never pass the cap.
func TestCurveHalvesRemainingEachHalfLife(t *testing.T) {
	const blocksPerYear = 10_000 // small so the test runs quickly
	supply := mainnetSupply
	start := types.SupplyCap.Sub(supply)

	for half := 1; half <= 3; half++ {
		for b := 0; b < 2*blocksPerYear; b++ { // one 2-year half-life
			next := supply.Add(mintBlock(supply, twoYears, blocksPerYear))
			require.True(t, next.GTE(supply), "supply decreased")
			require.True(t, next.LTE(types.SupplyCap), "supply passed cap")
			supply = next
		}
		remaining := toFloat(types.SupplyCap.Sub(supply))
		want := toFloat(start) / math.Pow(2, float64(half))
		// Stepping per block instead of continuously adds a relative error
		// of about λ/(2*blocksPerYear) per year.
		require.InEpsilon(t, want, remaining, 1e-4, "after %d half-lives", half)
	}
}

// Validate rejects half-lives under 1 year, but the rate function must still
// be safe for any input. The shortest possible half-life (1 second) must not
// push supply past the cap: λ is limited so one block mints at most all of the remaining tokens.
func TestNeverPassesCapWithExtremeHalfLife(t *testing.T) {
	for _, perYear := range []uint64{1, 2, 7, mainnetBlocksYear} {
		for _, remaining := range []int64{1, 7, 1_000, 1_000_000_000_000_000_000} {
			supply := types.SupplyCap.SubRaw(remaining)
			for b := 0; b < 5; b++ {
				supply = supply.Add(mintBlock(supply, 1, perYear))
				require.True(t, supply.LTE(types.SupplyCap),
					"blocks/year %d, remaining %d: supply %s passed cap", perYear, remaining, supply)
			}
		}
	}
}

// Near the cap, rounding down must leave supply at or below the cap.
func TestRoundingNearCap(t *testing.T) {
	for _, remaining := range []int64{1, 2, 10, 1_000_000} {
		supply := types.SupplyCap.SubRaw(remaining)
		for b := 0; b < 1_000; b++ {
			supply = supply.Add(mintBlock(supply, twoYears, mainnetBlocksYear))
		}
		require.True(t, supply.LTE(types.SupplyCap))
	}
}
