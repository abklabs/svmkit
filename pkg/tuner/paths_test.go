package tuner

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultTunerPaths(t *testing.T) {
	tp, err := NewDefaultTunerPaths(nil)
	require.NoError(t, err)
	require.NotNil(t, tp)
	require.NotNil(t, tp.TunerLogPath)
	assert.Equal(t, "/home/sol/svmkit-tuner.log", *tp.TunerLogPath)
	require.NotNil(t, tp.SysctlConfPath)
	assert.Equal(t, "/etc/sysctl.d/zzz-svmkit-tuner.conf", *tp.SysctlConfPath)
	require.NotNil(t, tp.RunBinPath)
	assert.Equal(t, "/usr/bin/run-tuner", *tp.RunBinPath)
}

func TestTunerPathsMerge(t *testing.T) {
	tp1, err := NewDefaultTunerPaths(nil)
	require.NoError(t, err)
	require.NotNil(t, tp1)

	customLog := "/custom/tuner.log"
	customConf := "/etc/sysctl.d/alt-tuner.conf"
	customBin := "/usr/bin/alt-run-tuner"

	tp2 := &TunerPaths{
		TunerLogPath:   &customLog,
		SysctlConfPath: &customConf,
		RunBinPath:     &customBin,
	}

	err = tp1.Merge(tp2)
	require.NoError(t, err)
	require.NotNil(t, tp1.TunerLogPath)
	assert.Equal(t, customLog, *tp1.TunerLogPath)
	require.NotNil(t, tp1.SysctlConfPath)
	assert.Equal(t, customConf, *tp1.SysctlConfPath)
	require.NotNil(t, tp1.RunBinPath)
	assert.Equal(t, customBin, *tp1.RunBinPath)
}

func TestTunerPathsCheck(t *testing.T) {
	tp, err := NewDefaultTunerPaths(nil)
	require.NoError(t, err)
	require.NotNil(t, tp)

	err = tp.Check()
	assert.NoError(t, err)

	tp.TunerLogPath = nil
	err = tp.Check()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "TunerLogPath is nil")

	logPath := "/home/sol/svmkit-tuner.log"
	tp.TunerLogPath = &logPath
	tp.SysctlConfPath = nil
	err = tp.Check()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SysctlConfPath is nil")
}
