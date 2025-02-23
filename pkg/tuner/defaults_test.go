package tuner

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDefaultTuner(t *testing.T) {
	// Create a new default Tuner (update this call if it now expects different arguments)
	tuner, err := NewDefaultTuner()
	require.NoError(t, err, "NewDefaultTuner should not return an error")
	require.NotNil(t, tuner, "NewDefaultTuner should return a non-nil Tuner struct")

	// --------------------------------------------------------------------
	// CPU Governor
	// --------------------------------------------------------------------
	require.NotNil(t, tuner.Params.CpuGovernor, "expected CpuGovernor to be a non-nil pointer")
	require.Equal(t, CpuGovernorPerformance, *tuner.Params.CpuGovernor,
		"cpuGovernor should default to 'performance'")

	// --------------------------------------------------------------------
	// Net defaults
	// --------------------------------------------------------------------
	require.NotNil(t, tuner.Params.Net, "tuner.Params.Net should not be nil")
	require.NotNil(t, tuner.Params.Net.NetIpv4TcpRmem)
	require.Equal(t, "10240 87380 12582912", *tuner.Params.Net.NetIpv4TcpRmem)
	require.NotNil(t, tuner.Params.Net.NetIpv4TcpWmem)
	require.Equal(t, "10240 87380 12582912", *tuner.Params.Net.NetIpv4TcpWmem)
	require.NotNil(t, tuner.Params.Net.NetIpv4TcpCongestionControl)
	require.Equal(t, "westwood", *tuner.Params.Net.NetIpv4TcpCongestionControl)
	require.NotNil(t, tuner.Params.Net.NetIpv4TcpFastopen)
	require.Equal(t, 3, *tuner.Params.Net.NetIpv4TcpFastopen)
	require.NotNil(t, tuner.Params.Net.NetIpv4TcpTimestamps)
	require.Equal(t, 0, *tuner.Params.Net.NetIpv4TcpTimestamps)
	require.NotNil(t, tuner.Params.Net.NetIpv4TcpSack)
	require.Equal(t, 1, *tuner.Params.Net.NetIpv4TcpSack)
	require.NotNil(t, tuner.Params.Net.NetIpv4TcpLowLatency)
	require.Equal(t, 1, *tuner.Params.Net.NetIpv4TcpLowLatency)
	require.NotNil(t, tuner.Params.Net.NetIpv4TcpTwReuse)
	require.Equal(t, 1, *tuner.Params.Net.NetIpv4TcpTwReuse)
	require.NotNil(t, tuner.Params.Net.NetIpv4TcpNoMetricsSave)
	require.Equal(t, 1, *tuner.Params.Net.NetIpv4TcpNoMetricsSave)
	require.NotNil(t, tuner.Params.Net.NetIpv4TcpModerateRcvbuf)
	require.Equal(t, 1, *tuner.Params.Net.NetIpv4TcpModerateRcvbuf)
	require.NotNil(t, tuner.Params.Net.NetCoreRmemMax)
	require.Equal(t, 134217728, *tuner.Params.Net.NetCoreRmemMax)
	require.NotNil(t, tuner.Params.Net.NetCoreRmemDefault)
	require.Equal(t, 134217728, *tuner.Params.Net.NetCoreRmemDefault)
	require.NotNil(t, tuner.Params.Net.NetCoreWmemMax)
	require.Equal(t, 134217728, *tuner.Params.Net.NetCoreWmemMax)
	require.NotNil(t, tuner.Params.Net.NetCoreWmemDefault)
	require.Equal(t, 134217728, *tuner.Params.Net.NetCoreWmemDefault)

	// --------------------------------------------------------------------
	// Kernel defaults
	// --------------------------------------------------------------------
	require.NotNil(t, tuner.Params.Kernel, "tuner.Params.Kernel should not be nil")
	require.NotNil(t, tuner.Params.Kernel.KernelTimerMigration)
	require.Equal(t, 0, *tuner.Params.Kernel.KernelTimerMigration)
	require.NotNil(t, tuner.Params.Kernel.KernelNmiWatchdog)
	require.Equal(t, 0, *tuner.Params.Kernel.KernelNmiWatchdog)
	require.NotNil(t, tuner.Params.Kernel.KernelSchedMinGranularityNs)
	require.Equal(t, 10000000, *tuner.Params.Kernel.KernelSchedMinGranularityNs)
	require.NotNil(t, tuner.Params.Kernel.KernelSchedWakeupGranularityNs)
	require.Equal(t, 15000000, *tuner.Params.Kernel.KernelSchedWakeupGranularityNs)
	require.NotNil(t, tuner.Params.Kernel.KernelHungTaskTimeoutSecs)
	require.Equal(t, 600, *tuner.Params.Kernel.KernelHungTaskTimeoutSecs)
	require.NotNil(t, tuner.Params.Kernel.KernelPidMax)
	require.Equal(t, 65536, *tuner.Params.Kernel.KernelPidMax)

	// --------------------------------------------------------------------
	// VM defaults
	// --------------------------------------------------------------------
	require.NotNil(t, tuner.Params.Vm, "tuner.Params.Vm should not be nil")
	require.NotNil(t, tuner.Params.Vm.VmSwappiness)
	require.Equal(t, 30, *tuner.Params.Vm.VmSwappiness)
	require.NotNil(t, tuner.Params.Vm.VmMaxMapCount)
	require.Equal(t, 700000, *tuner.Params.Vm.VmMaxMapCount)
	require.NotNil(t, tuner.Params.Vm.VmStatInterval)
	require.Equal(t, 10, *tuner.Params.Vm.VmStatInterval)
	require.NotNil(t, tuner.Params.Vm.VmDirtyRatio)
	require.Equal(t, 40, *tuner.Params.Vm.VmDirtyRatio)
	require.NotNil(t, tuner.Params.Vm.VmDirtyBackgroundRatio)
	require.Equal(t, 10, *tuner.Params.Vm.VmDirtyBackgroundRatio)
	require.NotNil(t, tuner.Params.Vm.VmMinFreeKbytes)
	require.Equal(t, 3000000, *tuner.Params.Vm.VmMinFreeKbytes)
	require.NotNil(t, tuner.Params.Vm.VmDirtyExpireCentisecs)
	require.Equal(t, 36000, *tuner.Params.Vm.VmDirtyExpireCentisecs)
	require.NotNil(t, tuner.Params.Vm.VmDirtyWritebackCentisecs)
	require.Equal(t, 3000, *tuner.Params.Vm.VmDirtyWritebackCentisecs)
	require.NotNil(t, tuner.Params.Vm.VmDirtytimeExpireSeconds)
	require.Equal(t, 43200, *tuner.Params.Vm.VmDirtytimeExpireSeconds)

	// --------------------------------------------------------------------
	// FS defaults
	// --------------------------------------------------------------------
	require.NotNil(t, tuner.Params.Fs, "tuner.Params.Fs should not be nil")
	require.NotNil(t, tuner.Params.Fs.FsNrOpen)
	require.Equal(t, 1000000, *tuner.Params.Fs.FsNrOpen)
}
