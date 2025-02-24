package solana

import (
	"errors"

	"github.com/abklabs/svmkit/pkg/runner"
)

type WithdrawArgs struct {
	Destination string `pulumi:"destination"`
}

type StakeAccountKeyPairs struct {
	StakeAccount      string  `pulumi:"stakeAccount" provider:"secret"`
	VoteAccount       string  `pulumi:"voteAccount" provider:"secret"`
	StakeAuthority    *string `pulumi:"stakeAuthority,optional"`
	WithdrawAuthority *string `pulumi:"withdrawAuthority,optional"`
}

type StakeAccountArgs struct {
	TransactionOptions   *TxnOptions          `pulumi:"transactionOptions"`
	StakeAccountKeyPairs StakeAccountKeyPairs `pulumi:"keyPairs"`
	Amount               float64              `pulumi:"amount"`
	WithdrawAddress      *string              `pulumi:"withdrawAddress,optional"`
	ForceDelete          bool                 `pulumi:"forceDelete"`
}

type StakeState int

const (
	StakeStateWarmup StakeState = iota
	StakeStateStaked
	StakeStateCooldown
	StakeStateUnstaked
)

type StakeAccountState struct {
	StakeState       StakeState       `pulumi:"stakeState"`
	StakeAccountArgs StakeAccountArgs `pulumi:"stakeAccountArgs"`
}

type StakeAccountClient struct{}

type StakeAccountCreate struct {
	StakeAccountArgs
}

type StakeAccountUpdate struct {
	StakeAccountArgs
	StakeAccountState
}

type StakeAccountDelete struct {
	StakeAccountArgs
	StakeAccountState
}

func (v *StakeAccountClient) Create(args StakeAccountArgs) runner.Command {
	return &StakeAccountCreate{
		args,
	}
}

func (v *StakeAccountClient) Update(oldState StakeAccountState, newArgs StakeAccountArgs) runner.Command {
	return &StakeAccountUpdate{
		StakeAccountArgs:  newArgs,
		StakeAccountState: oldState,
	}
}

func (v *StakeAccountClient) Delete(oldState StakeAccountState, newArgs StakeAccountArgs) runner.Command {
	return &StakeAccountDelete{
		StakeAccountArgs:  newArgs,
		StakeAccountState: oldState,
	}
}

func env(newArgs StakeAccountArgs) *runner.EnvBuilder {
	// Sets default env for all stake Commands
	b := runner.NewEnvBuilder()

	b.SetFloat64("STAKE_AMOUNT", newArgs.Amount)

	if opt := newArgs.TransactionOptions; opt != nil {
		cli := CLITxnOptions{*opt}
		b.SetArray("SOLANA_CLI_TXN_FLAGS", cli.Flags().Args())
	}

	return b
}

// ------------------------------------------------------------
// StakeAccount Create Command
// ------------------------------------------------------------

func (v *StakeAccountCreate) Check() error {
	return nil
}

func (v *StakeAccountCreate) Env() *runner.EnvBuilder {
	e := env(v.StakeAccountArgs)
	e.Set("STAKE_ACCOUNT_ACTION", "CREATE")

	return e
}

func (v *StakeAccountCreate) AddToPayload(p *runner.Payload) error {
	stakeAccountScript, err := assets.Open(assetsStakeAccountScript)

	if err != nil {
		return err
	}

	p.AddReader("steps.sh", stakeAccountScript)

	// 	p.AddString("stake_account.json", v.StakeAccountKeyPairs.StakeAccount)
	// 	p.AddString("vote_account.json", v.StakeAccountKeyPairs.VoteAccount)

	// 	if opt := v.TransactionOptions; opt != nil {
	// 		cli := CLITxnOptions{*opt}

	// 		err := cli.AddToPayload(p)

	// 		if err != nil {
	// 			return err
	// 		}
	// 	}

	return nil
}

// ------------------------------------------------------------
// StakeAccount Delete Command
// ------------------------------------------------------------

func (v *StakeAccountDelete) Check() error {
	if v.WithdrawAddress == nil && !v.ForceDelete {
		return errors.New("must provide withdraw address or set force delete to true")
	}
	return nil
}

func (v *StakeAccountDelete) Env() *runner.EnvBuilder {
	e := env(v.StakeAccountArgs)
	e.Set("STAKE_ACCOUNT_ACTION", "DELETE")

	// if v.StakeAccount.StakeAccountKeyPairs.WithdrawAuthority != nil {
	// 	e.SetBool("ADD_WITHDRAW_AUTHORITY", true)
	// }

	// if v.ForceDelete {
	// 	e.SetBool("FORCE_DELETE", true)
	// }

	return e
}

func (v *StakeAccountDelete) AddToPayload(p *runner.Payload) error {
	stakeAccountScript, err := assets.Open(assetsStakeAccountScript)

	if err != nil {
		return err
	}

	p.AddReader("steps.sh", stakeAccountScript)

	// p.AddString("stake_account.json", v.StakeAccountKeyPairs.StakeAccount)

	// if v.StakeAccountKeyPairs.WithdrawAuthority != nil {
	// 	p.AddString("withdraw_authority.json", *v.StakeAccountKeyPairs.WithdrawAuthority)
	// }

	// if opt := v.TransactionOptions; opt != nil {
	// 	cli := CLITxnOptions{*opt}

	// 	err := cli.AddToPayload(p)

	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

// ------------------------------------------------------------
// StakeAccount Update Command
// ------------------------------------------------------------

func (v *StakeAccountUpdate) Check() error {
	return nil
}

func (v *StakeAccountUpdate) Env() *runner.EnvBuilder {
	e := env(v.StakeAccountArgs)
	e.Set("STAKE_ACCOUNT_ACTION", "CREATE")

	return e
}

func (v *StakeAccountUpdate) AddToPayload(p *runner.Payload) error {
	stakeAccountScript, err := assets.Open(assetsStakeAccountScript)

	if err != nil {
		return err
	}

	p.AddReader("steps.sh", stakeAccountScript)

	// 	p.AddString("stake_account.json", v.StakeAccountKeyPairs.StakeAccount)
	// 	p.AddString("vote_account.json", v.StakeAccountKeyPairs.VoteAccount)

	// 	if opt := v.TransactionOptions; opt != nil {
	// 		cli := CLITxnOptions{*opt}

	// 		err := cli.AddToPayload(p)

	// 		if err != nil {
	// 			return err
	// 		}
	// 	}

	return nil
}
