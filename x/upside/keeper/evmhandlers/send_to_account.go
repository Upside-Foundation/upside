package evmhandler

import (
	"fmt"
	"math/big"

	upsidekeeper "github.com/0xUp5ide/upside-foundation/x/upside/keeper"
	"github.com/0xUp5ide/upside-foundation/x/upside/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.EvmLogHandler = SendToAccountHandler{}

const SendToAccountEventName = "__UpsideSendToAccount"

// SendToAccountEvent represent the signature of
// `event __UpsideSendToAccount(address recipient, uint256 amount)`
var SendToAccountEvent abi.Event

func init() {
	addressType, _ := abi.NewType("address", "", nil)
	uint256Type, _ := abi.NewType("uint256", "", nil)

	SendToAccountEvent = abi.NewEvent(
		SendToAccountEventName,
		SendToAccountEventName,
		false,
		abi.Arguments{abi.Argument{
			Name:    "recipient",
			Type:    addressType,
			Indexed: false,
		}, abi.Argument{
			Name:    "amount",
			Type:    uint256Type,
			Indexed: false,
		}},
	)
}

// SendToAccountHandler handles `__UpsideSendToAccount` log
type SendToAccountHandler struct {
	bankKeeper   types.BankKeeper
	upsideKeeper upsidekeeper.Keeper
}

func NewSendToAccountHandler(bankKeeper types.BankKeeper, upsideKeeper upsidekeeper.Keeper) *SendToAccountHandler {
	return &SendToAccountHandler{
		bankKeeper:   bankKeeper,
		upsideKeeper: upsideKeeper,
	}
}

func (h SendToAccountHandler) EventID() common.Hash {
	return SendToAccountEvent.ID
}

func (h SendToAccountHandler) Handle(
	ctx sdk.Context,
	contract common.Address,
	topics []common.Hash,
	data []byte,
	_ func(contractAddress common.Address, logSig common.Hash, logData []byte),
) error {
	unpacked, err := SendToAccountEvent.Inputs.Unpack(data)
	if err != nil {
		// log and ignore
		h.upsideKeeper.Logger(ctx).Error("log signature matches but failed to decode", "error", err)
		return nil
	}

	denom, found := h.upsideKeeper.GetDenomByContract(ctx, contract)
	if !found {
		return fmt.Errorf("contract %s is not connected to native token", contract)
	}

	contractAddr := sdk.AccAddress(contract.Bytes())
	recipient := sdk.AccAddress(unpacked[0].(common.Address).Bytes())
	coins := sdk.NewCoins(sdk.NewCoin(denom, sdkmath.NewIntFromBigInt(unpacked[1].(*big.Int))))
	err = h.bankKeeper.SendCoins(ctx, contractAddr, recipient, coins)
	if err != nil {
		return err
	}

	return nil
}
