package solana

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/abklabs/svmkit/pkg/runner"
	"golang.org/x/crypto/ssh"
)

type StakeAccountKeyPairs struct {
	StakeAccount      string  `pulumi:"stakeAccount" provider:"secret"`
	VoteAccount       *string `pulumi:"voteAccount,optional" provider:"secret"`
	StakeAuthority    *string `pulumi:"stakeAuthority,optional" provider:"secret"`
	WithdrawAuthority *string `pulumi:"withdrawAuthority,optional" provider:"secret"`
}

type StakeAccountLockup struct {
	EpochAvailable  int64  `pulumi:"epochAvailable"`
	CustodianPubkey string `pulumi:"custodianPubkey"`
}

type StakeAccount struct {
	StakeAccountKeyPairs StakeAccountKeyPairs `pulumi:"keyPairs"`
	Amount               float64              `pulumi:"amount"`
	WithdrawAddress      *string              `pulumi:"withdrawAddress,optional"`
	TransactionOptions   *TxnOptions          `pulumi:"transactionOptions"`
	ForceDelete          bool                 `pulumi:"forceDelete"`
	LockupArgs           *StakeAccountLockup  `pulumi:"lockupArgs,optional"`
}

type CliAuthorized struct {
	Staker     string `json:"staker"`
	Withdrawer string `json:"withdrawer"`
}

