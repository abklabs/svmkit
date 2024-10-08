package validator

import "github.com/abklabs/svmkit/pkg/runner"

// KeyPairs is a map of key pairs for the blockchain validator.
type KeyPairs map[string]string

// Version is the semantic version (semver) of the blockchain validator apt package.
type Version *string

// Client is an interface for managing the blockchain validator.
type Client interface {
	// Install returns a Command to install the blockchain validator.
	Install() runner.Command
}

// ClientFlags is an interface for client flags for the validator.
type ClientFlags interface {
	// ToArgs converts client flags to a slice of arguments.
	ToArgs() []string

	// S adds a string flag to the arguments.
	S(k string, v interface{}) string
	// B adds a boolean flag to the arguments.
	B(k string, v bool) string
}
