package solana

import (
	"github.com/abklabs/svmkit/pkg/runner"
)

type StakeAccountKeyPairs struct {
	StakeAccount string `pulumi:"stakeAccount" provider:"secret"`
	VoteAccount  string `pulumi:"voteAccount" provider:"secret"`
}

type StakeAccount struct {
	TransactionOptions   *TxnOptions          `pulumi:"transactionOptions"`
	StakeAccountKeyPairs StakeAccountKeyPairs `pulumi:"keyPairs"`
	Amount               float64              `pulumi:"amount"`
}

type WithdrawOptions struct {
	WithdrawAccount string `pulumi:"withdrawAccount" provider:"secret"`
	// Just withdraw full amount for now
	// Amount          float64 `pulumi:"amount"`
}

func (v *StakeAccount) Create() runner.Command {
	return &StakeAccountCreate{
		StakeAccount: *v,
	}
}

func (v *StakeAccount) Env() *runner.EnvBuilder {
	b := runner.NewEnvBuilder()

	b.SetFloat64("STAKE_AMOUNT", v.Amount)

	if opt := v.TransactionOptions; opt != nil {
		cli := CLITxnOptions{*opt}
		b.SetArray("SOLANA_CLI_TXN_FLAGS", cli.Flags().Args())
	}

	return b
}

type StakeAccountCreate struct {
	StakeAccount
}

func (v *StakeAccountCreate) Check() error {
	return nil
}

func (v *StakeAccountCreate) Env() *runner.EnvBuilder {
	e := v.StakeAccount.Env()
	e.Set("STAKE_ACCOUNT_ACTION", "CREATE")

	return e
}

func (v *StakeAccountCreate) AddToPayload(p *runner.Payload) error {
	stakeAccountScript, err := assets.Open(assetsStakeAccountScript)

	if err != nil {
		return err
	}

	p.AddReader("steps.sh", stakeAccountScript)

	p.AddString("stake_account.json", v.StakeAccountKeyPairs.StakeAccount)
	p.AddString("vote_account.json", v.StakeAccountKeyPairs.VoteAccount)

	if opt := v.TransactionOptions; opt != nil {
		cli := CLITxnOptions{*opt}

		err := cli.AddToPayload(p)

		if err != nil {
			return err
		}
	}

	return nil
}

type StakeAccountUpdate struct {
	StakeAccount
}

func (v *StakeAccountUpdate) Check() error {
	return nil
}

func (v *StakeAccountUpdate) Env() *runner.EnvBuilder {
	e := v.StakeAccount.Env()
	e.Set("STAKE_ACCOUNT_ACTION", "DEACTIVATE")

	return e
}

func (v *StakeAccountUpdate) AddToPayload(p *runner.Payload) error {
	stakeAccountScript, err := assets.Open(assetsStakeAccountScript)

	if err != nil {
		return err
	}

	p.AddReader("steps.sh", stakeAccountScript)

	p.AddString("stake_account.json", v.StakeAccount.StakeAccountKeyPairs.StakeAccount)

	if opt := v.TransactionOptions; opt != nil {
		cli := CLITxnOptions{*opt}

		err := cli.AddToPayload(p)

		if err != nil {
			return err
		}
	}

	return nil
}

type StakeAccountDelete struct {
	StakeAccount
	WithdrawOptions
}

func (v *StakeAccountDelete) Check() error {
	//check that stake_amount >= withdraw amount
	// Check that undelegated amount >= withdraw amount
	return nil
}

func (v *StakeAccountDelete) Env() *runner.EnvBuilder {
	e := v.StakeAccount.Env()
	e.Set("STAKE_ACCOUNT_ACTION", "WITHDRAW")
	e.SetFloat64("WITHDRAW_AMOUNT", v.StakeAccount.Amount)
	// e.SetFloat64("WITHDRAW_AMOUNT", v.WithdrawOptions.Amount)

	return e
}

func (v *StakeAccountDelete) AddToPayload(p *runner.Payload) error {
	stakeAccountScript, err := assets.Open(assetsStakeAccountScript)

	if err != nil {
		return err
	}

	p.AddReader("steps.sh", stakeAccountScript)

	p.AddString("stake_account.json", v.StakeAccount.StakeAccountKeyPairs.StakeAccount)
	p.AddString("withdraw_account.json", v.WithdrawOptions.WithdrawAccount)

	if opt := v.TransactionOptions; opt != nil {
		cli := CLITxnOptions{*opt}

		err := cli.AddToPayload(p)

		if err != nil {
			return err
		}
	}

	return nil
}
