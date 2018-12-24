package utils

import (
	"math/big"

	shardparams "github.com/ovcharovvladimir/Prysm/sness/params"
	"github.com/urfave/cli"
)

var (
	// DepositFlag defines whether a node will withdraw ESS from the user's account.
	DepositFlag = cli.BoolFlag{
		Name:  "deposit",
		Usage: "To become a voter in a sharding node, " + new(big.Int).Div(shardparams.DefaultConfig.NotaryDeposit, new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)).String() + " ETH will be deposited into VRC",
	}
	// ModeFlag defines the role of the sharding client. Either proposer, voter, or simulator.
	ModeFlag = cli.StringFlag{
		Name:  "mode",
		Usage: `use the --mode voter or --actor proposer to start a voter or proposer service in the sharding node. If omitted, the sharding node registers an Observer service that simply observes the activity in the sharded network`,
	}
	// ShardIDFlag specifies which shard to listen to.
	ShardIDFlag = cli.IntFlag{
		Name:  "shardid",
		Usage: `use the --shardid to determine which shard to start p2p server, listen for incoming transactions and perform proposer/observer duties`,
	}

	VrcContractFlag = cli.StringFlag{
		Name:  "vrcaddr",
		Usage: "use --vrcaddr to set vrc contract address",
	}

	PubKeyFlag = cli.StringFlag{
		Name:  "pkey",
		Usage: "use --pkey to set voter publik key",
	}

	Web3ProviderFlag = cli.StringFlag{
		Name:  "rpc",
		Usage: "use --rpc to set we3 rpc address",
	}
)
