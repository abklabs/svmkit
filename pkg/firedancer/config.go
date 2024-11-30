package firedancer

import (
	"github.com/BurntSushi/toml"
	"io"
)

type ConfigLog struct {
	Path         *string `toml:"path,omitempty" pulumi:"path,optional"`
	Colorize     *string `toml:"colorize,omitempty" pulumi:"colorize,optional"`
	LevelLogfile *string `toml:"level_logfile,omitempty" pulumi:"levelLogfile,optional"`
	LevelStderr  *string `toml:"level_stderr,omitempty" pulumi:"levelStderr,optional"`
	LevelFlush   *string `toml:"level_flush,omitempty" pulumi:"levelFlush,optional"`
}

type ConfigReporting struct {
	SolanaMetricsConfig *string `toml:"solana_metrics_config,omitempty" pulumi:"solanaMetricsConfig,optional"`
}

type ConfigLedger struct {
	Path                    *string   `toml:"path,omitempty" pulumi:"path,optional"`
	AccountsPath            *string   `toml:"accounts_path,omitempty" pulumi:"accountsPath,optional"`
	LimitSize               *int      `toml:"limit_size,omitempty" pulumi:"limitSize,optional"`
	AccountIndexes          *[]string `toml:"account_indexes,omitempty" pulumi:"accountIndexes,optional"`
	AccountIndexExcludeKeys *[]string `toml:"account_index_exclude_keys,omitempty" pulumi:"accountIndexExcludeKeys,optional"`
	AccountIndexIncludeKeys *[]string `toml:"account_index_include_keys,omitempty" pulumi:"accountIndexIncludeKeys,optional"`
	SnapshotArchiveFormat   *string   `toml:"snapshot_archive_format,omitempty" pulumi:"snapshotArchiveFormat,optional"`
	RequireTower            *bool     `toml:"require_tower,omitempty" pulumi:"requireTower,optional"`
}

type ConfigGossip struct {
	Entrypoints *[]string `toml:"entrypoints,omitempty" pulumi:"entrypoints,optional"`
	PortCheck   *bool     `toml:"port_check,omitempty" pulumi:"portCheck,optional"`
	Port        *int      `toml:"port,omitempty" pulumi:"port,optional"`
	Host        *string   `toml:"host,omitempty" pulumi:"host,optional"`
}

type ConfigRPC struct {
	Port                          *int  `toml:"port,omitempty" pulumi:"port,optional"`
	FullAPI                       *bool `toml:"full_api,omitempty" pulumi:"fullApi,optional"`
	Private                       *bool `toml:"private,omitempty" pulumi:"private,optional"`
	TransactionHistory            *bool `toml:"transaction_history,omitempty" pulumi:"transactionHistory,optional"`
	ExtendedTxMetadataStorage     *bool `toml:"extended_tx_metadata_storage,omitempty" pulumi:"extendedTxMetadataStorage,optional"`
	OnlyKnown                     *bool `toml:"only_known,omitempty" pulumi:"onlyKnown,optional"`
	PubsubEnableBlockSubscription *bool `toml:"pubsub_enable_block_subscription,omitempty" pulumi:"pubsubEnableBlockSubscription,optional"`
	PubsubEnableVoteSubscription  *bool `toml:"pubsub_enable_vote_subscription,omitempty" pulumi:"pubsubEnableVoteSubscription,optional"`
	BigtableLedgerStorage         *bool `toml:"bigtable_ledger_storage,omitempty" pulumi:"bigtableLedgerStorage,optional"`
}

type ConfigSnapshots struct {
	IncrementalSnapshots                *bool   `toml:"incremental_snapshots,omitempty" pulumi:"incrementalSnapshots,optional"`
	FullSnapshotIntervalSlots           *int    `toml:"full_snapshot_interval_slots,omitempty" pulumi:"fullSnapshotIntervalSlots,optional"`
	IncrementalSnapshotIntervalSlots    *int    `toml:"incremental_snapshot_interval_slots,omitempty" pulumi:"incrementalSnapshotIntervalSlots,optional"`
	MaximumFullSnapshotsToRetain        *int    `toml:"maximum_full_snapshots_to_retain,omitempty" pulumi:"maximumFullSnapshotsToRetain,optional"`
	MaximumIncrementalSnapshotsToRetain *int    `toml:"maximum_incremental_snapshots_to_retain,omitempty" pulumi:"maximumIncrementalSnapshotsToRetain,optional"`
	MinimumSnapshotDownloadSpeed        *int    `toml:"minimum_snapshot_download_speed,omitempty" pulumi:"minimumSnapshotDownloadSpeed,optional"`
	Path                                *string `toml:"path,omitempty" pulumi:"path,optional"`
	IncrementalPath                     *string `toml:"incremental_path,omitempty" pulumi:"incrementalPath,optional"`
}

