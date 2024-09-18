package svm

import (
	"github.com/abklabs/svmkit/pkg/module"
	"github.com/abklabs/svmkit/pkg/ssh"
)

// ValidatorArgs represents the arguments required to configure a validator.
type ValidatorArgs struct {
	// Connection holds the SSH connection details needed to access the validator.
	Connection ssh.Connection `pulumi:"connection"`

	// KeyPairs contains the key pairs used by the validator, marked as secret.
	KeyPairs module.KeyPairs `pulumi:"keyPairs" provider:"secret"`
}
