package solana

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/abklabs/svmkit/pkg/runner"
	"golang.org/x/crypto/ssh"
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

type StakeAccount struct {
	StakeAccountKeyPairs StakeAccountKeyPairs `pulumi:"keyPairs"`
	Amount               float64              `pulumi:"amount"`
	WithdrawAddress      *string              `pulumi:"withdrawAddress,optional"`
	TransactionOptions   *TxnOptions          `pulumi:"transactionOptions"`
	ForceDelete          bool                 `pulumi:"forceDelete"`
	LockupArgs           *StakeAccountLockup  `pulumi:"lockupArgs"`
}

// ------------------------------------------------------------
// StakeAccount Client
// ------------------------------------------------------------

type StakeAccountClient struct {
	operator StakeOperator
}

func (c *StakeAccountClient) Create(args StakeAccount) (StakeAccount, error) {
	err := c.operator.Create(args)
	if err != nil {
		return StakeAccount{}, err
	}
	return args, nil
}

func (c *StakeAccountClient) Read(args StakeAccount) (StakeAccount, error) {
	_, err := c.operator.GetStatus(args)

	if err != nil {
		return StakeAccount{}, err
	}

	//TODO: Take the status and convert it to a StakeAccount
	return StakeAccount{}, fmt.Errorf("not implemented")
}

// func (c *StakeAccountClient) Update(state StakeAccount, newArgs StakeAccount) (*StakeAccount, error) {
//   // TODO: Do a read with the string handler on the client first. Then do logic
// 	return nil, nil
// }

// func (c *StakeAccountClient) Delete(oldState StakeAccount) error {
// 	//TODO:
//    Run a read and then only withdraw all if it's in a state to do so
// }

// ------------------------------------------------------------
// StakeOperator Interface
// ------------------------------------------------------------

// StakeOperator defines interface for stake account operations
// It's agnostic to connection type and query method (RPC or bash commands)
type StakeOperator interface {
	GetStatus(stakeAccount StakeAccount) (CliStakeState, error)
	Create(account StakeAccount) error
}

// ------------------------------------------------------------
// StakeOperator CLI Implementation
// ------------------------------------------------------------

type CliStakeOperator struct {
	client *ssh.Client
	handle runner.DeployerHandler
	ctx    context.Context
}

func (op *CliStakeOperator) runCommand(cmd runner.Command, handler runner.DeployerHandler) error {
	if err := cmd.Check(); err != nil {
		return fmt.Errorf("failed to check component config: %w", err)
	}

	r := runner.NewRunner(op.client, cmd)
	if err := r.Run(op.ctx, handler); err != nil {
		return err
	}

	return nil
}

func (op *CliStakeOperator) Create(stakeAccount StakeAccount) error {
	cmd := &StakeAccountCreate{StakeAccount: stakeAccount}

	err := op.runCommand(cmd, op.handle)

	if err != nil {
		return err
	}
	return nil
}

func (op *CliStakeOperator) GetStatus(stakeAccount StakeAccount) (CliStakeState, error) {
	// This is the lone operation that doesn't use the provided handler and
	// uses a StringHandler so that the output can be parsed
	handler := &StringHandler{}

	cmd := &StakeAccountRead{StakeAccount: stakeAccount}

	err := op.runCommand(cmd, handler)
	if err != nil {
		return CliStakeState{}, err
	}

	output := handler.Output

	status, err := parseOutput(output)
	if err != nil {
		return CliStakeState{}, err
	}

	return status, nil

}

type CliAuthorized struct {
	Staker     string `json:"staker"`
	Withdrawer string `json:"withdrawer"`
}

type CliLockup struct {
	UnixTimestamp int64  `json:"unixTimestamp"` // Assuming UnixTimestamp is an int64
	Epoch         uint64 `json:"epoch"`         // Assuming Epoch is a uint64
	Custodian     string `json:"custodian"`
}
type CliStakeState struct {
	StakeType                   string         `json:"stakeType"`
	AccountBalance              uint64         `json:"accountBalance"`
	DelegatedStake              *uint64        `json:"delegatedStake,omitempty"`
	DelegatedVoteAccountAddress *string        `json:"delegatedVoteAccountAddress,omitempty"`
	ActivationEpoch             *uint64        `json:"activationEpoch,omitempty"`
	DeactivationEpoch           *uint64        `json:"deactivationEpoch,omitempty"`
	Authorized                  *CliAuthorized `json:",omitempty"`
	Lockup                      *CliLockup     `json:",omitempty"`
	ActiveStake                 *uint64        `json:"activeStake,omitempty"`
	ActivatingStake             *uint64        `json:"activatingStake,omitempty"`
	DeactivatingStake           *uint64        `json:"deactivatingStake,omitempty"`
}

