package solana

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/gagliardetto/solana-go"
	"golang.org/x/crypto/ssh"
)

// TODO: VoteAccount should be an address
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

// These values are go types converted directly from the rust types in the agave cli program
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

type CreateArgs struct {
	StakeAccountKeyPair      string
	Amount                   float64
	LockupArgs               *StakeAccountLockup
	StakeAuthorityAddress    *string
	WithdrawAuthorityAddress *string
}

type DelegateArgs struct {
	StakeAccountAddress   string
	VoteAccountAddress    string
	StakeAuthorityKeypair *string
}

type DeactivateArgs struct {
	StakeAccountAddress   string
	StakeAuthorityKeypair *string
}

type AuthorizeType int

const (
	AuthorizeStaker AuthorizeType = iota
	AuthorizeWithdrawer
)

type AuthorizeArgs struct {
	StakeAccountAddress string

	OldKeyPair *string
	NewAddress string
	AuthType   AuthorizeType

	// Only needed on withdrawer authorization change
	LockupKeypair *string
}

type WithdrawArgs struct {
	StakeAccountAddress      string
	Amount                   float64
	WithdrawAddress          string
	WithdrawAuthorityKeypair *string
	LockupKeypair            *string
}

type SetLockupArgs struct {
	StakeAccountAddress string
	LockupArgs          *StakeAccountLockup
}

// StakeOperator defines interface for stake account operations
// It's agnostic to connection type and query method (RPC or bash commands)
type StakeOperator interface {
	GetStatus(address string) (CliStakeState, error)
	Create(CreateArgs) error
	Delegate(DelegateArgs) error
	Deactivate(DeactivateArgs) error
	SetLockup(SetLockupArgs) error
	Authorize(AuthorizeArgs) error
	Withdraw(WithdrawArgs) error
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
	var mWithdrawAuthAddress *string
	var mStakeAuthAddress *string

	// Set withdraw authority if provided
	if args.StakeAccountKeyPairs.WithdrawAuthority != nil {
		withdrawAuthAddress, err := getPubkeyFromJson(*args.StakeAccountKeyPairs.WithdrawAuthority)
		if err != nil {
			return StakeAccount{}, err
		}
		mWithdrawAuthAddress = &withdrawAuthAddress
	}

	// Set stake authority if provided
	if args.StakeAccountKeyPairs.StakeAuthority != nil {
		stakeAuthAddress, err := getPubkeyFromJson(*args.StakeAccountKeyPairs.StakeAuthority)
		if err != nil {
			return StakeAccount{}, err
		}
		mStakeAuthAddress = &stakeAuthAddress
	}

	createArgs := CreateArgs{
		StakeAccountKeyPair:      args.StakeAccountKeyPairs.StakeAccount,
		Amount:                   args.Amount,
		LockupArgs:               args.LockupArgs,
		StakeAuthorityAddress:    mStakeAuthAddress,
		WithdrawAuthorityAddress: mWithdrawAuthAddress,
	}

	if err := c.operator.Create(createArgs); err != nil {
		return StakeAccount{}, err
	}

	stakeAccountAddress, err := getPubkeyFromJson(args.StakeAccountKeyPairs.StakeAccount)
	if err != nil {
		return StakeAccount{}, err
	}

	if args.StakeAccountKeyPairs.VoteAccount != nil {
		voteAccountAddress, err := getPubkeyFromJson(*args.StakeAccountKeyPairs.VoteAccount)
		if err != nil {
			return StakeAccount{}, err
		}
		delegateArgs := DelegateArgs{
			StakeAccountAddress:   stakeAccountAddress,
			VoteAccountAddress:    voteAccountAddress,
			StakeAuthorityKeypair: args.StakeAccountKeyPairs.StakeAuthority,
		}
		if err := c.operator.Delegate(delegateArgs); err != nil {
			return StakeAccount{}, err
		}
	}
	return args, nil
}

