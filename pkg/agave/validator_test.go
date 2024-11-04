package agave

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatorEnv(t *testing.T) {
	{
		v := ValidatorEnv{}

		if v.ToString() != "" {
			t.Error("empty validator did not return an empty stirng")
		}
	}

	{
		m := Metrics{
			URL:      "noproto://nowhere",
			Database: "nodb",
			User:     "notauser",
			Password: "notapassword",
		}

		v := ValidatorEnv{
			Metrics: &m,
		}

		if v.ToString() != `SOLANA_METRICS_CONFIG="host=noproto://nowhere,db=nodb,u=notauser,p=notapassword"` {
			t.Error("validator env output didn't match expectations")
		}
	}
}

func TestValidatorFlags(t *testing.T) {
	gossipHost := "somehost"
	tvuReceiveThreads := 88
	expectedGenesisHash := "asdkjasldjadk"
	fullRPCAPI := false
	noVoting := true
	allowPrivateAddr := true

	f := Flags{
		EntryPoint:                   &[]string{"hostA", "hostB"},
		KnownValidator:               &[]string{"validatorA", "validatorB"},
		UseSnapshotArchivesAtStartup: "whatever",
		RpcPort:                      38,
		PrivateRPC:                   false,
		OnlyKnownRPC:                 true,
		DynamicPortRange:             "8002-8003",
		GossipHost:                   &gossipHost,
		GossipPort:                   62,
		RpcBindAddress:               "1.2.3.4",
		WalRecoveryMode:              "someaddress",
		LimitLedgerSize:              62,
		BlockProductionMethod:        "whatever",
		TvuReceiveThreads:            &tvuReceiveThreads,
		NoWaitForVoteToStartLeader:   false,
		FullSnapshotIntervalSlots:    1000,
		ExpectedGenesisHash:          &expectedGenesisHash,
		FullRpcAPI:                   &fullRPCAPI,
		NoVoting:                     &noVoting,
		AllowPrivateAddr:             &allowPrivateAddr,
		ExtraFlags:                   &[]string{"--extra", "--flag"},
	}

	assert.Equal(t, f.ToArgs(), []string{"--identity /home/sol/validator-keypair.json",
		"--vote-account /home/sol/vote-account-keypair.json",
		"--entrypoint hostA",
		"--entrypoint hostB",
		"--known-validator validatorA",
		"--known-validator validatorB",
		"--expected-genesis-hash asdkjasldjadk",
		"--use-snapshot-archives-at-startup whatever",
		"--rpc-port 38",
		"--dynamic-port-range 8002-8003",
		"--gossip-host somehost",
		"--gossip-port 62",
		"--rpc-bind-address 1.2.3.4",
		"--wal-recovery-mode someaddress",
		"--log /home/sol/log",
		"--accounts /home/sol/accounts",
		"--ledger /home/sol/ledger",
		"--limit-ledger-size 62",
		"--block-production-method whatever",
		"--tvu-receive-threads 88",
		"--full-snapshot-interval-slots 1000",
		"",
		"--only-known-rpc",
		"",
		"",
		"--no-voting",
		"--allow-private-addr",
		"--extra",
		"--flag"})
}
