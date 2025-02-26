package solana

import (
	"errors"
	"fmt"
	"slices"

	"github.com/abklabs/svmkit/pkg/runner"
)

type StakeAccountKeyPairs struct {
	StakeAccount      string  `pulumi:"stakeAccount" provider:"secret"`
	VoteAccount       *string `pulumi:"voteAccount, optional" provider:"secret"`
	StakeAuthority    *string `pulumi:"stakeAuthority,optional" provider:"secret"`
	WithdrawAuthority *string `pulumi:"withdrawAuthority,optional" provider:"secret"`
}

type StakeAccountLockup struct {
	EpochAvailable  uint64 `pulumi:"epochAvailable"`
	CustodianPubkey string `pulumi:"custodianPubkey"`
}

type StakeAccountArgs struct {
	StakeAccountKeyPairs StakeAccountKeyPairs `pulumi:"keyPairs"`
	Amount               float64              `pulumi:"amount"`
	WithdrawAddress      *string              `pulumi:"withdrawAddress,optional"`
	TransactionOptions   *TxnOptions          `pulumi:"transactionOptions"`
	ForceDelete          bool                 `pulumi:"forceDelete"`
	LockupArgs           *StakeAccountLockup  `pulumi:"lockupArgs"`
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

func addKeyPairsToPayload(p *runner.Payload, keys StakeAccountKeyPairs) {
	p.AddString("stake_account.json", keys.StakeAccount)
	if keys.VoteAccount != nil {
		p.AddString("vote_account.json", *keys.VoteAccount)
	}
	if keys.StakeAuthority != nil {
		p.AddString("stake_authority.json", *keys.StakeAuthority)
	}
	if keys.WithdrawAuthority != nil {
		p.AddString("withdraw_authority.json", *keys.WithdrawAuthority)
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
	if v.StakeAccountKeyPairs.VoteAccount == nil {
		return errors.New("cannot create stake account without specifying delegate address")
	}
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
	addKeyPairsToPayload(p, v.StakeAccountKeyPairs)
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
	addKeyPairsToPayload(p, v.StakeAccountKeyPairs)
	return nil
}

// ------------------------------------------------------------
// StakeAccount Update Command
// ------------------------------------------------------------
type UpdateType int

const (
	UpdateTypeAuthority UpdateType = iota
	UpdateTypeDelegate
	UpdateTypeDeactivate
	UpdateTypeLock
)

func (v *StakeAccountUpdate) UpdatePlan() []UpdateType {
	var updates []UpdateType

	oldKps := v.state.StakeAccountKeyPairs
	newKps := v.newArgs.StakeAccountKeyPairs

	// Authority Diff
	if (oldKps.WithdrawAuthority != newKps.WithdrawAuthority) || (&oldKps.StakeAuthority != &newKps.StakeAuthority) {
		updates = append(updates, UpdateTypeAuthority)
	}

	// Lockup Args Diff
	if v.state.LockupArgs != v.newArgs.LockupArgs {
		updates = append(updates, UpdateTypeLock)
	}

	// Vote Account diff
	if oldKps.VoteAccount != newKps.VoteAccount {
		if oldKps.VoteAccount != nil && newKps.VoteAccount == nil {
			updates = append(updates, UpdateTypeDeactivate)
		} else {
			updates = append(updates, UpdateTypeDelegate)
		}
	}
	return updates
}

func (v *StakeAccountUpdate) Check() error {
	if v.state.StakeAccountKeyPairs.StakeAccount != v.newArgs.StakeAccountKeyPairs.StakeAccount {
		return errors.New("stake account address can not be rotated")
	}

	if v.newArgs.Amount != v.state.Amount {
		// This will trigger a split in the future
		return errors.New("cannot change stake amount; operation not currently supported")
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

	e.SetBool("STAKE_ACCOUNT_DEACTIVATE", false)
	e.SetBool("STAKE_ACCOUNT_DELEGATE", false)
	e.SetBool("STAKE_ACCOUNT_AUTHORITY", false)
	e.SetBool("STAKE_ACCOUNT_LOCKUP", false)
	e.SetBool("STAKE_AUTHORITY_UPDATE", false)
	e.SetBool("WITHDRAW_AUTHORITY_UPDATE", false)

	updates := v.UpdatePlan()
	if slices.Contains(updates, UpdateTypeDeactivate) {
		e.SetBool("STAKE_ACCOUNT_DEACTIVATE", true)
	}
	if slices.Contains(updates, UpdateTypeDelegate) {
		e.SetBool("STAKE_ACCOUNT_DELEGATE", true)
	}
	if slices.Contains(updates, UpdateTypeAuthority) {
		e.SetBool("STAKE_ACCOUNT_AUTHORITY", true)

		oldStakeAuth := v.state.StakeAccountKeyPairs.StakeAuthority
		newStakeAuth := v.newArgs.StakeAccountKeyPairs.StakeAuthority
		if (oldStakeAuth == nil) != (newStakeAuth == nil) || (oldStakeAuth != nil && newStakeAuth != nil && *oldStakeAuth != *newStakeAuth) {
			e.SetBool("STAKE_AUTHORITY_UPDATE", true)
			if oldStakeAuth != nil {
				e.SetBool("STAKE_AUTHORITY", true)
			}
		}

		oldWithdrawAuth := v.state.StakeAccountKeyPairs.WithdrawAuthority
		newWithdrawAuth := v.newArgs.StakeAccountKeyPairs.WithdrawAuthority
		if (oldWithdrawAuth == nil) != (newWithdrawAuth == nil) || (oldWithdrawAuth != nil && newWithdrawAuth != nil && *oldWithdrawAuth != *newWithdrawAuth) {
			e.SetBool("WITHDRAW_AUTHORITY_UPDATE", true)
			if oldWithdrawAuth != nil {
				e.SetBool("WITHDRAW_AUTHORITY", true)
			}
		}
	}
	if slices.Contains(updates, UpdateTypeLock) {
		e.SetBool("STAKE_ACCOUNT_LOCKUP", true)
		if v.newArgs.LockupArgs != nil {
			e.Set("EPOCH_AVAILABLE", fmt.Sprintf("%d", v.newArgs.LockupArgs.EpochAvailable))
			e.Set("CUSTODIAN_PUBKEY", v.newArgs.LockupArgs.CustodianPubkey)
		}
	}
	return e
}

func (v *StakeAccountUpdate) AddToPayload(p *runner.Payload) error {
	if err := setupPayload(p, v.newArgs.TransactionOptions); err != nil {
		return err
	}

	// Add the stake account keypair
	p.AddString("stake_account.json", v.state.StakeAccountKeyPairs.StakeAccount)

	updates := v.UpdatePlan()

	// Add old authority keypairs if needed for authorization
	if slices.Contains(updates, UpdateTypeAuthority) {
		oldStakeAuth := v.state.StakeAccountKeyPairs.StakeAuthority
		newStakeAuth := v.newArgs.StakeAccountKeyPairs.StakeAuthority
		if (oldStakeAuth == nil) != (newStakeAuth == nil) || (oldStakeAuth != nil && newStakeAuth != nil && *oldStakeAuth != *newStakeAuth) {
			if oldStakeAuth != nil {
				p.AddString("stake_authority.json", *oldStakeAuth)
			}
			if newStakeAuth != nil {
				p.AddString("new_stake_authority.json", *newStakeAuth)
			}
		}

		oldWithdrawAuth := v.state.StakeAccountKeyPairs.WithdrawAuthority
		newWithdrawAuth := v.newArgs.StakeAccountKeyPairs.WithdrawAuthority
		if (oldWithdrawAuth == nil) != (newWithdrawAuth == nil) || (oldWithdrawAuth != nil && newWithdrawAuth != nil && *oldWithdrawAuth != *newWithdrawAuth) {
			if oldWithdrawAuth != nil {
				p.AddString("withdraw_authority.json", *oldWithdrawAuth)
			}
			if newWithdrawAuth != nil {
				p.AddString("new_withdraw_authority.json", *newWithdrawAuth)
			}
		}
	}

	// For delegation updates, add the vote account
	if slices.Contains(updates, UpdateTypeDelegate) && v.newArgs.StakeAccountKeyPairs.VoteAccount != nil {
		p.AddString("new_vote_account.json", *v.newArgs.StakeAccountKeyPairs.VoteAccount)
	}

	return nil
}
