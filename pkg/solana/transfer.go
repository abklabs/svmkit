package solana

import (
	"github.com/abklabs/svmkit/pkg/runner"
)

type Transfer struct {
	PayerKeyPair           string  `pulumi:"payerKeyPair" provider:"secret"`
	Amount                 float64 `pulumi:"amount"`
	RecipientPubkey        string  `pulumi:"recipientPubkey"`
	AllowUnfundedRecipient *bool   `pulumi:"allowUnfundedRecipient,optional"`
}

func (v *Transfer) Create() runner.Command {
	return &TransferCreate{
		Transfer: *v,
	}
}

func (v *Transfer) Env() *runner.EnvBuilder {
	b := runner.NewEnvBuilder()

	b.SetMap(map[string]string{
		"RECIPIENT_PUBKEY": v.RecipientPubkey,
	})

	b.SetFloat64("AMOUNT", v.Amount)
	b.SetBoolP("ALLOW_UNFUNDED_RECIPIENT", v.AllowUnfundedRecipient)

	return b
}

type TransferCreate struct {
	Transfer
}

func (v *TransferCreate) Check() error {
	return nil
}

func (v *TransferCreate) Env() *runner.EnvBuilder {
	e := v.Transfer.Env()
	e.Set("TRANSFER_ACTION", "CREATE")

	return e
}

func (v *TransferCreate) AddToPayload(p *runner.Payload) error {
	p.AddString("steps.sh", TransferScript)

	p.AddString("payer.json", v.PayerKeyPair)

	return nil
}
