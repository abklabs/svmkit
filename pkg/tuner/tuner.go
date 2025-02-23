package tuner

import (
	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/abklabs/svmkit/pkg/runner/deb"

	"dario.cat/mergo"
	"github.com/pulumi/pulumi-go-provider/infer"
)

type TunerCommand struct {
	Tuner
}

func (cmd *TunerCommand) Env() *runner.EnvBuilder {
	tunerEnv := runner.NewEnvBuilder()

	if cmd.Params.CpuGovernor != nil {
		tunerEnv.Set("CPU_GOVERNOR", string(*cmd.Params.CpuGovernor))
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

	return nil
}

func (cmd *TunerCommand) AddToPayload(p *runner.Payload) error {
	if err := p.AddTemplate("steps.sh", tunerScriptTmpl, cmd); err != nil {
		return err
	}

	if err := p.AddTemplate("svmkit-tuner.conf", svmkitTunerConfTmpl, cmd.Params); err != nil {
		return err
	}

	if err := cmd.RunnerCommand.AddToPayload(p); err != nil {
		return err
	}

	return nil
}

type CpuGovernor string

const (
	CpuGovernorPerformance  CpuGovernor = "performance"
	CpuGovernorPowersave    CpuGovernor = "powersave"
	CpuGovernorOndemand     CpuGovernor = "ondemand"
	CpuGovernorConservative CpuGovernor = "conservative"
	CpuGovernorSchedutil    CpuGovernor = "schedutil"
	CpuGovernorUserspace    CpuGovernor = "userspace"
)

func (CpuGovernor) Values() []infer.EnumValue[CpuGovernor] {
	return []infer.EnumValue[CpuGovernor]{
		{
			Name:        "performance",
			Value:       CpuGovernorPerformance,
			Description: "The performance governor",
		},
		{
			Name:        "powersave",
			Value:       CpuGovernorPowersave,
			Description: "The powersave governor",
		},
		{
			Name:        "ondemand",
			Value:       CpuGovernorOndemand,
			Description: "The ondemand governor",
		},
		{
			Name:        "conservative",
			Value:       CpuGovernorConservative,
			Description: "The conservative governor",
		},
		{
			Name:        "schedutil",
			Value:       CpuGovernorSchedutil,
			Description: "The schedutil governor",
		},
		{
			Name:        "userspace",
			Value:       CpuGovernorUserspace,
			Description: "The userspace governor",
		},
	}
}

type TunerParams struct {
	CpuGovernor *CpuGovernor       `pulumi:"cpuGovernor,optional" toml:"cpuGovernor,omitempty"`
	Net         *TunerNetParams    `pulumi:"net,optional" toml:"net,omitempty"`
	Kernel      *TunerKernelParams `pulumi:"kernel,optional" toml:"kernel,omitempty"`
	Vm          *TunerVmParams     `pulumi:"vm,optional" toml:"vm,omitempty"`
	Fs          *TunerFsParams     `pulumi:"fs,optional" toml:"fs,omitempty"`
}

type Tuner struct {
	runner.RunnerCommand
	Params TunerParams `pulumi:"params" toml:"params"`
}

func (f *Tuner) Create() runner.Command {
	return &TunerCommand{
		Tuner: *f,
	}
}

type TunerNetParams struct {
	// net.ipv4.tcp_rmem => "10240 87380 12582912"
	NetIpv4TcpRmem *string `pulumi:"netIpv4TcpRmem,optional" toml:"netIpv4TcpRmem,omitempty"`

	// net.ipv4.tcp_wmem => "10240 87380 12582912"
	NetIpv4TcpWmem *string `pulumi:"netIpv4TcpWmem,optional" toml:"netIpv4TcpWmem,omitempty"`

	// net.ipv4.tcp_congestion_control => "westwood"
	NetIpv4TcpCongestionControl *string `pulumi:"netIpv4TcpCongestionControl,optional" toml:"netIpv4TcpCongestionControl,omitempty"`

	// net.ipv4.tcp_fastopen => 3
	NetIpv4TcpFastopen *int `pulumi:"netIpv4TcpFastopen,optional" toml:"netIpv4TcpFastopen,omitempty"`

	// net.ipv4.tcp_timestamps => 0
	NetIpv4TcpTimestamps *int `pulumi:"netIpv4TcpTimestamps,optional" toml:"netIpv4TcpTimestamps,omitempty"`

	// net.ipv4.tcp_sack => 1
	NetIpv4TcpSack *int `pulumi:"netIpv4TcpSack,optional" toml:"netIpv4TcpSack,omitempty"`

	// net.ipv4.tcp_low_latency => 1
	NetIpv4TcpLowLatency *int `pulumi:"netIpv4TcpLowLatency,optional" toml:"netIpv4TcpLowLatency,omitempty"`

	// net.ipv4.tcp_tw_reuse => 1
	NetIpv4TcpTwReuse *int `pulumi:"netIpv4TcpTwReuse,optional" toml:"netIpv4TcpTwReuse,omitempty"`

	// net.ipv4.tcp_no_metrics_save => 1
	NetIpv4TcpNoMetricsSave *int `pulumi:"netIpv4TcpNoMetricsSave,optional" toml:"netIpv4TcpNoMetricsSave,omitempty"`

	// net.ipv4.tcp_moderate_rcvbuf => 1
	NetIpv4TcpModerateRcvbuf *int `pulumi:"netIpv4TcpModerateRcvbuf,optional" toml:"netIpv4TcpModerateRcvbuf,omitempty"`

	// net.core.rmem_max => 134217728
	NetCoreRmemMax *int `pulumi:"netCoreRmemMax,optional" toml:"netCoreRmemMax,omitempty"`

	// net.core.rmem_default => 134217728
	NetCoreRmemDefault *int `pulumi:"netCoreRmemDefault,optional" toml:"netCoreRmemDefault,omitempty"`

	// net.core.wmem_max => 134217728
	NetCoreWmemMax *int `pulumi:"netCoreWmemMax,optional" toml:"netCoreWmemMax,omitempty"`

	// net.core.wmem_default => 134217728
	NetCoreWmemDefault *int `pulumi:"netCoreWmemDefault,optional" toml:"netCoreWmemDefault,omitempty"`
}

type TunerKernelParams struct {
	// kernel.timer_migration => 0
	KernelTimerMigration *int `pulumi:"kernelTimerMigration,optional" toml:"kernelTimerMigration,omitempty"`

	// kernel.nmi_watchdog => 0
	KernelNmiWatchdog *int `pulumi:"kernelNmiWatchdog,optional" toml:"kernelNmiWatchdog,omitempty"`

	// kernel.sched_min_granularity_ns => 10000000
	KernelSchedMinGranularityNs *int `pulumi:"kernelSchedMinGranularityNs,optional" toml:"kernelSchedMinGranularityNs,omitempty"`

	// kernel.sched_wakeup_granularity_ns => 15000000
	KernelSchedWakeupGranularityNs *int `pulumi:"kernelSchedWakeupGranularityNs,optional" toml:"kernelSchedWakeupGranularityNs,omitempty"`

	// kernel.hung_task_timeout_secs => 600
	KernelHungTaskTimeoutSecs *int `pulumi:"kernelHungTaskTimeoutSecs,optional" toml:"kernelHungTaskTimeoutSecs,omitempty"`

	// kernel.pid_max => 65536
	KernelPidMax *int `pulumi:"kernelPidMax,optional" toml:"kernelPidMax,omitempty"`
}

type TunerVmParams struct {
	// vm.swappiness => 30
	VmSwappiness *int `pulumi:"vmSwappiness,optional" toml:"vmSwappiness,omitempty"`

	// vm.max_map_count => 700000
	VmMaxMapCount *int `pulumi:"vmMaxMapCount,optional" toml:"vmMaxMapCount,omitempty"`

	// vm.stat_interval => 10
	VmStatInterval *int `pulumi:"vmStatInterval,optional" toml:"vmStatInterval,omitempty"`

	// vm.dirty_ratio => 40
	VmDirtyRatio *int `pulumi:"vmDirtyRatio,optional" toml:"vmDirtyRatio,omitempty"`

	// vm.dirty_background_ratio => 10
	VmDirtyBackgroundRatio *int `pulumi:"vmDirtyBackgroundRatio,optional" toml:"vmDirtyBackgroundRatio,omitempty"`

	// vm.min_free_kbytes => 3000000
	VmMinFreeKbytes *int `pulumi:"vmMinFreeKbytes,optional" toml:"vmMinFreeKbytes,omitempty"`

	// vm.dirty_expire_centisecs => 36000
	VmDirtyExpireCentisecs *int `pulumi:"vmDirtyExpireCentisecs,optional" toml:"vmDirtyExpireCentisecs,omitempty"`

	// vm.dirty_writeback_centisecs => 3000
	VmDirtyWritebackCentisecs *int `pulumi:"vmDirtyWritebackCentisecs,optional" toml:"vmDirtyWritebackCentisecs,omitempty"`

	// vm.dirtytime_expire_seconds => 43200
	VmDirtytimeExpireSeconds *int `pulumi:"vmDirtytimeExpireSeconds,optional" toml:"vmDirtytimeExpireSeconds,omitempty"`
}

type TunerFsParams struct {
	// fs.nr_open => 1000000
	FsNrOpen *int `pulumi:"fsNrOpen,optional" toml:"fsNrOpen,omitempty"`
}

func (t *Tuner) Merge(other *Tuner) error {
	if other == nil {
		return nil
	}
	return mergo.Merge(t, other, mergo.WithOverride)
}
