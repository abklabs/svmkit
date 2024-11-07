package solana

import (
	"github.com/abklabs/svmkit/pkg/utils"
)

type CLIConfig struct {
	URL *string
}

func (f *CLIConfig) ToFlags() *utils.FlagBuilder {
	b := utils.FlagBuilder{}

	b.AppendP("url", f.URL)

	return &b
}
