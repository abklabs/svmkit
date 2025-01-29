package tuner

import (
	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/abklabs/svmkit/pkg/runner/deb"
)

const (
	// TCP Buffer Sizes
	defaultNetIpv4TcpRmem = "10240 87380 12582912"
	defaultNetIpv4TcpWmem = "10240 87380 12582912"

	// TCP Optimization
	defaultNetIpv4TcpCongestionControl = "westwood"
	defaultNetIpv4TcpFastopen          = 3
	defaultNetIpv4TcpTimestamps        = 0
	defaultNetIpv4TcpSack              = 1
	defaultNetIpv4TcpLowLatency        = 1
	defaultNetIpv4TcpTwReuse           = 1
	defaultNetIpv4TcpNoMetricsSave     = 1
	defaultNetIpv4TcpModerateRcvbuf    = 1

	// Kernel Optimization
	defaultKernelTimerMigration           = 0
	defaultKernelNmiWatchdog              = 0
	defaultKernelSchedMinGranularityNs    = 10000000
	defaultKernelSchedWakeupGranularityNs = 15000000
	defaultKernelHungTaskTimeoutSecs      = 600
	defaultKernelPidMax                   = 65536

	// Virtual Memory Tuning
	defaultVmSwappiness              = 30
	defaultVmMaxMapCount             = 700000
	defaultVmStatInterval            = 10
	defaultVmDirtyRatio              = 40
	defaultVmDirtyBackgroundRatio    = 10
	defaultVmMinFreeKbytes           = 3000000
	defaultVmDirtyExpireCentisecs    = 36000
	defaultVmDirtyWritebackCentisecs = 3000
	defaultVmDirtytimeExpireSeconds  = 43200

	// Validator-Specific Networking
	defaultNetCoreRmemMax     = 134217728
	defaultNetCoreRmemDefault = 134217728
	defaultNetCoreWmemMax     = 134217728
	defaultNetCoreWmemDefault = 134217728
)

type CpuGovernor string

const (
	CpuGovernorPerformance  CpuGovernor = "performance"
	CpuGovernorPowersave    CpuGovernor = "powersave"
	CpuGovernorOndemand     CpuGovernor = "ondemand"
	CpuGovernorConservative CpuGovernor = "conservative"
	CpuGovernorSchedutil    CpuGovernor = "schedutil"
	CpuGovernorUserspace    CpuGovernor = "userspace"
)

type TunerCommand struct {
	Tuner
}

func (cmd *TunerCommand) Env() *runner.EnvBuilder {
	tunerEnv := runner.NewEnvBuilder()

	if cmd.CpuGovernor != nil {
		tunerEnv.Set("CPU_GOVERNOR", string(*cmd.CpuGovernor))
	} else {
		value := string(CpuGovernorPerformance)
		tunerEnv.Set("CPU_GOVERNOR", value)
	}

	tunerEnv.Merge(cmd.RunnerCommand.Env())

	return tunerEnv
}

func (cmd *TunerCommand) Check() error {
	cmd.RunnerCommand.SetConfigDefaults()

	pkgGrp := deb.Package{}.MakePackageGroup("cpufrequtils")

	if err := cmd.RunnerCommand.UpdatePackageGroup(pkgGrp); err != nil {
		return err
	}

	if cmd.Net == nil {
		cmd.Net = &TunerNetParams{}
	}

	cmd.Net.SetDefaults()

	if cmd.Kernel == nil {
		cmd.Kernel = &TunerKernelParams{}
	}

	cmd.Kernel.SetDefaults()

	if cmd.Vm == nil {
		cmd.Vm = &TunerVmParams{}
	}

	cmd.Vm.SetDefaults()

	return nil
}

func (cmd *TunerCommand) AddToPayload(p *runner.Payload) error {
	if err := p.AddTemplate("steps.sh", tunerScriptTmpl, cmd); err != nil {
		return err
	}

	if err := p.AddTemplate("svmkit-tuner.conf", svmkitTunerConfTmpl, cmd); err != nil {
		return err
	}

	if err := cmd.RunnerCommand.AddToPayload(p); err != nil {
		return err
	}

	return nil
}

type Tuner struct {
	runner.RunnerCommand
	Net         *TunerNetParams    `pulumi:"net,optional"`
	Kernel      *TunerKernelParams `pulumi:"kernel,optional"`
	Vm          *TunerVmParams     `pulumi:"vm,optional"`
	CpuGovernor *CpuGovernor       `pulumi:"cpuGovernor,optional"`
}

func (f *Tuner) Create() runner.Command {
	return &TunerCommand{
		Tuner: *f,
	}
}

