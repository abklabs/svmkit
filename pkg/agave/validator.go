package agave

import (
	_ "embed"
	"errors"
	"fmt"
	"strings"

	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/abklabs/svmkit/pkg/solana"
	"github.com/pulumi/pulumi-go-provider/infer"
)

const (
	accountsPath = "/home/sol/accounts"
	ledgerPath   = "/home/sol/ledger"
	logPath      = "/home/sol/log"

	identityKeyPairPath    = "/home/sol/validator-keypair.json"
	voteAccountKeyPairPath = "/home/sol/vote-account-keypair.json"
)

type Variant string

const (
	VariantSolana      Variant = "solana"
	VariantAgave       Variant = "agave"
	VariantPowerledger Variant = "powerledger"
	VariantJito        Variant = "jito"
	VariantPyth        Variant = "pyth"
	VariantMantis      Variant = "mantis"
	VariantXen         Variant = "xen"
)

func (Variant) Values() []infer.EnumValue[Variant] {
	return []infer.EnumValue[Variant]{
		{
			Name:        string(VariantSolana),
			Value:       VariantSolana,
			Description: "The Solana validator",
		},
		{
			Name:        string(VariantAgave),
			Value:       VariantAgave,
			Description: "The Agave validator",
		},
		{
			Name:        string(VariantPowerledger),
			Value:       VariantPowerledger,
			Description: "The Powerledger validator",
		},
		{
			Name:        string(VariantJito),
			Value:       VariantJito,
			Description: "The Jito validator",
		},
		{
			Name:        string(VariantPyth),
			Value:       VariantPyth,
			Description: "The Pyth validator",
		},
		{
			Name:        string(VariantMantis),
			Value:       VariantMantis,
			Description: "The Mantis validator",
		},
		{
			Name:        string(VariantXen),
			Value:       VariantXen,
			Description: "The Xen validator",
		},
	}
}

type KeyPairs struct {
	Identity    string `pulumi:"identity" provider:"secret"`
	VoteAccount string `pulumi:"voteAccount" provider:"secret"`
}

type Metrics struct {
	URL      string `pulumi:"url"`
	Database string `pulumi:"database"`
	User     string `pulumi:"user"`
	Password string `pulumi:"password"`
}

func (m *Metrics) Check() error {
	if m.URL == "" {
		return errors.New("metrics URL cannot be empty")
	}

	if m.Database == "" {
		return errors.New("metrics database cannot be empty")
	}

	if m.User == "" {
		return errors.New("metrics user cannot be empty")
	}

	return nil
}

// String constructs the Solana metrics configuration string from the separate fields
// and returns it as an environment variable string.
func (m *Metrics) String() string {

	// Note: We allow empty password as it might be a valid case in some scenarios
	configParts := []string{
		fmt.Sprintf("host=%s", m.URL),
		fmt.Sprintf("db=%s", m.Database),
		fmt.Sprintf("u=%s", m.User),
		fmt.Sprintf("p=%s", m.Password),
	}

	return strings.Join(configParts, ",")

}

type InstallCommand struct {
	Agave
}

func (cmd *InstallCommand) Check() error {
	if m := cmd.Metrics; m != nil {
		if err := m.Check(); err != nil {
			return fmt.Errorf("Warning: Invalid metrics URL: %v\n", err)
		}
	}

	return nil
}

