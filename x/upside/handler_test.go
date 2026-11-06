package upside_test

import (
	"errors"
	"testing"
	"time"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/0xUp5ide/upside-foundation/app"
	"github.com/0xUp5ide/upside-foundation/x/upside/keeper"
	"github.com/0xUp5ide/upside-foundation/x/upside/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type UpsideTestSuite struct {
	suite.Suite

	ctx     sdk.Context
	app     *app.App
	address sdk.AccAddress
}

func TestUpsideTestSuite(t *testing.T) {
	suite.Run(t, new(UpsideTestSuite))
}

func (suite *UpsideTestSuite) SetupTest() {
	checkTx := false
	privKey, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	suite.address = sdk.AccAddress(privKey.PubKey().Address())
	suite.app = app.Setup(suite.T(), suite.address.String())
	suite.ctx = suite.app.NewContext(checkTx).WithBlockHeader(tmproto.Header{Height: 1, ChainID: app.TestAppChainID, Time: time.Now().UTC()})
}

func (suite *UpsideTestSuite) TestMsgConvertVouchers() {
	testCases := []struct {
		name          string
		msg           *types.MsgConvertVouchers
		malleate      func()
		expectedError error
	}{
		{
			"Wrong address",
			types.NewMsgConvertVouchers("test", sdk.NewCoins(sdk.NewCoin("aup", sdkmath.NewInt(1)))),
			func() {},
			errors.New("decoding bech32 failed: invalid bech32 string length 4"),
		},
		{
			"Empty address",
			types.NewMsgConvertVouchers("", sdk.NewCoins(sdk.NewCoin("aup", sdkmath.NewInt(1)))),
			func() {},
			errors.New("empty address string is not allowed"),
		},
		{
			"Correct address with non supported coin denom",
			types.NewMsgConvertVouchers(suite.address.String(), sdk.NewCoins(sdk.NewCoin("fake", sdkmath.NewInt(1)))),
			func() {},
			errors.New("coin fake is not supported for conversion"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msgSrv := keeper.NewMsgServerImpl(suite.app.UpsideKeeper)
			_, err := msgSrv.ConvertVouchers(suite.ctx, tc.msg)
			if tc.expectedError != nil {
				suite.Require().EqualError(err, tc.expectedError.Error())
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *UpsideTestSuite) TestMsgTransferTokens() {
	testCases := []struct {
		name          string
		msg           *types.MsgTransferTokens
		malleate      func()
		expectedError error
	}{
		{
			"Wrong from address",
			types.NewMsgTransferTokens("test", "to", sdk.NewCoins(sdk.NewCoin("aup", sdkmath.NewInt(1)))),
			func() {},
			errors.New("decoding bech32 failed: invalid bech32 string length 4"),
		},
		{
			"Empty from address",
			types.NewMsgTransferTokens("", "to", sdk.NewCoins(sdk.NewCoin("aup", sdkmath.NewInt(1)))),
			func() {},
			errors.New("empty address string is not allowed"),
		},
		{
			"Empty to address",
			types.NewMsgTransferTokens(suite.address.String(), "", sdk.NewCoins(sdk.NewCoin("aup", sdkmath.NewInt(1)))),
			func() {},
			errors.New("to address cannot be empty"),
		},
		{
			"Correct address with non supported coin denom",
			types.NewMsgTransferTokens(suite.address.String(), "to", sdk.NewCoins(sdk.NewCoin("fake", sdkmath.NewInt(1)))),
			func() {},
			errors.New("the coin fake is neither an ibc voucher or a upside token"),
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msgSrv := keeper.NewMsgServerImpl(suite.app.UpsideKeeper)
			_, err := msgSrv.TransferTokens(suite.ctx, tc.msg)
			if tc.expectedError != nil {
				suite.Require().EqualError(err, tc.expectedError.Error())
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *UpsideTestSuite) TestUpdateTokenMapping() {
	suite.SetupTest()

	denom := "gravity0x6E7eef2b30585B2A4D45Ba9312015d5354FDB067"
	contract := "0x57f96e6B86CdeFdB3d412547816a82E3E0EbF9D2"

	msg := types.NewMsgUpdateTokenMapping(suite.address.String(), denom, contract, "", 0)
	err := suite.app.UpsideKeeper.RegisterOrUpdateTokenMapping(suite.ctx, msg)
	suite.Require().NoError(err)

	contractAddr, found := suite.app.UpsideKeeper.GetContractByDenom(suite.ctx, denom)
	suite.Require().True(found)
	suite.Require().Equal(contract, contractAddr.Hex())
}
