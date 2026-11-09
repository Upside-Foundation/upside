package keeper_test

import (
	upsidemodulekeeper "github.com/Upside-Foundation/upside/x/upside/keeper"
	"github.com/Upside-Foundation/upside/x/upside/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (suite *KeeperTestSuite) TestUpdateParams() {
	testCases := []struct {
		name      string
		req       *types.MsgUpdateParams
		expectErr bool
		expErrMsg string
	}{
		{
			name: "gov module account address as valid authority",
			req: &types.MsgUpdateParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params: types.Params{
					IbcUpsideDenom:       types.IbcUpsideDenomDefaultValue,
					IbcTimeout:           10,
					UpsideAdmin:          sdk.AccAddress(suite.address.Bytes()).String(),
					EnableAutoDeployment: true,
					MintHalfLifeSeconds:  types.MintHalfLifeSecondsDefaultValue,
				},
			},
			expectErr: false,
			expErrMsg: "",
		},
		{
			name: "set zero mint half-life",
			req: &types.MsgUpdateParams{
				Authority: suite.app.UpsideKeeper.GetAuthority(),
				Params: types.Params{
					IbcUpsideDenom:       types.IbcUpsideDenomDefaultValue,
					IbcTimeout:           10,
					UpsideAdmin:          sdk.AccAddress(suite.address.Bytes()).String(),
					EnableAutoDeployment: true,
				},
			},
			expectErr: true,
			expErrMsg: "mint half-life must be at least 31557600 seconds",
		},
		{
			name: "set mint half-life just under 1 year",
			req: &types.MsgUpdateParams{
				Authority: suite.app.UpsideKeeper.GetAuthority(),
				Params: types.Params{
					IbcUpsideDenom:       types.IbcUpsideDenomDefaultValue,
					IbcTimeout:           10,
					UpsideAdmin:          sdk.AccAddress(suite.address.Bytes()).String(),
					EnableAutoDeployment: true,
					MintHalfLifeSeconds:  types.MinMintHalfLifeSeconds - 1,
				},
			},
			expectErr: true,
			expErrMsg: "mint half-life must be at least 31557600 seconds",
		},
		{
			name: "set invalid authority",
			req: &types.MsgUpdateParams{
				Authority: "foo",
			},
			expectErr: true,
			expErrMsg: "invalid authority",
		},
		{
			name: "set invalid ibc up denomination",
			req: &types.MsgUpdateParams{
				Authority: suite.app.UpsideKeeper.GetAuthority(),
				Params: types.Params{
					IbcUpsideDenom:       "foo",
					IbcTimeout:           10,
					UpsideAdmin:          sdk.AccAddress(suite.address.Bytes()).String(),
					EnableAutoDeployment: true,
				},
			},
			expectErr: true,
			expErrMsg: "invalid ibc denom",
		},
		{
			name: "set invalid upside admin address",
			req: &types.MsgUpdateParams{
				Authority: suite.app.UpsideKeeper.GetAuthority(),
				Params: types.Params{
					IbcUpsideDenom:       types.IbcUpsideDenomDefaultValue,
					IbcTimeout:           10,
					UpsideAdmin:          "foo",
					EnableAutoDeployment: true,
				},
			},
			expectErr: true,
			expErrMsg: "invalid bech32 string",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msgServer := upsidemodulekeeper.NewMsgServerImpl(suite.app.UpsideKeeper)
			_, err := msgServer.UpdateParams(suite.ctx, tc.req)
			if tc.expectErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
