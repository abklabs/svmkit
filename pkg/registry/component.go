package registry

import (
	"github.com/abklabs/svmkit/pkg/runner"

	"github.com/abklabs/svmkit/pkg/agave"
	"github.com/abklabs/svmkit/pkg/firedancer"
	"github.com/abklabs/svmkit/pkg/firewall"
	"github.com/abklabs/svmkit/pkg/machine"
	"github.com/abklabs/svmkit/pkg/solana"
	"github.com/abklabs/svmkit/pkg/solana/explorer"
	"github.com/abklabs/svmkit/pkg/solana/faucet"
	"github.com/abklabs/svmkit/pkg/solana/genesis"
	"github.com/abklabs/svmkit/pkg/solana/watchtower"
	"github.com/abklabs/svmkit/pkg/tuner"
)

type Component int

const (
	ComponentAgave Component = iota
	ComponentFiredancer
	ComponentExplorer
	ComponentFaucet
	ComponentGenesis
	ComponentMachine
	ComponentStakeAccount
	ComponentTransfer
	ComponentTuner
	ComponentFirewall
	ComponentVoteAccount
	ComponentWatchtower
)

func (a Component) String() string {
	switch a {
	case ComponentAgave:
		return "agave"
	case ComponentFiredancer:
		return "firedancer"
	case ComponentExplorer:
		return "explorer"
	case ComponentFaucet:
		return "faucet"
	case ComponentGenesis:
		return "genesis"
	case ComponentMachine:
		return "machine"
	case ComponentStakeAccount:
		return "stakeAccount"
	case ComponentTransfer:
		return "transfer"
	case ComponentTuner:
		return "tuner"
	case ComponentFirewall:
		return "firewall"
	case ComponentVoteAccount:
		return "voteAccount"
	case ComponentWatchtower:
		return "watchtower"

	default:
		return "invalid"
	}
}

type ComponentOp struct {
	Action  Action
	Creator func() runner.Command
}

type ComponentCommand struct {
	Name    Component
	Summary string
	Op      []*ComponentOp
}

var Components = []ComponentCommand{
	{
		ComponentAgave,
		"Install and configure an Agave-variant validator.",
		[]*ComponentOp{
			{
				ActionCreate,
				func() runner.Command {
					return &agave.InstallCommand{}
				},
			},
		},
	},
	{
		ComponentFiredancer,
		"Install and configure an Firedancer-variant validator.",
		[]*ComponentOp{
			{
				ActionCreate,
				func() runner.Command {
					return &firedancer.InstallCommand{}
				},
			},
			{
				ActionDelete,
				func() runner.Command {
					return &firedancer.UninstallCommand{}
				},
			},
		},
	},
	{
		ComponentExplorer,
		"Install and configure the Solana Explorer.",
		[]*ComponentOp{
			{
				ActionCreate,
				func() runner.Command {
					return (&explorer.Explorer{}).Install()
				},
			},
		},
	},
	{
		ComponentFaucet,
		"Install and configure the Solana Faucet.",
		[]*ComponentOp{
			{
				ActionCreate,
				func() runner.Command {
					return (&faucet.Faucet{}).Install()
				},
			},
		},
	},
	{
		ComponentGenesis,
		"Instantiate a new ledger and accounts for a Solana network.",
		[]*ComponentOp{
			{
				ActionCreate,
				func() runner.Command {
					return &genesis.CreateCommand{}
				},
			},
		},
	},
	{
		ComponentMachine,
		"Configure a machine to be used by SVMKit.",
		[]*ComponentOp{
			{
				ActionCreate,
				func() runner.Command {
					return &machine.CreateCommand{}
				},
			},
		},
	},
	{
		ComponentStakeAccount,
		"Create a new Solana stake account.",
		[]*ComponentOp{
			{
				ActionCreate,
				func() runner.Command {
					return &solana.StakeAccountCreate{}
				},
			},
		},
	},
	{
		ComponentTransfer,
		"Perform a SOL transfer on a Solana network.",
		[]*ComponentOp{
			{
				ActionCreate,
				func() runner.Command {
					return &solana.StakeAccountCreate{}
				},
			},
		},
	},
	{
		ComponentTuner,
		"Configure the SVMKit systems tuner.",
		[]*ComponentOp{
			{
				ActionCreate,
				func() runner.Command {
					return (&tuner.Tuner{}).Create()
				},
			},
		},
	},
	{
		ComponentFirewall,
		"Configure the SVMKit firewall.",
		[]*ComponentOp{
			{
				ActionCreate,
				func() runner.Command {
					return (&firewall.Firewall{}).Create()
				},
			},
		},
	},
	{
		ComponentWatchtower,
		"Deploy the Solana Watchtower.",
		[]*ComponentOp{
			{
				ActionCreate,
				func() runner.Command {
					return &watchtower.InstallCommand{}
				},
			},
		},
	},
	{
		ComponentVoteAccount,
		"Create and manage a Solana vote account.",
		[]*ComponentOp{
			{
				ActionCreate,
				func() runner.Command {
					return &solana.VoteAccountCreate{}
				},
			},
			{
				ActionDelete,
				func() runner.Command {
					return &solana.VoteAccountDelete{}
				},
			},
		},
	},
}
