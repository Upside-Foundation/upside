package app

import (
	cmdcfg "github.com/0xUp5ide/upside-foundation/cmd/upsided/config"
	ethcfg "github.com/evmos/ethermint/cmd/config"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func SetConfig() {
	config := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(config)
	ethcfg.SetBip44CoinType(config)
	// Make sure address is compatible with ethereum
	config.SetAddressVerifier(VerifyAddressFormat)
	config.Seal()
}
