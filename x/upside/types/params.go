package types

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	// KeyIbcUpsideDenom is store's key for the IBC Upside denomination
	KeyIbcUpsideDenom = []byte("ibcUpsideDenom")
	// KeyIbcTimeout is store's key for the IBC Timeout
	KeyIbcTimeout = []byte("IbcTimeout")
	// KeyUpsideAdmin is store's key for the admin address
	KeyUpsideAdmin = []byte("UpsideAdmin")
	// KeyEnableAutoDeployment is store's key for the EnableAutoDeployment
	KeyEnableAutoDeployment = []byte("EnableAutoDeployment")
	// KeyMaxCallbackGas is store's key for the MaxCallbackGas
	KeyMaxCallbackGas = []byte("MaxCallbackGas")
)

const (
	IbcUpsideDenomDefaultValue    = "ibc/6B5A664BF0AF4F71B2F0BAA33141E2F1321242FBD5D19762F541EC971ACB0865"
	IbcTimeoutDefaultValue     = uint64(86400000000000) // 1 day
	MaxCallbackGasDefaultValue = uint64(50000)
	// MintHalfLifeSecondsDefaultValue is 2 Julian years.
	MintHalfLifeSecondsDefaultValue = uint64(2 * SecondsPerYear)
	// MinMintHalfLifeSeconds is 1 Julian year. With a shorter half-life one
	// governance proposal could mint most of the tokens below the cap within
	// days; with 1 second, the next block would mint all of them.
	MinMintHalfLifeSeconds = uint64(SecondsPerYear)
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new parameter configuration for the upside module
func NewParams(ibcUpsideDenom string, ibcTimeout uint64, upsideAdmin string, enableAutoDeployment bool, maxCallbackGas, mintHalfLifeSeconds uint64) Params {
	return Params{
		IbcUpsideDenom:          ibcUpsideDenom,
		IbcTimeout:           ibcTimeout,
		UpsideAdmin:          upsideAdmin,
		EnableAutoDeployment: enableAutoDeployment,
		MaxCallbackGas:       maxCallbackGas,
		MintHalfLifeSeconds:  mintHalfLifeSeconds,
	}
}

// DefaultParams is the default parameter configuration for the upside module
func DefaultParams() Params {
	return Params{
		IbcUpsideDenom:          IbcUpsideDenomDefaultValue,
		IbcTimeout:           IbcTimeoutDefaultValue,
		UpsideAdmin:          "",
		EnableAutoDeployment: false,
		MaxCallbackGas:       MaxCallbackGasDefaultValue,
		MintHalfLifeSeconds:  MintHalfLifeSecondsDefaultValue,
	}
}

// Validate all upside module parameters
func (p Params) Validate() error {
	if err := validateIsUint64(p.IbcTimeout); err != nil {
		return err
	}
	if err := validateIsIbcDenom(p.IbcUpsideDenom); err != nil {
		return err
	}
	if len(p.UpsideAdmin) > 0 {
		if _, err := sdk.AccAddressFromBech32(p.UpsideAdmin); err != nil {
			return err
		}
	}
	if err := validateIsUint64(p.MaxCallbackGas); err != nil {
		return err
	}
	if p.MintHalfLifeSeconds < MinMintHalfLifeSeconds {
		return fmt.Errorf("mint half-life must be at least %d seconds (1 year): %d", MinMintHalfLifeSeconds, p.MintHalfLifeSeconds)
	}
	return nil
}

// String implements the fmt.Stringer interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyIbcUpsideDenom, &p.IbcUpsideDenom, validateIsIbcDenom),
		paramtypes.NewParamSetPair(KeyIbcTimeout, &p.IbcTimeout, validateIsUint64),
		paramtypes.NewParamSetPair(KeyUpsideAdmin, &p.UpsideAdmin, validateIsAddress),
		paramtypes.NewParamSetPair(KeyEnableAutoDeployment, &p.EnableAutoDeployment, validateIsBool),
		paramtypes.NewParamSetPair(KeyMaxCallbackGas, &p.MaxCallbackGas, validateIsUint64),
	}
}

func validateIsIbcDenom(i interface{}) error {
	s, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !IsValidIBCDenom(s) {
		return fmt.Errorf("invalid ibc denom: %T", i)
	}
	return nil
}

func validateIsUint64(i interface{}) error {
	if _, ok := i.(uint64); !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateIsAddress(i interface{}) error {
	s, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if len(s) > 0 {
		if _, err := sdk.AccAddressFromBech32(s); err != nil {
			return err
		}
	}
	return nil
}

func validateIsBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
