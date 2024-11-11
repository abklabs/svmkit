package solana

import (
	"github.com/abklabs/svmkit/pkg/utils"
)

type CLIConfig struct {
	URL     *string
	KeyPair *string
}

func (f *CLIConfig) ToFlags() *utils.FlagBuilder {
	b := utils.FlagBuilder{}

	b.AppendP("url", f.URL)
	b.AppendP("keypair", f.KeyPair)

	return &b
}
