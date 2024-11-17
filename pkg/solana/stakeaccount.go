package solana

import (
	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/abklabs/svmkit/pkg/utils"
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

func (v *StakeAccount) Env() *utils.EnvBuilder {
	b := utils.NewEnvBuilder()

	b.SetMap(map[string]string{
		"STAKE_ACCOUNT_KEYPAIR": v.StakeAccountKeyPairs.StakeAccount,
		"VOTE_ACCOUNT_KEYPAIR":  v.StakeAccountKeyPairs.VoteAccount,
	})

	b.SetFloat64("STAKE_AMOUNT", v.Amount)

	return b
}

type StakeAccountCreate struct {
	StakeAccount
}

func (v *StakeAccountCreate) Check() error {
	return nil
}

func (v *StakeAccountCreate) Env() *utils.EnvBuilder {
	e := v.StakeAccount.Env()
	e.Set("STAKE_ACCOUNT_ACTION", "CREATE")

	return e
}

func (v *StakeAccountCreate) Script() string {
	return StakeAccountScript
}
