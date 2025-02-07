package tuner

import (
	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/abklabs/svmkit/pkg/runner/deb"
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

	if cmd.Variant != nil {
		tunerEnv.Set("TUNER_VARIANT", string(*cmd.Variant))
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
	Variant     *TunerVariant      `pulumi:"variant,optional"`
	Net         *TunerNetParams    `pulumi:"net,optional" toml:"net,omitempty"`
	Kernel      *TunerKernelParams `pulumi:"kernel,optional" toml:"kernel,omitempty"`
	Vm          *TunerVmParams     `pulumi:"vm,optional" toml:"vm,omitempty"`
	CpuGovernor *CpuGovernor       `pulumi:"cpuGovernor,optional" toml:"cpuGovernor,omitempty"`
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

func (t *Tuner) Merge(other *Tuner) {
	if other == nil {
		return
	}

	if other.Net != nil {
		if t.Net == nil {
			t.Net = &TunerNetParams{}
		}
		if other.Net.NetIpv4TcpRmem != nil {
			t.Net.NetIpv4TcpRmem = other.Net.NetIpv4TcpRmem
		}
		if other.Net.NetIpv4TcpWmem != nil {
			t.Net.NetIpv4TcpWmem = other.Net.NetIpv4TcpWmem
		}
		if other.Net.NetIpv4TcpCongestionControl != nil {
			t.Net.NetIpv4TcpCongestionControl = other.Net.NetIpv4TcpCongestionControl
		}
		if other.Net.NetIpv4TcpFastopen != nil {
			t.Net.NetIpv4TcpFastopen = other.Net.NetIpv4TcpFastopen
		}
		if other.Net.NetIpv4TcpTimestamps != nil {
			t.Net.NetIpv4TcpTimestamps = other.Net.NetIpv4TcpTimestamps
		}
		if other.Net.NetIpv4TcpSack != nil {
			t.Net.NetIpv4TcpSack = other.Net.NetIpv4TcpSack
		}
		if other.Net.NetIpv4TcpLowLatency != nil {
			t.Net.NetIpv4TcpLowLatency = other.Net.NetIpv4TcpLowLatency
		}
		if other.Net.NetIpv4TcpTwReuse != nil {
			t.Net.NetIpv4TcpTwReuse = other.Net.NetIpv4TcpTwReuse
		}
		if other.Net.NetIpv4TcpNoMetricsSave != nil {
			t.Net.NetIpv4TcpNoMetricsSave = other.Net.NetIpv4TcpNoMetricsSave
		}
		if other.Net.NetIpv4TcpModerateRcvbuf != nil {
			t.Net.NetIpv4TcpModerateRcvbuf = other.Net.NetIpv4TcpModerateRcvbuf
		}
		if other.Net.NetCoreRmemMax != nil {
			t.Net.NetCoreRmemMax = other.Net.NetCoreRmemMax
		}
		if other.Net.NetCoreRmemDefault != nil {
			t.Net.NetCoreRmemDefault = other.Net.NetCoreRmemDefault
		}
		if other.Net.NetCoreWmemMax != nil {
			t.Net.NetCoreWmemMax = other.Net.NetCoreWmemMax
		}
		if other.Net.NetCoreWmemDefault != nil {
			t.Net.NetCoreWmemDefault = other.Net.NetCoreWmemDefault
		}
	}

	if other.Kernel != nil {
		if t.Kernel == nil {
			t.Kernel = &TunerKernelParams{}
		}
		if other.Kernel.KernelTimerMigration != nil {
			t.Kernel.KernelTimerMigration = other.Kernel.KernelTimerMigration
		}
		if other.Kernel.KernelNmiWatchdog != nil {
			t.Kernel.KernelNmiWatchdog = other.Kernel.KernelNmiWatchdog
		}
		if other.Kernel.KernelSchedMinGranularityNs != nil {
			t.Kernel.KernelSchedMinGranularityNs = other.Kernel.KernelSchedMinGranularityNs
		}
		if other.Kernel.KernelSchedWakeupGranularityNs != nil {
			t.Kernel.KernelSchedWakeupGranularityNs = other.Kernel.KernelSchedWakeupGranularityNs
		}
		if other.Kernel.KernelHungTaskTimeoutSecs != nil {
			t.Kernel.KernelHungTaskTimeoutSecs = other.Kernel.KernelHungTaskTimeoutSecs
		}
		if other.Kernel.KernelPidMax != nil {
			t.Kernel.KernelPidMax = other.Kernel.KernelPidMax
		}
	}

	if other.Vm != nil {
		if t.Vm == nil {
			t.Vm = &TunerVmParams{}
		}
		if other.Vm.VmSwappiness != nil {
			t.Vm.VmSwappiness = other.Vm.VmSwappiness
		}
		if other.Vm.VmMaxMapCount != nil {
			t.Vm.VmMaxMapCount = other.Vm.VmMaxMapCount
		}
		if other.Vm.VmStatInterval != nil {
			t.Vm.VmStatInterval = other.Vm.VmStatInterval
		}
		if other.Vm.VmDirtyRatio != nil {
			t.Vm.VmDirtyRatio = other.Vm.VmDirtyRatio
		}
		if other.Vm.VmDirtyBackgroundRatio != nil {
			t.Vm.VmDirtyBackgroundRatio = other.Vm.VmDirtyBackgroundRatio
		}
		if other.Vm.VmMinFreeKbytes != nil {
			t.Vm.VmMinFreeKbytes = other.Vm.VmMinFreeKbytes
		}
		if other.Vm.VmDirtyExpireCentisecs != nil {
			t.Vm.VmDirtyExpireCentisecs = other.Vm.VmDirtyExpireCentisecs
		}
		if other.Vm.VmDirtyWritebackCentisecs != nil {
			t.Vm.VmDirtyWritebackCentisecs = other.Vm.VmDirtyWritebackCentisecs
		}
		if other.Vm.VmDirtytimeExpireSeconds != nil {
			t.Vm.VmDirtytimeExpireSeconds = other.Vm.VmDirtytimeExpireSeconds
		}
	}

	if other.CpuGovernor != nil {
		t.CpuGovernor = other.CpuGovernor
	}
}