func (cmd *InstallCommand) Env() *runner.EnvBuilder {
	validatorEnv := runner.NewEnvBuilder()

	if m := cmd.Metrics; m != nil {
		validatorEnv.Set("SOLANA_METRICS_CONFIG", m.String())
	}

	b := runner.NewEnvBuilder()

	b.SetMap(map[string]string{
		"VALIDATOR_FLAGS": strings.Join(cmd.Flags.ToArgs(), " "),
		"VALIDATOR_ENV":   validatorEnv.String(),
	})

	{
		s := identityKeyPairPath
		conf := solana.CLIConfig{
			KeyPair: &s,
		}

		if senv := cmd.Environment; senv != nil {
			conf.URL = senv.RPCURL
		}

		b.Set("SOLANA_CLI_CONFIG_FLAGS", conf.ToFlags().String())
	}

	if t := cmd.TimeoutConfig; t != nil {
		b.Merge(t.Env())
	}

	b.SetP("VALIDATOR_VERSION", cmd.Version)

	if cmd.Variant != nil {
		b.Set("VALIDATOR_VARIANT", string(*cmd.Variant))
	} else {
		b.Set("VALIDATOR_VARIANT", string(VariantAgave))
	}

	b.Set("RPC_BIND_ADDRESS", cmd.Flags.RpcBindAddress)
	b.SetInt("RPC_PORT", cmd.Flags.RpcPort)

	if cmd.Flags.FullRpcAPI != nil && *cmd.Flags.FullRpcAPI && cmd.StartupPolicy != nil {
		b.SetBoolP("WAIT_FOR_RPC_HEALTH", cmd.StartupPolicy.WaitForRPCHealth)
	}

	if i := cmd.Info; i != nil {
		b.Set("VALIDATOR_INFO_NAME", i.Name)
		b.SetP("VALIDATOR_INFO_WEBSITE", i.Website)
		b.SetP("VALIDATOR_INFO_ICON_URL", i.IconURL)
		b.SetP("VALIDATOR_INFO_DETAILS", i.Details)
	}

	if s := cmd.ShutdownPolicy; s != nil {
		b.SetArray("VALIDATOR_EXIT_FLAGS", s.ToFlags().ToArgs())
	}

	b.Set("LEDGER_PATH", ledgerPath)

	return b
}

func (cmd *InstallCommand) AddToPayload(p *runner.Payload) error {
	err := p.AddTemplate("steps.sh", installScriptTmpl, cmd)

	if err != nil {
		return err
	}

	if err := p.AddTemplate("check-validator", checkValidatorScriptTmpl, cmd); err != nil {
		return err
	}

	p.AddString("validator-keypair.json", cmd.KeyPairs.Identity)
	p.AddString("vote-account-keypair.json", cmd.KeyPairs.VoteAccount)

	return nil
}

type Agave struct {
	Environment    *solana.Environment   `pulumi:"environment,optional"`
	Version        *string               `pulumi:"version,optional"`
	Variant        *Variant              `pulumi:"variant,optional"`
	KeyPairs       KeyPairs              `pulumi:"keyPairs"`
	Flags          Flags                 `pulumi:"flags"`
	Metrics        *Metrics              `pulumi:"metrics,optional"`
	Info           *solana.ValidatorInfo `pulumi:"info,optional"`
	TimeoutConfig  *TimeoutConfig        `pulumi:"timeoutConfig,optional"`
	StartupPolicy  *StartupPolicy        `pulumi:"startupPolicy,optional"`
	ShutdownPolicy *ShutdownPolicy       `pulumi:"shutdownPolicy,optional"`
}

func (agave *Agave) Install() runner.Command {
	return &InstallCommand{
		Agave: *agave,
	}
}