type CliLockup struct {
	UnixTimestamp int64  `json:"unixTimestamp"`
	Epoch         uint64 `json:"epoch"`
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
// StakeOperator Interface
// ------------------------------------------------------------

// StakeOperator defines interface for stake account operations
// It's agnostic to connection type and query method (RPC or bash commands)
type StakeOperator interface {
	// TODO: Only need to pass in the stake account address; but rn we repr that as a keypair
	// and so we currently pass in the whole StakeAccount that contains the keypair. Eventually we
	// want to decouple this implementation detail
	GetStatus(stakeAccount StakeAccount) (CliStakeState, error)

	Create(StakeAccount) error
	Update(oldState StakeAccount, newArgs StakeAccount) error
	Delete(StakeAccount) error
}

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

// ------------------------------------------------------------
// StakeAccount Client
// ------------------------------------------------------------

type StakeAccountClient struct {
	operator StakeOperator
}

func NewStakeAccountClient(operator StakeOperator) *StakeAccountClient {
	return &StakeAccountClient{operator: operator}
}

func (c *StakeAccountClient) Create(args StakeAccount) (StakeAccount, error) {
	if args.StakeAccountKeyPairs.VoteAccount == nil {
		return StakeAccount{}, errors.New("cannot create stake account without specifying delegate address")
	}

	if err := c.operator.Create(args); err != nil {
		return StakeAccount{}, err
	}
	return args, nil
}

func (c *StakeAccountClient) Update(state StakeAccount, newArgs StakeAccount) (StakeAccount, error) {
	readState, err := c.operator.GetStatus(state)
	if err != nil {
		return StakeAccount{}, errors.New("failed to read stake account state from chain")
	}

	if state.StakeAccountKeyPairs.StakeAccount != newArgs.StakeAccountKeyPairs.StakeAccount {
		return StakeAccount{}, errors.New("stake account address can not be rotated")
	}
	if newArgs.Amount != state.Amount {
		return StakeAccount{}, errors.New("cannot change stake amount; operation not currently supported")
	}

	// Handle vote-account change
	currentVA := state.StakeAccountKeyPairs.VoteAccount
	newVA := newArgs.StakeAccountKeyPairs.VoteAccount

	// Compare vote accounts properly by value, not by pointer
	voteAccountChanged := false

	// Different nil status (one is nil, the other isn't)
	if (currentVA == nil) != (newVA == nil) {
		voteAccountChanged = true
	} else if currentVA != nil && newVA != nil {
		// Both non-nil, compare the actual string contents
		if *currentVA != *newVA {
			voteAccountChanged = true
		}
	}

	if voteAccountChanged {
		// Stake must be fully deactivated to redelegate per the stake program
		if currentVA != nil && newVA != nil && (*readState.DelegatedStake != 0 || *readState.DeactivatingStake != 0) {
			return StakeAccount{}, errors.New("cannot redelegate stake until it is fully deactivated")
		}
	}

	if err := c.operator.Update(state, newArgs); err != nil {
		return StakeAccount{}, err
	}
	return newArgs, nil
}

func (c *StakeAccountClient) Delete(state StakeAccount) error {
	readState, err := c.operator.GetStatus(state)
	if err != nil {
		return errors.New("failed to read stake account state from chain")
	}

	if state.WithdrawAddress == nil && !state.ForceDelete {
		return errors.New("must provide withdraw address or set force delete to true")
	}

	if state.WithdrawAddress != nil && state.ForceDelete {
		// You must not have a withdraw address set if you forcibly delete
		return errors.New("cannot provide withdraw address and set force delete to true")
	}

	if state.WithdrawAddress != nil && (*readState.DelegatedStake != 0 || *readState.DeactivatingStake != 0) {
		return errors.New("cannot withdraw stake until it is fully deactivated")
	}

	if err := c.operator.Delete(state); err != nil {
		return err
	}
	return nil
}

// ------------------------------------------------------------
// StakeOperator CLI Implementation
// ------------------------------------------------------------

type CliStakeOperator struct {
	client  *ssh.Client
	handler runner.DeployerHandler
	ctx     context.Context
}

func NewCliStakeOperator(client *ssh.Client, handler runner.DeployerHandler, ctx context.Context) *CliStakeOperator {
	return &CliStakeOperator{
		client:  client,
		handler: handler,
		ctx:     ctx,
	}
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

func (op *CliStakeOperator) Create(args StakeAccount) error {
	cmd := &StakeAccountCreate{StakeAccount: args}
	if err := op.runCommand(cmd, op.handler); err != nil {
		return err
	}
	return nil
}

func (op *CliStakeOperator) Update(oldState StakeAccount, newArgs StakeAccount) error {
	cmd := &StakeAccountUpdate{state: oldState, newArgs: newArgs}
	if err := op.runCommand(cmd, op.handler); err != nil {
		return err
	}
	return nil
}

func (op *CliStakeOperator) Delete(args StakeAccount) error {
	cmd := &StakeAccountDelete{StakeAccount: args}
	if err := op.runCommand(cmd, op.handler); err != nil {
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

// ------------------------------------------------------------
// Common Command Helper Functions
// ------------------------------------------------------------
func env(newArgs StakeAccount) *runner.EnvBuilder {
	// Sets default env for all stake Commands
	b := runner.NewEnvBuilder()

	// Set stake amount - always required
	b.SetFloat64("STAKE_AMOUNT", newArgs.Amount)

	// Set transaction flags if available
	if opt := newArgs.TransactionOptions; opt != nil {
		cli := CLITxnOptions{*opt}
		b.SetArray("SOLANA_CLI_TXN_FLAGS", cli.Flags().Args())
	}

	return b
}

func setupPayload(p *runner.Payload, opt *TxnOptions) error {
	// Add the stake account script
	stakeAccountScript, err := assets.Open(assetsStakeAccountScript)
	if err != nil {
		return fmt.Errorf("failed to open stake account script: %w", err)
	}
	p.AddReader("steps.sh", stakeAccountScript)

	// Add transaction options if provided
	if opt != nil {
		cli := CLITxnOptions{*opt}
		if err := cli.AddToPayload(p); err != nil {
			return fmt.Errorf("failed to add transaction options: %w", err)
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

// ------------------------------------------------------------
// StakeAccount Create Command
// ------------------------------------------------------------

type StakeAccountCreate struct {
	StakeAccount
}

func (v *StakeAccountCreate) Check() error {
	if v.Amount < 0 {
		return errors.New("stake amount cannot be negative")
	}
	return nil
}

func (v *StakeAccountCreate) Env() *runner.EnvBuilder {
	e := env(v.StakeAccount)
	e.Set("STAKE_ACCOUNT_ACTION", "CREATE")

	// No need to set authority flags - shell script will check file existence

	// Set lockup parameters if provided
	if v.LockupArgs != nil {
		e.Set("EPOCH_AVAILABLE", fmt.Sprintf("%d", v.LockupArgs.EpochAvailable))
		e.Set("CUSTODIAN_PUBKEY", v.LockupArgs.CustodianPubkey)
	}
	return e
}

func (v *StakeAccountCreate) AddToPayload(p *runner.Payload) error {
	// Add stake account script
	if err := setupPayload(p, v.TransactionOptions); err != nil {
		return err
	}

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
	return nil
}

// ------------------------------------------------------------
// StakeAccount Read Command
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
	if err := setupPayload(p, v.TransactionOptions); err != nil {
		return err
	}
	p.AddString("stake_account.json", v.StakeAccountKeyPairs.StakeAccount)
	return nil
}

// ------------------------------------------------------------
// StakeAccount Delete Command
// ------------------------------------------------------------
type StakeAccountDelete struct {
	StakeAccount
}

func (v *StakeAccountDelete) Check() error {
	return nil
}

func (v *StakeAccountDelete) Env() *runner.EnvBuilder {
	e := env(v.StakeAccount)
	e.Set("STAKE_ACCOUNT_ACTION", "DELETE")

	// Only set FORCE_DELETE if it's true
	if v.ForceDelete {
		e.SetBool("FORCE_DELETE", true)
	}

	// Add withdraw address if available
	if v.WithdrawAddress != nil {
		e.Set("WITHDRAW_ADDRESS", *v.WithdrawAddress)
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

type StakeAccountUpdate struct {
	state   StakeAccount
	newArgs StakeAccount
}

type UpdateType int

const (
	UpdateTypeAuthority UpdateType = iota
	UpdateTypeDelegate
	UpdateTypeDeactivate
	UpdateTypeLock
)

func (v *StakeAccountUpdate) updatePlan() []UpdateType {
	var updates []UpdateType

	oldKps := v.state.StakeAccountKeyPairs
	newKps := v.newArgs.StakeAccountKeyPairs

	// Authority Diff - compare actual values, not just pointers
	authorityChanged := false

	// Withdraw authority comparison
	if (oldKps.WithdrawAuthority == nil) != (newKps.WithdrawAuthority == nil) {
		// One is nil, the other isn't
		authorityChanged = true
	} else if oldKps.WithdrawAuthority != nil && newKps.WithdrawAuthority != nil {
		// Both non-nil, compare the actual string contents
		if *oldKps.WithdrawAuthority != *newKps.WithdrawAuthority {
			authorityChanged = true
		}
	}

	// Stake authority comparison
	if (oldKps.StakeAuthority == nil) != (newKps.StakeAuthority == nil) {
		// One is nil, the other isn't
		authorityChanged = true
	} else if oldKps.StakeAuthority != nil && newKps.StakeAuthority != nil {
		// Both non-nil, compare the actual string contents
		if *oldKps.StakeAuthority != *newKps.StakeAuthority {
			authorityChanged = true
		}
	}

	if authorityChanged {
		updates = append(updates, UpdateTypeAuthority)
	}

	// Lockup Args Diff
	if v.state.LockupArgs != v.newArgs.LockupArgs {
		updates = append(updates, UpdateTypeLock)
	}

	// Vote Account diff - compare actual vote account values, not just pointers
	voteAccountChanged := false

	// Different nil status (one is nil, the other isn't)
	if (oldKps.VoteAccount == nil) != (newKps.VoteAccount == nil) {
		voteAccountChanged = true
	} else if oldKps.VoteAccount != nil && newKps.VoteAccount != nil {
		// Both non-nil, compare the actual string contents
		if *oldKps.VoteAccount != *newKps.VoteAccount {
			voteAccountChanged = true
		}
	}

	if voteAccountChanged {
		if oldKps.VoteAccount != nil && newKps.VoteAccount == nil {
			updates = append(updates, UpdateTypeDeactivate)
		} else {
			updates = append(updates, UpdateTypeDelegate)
		}
	}
	return updates
}

func (v *StakeAccountUpdate) Check() error {
	return nil
}

func (v *StakeAccountUpdate) Env() *runner.EnvBuilder {
	e := env(v.newArgs)
	e.Set("STAKE_ACCOUNT_ACTION", "UPDATE")

	// Determine necessary operations
	updates := v.updatePlan()

	// Set operation type for deactivation
	if slices.Contains(updates, UpdateTypeDeactivate) {
		e.Set("OPERATION", "DEACTIVATE")
	}

	// Set lockup parameters if needed
	if slices.Contains(updates, UpdateTypeLock) && v.newArgs.LockupArgs != nil {
		e.Set("EPOCH_AVAILABLE", fmt.Sprintf("%d", v.newArgs.LockupArgs.EpochAvailable))
		e.Set("CUSTODIAN_PUBKEY", v.newArgs.LockupArgs.CustodianPubkey)
	}

	// Set flags for authority updates
	if slices.Contains(updates, UpdateTypeAuthority) {
		oldStakeAuth := v.state.StakeAccountKeyPairs.StakeAuthority
		newStakeAuth := v.newArgs.StakeAccountKeyPairs.StakeAuthority
		if (oldStakeAuth == nil) != (newStakeAuth == nil) || (oldStakeAuth != nil && newStakeAuth != nil && *oldStakeAuth != *newStakeAuth) {
			e.Set("UPDATE_STAKE_AUTHORITY", "true")
		}

		oldWithdrawAuth := v.state.StakeAccountKeyPairs.WithdrawAuthority
		newWithdrawAuth := v.newArgs.StakeAccountKeyPairs.WithdrawAuthority
		if (oldWithdrawAuth == nil) != (newWithdrawAuth == nil) || (oldWithdrawAuth != nil && newWithdrawAuth != nil && *oldWithdrawAuth != *newWithdrawAuth) {
			e.Set("UPDATE_WITHDRAW_AUTHORITY", "true")
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

	updates := v.updatePlan()

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

	// Validate that required files were added
	return nil
}
