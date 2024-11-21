package solana

import (
	"github.com/abklabs/svmkit/pkg/runner"
)

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
