package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

const (
	codeErrIbcUpsideDenomEmpty = uint32(iota) + 2 // NOTE: code 1 is reserved for internal errors
	codeErrIbcUpsideDenomInvalid
)

// x/upside module sentinel errors
var (
	ErrIbcUpsideDenomEmpty   = errors.Register(ModuleName, codeErrIbcUpsideDenomEmpty, "ibc upside denom is not set")
	ErrIbcUpsideDenomInvalid = errors.Register(ModuleName, codeErrIbcUpsideDenomInvalid, "ibc upside denom is invalid")
	// this line is used by starport scaffolding # ibc/errors
)
