package genesis

import (
	"github.com/abklabs/pulumi-svm/provider/pkg/runner"
)

// Genesis is an interface for managing the blockchain validator.
type Genesis interface {
	// Create returns a Command to create the blockchain validator.
	Create() runner.Command
}

type PrimorialEntry struct {
	Pubkey   string `pulumi:"pubkey"`
	Lamports string `pulumi:"lamports"`
}
