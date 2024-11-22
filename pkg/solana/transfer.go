package solana

import (
	"github.com/abklabs/svmkit/pkg/runner"
)

type Transfer struct {
	TransactionOptions     *TxnOptions `pulumi:"transactionOptions"`
	Amount                 float64     `pulumi:"amount"`
	RecipientPubkey        string      `pulumi:"recipientPubkey"`
	AllowUnfundedRecipient *bool       `pulumi:"allowUnfundedRecipient,optional"`
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

	if opt := v.TransactionOptions; opt != nil {
		cli := CLITxnOptions{*opt}
		b.SetArray("SOLANA_CLI_TXN_FLAGS", cli.ToFlags().ToArgs())
	}

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

	if opt := v.TransactionOptions; opt != nil {
		cli := CLITxnOptions{*opt}

		err := cli.AddToPayload(p)

		if err != nil {
			return err
		}
	}

	return nil
}