func (c *StakeAccountClient) Update(state StakeAccount, newArgs StakeAccount) (StakeAccount, error) {
	// Basic checks
	if state.StakeAccountKeyPairs.StakeAccount != newArgs.StakeAccountKeyPairs.StakeAccount {
		return StakeAccount{}, errors.New("stake account address can not be rotated")
	}
	if newArgs.Amount != state.Amount {
		return StakeAccount{}, errors.New("cannot change stake amount; operation not currently supported")
	}

	// Read state from chain
	stakeAccountAddress, err := getPubkeyFromJson(state.StakeAccountKeyPairs.StakeAccount)
	if err != nil {
		return StakeAccount{}, err
	}

	readState, err := c.operator.GetStatus(stakeAccountAddress)
	if err != nil {
		return StakeAccount{}, errors.New("failed to read stake account state from chain")
	}

	// Handle vote-account change
	currentVA := state.StakeAccountKeyPairs.VoteAccount
	newVA := newArgs.StakeAccountKeyPairs.VoteAccount
	vaChangeType := accountChange(currentVA, newVA)

	if vaChangeType == Added {
		if isFullyDeactivated(readState) {
			newVAAddress, err := getPubkeyFromJson(*newVA)
			if err != nil {
				return StakeAccount{}, err
			}
			// Delegate to new vote account
			delegateArgs := DelegateArgs{
				StakeAccountAddress: stakeAccountAddress,
				VoteAccountAddress:  newVAAddress,
				// Use old state since we haven't done authority updates yet
				StakeAuthorityKeypair: state.StakeAccountKeyPairs.StakeAuthority,
			}
			if err := c.operator.Delegate(delegateArgs); err != nil {
				return StakeAccount{}, err
			}
		} else {
			return StakeAccount{}, errors.New("cannot change vote account until stake is fully deactivated")
		}
	} else if vaChangeType == Removed {
		// TODO: Check readstate
		// Deactivate from current vote account
		deactivateArgs := DeactivateArgs{
			StakeAccountAddress:   stakeAccountAddress,
			StakeAuthorityKeypair: state.StakeAccountKeyPairs.StakeAuthority,
		}
		if err := c.operator.Deactivate(deactivateArgs); err != nil {
			return StakeAccount{}, err
		}
	} else if vaChangeType == Modified {
		return StakeAccount{}, errors.New("must fully deactivate stake before changing vote account")
	}

	currentWAuth := state.StakeAccountKeyPairs.WithdrawAuthority
	newWAuth := newArgs.StakeAccountKeyPairs.WithdrawAuthority
	wAuthChangeType := accountChange(currentWAuth, newWAuth)

	if wAuthChangeType == Added || wAuthChangeType == Modified {
		// Change withdraw authority
		newWAuthAddress, err := getPubkeyFromJson(*newWAuth)
		if err != nil {
			return StakeAccount{}, err
		}
		authorizeArgs := AuthorizeArgs{
			StakeAccountAddress: stakeAccountAddress,
			OldKeyPair:          state.StakeAccountKeyPairs.WithdrawAuthority,
			NewAddress:          newWAuthAddress,
			AuthType:            AuthorizeWithdrawer,
			// TODO: We currently only store a pubkey for lockup, but we need to store the full keypair
			// LockupKeypair:       state.LockupArgs.CustodianPubkey,
		}
		if err := c.operator.Authorize(authorizeArgs); err != nil {
			return StakeAccount{}, err
		}

	} else if wAuthChangeType == Removed {
		return StakeAccount{}, errors.New("cannot remove withdraw authority")
	}

	currentSAuth := state.StakeAccountKeyPairs.StakeAuthority
	newSAuth := newArgs.StakeAccountKeyPairs.StakeAuthority
	sAuthChangeType := accountChange(currentSAuth, newSAuth)

	if sAuthChangeType == Added || sAuthChangeType == Modified {
		// Change stake authority
		newSAuthAddress, err := getPubkeyFromJson(*newSAuth)
		if err != nil {
			return StakeAccount{}, err
		}
		authorizeArgs := AuthorizeArgs{
			StakeAccountAddress: stakeAccountAddress,
			OldKeyPair:          state.StakeAccountKeyPairs.StakeAuthority,
			NewAddress:          newSAuthAddress,
			AuthType:            AuthorizeStaker,
		}
		if err := c.operator.Authorize(authorizeArgs); err != nil {
			return StakeAccount{}, err
		}
	} else if sAuthChangeType == Removed {
		return StakeAccount{}, errors.New("cannot remove stake authority")
	}

	// TODO: Change lockup

	return newArgs, nil
}

type DiffStatus int

const (
	Unchanged DiffStatus = iota // both values are nil or equal
	Added                       // old is nil and new is non-nil
	Removed                     // old is non-nil and new is nil
	Modified                    // both non-nil but different
)

