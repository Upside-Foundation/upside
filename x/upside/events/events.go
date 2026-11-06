package events

import (
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	ica "github.com/0xUp5ide/upside-foundation/x/upside/events/bindings/cosmos/precompile/ica"
	relayer "github.com/0xUp5ide/upside-foundation/x/upside/events/bindings/cosmos/precompile/relayer"
	upsideeventstypes "github.com/0xUp5ide/upside-foundation/x/upside/events/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

var (
	RelayerEvents        map[string]*EventDescriptor
	IcaEvents            map[string]*EventDescriptor
	RelayerValueDecoders = ValueDecoders{
		channeltypes.AttributeKeyDataHex:             ConvertPacketData,
		sdk.AttributeKeyAmount:                       ConvertAmount,
		banktypes.AttributeKeyRecipient:              ConvertAccAddressFromBech32,
		banktypes.AttributeKeySpender:                ConvertAccAddressFromBech32,
		banktypes.AttributeKeyReceiver:               ConvertAccAddressFromBech32,
		banktypes.AttributeKeySender:                 ConvertAccAddressFromBech32,
		banktypes.AttributeKeyMinter:                 ConvertAccAddressFromBech32,
		banktypes.AttributeKeyBurner:                 ConvertAccAddressFromBech32,
		channeltypes.AttributeKeySequence:            ConvertUint64,
		channeltypes.AttributeKeySrcPort:             ReturnStringAsIs,
		upsideeventstypes.AttributeKeySrcPortInfo:    ReturnStringAsIs,
		channeltypes.AttributeKeySrcChannel:          ReturnStringAsIs,
		upsideeventstypes.AttributeKeySrcChannelInfo: ReturnStringAsIs,
		channeltypes.AttributeKeyDstPort:             ReturnStringAsIs,
		channeltypes.AttributeKeyDstChannel:          ReturnStringAsIs,
		channeltypes.AttributeKeyConnectionID:        ReturnStringAsIs,
		transfertypes.AttributeKeyDenom:              ReturnStringAsIs,
		transfertypes.AttributeKeyRefundReceiver:     ConvertAccAddressFromBech32,
		transfertypes.AttributeKeyRefundTokens:       ReturnStringAsIs,
	}
	IcaValueDecoders = ValueDecoders{
		upsideeventstypes.AttributeKeySeq:   ConvertUint64,
		channeltypes.AttributeKeySrcChannel: ReturnStringAsIs,
	}
)

func init() {
	var relayerABI abi.ABI
	if err := relayerABI.UnmarshalJSON([]byte(relayer.RelayerModuleMetaData.ABI)); err != nil {
		panic(err)
	}
	RelayerEvents = NewEventDescriptors(relayerABI)

	var icaABI abi.ABI
	if err := icaABI.UnmarshalJSON([]byte(ica.ICAModuleMetaData.ABI)); err != nil {
		panic(err)
	}
	IcaEvents = NewEventDescriptors(icaABI)
}

func RelayerConvertEvent(event sdk.Event) (*ethtypes.Log, error) {
	desc, ok := RelayerEvents[event.Type]
	if !ok {
		return nil, nil
	}
	replaceAttrs := map[string]string{
		upsideeventstypes.AttributeKeySrcPortInfo:    channeltypes.AttributeKeySrcPort,
		upsideeventstypes.AttributeKeySrcChannelInfo: channeltypes.AttributeKeySrcChannel,
	}
	return desc.ConvertEvent(event.Attributes, RelayerValueDecoders, replaceAttrs)
}

func IcaConvertEvent(event sdk.Event) (*ethtypes.Log, error) {
	desc, ok := IcaEvents[event.Type]
	if !ok {
		return nil, nil
	}
	return desc.ConvertEvent(event.Attributes, IcaValueDecoders, map[string]string{})
}
