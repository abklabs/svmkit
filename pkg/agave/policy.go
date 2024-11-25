package agave

import (
	"github.com/abklabs/svmkit/pkg/runner"
)

type StartupPolicy struct {
	WaitForRPCHealth *bool `pulumi:"waitForRPCHealth,optional"`
}

type ShutdownPolicy struct {
	Force                *bool `pulumi:"force,optional"`
	SkipHealthCheck      *bool `pulumi:"skipHealthCheck,optional"`
	SkipNewSnapshotCheck *bool `pulumi:"skipNewSnapshotCheck,optional"`
	MaxDelinquentStake   *int  `pulumi:"maxDelinquentStake,optional"`
	MinIdleTime          *int  `pulumi:"minIdleTime,optional"`
}

func (s *ShutdownPolicy) ToFlags() *runner.FlagBuilder {
	f := &runner.FlagBuilder{}

	f.AppendBoolP("force", s.Force)
	f.AppendBoolP("skip-health-check", s.SkipHealthCheck)
	f.AppendBoolP("skip-new-snapshot-check", s.SkipNewSnapshotCheck)
	f.AppendIntP("max-delinquent-stake", s.MaxDelinquentStake)
	f.AppendIntP("min-idle-time", s.MinIdleTime)

	return f
}
