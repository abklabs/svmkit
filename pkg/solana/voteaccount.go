package solana

import (
	"github.com/abklabs/svmkit/pkg/runner"
)

type VoteAccountKeyPairs struct {
	VoteAccount    string `pulumi:"voteAccount" provider:"secret"`
	Identity       string `pulumi:"identity" provider:"secret"`
	AuthWithdrawer string `pulumi:"authWithdrawer" provider:"secret"`
}

type VoteAccount struct {
	VoteAccountKeyPairs  VoteAccountKeyPairs `pulumi:"keyPairs"`
	AuthVoterPubkey      *string             `pulumi:"authVoterPubkey,optional"`
	CloseRecipientPubkey *string             `pulumi:"closeRecipientPubkey,optional"`
}

func (v *VoteAccount) Create() runner.Command {
	return &VoteAccountCreate{
		VoteAccount: *v,
	}
}

func (v *VoteAccount) Delete() runner.Command {
	return &VoteAccountDelete{
		VoteAccount: *v,
	}
}

func (v *VoteAccount) Env() *runner.EnvBuilder {
	b := runner.NewEnvBuilder()

	b.SetP("AUTH_VOTER_PUBKEY", v.AuthVoterPubkey)
	b.SetP("CLOSE_RECIPIENT_PUBKEY", v.CloseRecipientPubkey)

	return b
}

type VoteAccountCreate struct {
	VoteAccount
}

func (v *VoteAccountCreate) Check() error {
	return nil
}

func (v *VoteAccountCreate) Env() *runner.EnvBuilder {
	e := v.VoteAccount.Env()
	e.Set("VOTE_ACCOUNT_ACTION", "CREATE")

	return e
}

func (v *VoteAccountCreate) AddToPayload(p *runner.Payload) error {
	p.AddString("identity.json", v.VoteAccountKeyPairs.Identity)
	p.AddString("vote_account.json", v.VoteAccountKeyPairs.VoteAccount)
	p.AddString("auth_withdrawer.json", v.VoteAccountKeyPairs.AuthWithdrawer)

	p.AddString("steps.sh", VoteAccountScript)

	return nil
}

type VoteAccountDelete struct {
	VoteAccount
}

func (v *VoteAccountDelete) Check() error {
	return nil
}

func (v *VoteAccountDelete) Env() *runner.EnvBuilder {
	e := v.VoteAccount.Env()
	e.Set("VOTE_ACCOUNT_ACTION", "DELETE")

	return e
}

func (v *VoteAccountDelete) AddToPayload(p *runner.Payload) error {
	p.AddString("steps.sh", VoteAccountScript)

	return nil
}
