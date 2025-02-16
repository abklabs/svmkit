package firedancer

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ptr[T any](in T) *T {
	return &in
}

func TestConfigEncode(t *testing.T) {
	expectedTOML := `name = "fd1"
user = "sol"
scratch_directory = "/tmp/scratch"
dynamic_port_range = "8000-9000"

[log]
  path = "/var/log/test.log"
  colorize = "auto"
  level_logfile = "info"
  level_stderr = "warn"
  level_flush = "error"

[reporting]
  solana_metrics_config = "http://metrics.solana.com"

[ledger]
  path = "/home/sol/ledger"
  accounts_path = "/home/sol/accounts"
  limit_size = 1000
  account_indexes = ["index1", "index2"]
  account_index_exclude_keys = ["key1", "key2"]
  account_index_include_keys = ["key3", "key4"]
  snapshot_archive_format = "tar"
  require_tower = true

[gossip]
  entrypoints = ["entry1", "entry2"]
  port_check = true
  port = 8080
  host = "localhost"

[rpc]
  port = 8899
  full_api = true
  private = false
  transaction_history = true
  extended_tx_metadata_storage = false
  only_known = true
  pubsub_enable_block_subscription = true
  pubsub_enable_vote_subscription = false
  bigtable_ledger_storage = false

[snapshots]
  incremental_snapshots = true
  full_snapshot_interval_slots = 100
  incremental_snapshot_interval_slots = 50
  maximum_full_snapshots_to_retain = 5
  maximum_incremental_snapshots_to_retain = 10
  minimum_snapshot_download_speed = 100
  path = "/var/snapshots"
  incremental_path = "/var/incremental_snapshots"

[consensus]
  identity_path = "/var/identity"
  vote_account_path = "/var/vote_account"
  authorized_voter_paths = ["/var/voter1", "/var/voter2"]
  snapshot_fetch = true
  genesis_fetch = false
  poh_speed_test = true
  expected_genesis_hash = "hash123"
  wait_for_supermajority_at_slot = 1000
  expected_bank_hash = "bankhash123"
  expected_shred_version = 1
  wait_for_vote_to_start_leader = true
  os_network_limits_test = false
  hard_fork_at_slots = ["slot1", "slot2"]
  known_validators = ["validator1", "validator2"]

[layout]
  affinity = "affinity1"
  agave_affinity = "agave1"
  net_tile_count = 10
  quic_tile_count = 5
  resolv_tile_count = 3
  verify_tile_count = 7
  bank_tile_count = 2
  shred_tile_count = 4

[hugetlbfs]
  mount_path = "/mnt/hugetlbfs"

[development]
  [development.gossip]
    allow_private_address = true
`

	config := &Config{
		Name:             ptr("fd1"),
		User:             ptr("sol"),
		ScratchDirectory: ptr("/tmp/scratch"),
		DynamicPortRange: ptr("8000-9000"),
		Log: &ConfigLog{
			Path:         ptr("/var/log/test.log"),
			Colorize:     ptr("auto"),
			LevelLogfile: ptr("info"),
			LevelStderr:  ptr("warn"),
			LevelFlush:   ptr("error"),
		},
		Reporting: &ConfigReporting{
			SolanaMetricsConfig: ptr("http://metrics.solana.com"),
		},
		Ledger: &ConfigLedger{
			Path:                    ptr("/home/sol/ledger"),
			AccountsPath:            ptr("/home/sol/accounts"),
			LimitSize:               ptr(1000),
			AccountIndexes:          &[]string{"index1", "index2"},
			AccountIndexExcludeKeys: &[]string{"key1", "key2"},
			AccountIndexIncludeKeys: &[]string{"key3", "key4"},
			SnapshotArchiveFormat:   ptr("tar"),
			RequireTower:            ptr(true),
		},
		Gossip: &ConfigGossip{
			Entrypoints: &[]string{"entry1", "entry2"},
			PortCheck:   ptr(true),
			Port:        ptr(8080),
			Host:        ptr("localhost"),
		},
		RPC: &ConfigRPC{
			Port:                          ptr(8899),
			FullAPI:                       ptr(true),
			Private:                       ptr(false),
			TransactionHistory:            ptr(true),
			ExtendedTxMetadataStorage:     ptr(false),
			OnlyKnown:                     ptr(true),
			PubsubEnableBlockSubscription: ptr(true),
			PubsubEnableVoteSubscription:  ptr(false),
			BigtableLedgerStorage:         ptr(false),
		},
		Snapshots: &ConfigSnapshots{
			IncrementalSnapshots:                ptr(true),
			FullSnapshotIntervalSlots:           ptr(100),
			IncrementalSnapshotIntervalSlots:    ptr(50),
			MaximumFullSnapshotsToRetain:        ptr(5),
			MaximumIncrementalSnapshotsToRetain: ptr(10),
			MinimumSnapshotDownloadSpeed:        ptr(100),
			Path:                                ptr("/var/snapshots"),
			IncrementalPath:                     ptr("/var/incremental_snapshots"),
		},
		Consensus: &ConfigConsensus{
			IdentityPath:               ptr("/var/identity"),
			VoteAccountPath:            ptr("/var/vote_account"),
			AuthorizedVoterPaths:       &[]string{"/var/voter1", "/var/voter2"},
			SnapshotFetch:              ptr(true),
			GenesisFetch:               ptr(false),
			PohSpeedTest:               ptr(true),
			ExpectedGenesisHash:        ptr("hash123"),
			WaitForSupermajorityAtSlot: ptr(1000),
			ExpectedBankHash:           ptr("bankhash123"),
			ExpectedShredVersion:       ptr(1),
			WaitForVoteToStartLeader:   ptr(true),
			OsNetworkLimitsTest:        ptr(false),
			HardForkAtSlots:            &[]string{"slot1", "slot2"},
			KnownValidators:            &[]string{"validator1", "validator2"},
		},
		Layout: &ConfigLayout{
			Affinity:        ptr("affinity1"),
			AgaveAffinity:   ptr("agave1"),
			NetTileCount:    ptr(10),
			QuicTileCount:   ptr(5),
			ResolvTileCount: ptr(3),
			VerifyTileCount: ptr(7),
			BankTileCount:   ptr(2),
			ShredTileCount:  ptr(4),
		},
		HugeTLBFS: &ConfigHugeTLBFS{
			MountPath: ptr("/mnt/hugetlbfs"),
		},
		ExtraConfig: &[]string{
			`
[development]
  [development.gossip]
    allow_private_address = true
`,
		},
	}

	var buf bytes.Buffer
	err := config.Encode(&buf)
	assert.NoError(t, err, "Failed to encode config")

	assert.Equal(t, expectedTOML, buf.String(), "Encoded TOML does not match expected")
}
