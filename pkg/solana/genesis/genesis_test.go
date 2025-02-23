package genesis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenesisFlags(t *testing.T) {
	identityPubkey := "identity_key"
	votePubkey := "vote_key"
	stakePubkey := "stake_key"
	ledgerPath := "ledger_path"
	bootstrapStakeAuthorizedPubkey := "bootstrap_auth_key"
	bootstrapValidatorLamports := 5000000000
	bootstrapValidatorStakeLamports := 1000000000
	clusterType := "mainnet-beta"
	creationTime := "2024-01-01T00:00:00Z"
	deactivateFeatures := []string{"feature1", "feature2"}
	enableWarmupEpochs := true
	faucetPubkey := "faucet_key"
	faucetLamports := 1000000000
	feeBurnPercentage := 50
	hashesPerTick := "auto"
	inflation := "0.0"
	lamportsPerByteYear := 100000
	maxGenesisArchiveUnpackedSize := 1073741824
	rentBurnPercentage := 5
	rentExemptionThreshold := 10
	slotsPerEpoch := 432000
	targetLamportsPerSignature := 42
	targetSignaturesPerSlot := 250000
	targetTickDuration := 400
	ticksPerSlot := 64
	url := "http://127.0.0.1"
	voteCommissionPercentage := 10
	extraFlags := []string{"--extra-flag1", "--extra-flag2"}

	f := GenesisFlags{
		IdentityPubkey:                  identityPubkey,
		LedgerPath:                      ledgerPath,
		VotePubkey:                      votePubkey,
		StakePubkey:                     stakePubkey,
		BootstrapStakeAuthorizedPubkey:  &bootstrapStakeAuthorizedPubkey,
		BootstrapValidatorLamports:      &bootstrapValidatorLamports,
		BootstrapValidatorStakeLamports: &bootstrapValidatorStakeLamports,
		ClusterType:                     &clusterType,
		CreationTime:                    &creationTime,
		DeactivateFeatures:              &deactivateFeatures,
		EnableWarmupEpochs:              &enableWarmupEpochs,
		FaucetPubkey:                    &faucetPubkey,
		FaucetLamports:                  &faucetLamports,
		FeeBurnPercentage:               &feeBurnPercentage,
		HashesPerTick:                   &hashesPerTick,
		Inflation:                       &inflation,
		LamportsPerByteYear:             &lamportsPerByteYear,
		MaxGenesisArchiveUnpackedSize:   &maxGenesisArchiveUnpackedSize,
		RentBurnPercentage:              &rentBurnPercentage,
		RentExemptionThreshold:          &rentExemptionThreshold,
		SlotsPerEpoch:                   &slotsPerEpoch,
		TargetLamportsPerSignature:      &targetLamportsPerSignature,
		TargetSignaturesPerSlot:         &targetSignaturesPerSlot,
		TargetTickDuration:              &targetTickDuration,
		TicksPerSlot:                    &ticksPerSlot,
		Url:                             &url,
		VoteCommissionPercentage:        &voteCommissionPercentage,
		ExtraFlags:                      &extraFlags,
	}

	actualArgs := f.Args()

	// Construct the expected argument list
	expectedArgs := []string{
		"--primordial-accounts-file", primordialAccountPath,
		"--bootstrap-validator", identityPubkey, votePubkey, stakePubkey,
		"--ledger", ledgerPath,
		"--bootstrap-stake-authorized-pubkey", "bootstrap_auth_key",
		"--bootstrap-validator-lamports", "5000000000",
		"--bootstrap-validator-stake-lamports", "1000000000",
		"--cluster-type", "mainnet-beta",
		"--creation-time", "2024-01-01T00:00:00Z",
		"--deactivate-feature", "feature1",
		"--deactivate-feature", "feature2",
		"--enable-warmup-epochs",
		"--faucet-pubkey", "faucet_key",
		"--faucet-lamports", "1000000000",
		"--fee-burn-percentage", "50",
		"--hashes-per-tick", "auto",
		"--inflation", "0.0",
		"--lamports-per-byte-year", "100000",
		"--max-genesis-archive-unpacked-size", "1073741824",
		"--rent-burn-percentage", "5",
		"--rent-exemption-threshold", "10",
		"--slots-per-epoch", "432000",
		"--target-lamports-per-signature", "42",
		"--target-signatures-per-slot", "250000",
		"--target-tick-duration", "400",
		"--ticks-per-slot", "64",
		"--url", "http://127.0.0.1",
		"--vote-commission-percentage", "10",
		"--extra-flag1", "--extra-flag2",
	}

	assert.Equal(t, expectedArgs, actualArgs)
}
