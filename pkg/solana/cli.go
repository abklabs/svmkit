package solana

import (
	"strings"

	"github.com/abklabs/svmkit/pkg/runner"
)

func NewSecretPayload(path string, body string) runner.PayloadFile {
	return runner.PayloadFile{
		Path:   path,
		Reader: strings.NewReader(body),
		Mode:   0400,
	}
}

type CLIConfig struct {
	URL     *string
	KeyPair *string
}

func (f *CLIConfig) ToFlags() *runner.FlagBuilder {
	b := runner.FlagBuilder{}

	b.AppendP("url", f.URL)
	b.AppendP("keypair", f.KeyPair)

	return &b
}

type CLITxnOptions struct {
	TxnOptions
}

func (f *CLITxnOptions) ToFlags() *runner.FlagBuilder {
	b := &runner.FlagBuilder{}

	b.AppendP("blockhash", f.BlockHash)
	b.AppendFloat64P("with-compute-unit-price", f.WithComputeUnitPrice)
	b.AppendP("with-memo", f.WithMemo)
	b.AppendP("nonce", f.Nonce)
	b.AppendP("nonce-authority", f.NonceAuthority)

	if s := f.Signer; s != nil {
		for _, signer := range *s {
			b.AppendRaw("--signer", signer)
		}
	}

	if f.FeePayer != nil {
		b.AppendRaw("--fee-payer", "txn_feepayer.json")
	}

	if f.From != nil {
		b.AppendRaw("--from", "txn_from.json")
	}

	if f.KeyPair != nil {
		b.AppendRaw("--keypair", "txn_keypair.json")
	}

	b.AppendP("ws", f.WS)
	b.AppendP("url", f.URL)

	return b
}

func (c *CLITxnOptions) AddToPayload(p *runner.Payload) error {
	if c.KeyPair != nil {
		p.Add(NewSecretPayload("txn_keypair.json", *c.KeyPair))
	}

	if c.From != nil {
		p.Add(NewSecretPayload("txn_from.json", *c.From))
	}

	if c.FeePayer != nil {
		p.Add(NewSecretPayload("txn_fee_payer.json", *c.FeePayer))
	}

	return nil
}