func accountChange(oldVal *string, newVal *string) DiffStatus {
	// Both are nil => unchanged
	if oldVal == nil && newVal == nil {
		return Unchanged
	}
	// Added: old is nil, new is non-nil
	if oldVal == nil && newVal != nil {
		return Added
	}
	// Removed: old is non-nil, new is nil
	if oldVal != nil && newVal == nil {
		return Removed
	}
	// Both are non-nil. If equal, unchanged; if not, modified.
	if *oldVal == *newVal {
		return Unchanged
	}
	return Modified
}

func (c *StakeAccountClient) Delete(state StakeAccount) error {
	if state.ForceDelete {
		return nil
	}

	stakeAddress, err := getPubkeyFromJson(state.StakeAccountKeyPairs.StakeAccount)
	if err != nil {
		return err
	}

	readState, err := c.operator.GetStatus(stakeAddress)
	if err != nil {
		return errors.New("failed to read stake account state from chain")
	}

	//TODO: Is this extra forcedelete check necessary given we already checked it above?
	if state.WithdrawAddress == nil && !state.ForceDelete {
		return errors.New("must provide withdraw address or set force_delete to true")
	}

	if state.WithdrawAddress != nil && !isFullyDeactivated(readState) {
		return errors.New("cannot withdraw stake until it is fully deactivated")
	}

	// TODO: Check lockup state from read and ensure it's unlocked

	widrawArgs := WithdrawArgs{
		StakeAccountAddress:      stakeAddress,
		Amount:                   state.Amount,
		WithdrawAddress:          *state.WithdrawAddress,
		WithdrawAuthorityKeypair: state.StakeAccountKeyPairs.WithdrawAuthority,
		LockupKeypair:            state.StakeAccountKeyPairs.StakeAuthority,
	}

	if err := c.operator.Withdraw(widrawArgs); err != nil {
		return err
	}
	return nil
}

func isFullyDeactivated(state CliStakeState) bool {
	return (state.DelegatedStake == nil || *state.DelegatedStake == 0) &&
		(state.DeactivatingStake == nil || *state.DeactivatingStake == 0)
}

// ------------------------------------------------------------
// StakeOperator CLI Implementation
// ------------------------------------------------------------

type CliStakeOperator struct {
	client     *ssh.Client
	handler    runner.DeployerHandler
	ctx        context.Context
	txnOptions *TxnOptions
}

func NewCliStakeOperator(client *ssh.Client, handler runner.DeployerHandler, ctx context.Context) *CliStakeOperator {
	return &CliStakeOperator{
		client:     client,
		handler:    handler,
		ctx:        ctx,
		txnOptions: nil,
	}
}