type Flags struct {
	AccountIndex                        *[]string `pulumi:"accountIndex,optional"`
	AccountIndexExcludeKey              *[]string `pulumi:"accountIndexExcludeKey,optional"`
	AccountIndexIncludeKey              *[]string `pulumi:"accountIndexIncludeKey,optional"`
	AccountShrinkPath                   *[]string `pulumi:"accountShrinkPath,optional"`
	AccountsDbCacheLimitMb              *int      `pulumi:"accountsDbCacheLimitMb,optional"`
	AccountsDbTestHashCalculation       *bool     `pulumi:"accountsDbTestHashCalculation,optional"`
	AccountsHashCachePath               *string   `pulumi:"accountsHashCachePath,optional"`
	AccountsIndexBins                   *int      `pulumi:"accountsIndexBins,optional"`
	AccountsIndexPath                   *[]string `pulumi:"accountsIndexPath,optional"`
	AccountsIndexScanResultsLimitMb     *int      `pulumi:"accountsIndexScanResultsLimitMb,optional"`
	AccountsShrinkOptimizeTotalSpace    *bool     `pulumi:"accountsShrinkOptimizeTotalSpace,optional"`
	AccountsShrinkRatio                 *string   `pulumi:"accountsShrinkRatio,optional"`
	AllowPrivateAddr                    *bool     `pulumi:"allowPrivateAddr,optional"`
	AuthorizedVoter                     *[]string `pulumi:"authorizedVoter,optional"`
	BindAddress                         *string   `pulumi:"bindAddress,optional"`
	BlockProductionMethod               *string   `pulumi:"blockProductionMethod,optional"`
	BlockVerificationMethod             *string   `pulumi:"blockVerificationMethod,optional"`
	CheckVoteAccount                    *string   `pulumi:"checkVoteAccount,optional"`
	ContactDebugInterval                *int      `pulumi:"contactDebugInterval,optional"`
	Cuda                                *bool     `pulumi:"cuda,optional"`
	DebugKey                            *[]string `pulumi:"debugKey,optional"`
	DevHaltAtSlot                       *int      `pulumi:"devHaltAtSlot,optional"`
	DisableBankingTrace                 *bool     `pulumi:"disableBankingTrace,optional"`
	DynamicPortRange                    *string   `pulumi:"dynamicPortRange,optional"`
	EnableBankingTrace                  *int      `pulumi:"enableBankingTrace,optional"`
	EnableBigtableLedgerUpload          *bool     `pulumi:"enableBigtableLedgerUpload,optional"`
	EnableExtendedTxMetadataStorage     *bool     `pulumi:"enableExtendedTxMetadataStorage,optional"`
	EnableRpcBigtableLedgerStorage      *bool     `pulumi:"enableRpcBigtableLedgerStorage,optional"`
	EnableRpcTransactionHistory         *bool     `pulumi:"enableRpcTransactionHistory,optional"`
	EntryPoint                          *[]string `pulumi:"entryPoint,optional"`
	EtcdCacertFile                      *string   `pulumi:"etcdCacertFile,optional"`
	EtcdCertFile                        *string   `pulumi:"etcdCertFile,optional"`
	EtcdDomainName                      *string   `pulumi:"etcdDomainName,optional"`
	EtcdEndpoint                        *[]string `pulumi:"etcdEndpoint,optional"`
	EtcdKeyFile                         *string   `pulumi:"etcdKeyFile,optional"`
	ExpectedBankHash                    *string   `pulumi:"expectedBankHash,optional"`
	ExpectedGenesisHash                 *string   `pulumi:"expectedGenesisHash,optional"`
	ExpectedShredVersion                *int      `pulumi:"expectedShredVersion,optional"`
	ExtraFlags                          *[]string `pulumi:"extraFlags,optional"`
	FullRpcAPI                          *bool     `pulumi:"fullRpcAPI,optional"`
	FullSnapshotArchivePath             *string   `pulumi:"fullSnapshotArchivePath,optional"`
	FullSnapshotIntervalSlots           *int      `pulumi:"fullSnapshotIntervalSlots,optional"`
	GeyserPluginAlwaysEnabled           *bool     `pulumi:"geyserPluginAlwaysEnabled,optional"`
	GeyserPluginConfig                  *[]string `pulumi:"geyserPluginConfig,optional"`
	GossipHost                          *string   `pulumi:"gossipHost,optional"`
	GossipPort                          *int      `pulumi:"gossipPort,optional"`
	GossipValidator                     *[]string `pulumi:"gossipValidator,optional"`
	HardFork                            *[]int    `pulumi:"hardFork,optional"`
	HealthCheckSlotDistance             *int      `pulumi:"healthCheckSlotDistance,optional"`
	IncrementalSnapshotArchivePath      *string   `pulumi:"incrementalSnapshotArchivePath,optional"`
	InitCompleteFile                    *string   `pulumi:"initCompleteFile,optional"`
	KnownValidator                      *[]string `pulumi:"knownValidator,optional"`
	LimitLedgerSize                     *int      `pulumi:"limitLedgerSize,optional"`
	LogMessagesBytesLimit               *int      `pulumi:"logMessagesBytesLimit,optional"`
	MaxGenesisArchiveUnpackedSize       *int      `pulumi:"maxGenesisArchiveUnpackedSize,optional"`
	MaximumFullSnapshotsToRetain        *int      `pulumi:"maximumFullSnapshotsToRetain,optional"`
	MaximumIncrementalSnapshotsToRetain *int      `pulumi:"maximumIncrementalSnapshotsToRetain,optional"`
	MaximumLocalSnapshotAge             *int      `pulumi:"maximumLocalSnapshotAge,optional"`
	MaximumSnapshotDownloadAbort        *int      `pulumi:"maximumSnapshotDownloadAbort,optional"`
	MinimalSnapshotDownloadSpeed        *int      `pulumi:"minimalSnapshotDownloadSpeed,optional"`
	NoGenesisFetch                      *bool     `pulumi:"noGenesisFetch,optional"`
	NoIncrementalSnapshots              *bool     `pulumi:"noIncrementalSnapshots,optional"`
	NoSnapshotFetch                     *bool     `pulumi:"noSnapshotFetch,optional"`
	NoVoting                            *bool     `pulumi:"noVoting,optional"`
	NoWaitForVoteToStartLeader          bool      `pulumi:"noWaitForVoteToStartLeader"`
	OnlyKnownRPC                        *bool     `pulumi:"onlyKnownRPC,optional"`
	PrivateRPC                          *bool     `pulumi:"privateRPC,optional"`
	PublicRpcAddress                    *string   `pulumi:"publicRpcAddress,optional"`
	PublicTpuAddress                    *string   `pulumi:"publicTpuAddress,optional"`
	PublicTpuForwardsAddress            *string   `pulumi:"publicTpuForwardsAddress,optional"`
	RepairValidator                     *[]string `pulumi:"repairValidator,optional"`
	RequireTower                        *bool     `pulumi:"requireTower,optional"`
	RestrictedRepairOnlyMode            *bool     `pulumi:"restrictedRepairOnlyMode,optional"`
	RocksdbFifoShredStorageSize         *int      `pulumi:"rocksdbFifoShredStorageSize,optional"`
	RocksdbShredCompaction              *string   `pulumi:"rocksdbShredCompaction,optional"`
	RpcBigtableAppProfileId             *string   `pulumi:"rpcBigtableAppProfileId,optional"`
	RpcBigtableInstanceName             *string   `pulumi:"rpcBigtableInstanceName,optional"`
	RpcBigtableMaxMessageSize           *int      `pulumi:"rpcBigtableMaxMessageSize,optional"`
	RpcBigtableTimeout                  *int      `pulumi:"rpcBigtableTimeout,optional"`
	RpcBindAddress                      string    `pulumi:"rpcBindAddress"`
	RpcFaucetAddress                    *string   `pulumi:"rpcFaucetAddress,optional"`
	RpcMaxMultipleAccounts              *int      `pulumi:"rpcMaxMultipleAccounts,optional"`
	RpcMaxRequestBodySize               *int      `pulumi:"rpcMaxRequestBodySize,optional"`
	RpcNicenessAdjustment               *int      `pulumi:"rpcNicenessAdjustment,optional"`
	RpcPort                             int       `pulumi:"rpcPort"`
	RpcPubsubEnableBlockSubscription    *bool     `pulumi:"rpcPubsubEnableBlockSubscription,optional"`
	RpcPubsubEnableVoteSubscription     *bool     `pulumi:"rpcPubsubEnableVoteSubscription,optional"`
	RpcPubsubMaxActiveSubscriptions     *int      `pulumi:"rpcPubsubMaxActiveSubscriptions,optional"`
	RpcPubsubNotificationThreads        *int      `pulumi:"rpcPubsubNotificationThreads,optional"`
	RpcPubsubQueueCapacityBytes         *int      `pulumi:"rpcPubsubQueueCapacityBytes,optional"`
	RpcPubsubQueueCapacityItems         *int      `pulumi:"rpcPubsubQueueCapacityItems,optional"`
	RpcPubsubWorkerThreads              *int      `pulumi:"rpcPubsubWorkerThreads,optional"`
	RpcScanAndFixRoots                  *bool     `pulumi:"rpcScanAndFixRoots,optional"`
	RpcSendLeaderCount                  *int      `pulumi:"rpcSendLeaderCount,optional"`
	RpcSendRetryMs                      *int      `pulumi:"rpcSendRetryMs,optional"`
	RpcSendServiceMaxRetries            *int      `pulumi:"rpcSendServiceMaxRetries,optional"`
	RpcSendTransactionAlsoLeader        *bool     `pulumi:"rpcSendTransactionAlsoLeader,optional"`
	RpcSendTransactionRetryPoolMaxSize  *int      `pulumi:"rpcSendTransactionRetryPoolMaxSize,optional"`
	RpcSendTransactionTpuPeer           *[]string `pulumi:"rpcSendTransactionTpuPeer,optional"`
	RpcThreads                          *int      `pulumi:"rpcThreads,optional"`
	SkipPreflightHealthCheck            *bool     `pulumi:"skipPreflightHealthCheck,optional"`
	SkipSeedPhraseValidation            *bool     `pulumi:"skipSeedPhraseValidation,optional"`
	SkipStartupLedgerVerification       *bool     `pulumi:"skipStartupLedgerVerification,optional"`
	SnapshotArchiveFormat               *string   `pulumi:"snapshotArchiveFormat,optional"`
	SnapshotIntervalSlots               *int      `pulumi:"snapshotIntervalSlots,optional"`
	SnapshotPackagerNicenessAdjustment  *int      `pulumi:"snapshotPackagerNicenessAdjustment,optional"`
	SnapshotVersion                     *string   `pulumi:"snapshotVersion,optional"`
	StakedNodesOverrides                *string   `pulumi:"stakedNodesOverrides,optional"`
	TowerStorage                        *string   `pulumi:"towerStorage,optional"`
	TpuCoalesceMs                       *int      `pulumi:"tpuCoalesceMs,optional"`
	TpuConnectionPoolSize               *int      `pulumi:"tpuConnectionPoolSize,optional"`
	TpuDisableQuic                      *bool     `pulumi:"tpuDisableQuic,optional"`
	TpuEnableUdp                        *bool     `pulumi:"tpuEnableUdp,optional"`
	TvuReceiveThreads                   *int      `pulumi:"tvuReceiveThreads,optional"`
	UnifiedSchedulerHandlerThreads      *int      `pulumi:"unifiedSchedulerHandlerThreads,optional"`
	UseSnapshotArchivesAtStartup        *string   `pulumi:"useSnapshotArchivesAtStartup,optional"`
	WaitForSupermajority                *int      `pulumi:"waitForSupermajority,optional"`
	WalRecoveryMode                     string    `pulumi:"walRecoveryMode"`
}

