package solana

import (
	"github.com/abklabs/svmkit/pkg/runner"
)

type StakeAccountKeyPairs struct {
	StakeAccount string `pulumi:"stakeAccount" provider:"secret"`
	VoteAccount  string `pulumi:"voteAccount" provider:"secret"`
}

type StakeAccount struct {
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
	p.AddString("steps.sh", StakeAccountScript)

	p.AddString("stake_account.json", v.StakeAccountKeyPairs.StakeAccount)
	p.AddString("vote_account.json", v.StakeAccountKeyPairs.VoteAccount)

	return nil
}
