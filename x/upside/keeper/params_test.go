package keeper_test

import (
	"errors"

	upsidemodulekeeper "github.com/Upside-Foundation/upside/x/upside/keeper"
	keepertest "github.com/Upside-Foundation/upside/x/upside/keeper/mock"
	"github.com/Upside-Foundation/upside/x/upside/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (suite *KeeperTestSuite) TestGetSourceChannelID() {
	testCases := []struct {
		name          string
		ibcDenom      string
		expectedError error
		postCheck     func(channelID string)
	}{
		{
			"wrong ibc denom",
			"test",
			errors.New("test is invalid: ibc upside denom is invalid"),
			func(channelID string) {},
		},
		{
			"correct ibc denom",
			types.IbcUpsideDenomDefaultValue,
			nil,
			func(channelID string) {
				suite.Require().Equal(channelID, "channel-0")
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			// Create Upside Keeper with mock transfer keeper
			upsideKeeper := *upsidemodulekeeper.NewKeeper(
				suite.app.EncodingConfig().Codec,
				suite.app.GetKey(types.StoreKey),
				suite.app.GetKey(types.MemStoreKey),
				suite.app.BankKeeper,
				keepertest.IbcKeeperMock{},
				suite.app.EvmKeeper,
				suite.app.AccountKeeper,
				authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			)
			suite.app.UpsideKeeper = upsideKeeper

			channelID, err := suite.app.UpsideKeeper.GetSourceChannelID(suite.ctx, tc.ibcDenom)
			if tc.expectedError != nil {
				suite.Require().EqualError(err, tc.expectedError.Error())
			} else {
				suite.Require().NoError(err)
				tc.postCheck(channelID)
			}
		})
	}
}
