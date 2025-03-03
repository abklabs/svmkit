package explorer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultExplorerPaths(t *testing.T) {
	ep, err := NewDefaultExplorerPaths(nil)
	require.NoError(t, err)
	require.NotNil(t, ep)
	require.NotNil(t, ep.InstallPath)
	assert.Equal(t, "/opt/svmkit-solana-explorer", *ep.InstallPath)
	require.NotNil(t, ep.RunBinPath)
	assert.Equal(t, "/opt/svmkit-solana-explorer/run-explorer", *ep.RunBinPath)
}

func TestExplorerPathsMerge(t *testing.T) {
	ep1, err := NewDefaultExplorerPaths(nil)
	require.NoError(t, err)
	require.NotNil(t, ep1)

	customInstall := "/custom/explorer"
	customRunBin := "/custom/run-explorer"
	ep2 := &ExplorerPaths{
		InstallPath: &customInstall,
		RunBinPath:  &customRunBin,
	}

	err = ep1.Merge(ep2)
	require.NoError(t, err)
	require.NotNil(t, ep1.InstallPath)
	assert.Equal(t, customInstall, *ep1.InstallPath)
	require.NotNil(t, ep1.RunBinPath)
	assert.Equal(t, customRunBin, *ep1.RunBinPath)
}

func TestExplorerPathsCheck(t *testing.T) {
	ep, err := NewDefaultExplorerPaths(nil)
	require.NoError(t, err)
	require.NotNil(t, ep)

	err = ep.Check()
	assert.NoError(t, err)

	ep.InstallPath = nil
	err = ep.Check()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InstallPath is nil")

	installPath := "/opt/svmkit-solana-explorer"
	ep.InstallPath = &installPath
	ep.RunBinPath = nil
	err = ep.Check()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "RunBinPath is nil")
}
