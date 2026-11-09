package v2

import (
	"github.com/Upside-Foundation/upside/x/upside/exported"
	"github.com/Upside-Foundation/upside/x/upside/types"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migrate migrates the x/upside module state from the consensus version 1 to
// version 2. Specifically, it takes the parameters that are currently stored
// and managed by the x/params modules and stores them directly into the x/upside
// module state.
func Migrate(ctx sdk.Context, store storetypes.KVStore, legacySubspace exported.Subspace, cdc codec.BinaryCodec) error {
	var currParams types.Params
	legacySubspace.GetParamSetIfExists(ctx, &currParams)

	// The legacy subspace has no half-life key, so the value read is always
	// zero, which Validate rejects. Use the default.
	if currParams.MintHalfLifeSeconds == 0 {
		currParams.MintHalfLifeSeconds = types.MintHalfLifeSecondsDefaultValue
	}
	if err := currParams.Validate(); err != nil {
		return err
	}
	if currParams.GetMaxCallbackGas() == 0 {
		currParams.MaxCallbackGas = types.MaxCallbackGasDefaultValue
	}
	bz := cdc.MustMarshal(&currParams)
	store.Set(types.ParamsKey, bz)

	return nil
}
