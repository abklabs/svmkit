package agave

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultAgavePaths(t *testing.T) {
	ap, err := NewDefaultAgavePaths(nil)
	require.NoError(t, err)
	require.NotNil(t, ap)

	require.NotNil(t, ap.RunBinPath)
	assert.Equal(t, "/usr/bin/run-validator", *ap.RunBinPath)
	require.NotNil(t, ap.StopBinPath)
	assert.Equal(t, "/usr/bin/stop-validator", *ap.StopBinPath)
	require.NotNil(t, ap.CheckBinPath)
	assert.Equal(t, "/usr/bin/check-validator", *ap.CheckBinPath)
}

func TestAgavePathsMerge(t *testing.T) {
	ap1, err := NewDefaultAgavePaths(nil)
	require.NoError(t, err)
	require.NotNil(t, ap1)

	customRun := "/custom/run-validator"
	customStop := "/custom/stop-validator"
	customCheck := "/custom/check-validator"

	ap2 := &AgavePaths{
		RunBinPath:   &customRun,
		StopBinPath:  &customStop,
		CheckBinPath: &customCheck,
	}

	err = ap1.Merge(ap2)
	require.NoError(t, err)
	require.NotNil(t, ap1.RunBinPath)
	assert.Equal(t, customRun, *ap1.RunBinPath)
	require.NotNil(t, ap1.StopBinPath)
	assert.Equal(t, customStop, *ap1.StopBinPath)
	require.NotNil(t, ap1.CheckBinPath)
	assert.Equal(t, customCheck, *ap1.CheckBinPath)
}

func TestAgavePathsCheck(t *testing.T) {
	ap, err := NewDefaultAgavePaths(nil)
	require.NoError(t, err)
	require.NotNil(t, ap)

	err = ap.Check()
	assert.NoError(t, err)

	ap.RunBinPath = nil
	err = ap.Check()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "RunBinPath is nil")

	runBin := "/usr/bin/run-validator"
	ap.RunBinPath = &runBin
	ap.StopBinPath = nil
	err = ap.Check()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "StopBinPath is nil")

	stopBin := "/usr/bin/stop-validator"
	ap.StopBinPath = &stopBin
	ap.CheckBinPath = nil
	err = ap.Check()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CheckBinPath is nil")
}
