package faucet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultFaucetPaths(t *testing.T) {
	fp, err := NewDefaultFaucetPaths(nil)
	require.NoError(t, err)
	require.NotNil(t, fp)
	require.NotNil(t, fp.KeypairPath)
	assert.Equal(t, "/home/sol/faucet-keypair.json", *fp.KeypairPath)
	require.NotNil(t, fp.RunBinPath)
	assert.Equal(t, "/usr/bin/run-faucet", *fp.RunBinPath)
}

func TestFaucetPathsMerge(t *testing.T) {
	fp1, err := NewDefaultFaucetPaths(nil)
	require.NoError(t, err)
	require.NotNil(t, fp1)

	customKey := "/custom/faucet-keypair.json"
	customBin := "/custom/run-faucet"
	fp2 := &FaucetPaths{
		KeypairPath: &customKey,
		RunBinPath:  &customBin,
	}

	err = fp1.Merge(fp2)
	require.NoError(t, err)
	require.NotNil(t, fp1.KeypairPath)
	assert.Equal(t, customKey, *fp1.KeypairPath)
	require.NotNil(t, fp1.RunBinPath)
	assert.Equal(t, customBin, *fp1.RunBinPath)
}

func TestFaucetPathsCheck(t *testing.T) {
	fp, err := NewDefaultFaucetPaths(nil)
	require.NoError(t, err)
	require.NotNil(t, fp)

	err = fp.Check()
	assert.NoError(t, err)

	fp.KeypairPath = nil
	err = fp.Check()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "KeypairPath is nil")

	keyPath := "/home/sol/faucet-keypair.json"
	fp.KeypairPath = &keyPath
	fp.RunBinPath = nil
	err = fp.Check()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "RunBinPath is nil")
}
