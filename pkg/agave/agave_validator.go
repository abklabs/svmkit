package agave

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/abklabs/svmkit/pkg/module"
	"github.com/abklabs/svmkit/pkg/runner"
)

const (
	accountsPath = "/home/sol/accounts"
	ledgerPath   = "/home/sol/ledger"
	logPath      = "/home/sol/log"
)

type KeyPairs struct {
	Identity    string `pulumi:"identity" provider:"secret"`
	VoteAccount string `pulumi:"voteAccount" provider:"secret"`
}

type InstallCommand struct {
	runner.Command
	Flags    Flags
	KeyPairs KeyPairs
}

func (cmd *InstallCommand) Env() map[string]string {
	return map[string]string{
		"VALIDATOR_FLAGS":      strings.Join(cmd.Flags.toArgs(), " "),
		"IDENTITY_KEYPAIR":     cmd.KeyPairs.Identity,
		"VOTE_ACCOUNT_KEYPAIR": cmd.KeyPairs.VoteAccount,
	}
}

func (cmd *InstallCommand) Script() string {
	return InstallScript
}

type ValidatorPaths struct {
	Accounts string `pulumi:"accounts"`
	Ledger   string `pulumi:"ledger"`
	Log      string `pulumi:"log"`
}

type Agave struct {
	module.Validator
	KeyPairs KeyPairs
	Flags    Flags
}

func (agave *Agave) Install() runner.Command {
	return &InstallCommand{
		Flags:    agave.Flags,
		KeyPairs: agave.KeyPairs,
	}
}

type Flags struct {
	module.ValidatorFlags
	EntryPoint                   *[]string `pulumi:"entryPoint,optional"`
	KnownValidator               *[]string `pulumi:"knownValidator,optional"`
	UseSnapshotArchivesAtStartup string    `pulumi:"useSnapshotArchivesAtStartup"`
	RpcPort                      int       `pulumi:"rpcPort"`
	PrivateRPC                   bool      `pulumi:"privateRPC"`
	OnlyKnownRPC                 bool      `pulumi:"onlyKnownRPC"`
	DynamicPortRange             string    `pulumi:"dynamicPortRange"`
	GossipPort                   int       `pulumi:"gossipPort"`
	RpcBindAddress               string    `pulumi:"rpcBindAddress"`
	WalRecoveryMode              string    `pulumi:"walRecoveryMode"`
	LimitLedgerSize              int       `pulumi:"limitLedgerSize"`
	BlockProductionMethod        string    `pulumi:"blockProductionMethod"`
	TvuReceiveThreads            *int      `pulumi:"tvuReceiveThreads,optional"`
	NoWaitForVoteToStartLeader   bool      `pulumi:"noWaitForVoteToStartLeader"`
	FullSnapshotIntervalSlots    int       `pulumi:"fullSnapshotIntervalSlots"`
	ExpectedGenesisHash          *string   `pulumi:"expectedGenesisHash,optional"`
	FullRpcAPI                   *bool     `pulumi:"fullRpcAPI,optional"`
	NoVoting                     *bool     `pulumi:"noVoting,optional"`
}

func (f Flags) toArgs() []string {
	var l []string

	// Note: These locations are hard coded inside asset-builder.
	l = append(l, f.S("identity", "/home/sol/validator-keypair.json"))
	l = append(l, f.S("vote-account", "/home/sol/vote-account-keypair.json"))

	if f.EntryPoint != nil {
		for _, entrypoint := range *f.EntryPoint {
			l = append(l, f.S("entrypoint", entrypoint))
		}
	}

	if f.KnownValidator != nil {
		for _, knownValidator := range *f.KnownValidator {
			l = append(l, f.S("known-validator", knownValidator))
		}
	}

	if f.ExpectedGenesisHash != nil {
		l = append(l, f.S("expected-genesis-hash", *f.ExpectedGenesisHash))
	}
	l = append(l, f.S("use-snapshot-archives-at-startup", f.UseSnapshotArchivesAtStartup))
	l = append(l, f.S("rpc-port", f.RpcPort))
	l = append(l, f.S("dynamic-port-range", f.DynamicPortRange))
	l = append(l, f.S("gossip-port", f.GossipPort))
	l = append(l, f.S("rpc-bind-address", f.RpcBindAddress))
	l = append(l, f.S("wal-recovery-mode", f.WalRecoveryMode))
	l = append(l, f.S("log", logPath))
	l = append(l, f.S("accounts", accountsPath))
	l = append(l, f.S("ledger", ledgerPath))
	l = append(l, f.S("limit-ledger-size", f.LimitLedgerSize))
	l = append(l, f.S("block-production-method", f.BlockProductionMethod))
	if f.TvuReceiveThreads != nil {
		l = append(l, f.S("tvu-receive-threads", f.TvuReceiveThreads))
	}
	l = append(l, f.S("full-snapshot-interval-slots", f.FullSnapshotIntervalSlots))
	l = append(l, f.B("no-wait-for-vote-to-start-leader", f.NoWaitForVoteToStartLeader))
	l = append(l, f.B("only-known-rpc", f.OnlyKnownRPC))
	l = append(l, f.B("private-rpc", f.PrivateRPC))

	if f.FullRpcAPI != nil {
		l = append(l, f.B("full-rpc-api", *f.FullRpcAPI))
	}

	if f.NoVoting != nil {
		l = append(l, f.B("no-voting", *f.NoVoting))
	}

	return l
}

func (Flags) S(k string, v interface{}) string {
	return fmt.Sprintf("--%s %v", k, v)
}

func (Flags) B(k string, v bool) string {
	if v {
		return fmt.Sprintf("--%s", k)
	}
	return ""
}
