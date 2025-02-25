package solana

import (
	"errors"

	"github.com/abklabs/svmkit/pkg/runner"
)

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
	StakeState       `pulumi:"stakeState"`
	StakeAccountArgs `pulumi:"stakeAccountArgs"`
}

// ------------------------------------------------------------
// StakeAccount Command Types
// ------------------------------------------------------------

type StakeAccountCreate struct {
	StakeAccountArgs
}

type StakeAccountUpdate struct {
	newArgs StakeAccountArgs
	state   StakeAccountState
}

type StakeAccountDelete struct {
	StakeAccountState
}

// ------------------------------------------------------------
// StakeAccount Client
// ------------------------------------------------------------
type StakeAccountClient struct{}

func (v *StakeAccountClient) Create(args StakeAccountArgs) runner.Command {
	return &StakeAccountCreate{
		args,
	}
}

func (v *StakeAccountClient) Update(state StakeAccountState, newArgs StakeAccountArgs) runner.Command {
	return &StakeAccountUpdate{
		newArgs,
		state,
	}
}

func (v *StakeAccountClient) Delete(oldState StakeAccountState) runner.Command {
	return &StakeAccountDelete{
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
// Common Helper Functions
// ------------------------------------------------------------

func setupPayload(p *runner.Payload, opt *TxnOptions) error {
	stakeAccountScript, err := assets.Open(assetsStakeAccountScript)
	if err != nil {
		return err
	}
	p.AddReader("steps.sh", stakeAccountScript)

	if opt != nil {
		cli := CLITxnOptions{*opt}
		if err := cli.AddToPayload(p); err != nil {
			return err
		}
	}
	return nil
}

func addKeyPairsToPayload(p *runner.Payload, pairs StakeAccountKeyPairs, includeVoteAccount bool) {
	p.AddString("stake_account.json", pairs.StakeAccount)
	if includeVoteAccount {
		p.AddString("vote_account.json", pairs.VoteAccount)
	}
	if pairs.StakeAuthority != nil {
		p.AddString("stake_authority.json", *pairs.StakeAuthority)
	}
	if pairs.WithdrawAuthority != nil {
		p.AddString("withdraw_authority.json", *pairs.WithdrawAuthority)
	}
}

func setupEnvWithAuthorities(e *runner.EnvBuilder, pairs StakeAccountKeyPairs) {
	if pairs.StakeAuthority != nil {
		e.SetBool("STAKE_AUTHORITY", true)
	}
	if pairs.WithdrawAuthority != nil {
		e.SetBool("WITHDRAW_AUTHORITY", true)
	}
}

// ------------------------------------------------------------
// StakeAccount Create Command
// ------------------------------------------------------------

func (v *StakeAccountCreate) Check() error {
	if v.StakeAccountArgs.WithdrawAddress != nil {
		return errors.New("cannot withdraw on create")
	}
	return nil
}

func (v *StakeAccountCreate) Env() *runner.EnvBuilder {
	e := env(v.StakeAccountArgs)
	e.Set("STAKE_ACCOUNT_ACTION", "CREATE")
	setupEnvWithAuthorities(e, v.StakeAccountKeyPairs)
	return e
}

func (v *StakeAccountCreate) AddToPayload(p *runner.Payload) error {
	if err := setupPayload(p, v.TransactionOptions); err != nil {
		return err
	}
	addKeyPairsToPayload(p, v.StakeAccountKeyPairs, true)
	return nil
}

// ------------------------------------------------------------
// StakeAccount Delete Command
// ------------------------------------------------------------

func (v *StakeAccountDelete) Check() error {
	if v.WithdrawAddress == nil && !v.ForceDelete {
		return errors.New("must provide withdraw address or set force delete to true")
	}

	if v.WithdrawAddress != nil && v.ForceDelete {
		// You must not have a withdraw address set if you forcibly delete
		return errors.New("cannot provide withdraw address and set force delete to true")
	}

	if v.WithdrawAddress != nil && v.StakeState != StakeStateUnstaked {
		return errors.New("stake not fully deactivated, cannot delete")
	}

	return nil
}

func (v *StakeAccountDelete) Env() *runner.EnvBuilder {
	e := env(v.StakeAccountArgs)
	e.Set("STAKE_ACCOUNT_ACTION", "DELETE")

	if v.StakeAccountKeyPairs.WithdrawAuthority != nil {
		e.SetBool("ADD_WITHDRAW_AUTHORITY", true)
	}

	if v.WithdrawAddress != nil {
		e.Set("WITHDRAW_ADDRESS", *v.WithdrawAddress)
	}

	if v.ForceDelete {
		e.SetBool("FORCE_DELETE", true)
	}

	return e
}

func (v *StakeAccountDelete) AddToPayload(p *runner.Payload) error {
	if err := setupPayload(p, v.TransactionOptions); err != nil {
		return err
	}
	addKeyPairsToPayload(p, v.StakeAccountKeyPairs, false)
	return nil
}

// ------------------------------------------------------------
// StakeAccount Update Command
// ------------------------------------------------------------

func (v *StakeAccountUpdate) Check() error {
	if v.state.StakeAccountKeyPairs != v.newArgs.StakeAccountKeyPairs {
		// In future this is how redelgation will be signaled
		return errors.New("rotation of any keypair is not supported")
	}

	if v.newArgs.Amount != v.state.Amount {
		// This will trigger a split in the future
		return errors.New("cannot change stake amount; operation not currently supported")
	}

	if v.state.StakeState == StakeStateCooldown || v.state.StakeState == StakeStateWarmup {
		return errors.New("cannot update while in warmup/cooldown")
	}

	if v.state.WithdrawAddress == nil && v.newArgs.WithdrawAddress != nil {
		if v.state.StakeState != StakeStateUnstaked {
			return errors.New("cannot initiate a withdraw unless fully staked")
		}
	}

	return nil
}

func (v *StakeAccountUpdate) Env() *runner.EnvBuilder {
	e := env(v.newArgs)
	e.Set("STAKE_ACCOUNT_ACTION", "UPDATE")

	if v.state.WithdrawAddress == nil && v.newArgs.WithdrawAddress != nil {
		e.Set("STAKE_ACCOUNT_UPDATE_ACTION", "DEACTIVATE")
	}

	setupEnvWithAuthorities(e, v.state.StakeAccountKeyPairs)
	return e
}

func (v *StakeAccountUpdate) AddToPayload(p *runner.Payload) error {
	if err := setupPayload(p, v.newArgs.TransactionOptions); err != nil {
		return err
	}
	addKeyPairsToPayload(p, v.state.StakeAccountKeyPairs, true)
	return nil
}