type TunerNetParams struct {
	// net.ipv4.tcp_rmem => "10240 87380 12582912"
	NetIpv4TcpRmem *string `pulumi:"netIpv4TcpRmem,optional"`

	// net.ipv4.tcp_wmem => "10240 87380 12582912"
	NetIpv4TcpWmem *string `pulumi:"netIpv4TcpWmem,optional"`

	// net.ipv4.tcp_congestion_control => "westwood"
	NetIpv4TcpCongestionControl *string `pulumi:"netIpv4TcpCongestionControl,optional"`

	// net.ipv4.tcp_fastopen => 3
	NetIpv4TcpFastopen *int `pulumi:"netIpv4TcpFastopen,optional"`

	// net.ipv4.tcp_timestamps => 0
	NetIpv4TcpTimestamps *int `pulumi:"netIpv4TcpTimestamps,optional"`

	// net.ipv4.tcp_sack => 1
	NetIpv4TcpSack *int `pulumi:"netIpv4TcpSack,optional"`

	// net.ipv4.tcp_low_latency => 1
	NetIpv4TcpLowLatency *int `pulumi:"netIpv4TcpLowLatency,optional"`

	// net.ipv4.tcp_tw_reuse => 1
	NetIpv4TcpTwReuse *int `pulumi:"netIpv4TcpTwReuse,optional"`

	// net.ipv4.tcp_no_metrics_save => 1
	NetIpv4TcpNoMetricsSave *int `pulumi:"netIpv4TcpNoMetricsSave,optional"`

	// net.ipv4.tcp_moderate_rcvbuf => 1
	NetIpv4TcpModerateRcvbuf *int `pulumi:"netIpv4TcpModerateRcvbuf,optional"`

	// net.core.rmem_max => 134217728
	NetCoreRmemMax *int `pulumi:"netCoreRmemMax,optional"`

	// net.core.rmem_default => 134217728
	NetCoreRmemDefault *int `pulumi:"netCoreRmemDefault,optional"`

	// net.core.wmem_max => 134217728
	NetCoreWmemMax *int `pulumi:"netCoreWmemMax,optional"`

	// net.core.wmem_default => 134217728
	NetCoreWmemDefault *int `pulumi:"netCoreWmemDefault,optional"`
}

type TunerKernelParams struct {
	// kernel.timer_migration => 0
	KernelTimerMigration *int `pulumi:"kernelTimerMigration,optional"`

	// kernel.nmi_watchdog => 0
	KernelNmiWatchdog *int `pulumi:"kernelNmiWatchdog,optional"`

	// kernel.sched_min_granularity_ns => 10000000
	KernelSchedMinGranularityNs *int `pulumi:"kernelSchedMinGranularityNs,optional"`

	// kernel.sched_wakeup_granularity_ns => 15000000
	KernelSchedWakeupGranularityNs *int `pulumi:"kernelSchedWakeupGranularityNs,optional"`

	// kernel.hung_task_timeout_secs => 600
	KernelHungTaskTimeoutSecs *int `pulumi:"kernelHungTaskTimeoutSecs,optional"`

	// kernel.pid_max => 65536
	KernelPidMax *int `pulumi:"kernelPidMax,optional"`
}

type TunerVmParams struct {
	// vm.swappiness => 30
	VmSwappiness *int `pulumi:"vmSwappiness,optional"`

	// vm.max_map_count => 700000
	VmMaxMapCount *int `pulumi:"vmMaxMapCount,optional"`

	// vm.stat_interval => 10
	VmStatInterval *int `pulumi:"vmStatInterval,optional"`

	// vm.dirty_ratio => 40
	VmDirtyRatio *int `pulumi:"vmDirtyRatio,optional"`

	// vm.dirty_background_ratio => 10
	VmDirtyBackgroundRatio *int `pulumi:"vmDirtyBackgroundRatio,optional"`

	// vm.min_free_kbytes => 3000000
	VmMinFreeKbytes *int `pulumi:"vmMinFreeKbytes,optional"`

	// vm.dirty_expire_centisecs => 36000
	VmDirtyExpireCentisecs *int `pulumi:"vmDirtyExpireCentisecs,optional"`

	// vm.dirty_writeback_centisecs => 3000
	VmDirtyWritebackCentisecs *int `pulumi:"vmDirtyWritebackCentisecs,optional"`

	// vm.dirtytime_expire_seconds => 43200
	VmDirtytimeExpireSeconds *int `pulumi:"vmDirtytimeExpireSeconds,optional"`
}

