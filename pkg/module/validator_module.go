package module

import "github.com/abklabs/svmkit/pkg/runner"

type KeyPairs map[string]string

// Validator is an interface for managing the blockchain validator.
type Validator interface {
	// Install returns a Command to install the blockchain validator.
	Install() runner.Command
}

// ClientFlags is an interface for client flags for the validator.
type ValidatorFlags interface {
	// ToArgs converts client flags to a slice of arguments.
	ToArgs() []string

	// S adds a string flag to the arguments.
	S(k string, v interface{}) string
	// B adds a boolean flag to the arguments.
	B(k string, v bool) string
}
