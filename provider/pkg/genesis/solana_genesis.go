package genesis

import (
	"context"

	_ "embed"
	"fmt"

	"github.com/abklabs/svmkit/pkg/genesis"
	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/abklabs/svmkit/pkg/solana"
	"github.com/abklabs/svmkit/pkg/ssh"
	"github.com/abklabs/svmkit/provider/pkg/svm"
)

// Solana represents a Pulumi resource for building the genesis ledger for the Solana network.
type Solana struct{}

// SolanaArgs represents the input arguments required to create a Solana genesis resource.
type SolanaArgs struct {
	svm.GenesisArgs
	// Flags contains the configuration flags for the Solana genesis setup.
	Flags      solana.GenesisFlags      `pulumi:"flags"`
	Primordial []genesis.PrimorialEntry `pulumi:"primordial"`
}

// SolanaState represents the state of a Solana genesis resource.
type SolanaState struct {
	SolanaArgs
	// Hash is the genesis hash of the Solana ledger.
	Hash string `pulumi:"genesisHash"`
}

// Create sets up the Solana genesis ledger on the specified machine using the provided connection and flags.
//
// Parameters:
// - ctx: The context for the creation operation.
// - name: The name of the resource.
// - input: The input arguments for creating the resource.
// - preview: A boolean indicating whether this is a preview operation.
//
// Returns:
// - The name of the created resource.
// - The state of the created resource.
// - An error if the creation fails.
func (Solana) Create(ctx context.Context, name string, input SolanaArgs, preview bool) (string, SolanaState, error) {
	state := SolanaState{SolanaArgs: input}

	if preview {
		return name, state, nil
	}

	genesis := &solana.Genesis{
		Flags:      input.Flags,
		Primordial: input.Primordial,
	}
	command := genesis.Create()

	r := runner.Machine(input.Connection).
		Env(command.Env()).
		Script(command.Script())

	if err := r.Run(ctx); err != nil {
		return "", SolanaState{}, fmt.Errorf("failed to setup Solana genesis: %w", err)
	}

	// Establish SSH connection
	connection, err := input.Connection.Dial(ctx)
	if err != nil {
		return "", SolanaState{}, fmt.Errorf("failed to establish SSH connection: %w", err)
	}
	defer connection.Close()

	// Execute the command on the remote machine
	stdout, stderr, err := ssh.Exec(ctx, connection, "solana genesis-hash")
	if err != nil {
		return "", SolanaState{}, fmt.Errorf("failed to execute Solana genesis command: %w, stderr: %s", err, stderr)
	}

	state.Hash = stdout

	return name, state, nil
}
