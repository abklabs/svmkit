package genesis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultGenesisPaths(t *testing.T) {
	gp, err := NewDefaultGenesisPaths(nil)
	require.NoError(t, err)
	require.NotNil(t, gp)

	require.NotNil(t, gp.PrimordialAccountsPath)
	assert.Equal(t, "/home/sol/primordial.yaml", *gp.PrimordialAccountsPath)

	require.NotNil(t, gp.ValidatorAccountsPath)
	assert.Equal(t, "/home/sol/validator_accounts.yaml", *gp.ValidatorAccountsPath)

	require.NotNil(t, gp.SolanaSplCachePath)
	assert.Equal(t, "~/.cache/solana-spl", *gp.SolanaSplCachePath)
}

func TestGenesisPathsMerge(t *testing.T) {
	gp1, err := NewDefaultGenesisPaths(nil)
	require.NoError(t, err)
	require.NotNil(t, gp1)

	primordial := "/custom/primordial.yaml"
	valAccounts := "/custom/validator_accounts.yaml"
	splCache := "/custom/solana-spl-cache"

	gp2 := &GenesisPaths{
		PrimordialAccountsPath: &primordial,
		ValidatorAccountsPath:  &valAccounts,
		SolanaSplCachePath:     &splCache,
	}

	err = gp1.Merge(gp2)
	require.NoError(t, err)
	require.NotNil(t, gp1.PrimordialAccountsPath)
	assert.Equal(t, primordial, *gp1.PrimordialAccountsPath)

	require.NotNil(t, gp1.ValidatorAccountsPath)
	assert.Equal(t, valAccounts, *gp1.ValidatorAccountsPath)

	require.NotNil(t, gp1.SolanaSplCachePath)
	assert.Equal(t, splCache, *gp1.SolanaSplCachePath)
}

func TestGenesisPathsCheck(t *testing.T) {
	gp, err := NewDefaultGenesisPaths(nil)
	require.NoError(t, err)
	require.NotNil(t, gp)

	err = gp.Check()
	assert.NoError(t, err)

	gp.PrimordialAccountsPath = nil
	err = gp.Check()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "PrimordialAccountsPath is nil")

	primordial := "/home/sol/primordial.yaml"
	gp.PrimordialAccountsPath = &primordial
	gp.ValidatorAccountsPath = nil
	err = gp.Check()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ValidatorAccountsPath is nil")

	valAccounts := "/home/sol/validator_accounts.yaml"
	gp.ValidatorAccountsPath = &valAccounts
	gp.SolanaSplCachePath = nil
	err = gp.Check()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SolanaSplCachePath is nil")
}
