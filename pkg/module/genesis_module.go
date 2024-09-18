package module

import (
	"github.com/abklabs/svmkit/pkg/runner"
)

// Genesis is an interface for generating the genesis ledger for a network.
type Genesis interface {
	// Create returns a Command to create a genesis ledger.
	Create() runner.Command
}

// PrimordialEntry represents a seeding of lamports to an account at the genesis of the network.
type PrimorialEntry struct {
	Pubkey   string `pulumi:"pubkey"`
	Lamports string `pulumi:"lamports"`
}
