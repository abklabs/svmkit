package solana

import (
	"github.com/abklabs/svmkit/pkg/runner"
)

type StakeAccountKeyPairs struct {
	StakeAccount string `pulumi:"stakeAccount" provider:"secret"`
	VoteAccount  string `pulumi:"voteAccount" provider:"secret"`
}

type StakeAccount struct {
	runner.RunnerCommand

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
	v.SetConfigDefaults()

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
