package agave

import "fmt"

type ValidatorKeyPairs struct {
	Identity    string `pulumi:"identity" provider:"secret"`
	VoteAccount string `pulumi:"voteAccount" provider:"secret"`
}

type ValidatorPaths struct {
	Accounts string `pulumi:"accounts"`
	Ledger   string `pulumi:"ledger"`
	Log      string `pulumi:"log"`
}

type ValidatorFlags struct {
	EntryPoint                   *[]string      `pulumi:"entryPoint,optional"`
	KnownValidator               *[]string      `pulumi:"knownValidator,optional"`
	UseSnapshotArchivesAtStartup string         `pulumi:"useSnapshotArchivesAtStartup"`
	RpcPort                      int            `pulumi:"rpcPort"`
	PrivateRPC                   bool           `pulumi:"privateRPC"`
	OnlyKnownRPC                 bool           `pulumi:"onlyKnownRPC"`
	DynamicPortRange             string         `pulumi:"dynamicPortRange"`
	GossipPort                   int            `pulumi:"gossipPort"`
	RpcBindAddress               string         `pulumi:"rpcBindAddress"`
	WalRecoveryMode              string         `pulumi:"walRecoveryMode"`
	LimitLedgerSize              int            `pulumi:"limitLedgerSize"`
	BlockProductionMethod        string         `pulumi:"blockProductionMethod"`
	TvuReceiveThreads            *int           `pulumi:"tvuReceiveThreads,optional"`
	NoWaitForVoteToStartLeader   bool           `pulumi:"noWaitForVoteToStartLeader"`
	FullSnapshotIntervalSlots    int            `pulumi:"fullSnapshotIntervalSlots"`
	ExpectedGenesisHash          *string        `pulumi:"expectedGenesisHash,optional"`
	FullRpcAPI                   *bool          `pulumi:"fullRpcAPI,optional"`
	NoVoting                     *bool          `pulumi:"noVoting,optional"`
	Paths                        ValidatorPaths `pulumi:"paths"`
}

func Flags(flags ValidatorFlags) []string {
	var l []string

	b := func(k string, v bool) {
		if v {
			l = append(l, "--"+k)
		}
	}

	s := func(k string, v interface{}) {
		if v != nil {
			l = append(l, "--"+k, fmt.Sprintf("%v", v))
		}
	}

	// Note: These locations are hard coded inside asset-builder.
	s("identity", "/home/sol/validator-keypair.json")
	s("vote-account", "/home/sol/vote-account-keypair.json")

	if flags.EntryPoint != nil {
		for _, entrypoint := range *flags.EntryPoint {
			s("entrypoint", entrypoint)
		}

	}

	if flags.KnownValidator != nil {
		for _, knownValidator := range *flags.KnownValidator {
			s("known-validator", knownValidator)
		}
	}

	s("expected-genesis-hash", flags.ExpectedGenesisHash)
	s("use-snapshot-archives-at-startup", flags.UseSnapshotArchivesAtStartup)
	s("rpc-port", flags.RpcPort)
	s("dynamic-port-range", flags.DynamicPortRange)
	s("gossip-port", flags.GossipPort)
	s("rpc-bind-address", flags.RpcBindAddress)
	s("wal-recovery-mode", flags.WalRecoveryMode)
	s("log", flags.Paths.Log)
	s("accounts", flags.Paths.Accounts)
	s("ledger", flags.Paths.Ledger)
	s("limit-ledger-size", flags.LimitLedgerSize)
	s("block-production-method", flags.BlockProductionMethod)
	s("tvu-receive-threads", flags.TvuReceiveThreads)
	s("full-snapshot-interval-slots", flags.FullSnapshotIntervalSlots)

	b("no-wait-for-vote-to-start-leader", flags.NoWaitForVoteToStartLeader)
	b("only-known-rpc", flags.OnlyKnownRPC)
	b("private-rpc", flags.PrivateRPC)
	b("full-rpc-api", *flags.FullRpcAPI)
	b("no-voting", *flags.NoVoting)

	return l
}