func (f *TunerNetParams) SetDefaults() {
	// net.ipv4.tcp_rmem
	if f.NetIpv4TcpRmem == nil {
		value := defaultNetIpv4TcpRmem
		f.NetIpv4TcpRmem = &value
	}

	// net.ipv4.tcp_wmem
	if f.NetIpv4TcpWmem == nil {
		value := defaultNetIpv4TcpWmem
		f.NetIpv4TcpWmem = &value
	}

	// net.ipv4.tcp_congestion_control
	if f.NetIpv4TcpCongestionControl == nil {
		value := defaultNetIpv4TcpCongestionControl
		f.NetIpv4TcpCongestionControl = &value
	}

	// net.ipv4.tcp_fastopen
	if f.NetIpv4TcpFastopen == nil {
		value := defaultNetIpv4TcpFastopen
		f.NetIpv4TcpFastopen = &value
	}

	// net.ipv4.tcp_timestamps
	if f.NetIpv4TcpTimestamps == nil {
		value := defaultNetIpv4TcpTimestamps
		f.NetIpv4TcpTimestamps = &value
	}

	// net.ipv4.tcp_sack
	if f.NetIpv4TcpSack == nil {
		value := defaultNetIpv4TcpSack
		f.NetIpv4TcpSack = &value
	}

	// net.ipv4.tcp_low_latency
	if f.NetIpv4TcpLowLatency == nil {
		value := defaultNetIpv4TcpLowLatency
		f.NetIpv4TcpLowLatency = &value
	}

	// net.ipv4.tcp_tw_reuse
	if f.NetIpv4TcpTwReuse == nil {
		value := defaultNetIpv4TcpTwReuse
		f.NetIpv4TcpTwReuse = &value
	}

	// net.ipv4.tcp_no_metrics_save
	if f.NetIpv4TcpNoMetricsSave == nil {
		value := defaultNetIpv4TcpNoMetricsSave
		f.NetIpv4TcpNoMetricsSave = &value
	}

	// net.ipv4.tcp_moderate_rcvbuf
	if f.NetIpv4TcpModerateRcvbuf == nil {
		value := defaultNetIpv4TcpModerateRcvbuf
		f.NetIpv4TcpModerateRcvbuf = &value
	}

	// net.core.rmem_max
	if f.NetCoreRmemMax == nil {
		value := defaultNetCoreRmemMax
		f.NetCoreRmemMax = &value
	}

	// net.core.rmem_default
	if f.NetCoreRmemDefault == nil {
		value := defaultNetCoreRmemDefault
		f.NetCoreRmemDefault = &value
	}

	// net.core.wmem_max
	if f.NetCoreWmemMax == nil {
		value := defaultNetCoreWmemMax
		f.NetCoreWmemMax = &value
	}

	// net.core.wmem_default
	if f.NetCoreWmemDefault == nil {
		value := defaultNetCoreWmemDefault
		f.NetCoreWmemDefault = &value
	}
}

func (f *TunerKernelParams) SetDefaults() {
	// kernel.timer_migration
	if f.KernelTimerMigration == nil {
		value := defaultKernelTimerMigration
		f.KernelTimerMigration = &value
	}

	// kernel.nmi_watchdog
	if f.KernelNmiWatchdog == nil {
		value := defaultKernelNmiWatchdog
		f.KernelNmiWatchdog = &value
	}

	// kernel.sched_min_granularity_ns
	if f.KernelSchedMinGranularityNs == nil {
		value := defaultKernelSchedMinGranularityNs
		f.KernelSchedMinGranularityNs = &value
	}

	// kernel.sched_wakeup_granularity_ns
	if f.KernelSchedWakeupGranularityNs == nil {
		value := defaultKernelSchedWakeupGranularityNs
		f.KernelSchedWakeupGranularityNs = &value
	}

	// kernel.hung_task_timeout_secs
	if f.KernelHungTaskTimeoutSecs == nil {
		value := defaultKernelHungTaskTimeoutSecs
		f.KernelHungTaskTimeoutSecs = &value
	}

	// kernel.pid_max
	if f.KernelPidMax == nil {
		value := defaultKernelPidMax
		f.KernelPidMax = &value
	}
}

func (f *TunerVmParams) SetDefaults() {
	// vm.swappiness
	if f.VmSwappiness == nil {
		value := defaultVmSwappiness
		f.VmSwappiness = &value
	}

	// vm.max_map_count
	if f.VmMaxMapCount == nil {
		value := defaultVmMaxMapCount
		f.VmMaxMapCount = &value
	}

	// vm.stat_interval
	if f.VmStatInterval == nil {
		value := defaultVmStatInterval
		f.VmStatInterval = &value
	}

	// vm.dirty_ratio
	if f.VmDirtyRatio == nil {
		value := defaultVmDirtyRatio
		f.VmDirtyRatio = &value
	}

	// vm.dirty_background_ratio
	if f.VmDirtyBackgroundRatio == nil {
		value := defaultVmDirtyBackgroundRatio
		f.VmDirtyBackgroundRatio = &value
	}

	// vm.min_free_kbytes
	if f.VmMinFreeKbytes == nil {
		value := defaultVmMinFreeKbytes
		f.VmMinFreeKbytes = &value
	}

	// vm.dirty_expire_centisecs
	if f.VmDirtyExpireCentisecs == nil {
		value := defaultVmDirtyExpireCentisecs
		f.VmDirtyExpireCentisecs = &value
	}

	// vm.dirty_writeback_centisecs
	if f.VmDirtyWritebackCentisecs == nil {
		value := defaultVmDirtyWritebackCentisecs
		f.VmDirtyWritebackCentisecs = &value
	}

	// vm.dirtytime_expire_seconds
	if f.VmDirtytimeExpireSeconds == nil {
		value := defaultVmDirtytimeExpireSeconds
		f.VmDirtytimeExpireSeconds = &value
	}
}