func parseOutput(output string) (CliStakeState, error) {
	output = strings.TrimSpace(output)
	var state CliStakeState
	err := json.Unmarshal([]byte(output), &state)
	if err != nil {
		return CliStakeState{}, err
	}
	return state, nil
}

// ------------------------------------------------------------
// Common Helper Functions
// ------------------------------------------------------------
func env(newArgs StakeAccount) *runner.EnvBuilder {
	// Sets default env for all stake Commands
	b := runner.NewEnvBuilder()

	b.SetFloat64("STAKE_AMOUNT", newArgs.Amount)

	if opt := newArgs.TransactionOptions; opt != nil {
		cli := CLITxnOptions{*opt}
		b.SetArray("SOLANA_CLI_TXN_FLAGS", cli.Flags().Args())
	}

	return b
}

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
// Create
// ------------------------------------------------------------

type StakeAccountCreate struct {
	StakeAccount
}

func (v *StakeAccountCreate) Check() error {
	if v.StakeAccountKeyPairs.VoteAccount == nil {
		return errors.New("cannot create stake account without specifying delegate address")
	}
	if v.StakeAccount.WithdrawAddress != nil {
		return errors.New("cannot withdraw on create")
	}

	return nil
}

func (v *StakeAccountCreate) Env() *runner.EnvBuilder {
	e := env(v.StakeAccount)
	e.Set("STAKE_ACCOUNT_ACTION", "CREATE")
	if v.StakeAccountKeyPairs.StakeAuthority != nil {
		e.SetBool("STAKE_AUTHORITY", true)
	}
	if v.StakeAccountKeyPairs.WithdrawAuthority != nil {
		e.SetBool("WITHDRAW_AUTHORITY", true)
	}
	return e
}

func (v *StakeAccountCreate) AddToPayload(p *runner.Payload) error {
	// Add stake account script
	script, err := assets.Open(assetsStakeAccountScript)
	if err != nil {
		return fmt.Errorf("failed to open stake account script: %w", err)
	}
	p.AddReader("steps.sh", script)

	// Add keypair files
	p.AddString("stake_account.json", v.StakeAccountKeyPairs.StakeAccount)
	if v.StakeAccountKeyPairs.VoteAccount != nil {
		p.AddString("vote_account.json", *v.StakeAccountKeyPairs.VoteAccount)
	}
	if v.StakeAccountKeyPairs.StakeAuthority != nil {
		p.AddString("stake_authority.json", *v.StakeAccountKeyPairs.StakeAuthority)
	}
	if v.StakeAccountKeyPairs.WithdrawAuthority != nil {
		p.AddString("withdraw_authority.json", *v.StakeAccountKeyPairs.WithdrawAuthority)
	}

	// Add transaction options if present
	if v.TransactionOptions != nil {
		cli := CLITxnOptions{*v.TransactionOptions}
		if err := cli.AddToPayload(p); err != nil {
			return fmt.Errorf("failed to add transaction options: %w", err)
		}
	}

	return nil
}

// ------------------------------------------------------------
// Read
// ------------------------------------------------------------

type StakeAccountRead struct {
	StakeAccount
}

func (v *StakeAccountRead) Check() error {
	return nil
}

func (v *StakeAccountRead) Env() *runner.EnvBuilder {
	e := env(v.StakeAccount)
	e.Set("STAKE_ACCOUNT_ACTION", "READ")
	return e
}

func (v *StakeAccountRead) AddToPayload(p *runner.Payload) error {
	setupPayload(p, v.TransactionOptions)
	p.AddString("stake_account.json", v.StakeAccountKeyPairs.StakeAccount)
	return nil
}

// type StakeAccountUpdate struct {
//   old StakeAccount
//   new StakeAccount
// }

// type StakeAccountDelete struct {
// 	StakeAccount
// }

// // ------------------------------------------------------------
// // StakeAccount Delete Command
// // ------------------------------------------------------------

// func (v *StakeAccountDelete) Check() error {
// 	if v.WithdrawAddress == nil && v.ForceDelete {
// 		return errors.New("must provide withdraw address or set force delete to true")
// 	}

// 	if v.WithdrawAddress != nil && v.ForceDelete {
// 		// You must not have a withdraw address set if you forcibly delete
// 		return errors.New("cannot provide withdraw address and set force delete to true")
// 	}

