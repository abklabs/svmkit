package solana

import (
	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/abklabs/svmkit/pkg/utils"
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

func (v *Transfer) Env() *utils.EnvBuilder {
	b := utils.NewEnvBuilder()

	b.SetMap(map[string]string{
		"PAYER_KEYPAIR":    v.PayerKeyPair,
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

func (v *TransferCreate) Env() *utils.EnvBuilder {
	e := v.Transfer.Env()
	e.Set("TRANSFER_ACTION", "CREATE")

	return e
}

func (v *TransferCreate) Script() string {
	return TransferScript
}
