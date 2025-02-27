package agave

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatorEnv(t *testing.T) {
	m := Metrics{
		URL:      "noproto://nowhere",
		Database: "nodb",
		User:     "notauser",
		Password: "notapassword",
	}

	if m.String() != `host=noproto://nowhere,db=nodb,u=notauser,p=notapassword` {
		t.Error("validator env output didn't match expectations")
	}
}

func TestValidatorFlags(t *testing.T) {

	paths, err := NewDefaultAgavePaths(nil)

	if err != nil {
		t.Fatal(err)
	}

	accountIndex := []string{"program-id", "spl-token-owner"}
	accountIndexExcludeKey := []string{"excludeKey1", "excludeKey2"}
	accountIndexIncludeKey := []string{"includeKey1", "includeKey2"}
	accountShrinkPath := []string{"/path/to/shrink1", "/path/to/shrink2"}
	accountsDbCacheLimitMb := 1024
	accountsDbTestHashCalculation := true
	accountsHashCachePath := "/path/to/accounts_hash_cache"
	accountsIndexBins := 64
	accountsIndexPath := []string{"/path/to/accounts_index1", "/path/to/accounts_index2"}
	accountsIndexScanResultsLimitMb := 512
	accountsShrinkOptimizeTotalSpace := true
	accountsShrinkRatio := "0.8"
	allowPrivateAddr := true
	authorizedVoter := []string{"voterPubkey1", "voterPubkey2"}
	bindAddress := "0.0.0.0"
	blockProductionMethod := "whatever"
	blockVerificationMethod := "unified-scheduler"
	checkVoteAccount := "http://localhost:8899"
	contactDebugInterval := 120000
	cuda := true
	debugKey := []string{"debugKey1", "debugKey2"}
	devHaltAtSlot := 1000000
	disableBankingTrace := true
	dynamicPortRange := "8002-8003"
	enableBankingTrace := 15032385536
	enableBigtableLedgerUpload := true
	enableExtendedTxMetadataStorage := true
	enableRpcBigtableLedgerStorage := true
	enableRpcTransactionHistory := true
	entryPoint := []string{"hostA", "hostB"}
	etcdCacertFile := "/path/to/ca.crt"
	etcdCertFile := "/path/to/client.crt"
	etcdDomainName := "example.com"
	etcdEndpoint := []string{"etcd1.example.com:2379", "etcd2.example.com:2379"}
	etcdKeyFile := "/path/to/client.key"
	expectedBankHash := "expectedBankHashValue"
	expectedGenesisHash := "asdkjasldjadk"
	expectedShredVersion := 1234
	extraFlags := []string{"--extra", "--flag"}
	fullRpcAPI := true
	fullSnapshotArchivePath := "/path/to/full_snapshot"
	fullSnapshotIntervalSlots := 1000
	geyserPluginAlwaysEnabled := true
	geyserPluginConfig := []string{"/path/to/geyser_plugin1.yaml", "/path/to/geyser_plugin2.yaml"}
	gossipHost := "somehost"
	gossipPort := 62
	gossipValidator := []string{"validatorPubkey1", "validatorPubkey2"}
	hardFork := []int{500000, 1000000}
	healthCheckSlotDistance := 128
	incrementalSnapshotArchivePath := "/path/to/incremental_snapshot"
	initCompleteFile := "/path/to/init_complete"
	knownValidator := []string{"validatorA", "validatorB"}
	limitLedgerSize := 62
	log := "/home/sol/log"
	logMessagesBytesLimit := 1048576
	maxGenesisArchiveUnpackedSize := 10485760
	maximumFullSnapshotsToRetain := 2
	maximumIncrementalSnapshotsToRetain := 4
	maximumLocalSnapshotAge := 2500
	maximumSnapshotDownloadAbort := 5
	minimalSnapshotDownloadSpeed := 10485760
	noGenesisFetch := true
	noIncrementalSnapshots := true
	noSnapshotFetch := true
	noVoting := true
	noWaitForVoteToStartLeader := true
	onlyKnownRPC := true
	privateRPC := true
	publicRpcAddress := "0.0.0.0:8899"
	publicTpuAddress := "0.0.0.0:8001"
	publicTpuForwardsAddress := "0.0.0.0:8002"
	repairValidator := []string{"repairValidator1", "repairValidator2"}
	requireTower := true
	restrictedRepairOnlyMode := true
	rocksdbFifoShredStorageSize := 104857600
	rocksdbShredCompaction := "fifo"
	rpcBigtableAppProfileId := "default"
	rpcBigtableInstanceName := "solana-ledger"
	rpcBigtableMaxMessageSize := 67108864
	rpcBigtableTimeout := 30
	rpcBindAddress := "1.2.3.4"
	rpcFaucetAddress := "127.0.0.1:9900"
	rpcMaxMultipleAccounts := 100
	rpcMaxRequestBodySize := 51200
	rpcNicenessAdjustment := 0
	rpcPort := 38
	rpcPubsubEnableBlockSubscription := true
	rpcPubsubEnableVoteSubscription := true
	rpcPubsubMaxActiveSubscriptions := 1000000
	rpcPubsubNotificationThreads := 4
	rpcPubsubQueueCapacityBytes := 268435456
	rpcPubsubQueueCapacityItems := 10000000
	rpcPubsubWorkerThreads := 4
	rpcScanAndFixRoots := true
	rpcSendLeaderCount := 2
	rpcSendRetryMs := 2000
	rpcSendServiceMaxRetries := 5
	rpcSendTransactionAlsoLeader := true
	rpcSendTransactionRetryPoolMaxSize := 10000
	rpcSendTransactionTpuPeer := []string{"peer1:8000", "peer2:8000"}
	rpcThreads := 10
	skipPreflightHealthCheck := true
	skipSeedPhraseValidation := true
	skipStartupLedgerVerification := true
	snapshotArchiveFormat := "zstd"
	snapshotIntervalSlots := 100
	snapshotPackagerNicenessAdjustment := 0
	snapshotVersion := "1.2.0"
	stakedNodesOverrides := "/path/to/staked_nodes_overrides.yaml"
	towerStorage := "file"
	tpuCoalesceMs := 5
	tpuConnectionPoolSize := 10
	tpuDisableQuic := true
	tpuEnableUdp := true
	tvuReceiveThreads := 88
	unifiedSchedulerHandlerThreads := 2
	useSnapshotArchivesAtStartup := "whatever"
	waitForSupermajority := 1000
	walRecoveryMode := "someaddress"

	f := AgaveFlags{
		AccountIndex:                        &accountIndex,
		AccountIndexExcludeKey:              &accountIndexExcludeKey,
		AccountIndexIncludeKey:              &accountIndexIncludeKey,
		AccountShrinkPath:                   &accountShrinkPath,
		AccountsDbCacheLimitMb:              &accountsDbCacheLimitMb,
		AccountsDbTestHashCalculation:       &accountsDbTestHashCalculation,
		AccountsHashCachePath:               &accountsHashCachePath,
		AccountsIndexBins:                   &accountsIndexBins,
		AccountsIndexPath:                   &accountsIndexPath,
		AccountsIndexScanResultsLimitMb:     &accountsIndexScanResultsLimitMb,
		AccountsShrinkOptimizeTotalSpace:    &accountsShrinkOptimizeTotalSpace,
		AccountsShrinkRatio:                 &accountsShrinkRatio,
		AllowPrivateAddr:                    &allowPrivateAddr,
		AuthorizedVoter:                     &authorizedVoter,
		BindAddress:                         &bindAddress,
		BlockProductionMethod:               &blockProductionMethod,
		BlockVerificationMethod:             &blockVerificationMethod,
		CheckVoteAccount:                    &checkVoteAccount,
		ContactDebugInterval:                &contactDebugInterval,
		Cuda:                                &cuda,
		DebugKey:                            &debugKey,
		DevHaltAtSlot:                       &devHaltAtSlot,
		DisableBankingTrace:                 &disableBankingTrace,
		DynamicPortRange:                    &dynamicPortRange,
		EnableBankingTrace:                  &enableBankingTrace,
		EnableBigtableLedgerUpload:          &enableBigtableLedgerUpload,
		EnableExtendedTxMetadataStorage:     &enableExtendedTxMetadataStorage,
		EnableRpcBigtableLedgerStorage:      &enableRpcBigtableLedgerStorage,
		EnableRpcTransactionHistory:         &enableRpcTransactionHistory,
		EntryPoint:                          &entryPoint,
		EtcdCacertFile:                      &etcdCacertFile,
		EtcdCertFile:                        &etcdCertFile,
		EtcdDomainName:                      &etcdDomainName,
		EtcdEndpoint:                        &etcdEndpoint,
		EtcdKeyFile:                         &etcdKeyFile,
		ExpectedBankHash:                    &expectedBankHash,
		ExpectedGenesisHash:                 &expectedGenesisHash,
		ExpectedShredVersion:                &expectedShredVersion,
		ExtraFlags:                          &extraFlags,
		FullRpcAPI:                          &fullRpcAPI,
		FullSnapshotArchivePath:             &fullSnapshotArchivePath,
		FullSnapshotIntervalSlots:           &fullSnapshotIntervalSlots,
		GeyserPluginAlwaysEnabled:           &geyserPluginAlwaysEnabled,
		GeyserPluginConfig:                  &geyserPluginConfig,
		GossipHost:                          &gossipHost,
		GossipPort:                          &gossipPort,
		GossipValidator:                     &gossipValidator,
		HardFork:                            &hardFork,
		HealthCheckSlotDistance:             &healthCheckSlotDistance,
		IncrementalSnapshotArchivePath:      &incrementalSnapshotArchivePath,
		InitCompleteFile:                    &initCompleteFile,
		KnownValidator:                      &knownValidator,
		LimitLedgerSize:                     &limitLedgerSize,
		Log:                                 &log,
		LogMessagesBytesLimit:               &logMessagesBytesLimit,
		MaxGenesisArchiveUnpackedSize:       &maxGenesisArchiveUnpackedSize,
		MaximumFullSnapshotsToRetain:        &maximumFullSnapshotsToRetain,
		MaximumIncrementalSnapshotsToRetain: &maximumIncrementalSnapshotsToRetain,
		MaximumLocalSnapshotAge:             &maximumLocalSnapshotAge,
		MaximumSnapshotDownloadAbort:        &maximumSnapshotDownloadAbort,
		MinimalSnapshotDownloadSpeed:        &minimalSnapshotDownloadSpeed,
		NoGenesisFetch:                      &noGenesisFetch,
		NoIncrementalSnapshots:              &noIncrementalSnapshots,
		NoSnapshotFetch:                     &noSnapshotFetch,
		NoVoting:                            &noVoting,
		NoWaitForVoteToStartLeader:          noWaitForVoteToStartLeader,
		OnlyKnownRPC:                        &onlyKnownRPC,
		PrivateRPC:                          &privateRPC,
		PublicRpcAddress:                    &publicRpcAddress,
		PublicTpuAddress:                    &publicTpuAddress,
		PublicTpuForwardsAddress:            &publicTpuForwardsAddress,
		RepairValidator:                     &repairValidator,
		RequireTower:                        &requireTower,
		RestrictedRepairOnlyMode:            &restrictedRepairOnlyMode,
		RocksdbFifoShredStorageSize:         &rocksdbFifoShredStorageSize,
		RocksdbShredCompaction:              &rocksdbShredCompaction,
		RpcBigtableAppProfileId:             &rpcBigtableAppProfileId,
		RpcBigtableInstanceName:             &rpcBigtableInstanceName,
		RpcBigtableMaxMessageSize:           &rpcBigtableMaxMessageSize,
		RpcBigtableTimeout:                  &rpcBigtableTimeout,
		RpcBindAddress:                      rpcBindAddress,
		RpcFaucetAddress:                    &rpcFaucetAddress,
		RpcMaxMultipleAccounts:              &rpcMaxMultipleAccounts,
		RpcMaxRequestBodySize:               &rpcMaxRequestBodySize,
		RpcNicenessAdjustment:               &rpcNicenessAdjustment,
		RpcPort:                             rpcPort,
		RpcPubsubEnableBlockSubscription:    &rpcPubsubEnableBlockSubscription,
		RpcPubsubEnableVoteSubscription:     &rpcPubsubEnableVoteSubscription,
		RpcPubsubMaxActiveSubscriptions:     &rpcPubsubMaxActiveSubscriptions,
		RpcPubsubNotificationThreads:        &rpcPubsubNotificationThreads,
		RpcPubsubQueueCapacityBytes:         &rpcPubsubQueueCapacityBytes,
		RpcPubsubQueueCapacityItems:         &rpcPubsubQueueCapacityItems,
		RpcPubsubWorkerThreads:              &rpcPubsubWorkerThreads,
		RpcScanAndFixRoots:                  &rpcScanAndFixRoots,
		RpcSendLeaderCount:                  &rpcSendLeaderCount,
		RpcSendRetryMs:                      &rpcSendRetryMs,
		RpcSendServiceMaxRetries:            &rpcSendServiceMaxRetries,
		RpcSendTransactionAlsoLeader:        &rpcSendTransactionAlsoLeader,
		RpcSendTransactionRetryPoolMaxSize:  &rpcSendTransactionRetryPoolMaxSize,
		RpcSendTransactionTpuPeer:           &rpcSendTransactionTpuPeer,
		RpcThreads:                          &rpcThreads,
		SkipPreflightHealthCheck:            &skipPreflightHealthCheck,
		SkipSeedPhraseValidation:            &skipSeedPhraseValidation,
		SkipStartupLedgerVerification:       &skipStartupLedgerVerification,
		SnapshotArchiveFormat:               &snapshotArchiveFormat,
		SnapshotIntervalSlots:               &snapshotIntervalSlots,
		SnapshotPackagerNicenessAdjustment:  &snapshotPackagerNicenessAdjustment,
		SnapshotVersion:                     &snapshotVersion,
		StakedNodesOverrides:                &stakedNodesOverrides,
		TowerStorage:                        &towerStorage,
		TpuCoalesceMs:                       &tpuCoalesceMs,
		TpuConnectionPoolSize:               &tpuConnectionPoolSize,
		TpuDisableQuic:                      &tpuDisableQuic,
		TpuEnableUdp:                        &tpuEnableUdp,
		TvuReceiveThreads:                   &tvuReceiveThreads,
		UnifiedSchedulerHandlerThreads:      &unifiedSchedulerHandlerThreads,
		UseSnapshotArchivesAtStartup:        &useSnapshotArchivesAtStartup,
		WaitForSupermajority:                &waitForSupermajority,
		WalRecoveryMode:                     walRecoveryMode,
	}

	expectedArgs := []string{
		"--identity", *paths.ValidatorIdentityKeypairPath,
		"--vote-account", *paths.ValidatorVoteAccountKeypairPath,
		"--log", *paths.LogPath,
		"--accounts", *paths.AccountsPath,
		"--ledger", *paths.LedgerPath,
		"--account-index", "program-id",
		"--account-index", "spl-token-owner",
		"--account-index-exclude-key", "excludeKey1",
		"--account-index-exclude-key", "excludeKey2",
		"--account-index-include-key", "includeKey1",
		"--account-index-include-key", "includeKey2",
		"--account-shrink-path", "/path/to/shrink1",
		"--account-shrink-path", "/path/to/shrink2",
		"--accounts-db-cache-limit-mb", "1024",
		"--accounts-db-test-hash-calculation",
		"--accounts-hash-cache-path", "/path/to/accounts_hash_cache",
		"--accounts-index-bins", "64",
		"--accounts-index-path", "/path/to/accounts_index1",
		"--accounts-index-path", "/path/to/accounts_index2",
		"--accounts-index-scan-results-limit-mb", "512",
		"--accounts-shrink-optimize-total-space",
		"--accounts-shrink-ratio", "0.8",
		"--allow-private-addr",
		"--authorized-voter", "voterPubkey1",
		"--authorized-voter", "voterPubkey2",
		"--bind-address", "0.0.0.0",
		"--block-production-method", "whatever",
		"--block-verification-method", "unified-scheduler",
		"--check-vote-account", "http://localhost:8899",
		"--contact-debug-interval", "120000",
		"--cuda",
		"--debug-key", "debugKey1",
		"--debug-key", "debugKey2",
		"--dev-halt-at-slot", "1000000",
		"--disable-banking-trace",
		"--dynamic-port-range", "8002-8003",
		"--enable-banking-trace", "15032385536",
		"--enable-bigtable-ledger-upload",
		"--enable-extended-tx-metadata-storage",
		"--enable-rpc-bigtable-ledger-storage",
		"--enable-rpc-transaction-history",
		"--entrypoint", "hostA",
		"--entrypoint", "hostB",
		"--etcd-cacert-file", "/path/to/ca.crt",
		"--etcd-cert-file", "/path/to/client.crt",
		"--etcd-domain-name", "example.com",
		"--etcd-endpoint", "etcd1.example.com:2379",
		"--etcd-endpoint", "etcd2.example.com:2379",
		"--etcd-key-file", "/path/to/client.key",
		"--expected-bank-hash", "expectedBankHashValue",
		"--expected-genesis-hash", "asdkjasldjadk",
		"--expected-shred-version", "1234",
		"--extra", "--flag",
		"--full-rpc-api",
		"--full-snapshot-archive-path", "/path/to/full_snapshot",
		"--full-snapshot-interval-slots", "1000",
		"--geyser-plugin-always-enabled",
		"--geyser-plugin-config", "/path/to/geyser_plugin1.yaml",
		"--geyser-plugin-config", "/path/to/geyser_plugin2.yaml",
		"--gossip-host", "somehost",
		"--gossip-port", "62",
		"--gossip-validator", "validatorPubkey1",
		"--gossip-validator", "validatorPubkey2",
		"--hard-fork", "500000",
		"--hard-fork", "1000000",
		"--health-check-slot-distance", "128",
		"--incremental-snapshot-archive-path", "/path/to/incremental_snapshot",
		"--init-complete-file", "/path/to/init_complete",
		"--known-validator", "validatorA",
		"--known-validator", "validatorB",
		"--limit-ledger-size", "62",
		"--log-messages-bytes-limit", "1048576",
		"--max-genesis-archive-unpacked-size", "10485760",
		"--maximum-full-snapshots-to-retain", "2",
		"--maximum-incremental-snapshots-to-retain", "4",
		"--maximum-local-snapshot-age", "2500",
		"--maximum-snapshot-download-abort", "5",
		"--minimal-snapshot-download-speed", "10485760",
		"--no-genesis-fetch",
		"--no-incremental-snapshots",
		"--no-snapshot-fetch",
		"--no-voting",
		"--no-wait-for-vote-to-start-leader",
		"--only-known-rpc",
		"--private-rpc",
		"--public-rpc-address", "0.0.0.0:8899",
		"--public-tpu-address", "0.0.0.0:8001",
		"--public-tpu-forwards-address", "0.0.0.0:8002",
		"--repair-validator", "repairValidator1",
		"--repair-validator", "repairValidator2",
		"--require-tower",
		"--restricted-repair-only-mode",
		"--rocksdb-fifo-shred-storage-size", "104857600",
		"--rocksdb-shred-compaction", "fifo",
		"--rpc-bigtable-app-profile-id", "default",
		"--rpc-bigtable-instance-name", "solana-ledger",
		"--rpc-bigtable-max-message-size", "67108864",
		"--rpc-bigtable-timeout", "30",
		"--rpc-bind-address", "1.2.3.4",
		"--rpc-faucet-address", "127.0.0.1:9900",
		"--rpc-max-multiple-accounts", "100",
		"--rpc-max-request-body-size", "51200",
		"--rpc-niceness-adjustment", "0",
		"--rpc-port", "38",
		"--rpc-pubsub-enable-block-subscription",
		"--rpc-pubsub-enable-vote-subscription",
		"--rpc-pubsub-max-active-subscriptions", "1000000",
		"--rpc-pubsub-notification-threads", "4",
		"--rpc-pubsub-queue-capacity-bytes", "268435456",
		"--rpc-pubsub-queue-capacity-items", "10000000",
		"--rpc-pubsub-worker-threads", "4",
		"--rpc-scan-and-fix-roots",
		"--rpc-send-leader-count", "2",
		"--rpc-send-retry-ms", "2000",
		"--rpc-send-service-max-retries", "5",
		"--rpc-send-transaction-also-leader",
		"--rpc-send-transaction-retry-pool-max-size", "10000",
		"--rpc-send-transaction-tpu-peer", "peer1:8000",
		"--rpc-send-transaction-tpu-peer", "peer2:8000",
		"--rpc-threads", "10",
		"--skip-preflight-health-check",
		"--skip-seed-phrase-validation",
		"--skip-startup-ledger-verification",
		"--snapshot-archive-format", "zstd",
		"--snapshot-interval-slots", "100",
		"--snapshot-packager-niceness-adjustment", "0",
		"--snapshot-version", "1.2.0",
		"--staked-nodes-overrides", "/path/to/staked_nodes_overrides.yaml",
		"--tower-storage", "file",
		"--tpu-coalesce-ms", "5",
		"--tpu-connection-pool-size", "10",
		"--tpu-disable-quic",
		"--tpu-enable-udp",
		"--tvu-receive-threads", "88",
		"--unified-scheduler-handler-threads", "2",
		"--use-snapshot-archives-at-startup", "whatever",
		"--wait-for-supermajority", "1000",
		"--wal-recovery-mode", "someaddress",
	}

	actualArgs := f.Args(*paths)

	assert.Equal(t, expectedArgs, actualArgs)
}
