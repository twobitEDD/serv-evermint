package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/EscanBE/evermint/v12/x/incentives/client/cli"
)

var (
	RegisterIncentiveProposalHandler = govclient.NewProposalHandler(cli.NewRegisterIncentiveProposalCmd)
	CancelIncentiveProposalHandler   = govclient.NewProposalHandler(cli.NewCancelIncentiveProposalCmd)
)