type ConfigConsensus struct {
	IdentityPath               *string   `toml:"identity_path,omitempty" pulumi:"identityPath,optional"`
	VoteAccountPath            *string   `toml:"vote_account_path,omitempty" pulumi:"voteAccountPath,optional"`
	AuthorizedVoterPaths       *[]string `toml:"authorized_voter_paths,omitempty" pulumi:"authorizedVoterPaths,optional"`
	SnapshotFetch              *bool     `toml:"snapshot_fetch,omitempty" pulumi:"snapshotFetch,optional"`
	GenesisFetch               *bool     `toml:"genesis_fetch,omitempty" pulumi:"genesisFetch,optional"`
	PohSpeedTest               *bool     `toml:"poh_speed_test,omitempty" pulumi:"pohSpeedTest,optional"`
	ExpectedGenesisHash        *string   `toml:"expected_genesis_hash,omitempty" pulumi:"expectedGenesisHash,optional"`
	WaitForSupermajorityAtSlot *int      `toml:"wait_for_supermajority_at_slot,omitempty" pulumi:"waitForSupermajorityAtSlot,optional"`
	ExpectedBankHash           *string   `toml:"expected_bank_hash,omitempty" pulumi:"expectedBankHash,optional"`
	ExpectedShredVersion       *int      `toml:"expected_shred_version,omitempty" pulumi:"expectedShredVersion,optional"`
	WaitForVoteToStartLeader   *bool     `toml:"wait_for_vote_to_start_leader,omitempty" pulumi:"waitForVoteToStartLeader,optional"`
	OsNetworkLimitsTest        *bool     `toml:"os_network_limits_test,omitempty" pulumi:"osNetworkLimitsTest,optional"`
	HardForkAtSlots            *[]string `toml:"hard_fork_at_slots,omitempty" pulumi:"hardForkAtSlots,optional"`
	KnownValidators            *[]string `toml:"known_validators,omitempty" pulumi:"knownValidators,optional"`
}

type ConfigLayout struct {
	Affinity        *string `toml:"affinity,omitempty" pulumi:"affinity,optional"`
	AgaveAffinity   *string `toml:"agave_affinity,omitempty" pulumi:"agaveAffinity,optional"`
	NetTileCount    *int    `toml:"net_tile_count,omitempty" pulumi:"netTileCount,optional"`
	QuicTileCount   *int    `toml:"quic_tile_count,omitempty" pulumi:"quicTileCount,optional"`
	ResolvTileCount *int    `toml:"resolv_tile_count,omitempty" pulumi:"resolvTileCount,optional"`
	VerifyTileCount *int    `toml:"verify_tile_count,omitempty" pulumi:"verifyTileCount,optional"`
	BankTileCount   *int    `toml:"bank_tile_count,omitempty" pulumi:"bankTileCount,optional"`
	ShredTileCount  *int    `toml:"shred_tile_count,omitempty" pulumi:"shredTileCount,optional"`
}

type ConfigHugeTLBFS struct {
	MountPath *string `toml:"mount_path,omitempty" pulumi:"mountPath,optional"`
}

type Config struct {
	Name             *string `toml:"name,omitempty" pulumi:"name,optional"`
	User             *string `toml:"user,omitempty" pulumi:"user,optional"`
	ScratchDirectory *string `toml:"scratch_directory,omitempty" pulumi:"scratchDirectory,optional"`
	DynamicPortRange *string `toml:"dynamic_port_range,omitempty" pulumi:"dynamicPortRange,optional"`

	Log       *ConfigLog       `toml:"log,omitempty" pulumi:"log,optional"`
	Reporting *ConfigReporting `toml:"reporting,omitempty" pulumi:"reporting,optional"`
	Ledger    *ConfigLedger    `toml:"ledger,omitempty" pulumi:"ledger,optional"`
	Gossip    *ConfigGossip    `toml:"gossip,omitempty" pulumi:"gossip,optional"`
	RPC       *ConfigRPC       `toml:"rpc,omitempty" pulumi:"rpc,optional"`
	Snapshots *ConfigSnapshots `toml:"snapshots,omitempty" pulumi:"snapshots,optional"`
	Consensus *ConfigConsensus `toml:"consensus,omitempty" pulumi:"consensus,optional"`
	Layout    *ConfigLayout    `toml:"layout,omitempty" pulumi:"layout,optional"`
	HugeTLBFS *ConfigHugeTLBFS `toml:"hugetlbfs,omitempty" pulumi:"hugetlbfs,optional"`

	ExtraConfig *[]string `pulumi:"extraConfig,optional"`
}

func (c *Config) Encode(w io.Writer) error {
	if err := toml.NewEncoder(w).Encode(c); err != nil {
		return err
	}

	if c.ExtraConfig != nil {
		for _, v := range *c.ExtraConfig {
			if _, err := w.Write([]byte(v)); err != nil {
				return err
			}
		}
	}

	return nil
}
