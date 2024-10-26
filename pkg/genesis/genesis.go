package genesis

import (
	"github.com/abklabs/svmkit/pkg/runner"
)

// Version is the semantic version (semver) of the blockchain validator apt package.
type Version *string

// Genesis is an interface for managing the blockchain validator.
type Genesis interface {
	// Create returns a Command to create the blockchain validator.
	Create() runner.Command
}

type PrimorialEntry struct {
	Pubkey   string `pulumi:"pubkey"`
	Lamports string `pulumi:"lamports"`
}
