package firedancer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultFiredancerPaths(t *testing.T) {
	fp, err := NewDefaultFiredancerPaths(nil)
	require.NoError(t, err)
	require.NotNil(t, fp)
	require.NotNil(t, fp.ConfigPath)
	assert.Equal(t, "/home/sol/config.toml", *fp.ConfigPath)
}

func TestFiredancerPathsMerge(t *testing.T) {
	fp1, err := NewDefaultFiredancerPaths(nil)
	require.NoError(t, err)
	require.NotNil(t, fp1)

	customConfig := "/custom/config.toml"
	fp2 := &FiredancerPaths{
		ConfigPath: &customConfig,
	}

	err = fp1.Merge(fp2)
	require.NoError(t, err)
	require.NotNil(t, fp1.ConfigPath)
	assert.Equal(t, customConfig, *fp1.ConfigPath)
}

func TestFiredancerPathsCheck(t *testing.T) {
	fp, err := NewDefaultFiredancerPaths(nil)
	require.NoError(t, err)
	require.NotNil(t, fp)

	err = fp.Check()
	assert.NoError(t, err)

	fp.ConfigPath = nil
	err = fp.Check()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ConfigPath is nil")
}
