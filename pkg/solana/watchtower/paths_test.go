package watchtower

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultWatchtowerPaths(t *testing.T) {
	wp, err := NewDefaultWatchtowerPaths(nil)
	require.NoError(t, err)
	require.NotNil(t, wp)
	require.NotNil(t, wp.RunBinPath)
	assert.Equal(t, "/usr/bin/run-watchtower", *wp.RunBinPath)
}

func TestWatchtowerPathsMerge(t *testing.T) {
	wp1, err := NewDefaultWatchtowerPaths(nil)
	require.NoError(t, err)
	require.NotNil(t, wp1)

	customRun := "/custom/run-watchtower"
	wp2 := &WatchtowerPaths{
		RunBinPath: &customRun,
	}

	err = wp1.Merge(wp2)
	require.NoError(t, err)
	require.NotNil(t, wp1.RunBinPath)
	assert.Equal(t, customRun, *wp1.RunBinPath)
}

func TestWatchtowerPathsCheck(t *testing.T) {
	wp, err := NewDefaultWatchtowerPaths(nil)
	require.NoError(t, err)
	require.NotNil(t, wp)

	err = wp.Check()
	assert.NoError(t, err)

	wp.RunBinPath = nil
	err = wp.Check()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "RunBinPath is nil")
}
