package utils

import (
	flag "github.com/spf13/pflag"
	"io"
	"strings"
)

func AddScriptFlags(flags *flag.FlagSet) {
	flags.String("script-pre", "", "opsh code to inject before the script")
	flags.String("script-post", "", "opsh code to inject after the script")
}

func AssembleScript(flags *flag.FlagSet, readers ...io.Reader) (io.Reader, error) {
	pre, err := flags.GetString("script-pre")

	if err != nil {
		return nil, err
	}

	if pre != "" {
		readers = append([]io.Reader{strings.NewReader(pre + "\n")}, readers...)
	}

	post, err := flags.GetString("script-post")

	if err != nil {
		return nil, err
	}

	if post != "" {
		readers = append(readers, strings.NewReader(post+"\n"))
	}

	return io.MultiReader(readers...), nil
}
