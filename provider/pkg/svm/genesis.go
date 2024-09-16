package svm

import (
	"github.com/abklabs/pulumi-svm/provider/pkg/ssh"
)

// GenesisArgs represents the arguments required to configure a genesis ledger.
type GenesisArgs struct {
	// Connection holds the SSH connection details needed to access the genesis ledger.
	Connection ssh.Connection `pulumi:"connection"`
}