// 	if v.WithdrawAddress != nil && v.StakeState != StakeStateUnstaked {
// 		return errors.New("stake not fully deactivated, cannot delete")
// 	}

// 	return nil
// }

// func (v *StakeAccountDelete) Env() *runner.EnvBuilder {
// 	e := env(v.StakeAccount)
// 	e.Set("STAKE_ACCOUNT_ACTION", "DELETE")

// 	if v.StakeAccountKeyPairs.WithdrawAuthority != nil {
// 		e.SetBool("ADD_WITHDRAW_AUTHORITY", true)
// 	}

// 	if v.WithdrawAddress != nil {
// 		e.Set("WITHDRAW_ADDRESS", *v.WithdrawAddress)
// 	}

// 	if v.ForceDelete {
// 		e.SetBool("FORCE_DELETE", true)
// 	}

// 	return e
// }

// func (v *StakeAccountDelete) AddToPayload(p *runner.Payload) error {
// 	if err := setupPayload(p, v.TransactionOptions); err != nil {
// 		return err
// 	}
// 	addKeyPairsToPayload(p, v.StakeAccountKeyPairs)
// 	return nil
// }

// ------------------------------------------------------------
// StakeAccount Update Command
// ------------------------------------------------------------
// type UpdateType int

// const (
// 	UpdateTypeAuthority UpdateType = iota
// 	UpdateTypeDelegate
// 	UpdateTypeDeactivate
// 	UpdateTypeLock
// )

// func (v *StakeAccountUpdate) UpdatePlan() []UpdateType {
// 	var updates []UpdateType

// 	oldKps := v.state.StakeAccountKeyPairs
// 	newKps := v.newArgs.StakeAccountKeyPairs

// 	// Authority Diff
// 	if (oldKps.WithdrawAuthority != newKps.WithdrawAuthority) || (&oldKps.StakeAuthority != &newKps.StakeAuthority) {
// 		updates = append(updates, UpdateTypeAuthority)
// 	}

// 	// Lockup Args Diff
// 	if v.state.LockupArgs != v.newArgs.LockupArgs {
// 		updates = append(updates, UpdateTypeLock)
// 	}

// 	// Vote Account diff
// 	if oldKps.VoteAccount != newKps.VoteAccount {
// 		if oldKps.VoteAccount != nil && newKps.VoteAccount == nil {
// 			updates = append(updates, UpdateTypeDeactivate)
// 		} else {
// 			updates = append(updates, UpdateTypeDelegate)
// 		}
// 	}
// 	return updates
// }

// func (v *StakeAccountUpdate) Check() error {
// 	if v.state.StakeAccountKeyPairs.StakeAccount != v.newArgs.StakeAccountKeyPairs.StakeAccount {
// 		return errors.New("stake account address can not be rotated")
// 	}

// 	if v.newArgs.Amount != v.state.Amount {
// 		// This will trigger a split in the future
// 		return errors.New("cannot change stake amount; operation not currently supported")
// 	}

// 	if v.state.WithdrawAddress == nil && v.newArgs.WithdrawAddress != nil {
// 		if v.state.StakeState != StakeStateUnstaked {
// 			return errors.New("cannot initiate a withdraw unless fully staked")
// 		}
// 	}

// 	return nil
// }

// func (v *StakeAccountUpdate) Env() *runner.EnvBuilder {
// 	e := env(v.newArgs)
// 	e.Set("STAKE_ACCOUNT_ACTION", "UPDATE")

// 	e.SetBool("STAKE_ACCOUNT_DEACTIVATE", false)
// 	e.SetBool("STAKE_ACCOUNT_DELEGATE", false)
// 	e.SetBool("STAKE_ACCOUNT_AUTHORITY", false)
// 	e.SetBool("STAKE_ACCOUNT_LOCKUP", false)
// 	e.SetBool("STAKE_AUTHORITY_UPDATE", false)
// 	e.SetBool("WITHDRAW_AUTHORITY_UPDATE", false)

// 	updates := v.UpdatePlan()
// 	if slices.Contains(updates, UpdateTypeDeactivate) {
// 		e.SetBool("STAKE_ACCOUNT_DEACTIVATE", true)
// 	}
// 	if slices.Contains(updates, UpdateTypeDelegate) {
// 		e.SetBool("STAKE_ACCOUNT_DELEGATE", true)
// 	}
// 	if slices.Contains(updates, UpdateTypeAuthority) {
// 		e.SetBool("STAKE_ACCOUNT_AUTHORITY", true)

