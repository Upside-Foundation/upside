package simulation

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/Upside-Foundation/upside/x/upside/types"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

const (
	ibcUpsideDenomKey       = "ibc_upside_denom"
	ibcTimeoutKey           = "ibc_timeout"
	upsideAdminKey          = "upside_admin"
	enableAutoDeploymentKey = "enable_auto_deployment"
	maxCallbackGasKey       = "max_callback_gas"
)

func GenIbcUpsideDenom(r *rand.Rand) string {
	randDenom := make([]byte, 32)
	r.Read(randDenom)
	return fmt.Sprintf("ibc/%s", hex.EncodeToString(randDenom))
}

func GenIbcTimeout(r *rand.Rand) uint64 {
	timeout := r.Uint64()
	return timeout
}

func GenUpsideAdmin(r *rand.Rand, simState *module.SimulationState) string {
	adminAccount, _ := simtypes.RandomAcc(r, simState.Accounts)
	return adminAccount.Address.String()
}

func GenEnableAutoDeployment(r *rand.Rand) bool {
	return r.Intn(2) > 0
}

func GenMaxCallbackGas(r *rand.Rand) uint64 {
	maxCallbackGas := r.Uint64()
	return maxCallbackGas
}

// RandomizedGenState generates a random GenesisState for the upside module
func RandomizedGenState(simState *module.SimulationState) {
	// upside params
	var (
		ibcUpsideDenom       string
		ibcTimeout           uint64
		upsideAdmin          string
		enableAutoDeployment bool
		maxCallbackGas       uint64
	)

	simState.AppParams.GetOrGenerate(
		ibcUpsideDenomKey, &ibcUpsideDenom, simState.Rand,
		func(r *rand.Rand) { ibcUpsideDenom = GenIbcUpsideDenom(r) },
	)

	simState.AppParams.GetOrGenerate(
		ibcTimeoutKey, &ibcTimeout, simState.Rand,
		func(r *rand.Rand) { ibcTimeout = GenIbcTimeout(r) },
	)

	simState.AppParams.GetOrGenerate(
		upsideAdminKey, &upsideAdmin, simState.Rand,
		func(r *rand.Rand) { upsideAdmin = GenUpsideAdmin(r, simState) },
	)

	simState.AppParams.GetOrGenerate(
		enableAutoDeploymentKey, &enableAutoDeployment, simState.Rand,
		func(r *rand.Rand) { enableAutoDeployment = GenEnableAutoDeployment(r) },
	)

	simState.AppParams.GetOrGenerate(
		maxCallbackGasKey, &ibcTimeout, simState.Rand,
		func(r *rand.Rand) { maxCallbackGas = GenIbcTimeout(r) },
	)

	params := types.NewParams(ibcUpsideDenom, ibcTimeout, upsideAdmin, enableAutoDeployment, maxCallbackGas, types.MintHalfLifeSecondsDefaultValue)
	upsideGenesis := &types.GenesisState{
		Params:            params,
		ExternalContracts: nil,
		AutoContracts:     nil,
	}

	bz, err := json.MarshalIndent(upsideGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", types.ModuleName, bz)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(upsideGenesis)
}
