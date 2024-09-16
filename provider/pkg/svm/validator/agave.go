package validator

import (
	"context"

	_ "embed"
	"fmt"

	"github.com/abklabs/pulumi-svm/provider/pkg/agave"
	"github.com/abklabs/pulumi-svm/provider/pkg/runner"
)

// Agave represents a Pulumi resource for managing an Agave validator.
type Agave struct{}

// AgaveArgs represents the input arguments required to create an Agave resource.
type AgaveArgs struct {
	ValidatorArgs
	// Flags contains the configuration flags for the Agave validator.
	Flags    agave.Flags    `pulumi:"flags"`
	KeyPairs agave.KeyPairs `pulumi:"keyPairs" provider:"secret"`
}

// AgaveState represents the state of an Agave resource.
type AgaveState struct {
	AgaveArgs
}

// Create is the method that Pulumi calls to create an Agave resource.
// It sets up the Agave validator on the specified machine using the provided connection and flags.
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
func (Agave) Create(ctx context.Context, name string, input AgaveArgs, preview bool) (string, AgaveState, error) {
	state := AgaveState{AgaveArgs: input}

	if preview {
		return name, state, nil
	}

	client := &agave.Agave{
		Flags:    input.Flags,
		KeyPairs: input.KeyPairs,
	}
	command := client.Install()

	r := runner.Machine(input.Connection).
		Env(command.Env()).
		Script(command.Script())

	if err := r.Run(ctx); err != nil {
		return "", AgaveState{}, fmt.Errorf("failed to install validator: %w", err)
	}
	return name, state, nil
}