func (f Flags) ToArgs() []string {
	b := runner.FlagBuilder{}

	// Note: These locations are hard coded inside asset-builder.
	b.Append("--identity", identityKeyPairPath)
	b.Append("--vote-account", voteAccountKeyPairPath)
	b.Append("--log", logPath)
	b.Append("--accounts", accountsPath)
	b.Append("--ledger", ledgerPath)

	if f.AccountIndex != nil {
		for _, index := range *f.AccountIndex {
			b.AppendP("account-index", &index)
		}
	}

	if f.AccountIndexExcludeKey != nil {
		for _, key := range *f.AccountIndexExcludeKey {
			b.AppendP("account-index-exclude-key", &key)
		}
	}

	if f.AccountIndexIncludeKey != nil {
		for _, key := range *f.AccountIndexIncludeKey {
			b.AppendP("account-index-include-key", &key)
		}
	}

	if f.AccountShrinkPath != nil {
		for _, path := range *f.AccountShrinkPath {
			b.AppendP("account-shrink-path", &path)
		}
	}

	b.AppendIntP("accounts-db-cache-limit-mb", f.AccountsDbCacheLimitMb)
	b.AppendBoolP("accounts-db-test-hash-calculation", f.AccountsDbTestHashCalculation)
	b.AppendP("accounts-hash-cache-path", f.AccountsHashCachePath)
	b.AppendIntP("accounts-index-bins", f.AccountsIndexBins)

	if f.AccountsIndexPath != nil {
		for _, path := range *f.AccountsIndexPath {
			b.AppendP("accounts-index-path", &path)
		}
	}

	b.AppendIntP("accounts-index-scan-results-limit-mb", f.AccountsIndexScanResultsLimitMb)
	b.AppendBoolP("accounts-shrink-optimize-total-space", f.AccountsShrinkOptimizeTotalSpace)
	b.AppendP("accounts-shrink-ratio", f.AccountsShrinkRatio)

	b.AppendBoolP("allow-private-addr", f.AllowPrivateAddr)

	if f.AuthorizedVoter != nil {
		for _, voter := range *f.AuthorizedVoter {
			b.AppendP("authorized-voter", &voter)
		}
	}

	b.AppendP("bind-address", f.BindAddress)
	b.AppendP("block-production-method", f.BlockProductionMethod)
	b.AppendP("block-verification-method", f.BlockVerificationMethod)
	b.AppendP("check-vote-account", f.CheckVoteAccount)
	b.AppendIntP("contact-debug-interval", f.ContactDebugInterval)
	b.AppendBoolP("cuda", f.Cuda)

	if f.DebugKey != nil {
		for _, key := range *f.DebugKey {
			b.AppendP("debug-key", &key)
		}
	}

	b.AppendIntP("dev-halt-at-slot", f.DevHaltAtSlot)
	b.AppendBoolP("disable-banking-trace", f.DisableBankingTrace)
	b.AppendP("dynamic-port-range", f.DynamicPortRange)
	b.AppendIntP("enable-banking-trace", f.EnableBankingTrace)
	b.AppendBoolP("enable-bigtable-ledger-upload", f.EnableBigtableLedgerUpload)
	b.AppendBoolP("enable-extended-tx-metadata-storage", f.EnableExtendedTxMetadataStorage)
	b.AppendBoolP("enable-rpc-bigtable-ledger-storage", f.EnableRpcBigtableLedgerStorage)
	b.AppendBoolP("enable-rpc-transaction-history", f.EnableRpcTransactionHistory)

	if f.EntryPoint != nil {
		for _, entrypoint := range *f.EntryPoint {
			b.AppendP("entrypoint", &entrypoint)
		}
	}

	b.AppendP("etcd-cacert-file", f.EtcdCacertFile)
	b.AppendP("etcd-cert-file", f.EtcdCertFile)
	b.AppendP("etcd-domain-name", f.EtcdDomainName)

	if f.EtcdEndpoint != nil {
		for _, endpoint := range *f.EtcdEndpoint {
			b.AppendP("etcd-endpoint", &endpoint)
		}
	}

	b.AppendP("etcd-key-file", f.EtcdKeyFile)
	b.AppendP("expected-bank-hash", f.ExpectedBankHash)
	b.AppendP("expected-genesis-hash", f.ExpectedGenesisHash)
	b.AppendIntP("expected-shred-version", f.ExpectedShredVersion)

	if f.ExtraFlags != nil {
		b.Append(*f.ExtraFlags...)
	}

	b.AppendBoolP("full-rpc-api", f.FullRpcAPI)
	b.AppendP("full-snapshot-archive-path", f.FullSnapshotArchivePath)
	b.AppendIntP("full-snapshot-interval-slots", f.FullSnapshotIntervalSlots)
	b.AppendBoolP("geyser-plugin-always-enabled", f.GeyserPluginAlwaysEnabled)

	if f.GeyserPluginConfig != nil {
		for _, config := range *f.GeyserPluginConfig {
			b.AppendP("geyser-plugin-config", &config)
		}
	}

	b.AppendP("gossip-host", f.GossipHost)
	b.AppendIntP("gossip-port", f.GossipPort)

	if f.GossipValidator != nil {
		for _, validator := range *f.GossipValidator {
			b.AppendP("gossip-validator", &validator)
		}
	}

	if f.HardFork != nil {
		for _, slot := range *f.HardFork {
			b.AppendIntP("hard-fork", &slot)
		}
	}

	b.AppendIntP("health-check-slot-distance", f.HealthCheckSlotDistance)
	b.AppendP("incremental-snapshot-archive-path", f.IncrementalSnapshotArchivePath)
	b.AppendP("init-complete-file", f.InitCompleteFile)

	if f.KnownValidator != nil {
		for _, knownValidator := range *f.KnownValidator {
			b.AppendP("known-validator", &knownValidator)
		}
	}

	b.AppendIntP("limit-ledger-size", f.LimitLedgerSize)
	b.AppendIntP("log-messages-bytes-limit", f.LogMessagesBytesLimit)
	b.AppendIntP("max-genesis-archive-unpacked-size", f.MaxGenesisArchiveUnpackedSize)
	b.AppendIntP("maximum-full-snapshots-to-retain", f.MaximumFullSnapshotsToRetain)
	b.AppendIntP("maximum-incremental-snapshots-to-retain", f.MaximumIncrementalSnapshotsToRetain)
	b.AppendIntP("maximum-local-snapshot-age", f.MaximumLocalSnapshotAge)
	b.AppendIntP("maximum-snapshot-download-abort", f.MaximumSnapshotDownloadAbort)
	b.AppendIntP("minimal-snapshot-download-speed", f.MinimalSnapshotDownloadSpeed)
	b.AppendBoolP("no-genesis-fetch", f.NoGenesisFetch)
	b.AppendBoolP("no-incremental-snapshots", f.NoIncrementalSnapshots)
	b.AppendBoolP("no-snapshot-fetch", f.NoSnapshotFetch)
	b.AppendBoolP("no-voting", f.NoVoting)

	// Note: This flag is not documented in the Agave validator documentation, but it is
	// present in the source code.
	b.AppendBoolP("no-wait-for-vote-to-start-leader", &f.NoWaitForVoteToStartLeader)

	b.AppendBoolP("only-known-rpc", f.OnlyKnownRPC)
	b.AppendBoolP("private-rpc", f.PrivateRPC)
	b.AppendP("public-rpc-address", f.PublicRpcAddress)
	b.AppendP("public-tpu-address", f.PublicTpuAddress)
	b.AppendP("public-tpu-forwards-address", f.PublicTpuForwardsAddress)

	if f.RepairValidator != nil {
		for _, validator := range *f.RepairValidator {
			b.AppendP("repair-validator", &validator)
		}
	}

	b.AppendBoolP("require-tower", f.RequireTower)
	b.AppendBoolP("restricted-repair-only-mode", f.RestrictedRepairOnlyMode)
	b.AppendIntP("rocksdb-fifo-shred-storage-size", f.RocksdbFifoShredStorageSize)
	b.AppendP("rocksdb-shred-compaction", f.RocksdbShredCompaction)
	b.AppendP("rpc-bigtable-app-profile-id", f.RpcBigtableAppProfileId)
	b.AppendP("rpc-bigtable-instance-name", f.RpcBigtableInstanceName)
	b.AppendIntP("rpc-bigtable-max-message-size", f.RpcBigtableMaxMessageSize)
	b.AppendIntP("rpc-bigtable-timeout", f.RpcBigtableTimeout)
	b.AppendP("rpc-bind-address", &f.RpcBindAddress)
	b.AppendP("rpc-faucet-address", f.RpcFaucetAddress)
	b.AppendIntP("rpc-max-multiple-accounts", f.RpcMaxMultipleAccounts)
	b.AppendIntP("rpc-max-request-body-size", f.RpcMaxRequestBodySize)
	b.AppendIntP("rpc-niceness-adjustment", f.RpcNicenessAdjustment)
	b.AppendIntP("rpc-port", &f.RpcPort)
	b.AppendBoolP("rpc-pubsub-enable-block-subscription", f.RpcPubsubEnableBlockSubscription)
	b.AppendBoolP("rpc-pubsub-enable-vote-subscription", f.RpcPubsubEnableVoteSubscription)
	b.AppendIntP("rpc-pubsub-max-active-subscriptions", f.RpcPubsubMaxActiveSubscriptions)
	b.AppendIntP("rpc-pubsub-notification-threads", f.RpcPubsubNotificationThreads)
	b.AppendIntP("rpc-pubsub-queue-capacity-bytes", f.RpcPubsubQueueCapacityBytes)
	b.AppendIntP("rpc-pubsub-queue-capacity-items", f.RpcPubsubQueueCapacityItems)
	b.AppendIntP("rpc-pubsub-worker-threads", f.RpcPubsubWorkerThreads)
	b.AppendBoolP("rpc-scan-and-fix-roots", f.RpcScanAndFixRoots)
	b.AppendIntP("rpc-send-leader-count", f.RpcSendLeaderCount)
	b.AppendIntP("rpc-send-retry-ms", f.RpcSendRetryMs)
	b.AppendIntP("rpc-send-service-max-retries", f.RpcSendServiceMaxRetries)
	b.AppendBoolP("rpc-send-transaction-also-leader", f.RpcSendTransactionAlsoLeader)
	b.AppendIntP("rpc-send-transaction-retry-pool-max-size", f.RpcSendTransactionRetryPoolMaxSize)

	if f.RpcSendTransactionTpuPeer != nil {
		for _, peer := range *f.RpcSendTransactionTpuPeer {
			b.AppendP("rpc-send-transaction-tpu-peer", &peer)
		}
	}

	b.AppendIntP("rpc-threads", f.RpcThreads)
	b.AppendBoolP("skip-preflight-health-check", f.SkipPreflightHealthCheck)
	b.AppendBoolP("skip-seed-phrase-validation", f.SkipSeedPhraseValidation)
	b.AppendBoolP("skip-startup-ledger-verification", f.SkipStartupLedgerVerification)
	b.AppendP("snapshot-archive-format", f.SnapshotArchiveFormat)
	b.AppendIntP("snapshot-interval-slots", f.SnapshotIntervalSlots)
	b.AppendIntP("snapshot-packager-niceness-adjustment", f.SnapshotPackagerNicenessAdjustment)
	b.AppendP("snapshot-version", f.SnapshotVersion)
	b.AppendP("staked-nodes-overrides", f.StakedNodesOverrides)
	b.AppendP("tower-storage", f.TowerStorage)
	b.AppendIntP("tpu-coalesce-ms", f.TpuCoalesceMs)
	b.AppendIntP("tpu-connection-pool-size", f.TpuConnectionPoolSize)
	b.AppendBoolP("tpu-disable-quic", f.TpuDisableQuic)
	b.AppendBoolP("tpu-enable-udp", f.TpuEnableUdp)
	b.AppendIntP("tvu-receive-threads", f.TvuReceiveThreads)
	b.AppendIntP("unified-scheduler-handler-threads", f.UnifiedSchedulerHandlerThreads)
	b.AppendP("use-snapshot-archives-at-startup", f.UseSnapshotArchivesAtStartup)
	b.AppendIntP("wait-for-supermajority", f.WaitForSupermajority)
	b.AppendP("wal-recovery-mode", &f.WalRecoveryMode)

	return b.ToArgs()
}