// 		oldStakeAuth := v.state.StakeAccountKeyPairs.StakeAuthority
// 		newStakeAuth := v.newArgs.StakeAccountKeyPairs.StakeAuthority
// 		if (oldStakeAuth == nil) != (newStakeAuth == nil) || (oldStakeAuth != nil && newStakeAuth != nil && *oldStakeAuth != *newStakeAuth) {
// 			e.SetBool("STAKE_AUTHORITY_UPDATE", true)
// 			if oldStakeAuth != nil {
// 				e.SetBool("STAKE_AUTHORITY", true)
// 			}
// 		}

// 		oldWithdrawAuth := v.state.StakeAccountKeyPairs.WithdrawAuthority
// 		newWithdrawAuth := v.newArgs.StakeAccountKeyPairs.WithdrawAuthority
// 		if (oldWithdrawAuth == nil) != (newWithdrawAuth == nil) || (oldWithdrawAuth != nil && newWithdrawAuth != nil && *oldWithdrawAuth != *newWithdrawAuth) {
// 			e.SetBool("WITHDRAW_AUTHORITY_UPDATE", true)
// 			if oldWithdrawAuth != nil {
// 				e.SetBool("WITHDRAW_AUTHORITY", true)
// 			}
// 		}
// 	}
// 	if slices.Contains(updates, UpdateTypeLock) {
// 		e.SetBool("STAKE_ACCOUNT_LOCKUP", true)
// 		if v.newArgs.LockupArgs != nil {
// 			e.Set("EPOCH_AVAILABLE", fmt.Sprintf("%d", v.newArgs.LockupArgs.EpochAvailable))
// 			e.Set("CUSTODIAN_PUBKEY", v.newArgs.LockupArgs.CustodianPubkey)
// 		}
// 	}
// 	return e
// }

// func (v *StakeAccountUpdate) AddToPayload(p *runner.Payload) error {
// 	if err := setupPayload(p, v.newArgs.TransactionOptions); err != nil {
// 		return err
// 	}

// 	// Add the stake account keypair
// 	p.AddString("stake_account.json", v.state.StakeAccountKeyPairs.StakeAccount)

// 	updates := v.UpdatePlan()

// 	// Add old authority keypairs if needed for authorization
// 	if slices.Contains(updates, UpdateTypeAuthority) {
// 		oldStakeAuth := v.state.StakeAccountKeyPairs.StakeAuthority
// 		newStakeAuth := v.newArgs.StakeAccountKeyPairs.StakeAuthority
// 		if (oldStakeAuth == nil) != (newStakeAuth == nil) || (oldStakeAuth != nil && newStakeAuth != nil && *oldStakeAuth != *newStakeAuth) {
// 			if oldStakeAuth != nil {
// 				p.AddString("stake_authority.json", *oldStakeAuth)
// 			}
// 			if newStakeAuth != nil {
// 				p.AddString("new_stake_authority.json", *newStakeAuth)
// 			}
// 		}

// 		oldWithdrawAuth := v.state.StakeAccountKeyPairs.WithdrawAuthority
// 		newWithdrawAuth := v.newArgs.StakeAccountKeyPairs.WithdrawAuthority
// 		if (oldWithdrawAuth == nil) != (newWithdrawAuth == nil) || (oldWithdrawAuth != nil && newWithdrawAuth != nil && *oldWithdrawAuth != *newWithdrawAuth) {
// 			if oldWithdrawAuth != nil {
// 				p.AddString("withdraw_authority.json", *oldWithdrawAuth)
// 			}
// 			if newWithdrawAuth != nil {
// 				p.AddString("new_withdraw_authority.json", *newWithdrawAuth)
// 			}
// 		}
// 	}

// 	// For delegation updates, add the vote account
// 	if slices.Contains(updates, UpdateTypeDelegate) && v.newArgs.StakeAccountKeyPairs.VoteAccount != nil {
// 		p.AddString("new_vote_account.json", *v.newArgs.StakeAccountKeyPairs.VoteAccount)
// 	}

// 	return nil
// }

// ------------------------------------------------------------
// String Handler
// ------------------------------------------------------------

type StringHandler struct {
	Output string
}

func (h *StringHandler) IngestReaders(done chan<- struct{}, stdout io.Reader, stderr io.Reader) error {
	var buf strings.Builder

	// Copy stdout to our buffer
	if _, err := io.Copy(&buf, stdout); err != nil {
		return err
	}

	// Store the output
	h.Output = buf.String()

	// Signal completion
	close(done)
	return nil
}

func (h *StringHandler) AugmentError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("command failed: %w\nOutput: %s", err, h.Output)
}
