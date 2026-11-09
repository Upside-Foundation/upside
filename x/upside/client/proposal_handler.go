package client

import (
	"github.com/Upside-Foundation/upside/x/upside/client/cli"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

// ProposalHandler is the token mapping change proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.NewSubmitTokenMappingChangeProposalTxCmd)
