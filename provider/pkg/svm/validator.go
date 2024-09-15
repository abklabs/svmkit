package svm

import (
	"context"

	"github.com/abklabs/pulumi-svm/provider/pkg/agave"
)

type Validator struct{}

type ValidatorArgs struct {
	Variant  *string                 `pulumi:"variant,optional"`
	Flags    agave.ValidatorFlags    `pulumi:"flags"`
	KeyPairs agave.ValidatorKeyPairs `pulumi:"keyPairs" provider:"secret"`
}

type ValidatorState struct {
	ValidatorArgs
}

func (Validator) Create(ctx context.Context, name string, input ValidatorArgs, preview bool) (string, ValidatorState, error) {
	state := ValidatorState{ValidatorArgs: input}

	if preview {
		return name, state, nil
	}

	return name, state, nil
}
