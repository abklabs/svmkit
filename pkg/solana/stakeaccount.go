package solana

import (
	"fmt"

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

type StakeAccountDeactivate struct {
	StakeAccount
}

func (v *StakeAccountDeactivate) Check() error {
	if v.Amount != 0 {
		return fmt.Errorf("stake amount must be 0 for deactivation, got %f", v.Amount)
	}
	return nil
}

func (v *StakeAccountDeactivate) Env() *runner.EnvBuilder {
	e := v.StakeAccount.Env()
	e.Set("STAKE_ACCOUNT_ACTION", "DEACTIVATE")

	return e
}

func (v *StakeAccountDeactivate) AddToPayload(p *runner.Payload) error {
	stakeAccountScript, err := assets.Open(assetsStakeAccountScript)

	if err != nil {
		return err
	}

	p.AddReader("steps.sh", stakeAccountScript)

	p.AddString("stake_account.json", v.StakeAccountKeyPairs.StakeAccount)

	if opt := v.TransactionOptions; opt != nil {
		cli := CLITxnOptions{*opt}

		err := cli.AddToPayload(p)

		if err != nil {
			return err
		}
	}

	return nil
}
