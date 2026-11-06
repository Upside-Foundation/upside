package types

import (
	sdkmath "cosmossdk.io/math"
)

// Supply curve for the native token (aup, 18 decimals).
//
// The whitepaper curve is S(t) = S_cap * (1 - e^(-λt)) with λ = ln 2 / T_half.
// Its slope is dS/dt = λ * (S_cap - S): at every moment, a fixed fraction λ
// of the tokens still below the cap is minted. That form needs no start
// time, so the chain can join the curve at whatever its supply is today.
//
// x/mint mints inflation * supply per year, spread evenly over BlocksPerYear
// blocks. Setting
//
//	inflation = λ * (S_cap - S) / S
//
// makes each block mint λ/BlocksPerYear of the remaining tokens, which is the
// curve above taken in steps of one block. λ is limited to BlocksPerYear, so
// one block mints at most all of the remaining tokens, and every division
// rounds down, so supply never passes the cap through rounding.

var (
	// SupplyCap is 2,100,000,000 UP in aup.
	SupplyCap = sdkmath.NewIntWithDecimal(2_100_000_000, 18)

	// ln2 to 18 decimal places, rounded down.
	ln2 = sdkmath.LegacyMustNewDecFromStr("0.693147180559945309")
)

// MinBlocksPerYear and MaxBlocksPerYear bound x/mint's blocks_per_year param.
// x/mint spreads a year's minting over blocks_per_year blocks, so the real
// half-life is mint_half_life_seconds scaled by blocks_per_year over the real
// number of blocks per year. x/mint accepts any value from 1, and at 1 a
// single block mints about a third of the tokens below the cap. The bounds
// are half and double mainnet's measured 5,783,000 (blocks of 10.9 s and
// 2.7 s); outside them the inflation function returns zero.
const (
	MinBlocksPerYear = 2_891_500
	MaxBlocksPerYear = 11_566_000
)

// SecondsPerYear is a Julian year (365.25 days), the same year length x/mint's
// default BlocksPerYear uses.
const SecondsPerYear = 31_557_600

// MintInflationRate returns the annual inflation rate that keeps supply on the
// capped half-life curve. It returns zero when supply is at or above the cap,
// when supply is not positive, or when halfLifeSeconds or blocksPerYear is
// zero.
func MintInflationRate(supply sdkmath.Int, halfLifeSeconds, blocksPerYear uint64) sdkmath.LegacyDec {
	if halfLifeSeconds == 0 || blocksPerYear == 0 || !supply.IsPositive() || supply.GTE(SupplyCap) {
		return sdkmath.LegacyZeroDec()
	}
	remaining := SupplyCap.Sub(supply)

	// λ in units of 1/year: ln2 * SecondsPerYear / halfLifeSeconds.
	lambda := ln2.MulInt64(SecondsPerYear).QuoInt(sdkmath.NewIntFromUint64(halfLifeSeconds))

	// One block mints λ/blocksPerYear of the remaining tokens; keep that
	// fraction at or below 1.
	maxLambda := sdkmath.LegacyNewDecFromInt(sdkmath.NewIntFromUint64(blocksPerYear))
	if lambda.GT(maxLambda) {
		lambda = maxLambda
	}

	// λ * remaining / supply, rounded down.
	return lambda.MulInt(remaining).QuoInt(supply)
}
