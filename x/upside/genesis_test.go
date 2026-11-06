package upside_test

import (
	"github.com/0xUp5ide/upside-foundation/x/upside"
	"github.com/0xUp5ide/upside-foundation/x/upside/types"
)

func (suite *UpsideTestSuite) TestInitGenesis() {
	testCases := []struct {
		name     string
		malleate func()
		genState *types.GenesisState
		expPanic bool
	}{
		{
			"default",
			func() {},
			types.DefaultGenesis(),
			false,
		},
		{
			"Wrong ibcUpsideDenom length",
			func() {},
			&types.GenesisState{
				Params: types.Params{
					IbcUpsideDenom: "ibc/6B5A664BF0AF4F71B2F0BAA33141E2F1321242FBD5D19762F541EC971ACB086534",
				},
			},
			true,
		},
		{
			"Wrong ibcUpsideDenom prefix",
			func() {},
			&types.GenesisState{
				Params: types.Params{
					IbcUpsideDenom: "aaa/6B5A664BF0AF4F71B2F0BAA33141E2F1321242FBD5D19762F541EC971ACB0865",
				},
			},
			true,
		},
		{
			"Wrong denom in external token mapping",
			func() {},
			&types.GenesisState{
				ExternalContracts: []types.TokenMapping{
					{
						Denom:    "aaa/6B5A664BF0AF4F71B2F0BAA33141E2F1321242FBD5D19762F541EC971ACB0865",
						Contract: "0x0000000000000000000000000000000000000000",
					},
				},
			},
			true,
		},
		{
			"Wrong denom in auto token mapping",
			func() {},
			&types.GenesisState{
				AutoContracts: []types.TokenMapping{
					{
						Denom:    "aaa/6B5A664BF0AF4F71B2F0BAA33141E2F1321242FBD5D19762F541EC971ACB0865",
						Contract: "0x0000000000000000000000000000000000000000",
					},
				},
			},
			true,
		},
		{
			"Wrong contract in external token mapping",
			func() {},
			&types.GenesisState{
				ExternalContracts: []types.TokenMapping{
					{
						Denom:    "ibc/6B5A664BF0AF4F71B2F0BAA33141E2F1321242FBD5D19762F541EC971ACB0865",
						Contract: "0x00000000000000000000000000000000000000",
					},
				},
			},
			true,
		},
		{
			"Wrong contract in auto token mapping",
			func() {},
			&types.GenesisState{
				AutoContracts: []types.TokenMapping{
					{
						Denom:    "ibc/6B5A664BF0AF4F71B2F0BAA33141E2F1321242FBD5D19762F541EC971ACB0865",
						Contract: "0x00000000000000000000000000000000000000",
					},
				},
			},
			true,
		},
		{
			"Correct token mapping",
			func() {},
			&types.GenesisState{
				Params: types.DefaultParams(),
				ExternalContracts: []types.TokenMapping{
					{
						Denom:    "ibc/6B5A664BF0AF4F71B2F0BAA33141E2F1321242FBD5D19762F541EC971ACB0865",
						Contract: "0x0000000000000000000000000000000000000000",
					},
				},
				AutoContracts: []types.TokenMapping{
					{
						Denom:    "gravity0x0000000000000000000000000000000000000000",
						Contract: "0x0000000000000000000000000000000000000000",
					},
				},
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.malleate()
			if tc.expPanic {
				suite.Require().Panics(
					func() {
						upside.InitGenesis(suite.ctx, suite.app.UpsideKeeper, *tc.genState)
					},
				)
			} else {
				suite.Require().NotPanics(
					func() {
						upside.InitGenesis(suite.ctx, suite.app.UpsideKeeper, *tc.genState)
					},
				)
			}
		})
	}
}

func (suite *UpsideTestSuite) TestExportGenesis() {
	genesisState := upside.ExportGenesis(suite.ctx, suite.app.UpsideKeeper)
	suite.Require().Equal(genesisState.Params.IbcUpsideDenom, types.DefaultParams().IbcUpsideDenom)
}