func (op *CliStakeOperator) SetTxnOptions(opt *TxnOptions) {
	op.txnOptions = opt
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

func (op *CliStakeOperator) Create(args CreateArgs) error {
	cmd := &StakeAccountCreate{args, op.txnOptions}
	if err := op.runCommand(cmd, op.handler); err != nil {
		return err
	}
	return nil
}

func (op *CliStakeOperator) Delegate(args DelegateArgs) error {
	cmd := &StakeAccountDelegate{args, op.txnOptions}
	if err := op.runCommand(cmd, op.handler); err != nil {
		return err
	}
	return nil
}

func (op *CliStakeOperator) Deactivate(args DeactivateArgs) error {
	cmd := &StakeAccountDeactivate{args, op.txnOptions}
	if err := op.runCommand(cmd, op.handler); err != nil {
		return err
	}
	return nil
}

func (op *CliStakeOperator) Authorize(args AuthorizeArgs) error {
	cmd := &StakeAccountAuthorize{args, op.txnOptions}
	if err := op.runCommand(cmd, op.handler); err != nil {
		return err
	}
	return nil
}

func (op *CliStakeOperator) Withdraw(args WithdrawArgs) error {
	cmd := &StakeAccountWithdraw{args, op.txnOptions}
	if err := op.runCommand(cmd, op.handler); err != nil {
		return err
	}
	return nil
}

func (op *CliStakeOperator) GetStatus(stakeAddress string) (CliStakeState, error) {
	// This is the lone operation that doesn't use the provided handler and
	// uses a StringHandler so that the output can be parsed
	handler := &StringHandler{}

	cmd := &StakeAccountRead{stakeAddress, op.txnOptions}

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

func (op *CliStakeOperator) SetLockup(args SetLockupArgs) error {
	// TODO
	return nil
}

// ------------------------------------------------------------
// Common Command Helper Functions
// ------------------------------------------------------------
func envWithOptions(txnOptions *TxnOptions) *runner.EnvBuilder {
	// Sets default env for all stake Commands
	b := runner.NewEnvBuilder()

	// Set transaction flags if available
	if opt := txnOptions; opt != nil {
		cli := CLITxnOptions{*opt}
		b.SetArray("SOLANA_CLI_TXN_FLAGS", cli.Flags().Args())
	}
	return b
}

func getPubkeyFromJson(jsonStr string) (string, error) {
	// Parse the JSON string into a slice of ints
	var privateKeyInts []int
	err := json.Unmarshal([]byte(jsonStr), &privateKeyInts)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Convert back to a byte slice
	// NOTE: We only do this switch from ints => bytes bc pulumi doesn't support uint8
	privateKeyBytes := make([]byte, len(privateKeyInts))
	for i, v := range privateKeyInts {
		privateKeyBytes[i] = byte(v)
	}
	pubKey := solana.PrivateKey(privateKeyBytes).PublicKey().String()
	return pubKey, nil
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

// ------------------------------------------------------------
// StakeAccount Create Command
// ------------------------------------------------------------

type StakeAccountCreate struct {
	CreateArgs
	TxnOptions *TxnOptions
}

func (a *StakeAccountCreate) Check() error {
	if a.Amount < 0 {
		return errors.New("stake amount cannot be negative")
	}
	if a.StakeAccountKeyPair == "" {
		return errors.New("stake account keypair is required")
	}
	return nil
}

func (a *StakeAccountCreate) Env() *runner.EnvBuilder {
	e := envWithOptions(a.TxnOptions)
	e.Set("STAKE_ACCOUNT_ACTION", "CREATE")
	e.SetFloat64("STAKE_AMOUNT", a.Amount)

	// Set lockup parameters if provided
	if a.LockupArgs != nil {
		e.Set("EPOCH_AVAILABLE", fmt.Sprintf("%d", a.LockupArgs.EpochAvailable))
		e.Set("CUSTODIAN_PUBKEY", a.LockupArgs.CustodianPubkey)
	}
	if a.StakeAuthorityAddress != nil {
		e.Set("STAKE_AUTHORITY_ADDRESS", *a.StakeAuthorityAddress)
	}
	if a.WithdrawAuthorityAddress != nil {
		e.Set("WITHDRAW_AUTHORITY_ADDRESS", *a.WithdrawAuthorityAddress)
	}
	return e
}

func (a *StakeAccountCreate) AddToPayload(p *runner.Payload) error {
	// Add stake account script
	if err := setupPayload(p, a.TxnOptions); err != nil {
		return err
	}
	// Add keypair files
	p.AddString("stake_account.json", a.StakeAccountKeyPair)

	return nil
}

// ------------------------------------------------------------
// StakeAccount Delegate Command
// ------------------------------------------------------------
type StakeAccountDelegate struct {
	DelegateArgs
	TxnOptions *TxnOptions
}

func (v *StakeAccountDelegate) Check() error {
	if v.StakeAccountAddress == "" {
		return errors.New("stake account address is required")
	}
	if v.VoteAccountAddress == "" {
		return errors.New("vote account address is required")
	}
	return nil
}

func (v *StakeAccountDelegate) Env() *runner.EnvBuilder {
	e := envWithOptions(v.TxnOptions)
	e.Set("STAKE_ACCOUNT_ACTION", "DELEGATE")
	e.Set("STAKE_ACCOUNT_ADDRESS", v.StakeAccountAddress)
	e.Set("VOTE_ACCOUNT_ADDRESS", v.VoteAccountAddress)
	return e
}

func (v *StakeAccountDelegate) AddToPayload(p *runner.Payload) error {
	if err := setupPayload(p, v.TxnOptions); err != nil {
		return err
	}
	if v.StakeAuthorityKeypair != nil {
		p.AddString("stake_authority.json", *v.StakeAuthorityKeypair)
	}
	return nil
}

// ------------------------------------------------------------
// StakeAccount Deactivate Command
// ------------------------------------------------------------
type StakeAccountDeactivate struct {
	DeactivateArgs
	TxnOptions *TxnOptions
}

func (a *StakeAccountDeactivate) Check() error {
	if a.StakeAccountAddress == "" {
		return errors.New("stake account address is required")
	}
	return nil
}

func (a *StakeAccountDeactivate) Env() *runner.EnvBuilder {
	e := envWithOptions(a.TxnOptions)
	e.Set("STAKE_ACCOUNT_ACTION", "DEACTIVATE")
	e.Set("STAKE_ACCOUNT_ADDRESS", a.StakeAccountAddress)
	return e
}

func (a *StakeAccountDeactivate) AddToPayload(p *runner.Payload) error {
	if err := setupPayload(p, a.TxnOptions); err != nil {
		return err
	}
	if a.StakeAuthorityKeypair != nil {
		p.AddString("stake_authority.json", *a.StakeAuthorityKeypair)
	}
	return nil
}

// ------------------------------------------------------------
// StakeAccount Authorize Command
// ------------------------------------------------------------

type StakeAccountAuthorize struct {
	AuthorizeArgs
	TxnOptions *TxnOptions
}

func (a *StakeAccountAuthorize) Check() error {
	if a.AuthType == AuthorizeStaker && a.LockupKeypair != nil {
		return errors.New("staker authorization should not have a lockup keypair")
	}
	return nil
}

func (v *StakeAccountAuthorize) Env() *runner.EnvBuilder {
	e := envWithOptions(v.TxnOptions)
	e.Set("STAKE_ACCOUNT_ACTION", "AUTHORIZE")
	e.Set("STAKE_ACCOUNT_ADDRESS", v.StakeAccountAddress)
	e.Set("NEW_AUTHORITY_ADDRESS", v.NewAddress)

	if v.AuthType == AuthorizeStaker {
		e.Set("AUTH_TYPE", "STAKER")
	} else {
		e.Set("AUTH_TYPE", "WITHDRAWER")
	}
	return e
}

func (v *StakeAccountAuthorize) AddToPayload(p *runner.Payload) error {
	if err := setupPayload(p, v.TxnOptions); err != nil {
		return err
	}
	if v.OldKeyPair != nil {
		p.AddString("old_authority.json", *v.OldKeyPair)
	}
	if v.LockupKeypair != nil {
		p.AddString("lockup_keypair.json", *v.LockupKeypair)
	}
	return nil
}

// ------------------------------------------------------------
// StakeAccount Read Command
// ------------------------------------------------------------

type StakeAccountRead struct {
	StakeAddress string
	TxnOptions   *TxnOptions
}

func (v *StakeAccountRead) Check() error {
	return nil
}

func (v *StakeAccountRead) Env() *runner.EnvBuilder {
	e := envWithOptions(nil)
	e.Set("STAKE_ACCOUNT_ACTION", "READ")
	e.Set("STAKE_ACCOUNT_ADDRESS", v.StakeAddress)
	return e
}

func (v *StakeAccountRead) AddToPayload(p *runner.Payload) error {
	if err := setupPayload(p, nil); err != nil {
		return err
	}
	return nil
}

// ------------------------------------------------------------
// StakeAccount Withdraw Command
// ------------------------------------------------------------
type StakeAccountWithdraw struct {
	WithdrawArgs
	TxnOptions *TxnOptions
}

func (a *StakeAccountWithdraw) Check() error {
	if a.StakeAccountAddress == "" {
		return errors.New("stake account address is required")
	}
	if a.WithdrawAddress == "" {
		return errors.New("withdraw address is required")
	}
	if a.Amount < 0 {
		return errors.New("withdraw amount cannot be negative")
	}
	return nil
}

func (a *StakeAccountWithdraw) Env() *runner.EnvBuilder {
	e := envWithOptions(a.TxnOptions)
	e.Set("STAKE_ACCOUNT_ACTION", "WITHDRAW")
	e.SetFloat64("WITHDRAW_AMOUNT", a.Amount)
	e.Set("STAKE_ACCOUNT_ADDRESS", a.StakeAccountAddress)
	e.Set("WITHDRAW_ADDRESS", a.WithdrawAddress)
	return e
}

func (a *StakeAccountWithdraw) AddToPayload(p *runner.Payload) error {
	if err := setupPayload(p, a.TxnOptions); err != nil {
		return err
	}
	if a.WithdrawAuthorityKeypair != nil {
		p.AddString("withdraw_authority.json", *a.WithdrawAuthorityKeypair)
	}
	if a.LockupKeypair != nil {
		p.AddString("lockup_keypair.json", *a.LockupKeypair)
	}
	return nil
}
