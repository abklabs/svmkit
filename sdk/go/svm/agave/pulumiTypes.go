// Code generated by pulumi-language-go DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package agave

import (
	"context"
	"reflect"

	"example.com/pulumi-svm/sdk/go/svm/internal"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var _ = internal.GetEnvOrDefault

type ValidatorFlags struct {
	BlockProductionMethod        string         `pulumi:"blockProductionMethod"`
	DynamicPortRange             string         `pulumi:"dynamicPortRange"`
	EntryPoint                   []string       `pulumi:"entryPoint"`
	ExpectedGenesisHash          string         `pulumi:"expectedGenesisHash"`
	FullRpcAPI                   bool           `pulumi:"fullRpcAPI"`
	FullSnapshotIntervalSlots    int            `pulumi:"fullSnapshotIntervalSlots"`
	GossipPort                   int            `pulumi:"gossipPort"`
	KnownValidator               []string       `pulumi:"knownValidator"`
	LimitLedgerSize              int            `pulumi:"limitLedgerSize"`
	NoVoting                     bool           `pulumi:"noVoting"`
	NoWaitForVoteToStartLeader   bool           `pulumi:"noWaitForVoteToStartLeader"`
	OnlyKnownRPC                 bool           `pulumi:"onlyKnownRPC"`
	Paths                        ValidatorPaths `pulumi:"paths"`
	PrivateRPC                   bool           `pulumi:"privateRPC"`
	RpcBindAddress               string         `pulumi:"rpcBindAddress"`
	RpcPort                      int            `pulumi:"rpcPort"`
	TvuReceiveThreads            int            `pulumi:"tvuReceiveThreads"`
	UseSnapshotArchivesAtStartup string         `pulumi:"useSnapshotArchivesAtStartup"`
	WalRecoveryMode              string         `pulumi:"walRecoveryMode"`
}

// ValidatorFlagsInput is an input type that accepts ValidatorFlagsArgs and ValidatorFlagsOutput values.
// You can construct a concrete instance of `ValidatorFlagsInput` via:
//
//	ValidatorFlagsArgs{...}
type ValidatorFlagsInput interface {
	pulumi.Input

	ToValidatorFlagsOutput() ValidatorFlagsOutput
	ToValidatorFlagsOutputWithContext(context.Context) ValidatorFlagsOutput
}

type ValidatorFlagsArgs struct {
	BlockProductionMethod        pulumi.StringInput      `pulumi:"blockProductionMethod"`
	DynamicPortRange             pulumi.StringInput      `pulumi:"dynamicPortRange"`
	EntryPoint                   pulumi.StringArrayInput `pulumi:"entryPoint"`
	ExpectedGenesisHash          pulumi.StringInput      `pulumi:"expectedGenesisHash"`
	FullRpcAPI                   pulumi.BoolInput        `pulumi:"fullRpcAPI"`
	FullSnapshotIntervalSlots    pulumi.IntInput         `pulumi:"fullSnapshotIntervalSlots"`
	GossipPort                   pulumi.IntInput         `pulumi:"gossipPort"`
	KnownValidator               pulumi.StringArrayInput `pulumi:"knownValidator"`
	LimitLedgerSize              pulumi.IntInput         `pulumi:"limitLedgerSize"`
	NoVoting                     pulumi.BoolInput        `pulumi:"noVoting"`
	NoWaitForVoteToStartLeader   pulumi.BoolInput        `pulumi:"noWaitForVoteToStartLeader"`
	OnlyKnownRPC                 pulumi.BoolInput        `pulumi:"onlyKnownRPC"`
	Paths                        ValidatorPathsInput     `pulumi:"paths"`
	PrivateRPC                   pulumi.BoolInput        `pulumi:"privateRPC"`
	RpcBindAddress               pulumi.StringInput      `pulumi:"rpcBindAddress"`
	RpcPort                      pulumi.IntInput         `pulumi:"rpcPort"`
	TvuReceiveThreads            pulumi.IntInput         `pulumi:"tvuReceiveThreads"`
	UseSnapshotArchivesAtStartup pulumi.StringInput      `pulumi:"useSnapshotArchivesAtStartup"`
	WalRecoveryMode              pulumi.StringInput      `pulumi:"walRecoveryMode"`
}

func (ValidatorFlagsArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*ValidatorFlags)(nil)).Elem()
}

func (i ValidatorFlagsArgs) ToValidatorFlagsOutput() ValidatorFlagsOutput {
	return i.ToValidatorFlagsOutputWithContext(context.Background())
}

func (i ValidatorFlagsArgs) ToValidatorFlagsOutputWithContext(ctx context.Context) ValidatorFlagsOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ValidatorFlagsOutput)
}

type ValidatorFlagsOutput struct{ *pulumi.OutputState }

func (ValidatorFlagsOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*ValidatorFlags)(nil)).Elem()
}

func (o ValidatorFlagsOutput) ToValidatorFlagsOutput() ValidatorFlagsOutput {
	return o
}

func (o ValidatorFlagsOutput) ToValidatorFlagsOutputWithContext(ctx context.Context) ValidatorFlagsOutput {
	return o
}

func (o ValidatorFlagsOutput) BlockProductionMethod() pulumi.StringOutput {
	return o.ApplyT(func(v ValidatorFlags) string { return v.BlockProductionMethod }).(pulumi.StringOutput)
}

func (o ValidatorFlagsOutput) DynamicPortRange() pulumi.StringOutput {
	return o.ApplyT(func(v ValidatorFlags) string { return v.DynamicPortRange }).(pulumi.StringOutput)
}

func (o ValidatorFlagsOutput) EntryPoint() pulumi.StringArrayOutput {
	return o.ApplyT(func(v ValidatorFlags) []string { return v.EntryPoint }).(pulumi.StringArrayOutput)
}

func (o ValidatorFlagsOutput) ExpectedGenesisHash() pulumi.StringOutput {
	return o.ApplyT(func(v ValidatorFlags) string { return v.ExpectedGenesisHash }).(pulumi.StringOutput)
}

func (o ValidatorFlagsOutput) FullRpcAPI() pulumi.BoolOutput {
	return o.ApplyT(func(v ValidatorFlags) bool { return v.FullRpcAPI }).(pulumi.BoolOutput)
}

func (o ValidatorFlagsOutput) FullSnapshotIntervalSlots() pulumi.IntOutput {
	return o.ApplyT(func(v ValidatorFlags) int { return v.FullSnapshotIntervalSlots }).(pulumi.IntOutput)
}

func (o ValidatorFlagsOutput) GossipPort() pulumi.IntOutput {
	return o.ApplyT(func(v ValidatorFlags) int { return v.GossipPort }).(pulumi.IntOutput)
}

func (o ValidatorFlagsOutput) KnownValidator() pulumi.StringArrayOutput {
	return o.ApplyT(func(v ValidatorFlags) []string { return v.KnownValidator }).(pulumi.StringArrayOutput)
}

func (o ValidatorFlagsOutput) LimitLedgerSize() pulumi.IntOutput {
	return o.ApplyT(func(v ValidatorFlags) int { return v.LimitLedgerSize }).(pulumi.IntOutput)
}

func (o ValidatorFlagsOutput) NoVoting() pulumi.BoolOutput {
	return o.ApplyT(func(v ValidatorFlags) bool { return v.NoVoting }).(pulumi.BoolOutput)
}

func (o ValidatorFlagsOutput) NoWaitForVoteToStartLeader() pulumi.BoolOutput {
	return o.ApplyT(func(v ValidatorFlags) bool { return v.NoWaitForVoteToStartLeader }).(pulumi.BoolOutput)
}

func (o ValidatorFlagsOutput) OnlyKnownRPC() pulumi.BoolOutput {
	return o.ApplyT(func(v ValidatorFlags) bool { return v.OnlyKnownRPC }).(pulumi.BoolOutput)
}

func (o ValidatorFlagsOutput) Paths() ValidatorPathsOutput {
	return o.ApplyT(func(v ValidatorFlags) ValidatorPaths { return v.Paths }).(ValidatorPathsOutput)
}

func (o ValidatorFlagsOutput) PrivateRPC() pulumi.BoolOutput {
	return o.ApplyT(func(v ValidatorFlags) bool { return v.PrivateRPC }).(pulumi.BoolOutput)
}

func (o ValidatorFlagsOutput) RpcBindAddress() pulumi.StringOutput {
	return o.ApplyT(func(v ValidatorFlags) string { return v.RpcBindAddress }).(pulumi.StringOutput)
}

func (o ValidatorFlagsOutput) RpcPort() pulumi.IntOutput {
	return o.ApplyT(func(v ValidatorFlags) int { return v.RpcPort }).(pulumi.IntOutput)
}

func (o ValidatorFlagsOutput) TvuReceiveThreads() pulumi.IntOutput {
	return o.ApplyT(func(v ValidatorFlags) int { return v.TvuReceiveThreads }).(pulumi.IntOutput)
}

func (o ValidatorFlagsOutput) UseSnapshotArchivesAtStartup() pulumi.StringOutput {
	return o.ApplyT(func(v ValidatorFlags) string { return v.UseSnapshotArchivesAtStartup }).(pulumi.StringOutput)
}

func (o ValidatorFlagsOutput) WalRecoveryMode() pulumi.StringOutput {
	return o.ApplyT(func(v ValidatorFlags) string { return v.WalRecoveryMode }).(pulumi.StringOutput)
}

type ValidatorKeyPairs struct {
	Identity    string `pulumi:"identity"`
	VoteAccount string `pulumi:"voteAccount"`
}

// ValidatorKeyPairsInput is an input type that accepts ValidatorKeyPairsArgs and ValidatorKeyPairsOutput values.
// You can construct a concrete instance of `ValidatorKeyPairsInput` via:
//
//	ValidatorKeyPairsArgs{...}
type ValidatorKeyPairsInput interface {
	pulumi.Input

	ToValidatorKeyPairsOutput() ValidatorKeyPairsOutput
	ToValidatorKeyPairsOutputWithContext(context.Context) ValidatorKeyPairsOutput
}

type ValidatorKeyPairsArgs struct {
	Identity    pulumi.StringInput `pulumi:"identity"`
	VoteAccount pulumi.StringInput `pulumi:"voteAccount"`
}

func (ValidatorKeyPairsArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*ValidatorKeyPairs)(nil)).Elem()
}

func (i ValidatorKeyPairsArgs) ToValidatorKeyPairsOutput() ValidatorKeyPairsOutput {
	return i.ToValidatorKeyPairsOutputWithContext(context.Background())
}

func (i ValidatorKeyPairsArgs) ToValidatorKeyPairsOutputWithContext(ctx context.Context) ValidatorKeyPairsOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ValidatorKeyPairsOutput)
}

type ValidatorKeyPairsOutput struct{ *pulumi.OutputState }

func (ValidatorKeyPairsOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*ValidatorKeyPairs)(nil)).Elem()
}

func (o ValidatorKeyPairsOutput) ToValidatorKeyPairsOutput() ValidatorKeyPairsOutput {
	return o
}

func (o ValidatorKeyPairsOutput) ToValidatorKeyPairsOutputWithContext(ctx context.Context) ValidatorKeyPairsOutput {
	return o
}

func (o ValidatorKeyPairsOutput) Identity() pulumi.StringOutput {
	return o.ApplyT(func(v ValidatorKeyPairs) string { return v.Identity }).(pulumi.StringOutput)
}

func (o ValidatorKeyPairsOutput) VoteAccount() pulumi.StringOutput {
	return o.ApplyT(func(v ValidatorKeyPairs) string { return v.VoteAccount }).(pulumi.StringOutput)
}

type ValidatorPaths struct {
	Accounts string `pulumi:"accounts"`
	Ledger   string `pulumi:"ledger"`
	Log      string `pulumi:"log"`
}

// ValidatorPathsInput is an input type that accepts ValidatorPathsArgs and ValidatorPathsOutput values.
// You can construct a concrete instance of `ValidatorPathsInput` via:
//
//	ValidatorPathsArgs{...}
type ValidatorPathsInput interface {
	pulumi.Input

	ToValidatorPathsOutput() ValidatorPathsOutput
	ToValidatorPathsOutputWithContext(context.Context) ValidatorPathsOutput
}

type ValidatorPathsArgs struct {
	Accounts pulumi.StringInput `pulumi:"accounts"`
	Ledger   pulumi.StringInput `pulumi:"ledger"`
	Log      pulumi.StringInput `pulumi:"log"`
}

func (ValidatorPathsArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*ValidatorPaths)(nil)).Elem()
}

func (i ValidatorPathsArgs) ToValidatorPathsOutput() ValidatorPathsOutput {
	return i.ToValidatorPathsOutputWithContext(context.Background())
}

func (i ValidatorPathsArgs) ToValidatorPathsOutputWithContext(ctx context.Context) ValidatorPathsOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ValidatorPathsOutput)
}

type ValidatorPathsOutput struct{ *pulumi.OutputState }

func (ValidatorPathsOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*ValidatorPaths)(nil)).Elem()
}

func (o ValidatorPathsOutput) ToValidatorPathsOutput() ValidatorPathsOutput {
	return o
}

func (o ValidatorPathsOutput) ToValidatorPathsOutputWithContext(ctx context.Context) ValidatorPathsOutput {
	return o
}

func (o ValidatorPathsOutput) Accounts() pulumi.StringOutput {
	return o.ApplyT(func(v ValidatorPaths) string { return v.Accounts }).(pulumi.StringOutput)
}

func (o ValidatorPathsOutput) Ledger() pulumi.StringOutput {
	return o.ApplyT(func(v ValidatorPaths) string { return v.Ledger }).(pulumi.StringOutput)
}

func (o ValidatorPathsOutput) Log() pulumi.StringOutput {
	return o.ApplyT(func(v ValidatorPaths) string { return v.Log }).(pulumi.StringOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*ValidatorFlagsInput)(nil)).Elem(), ValidatorFlagsArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*ValidatorKeyPairsInput)(nil)).Elem(), ValidatorKeyPairsArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*ValidatorPathsInput)(nil)).Elem(), ValidatorPathsArgs{})
	pulumi.RegisterOutputType(ValidatorFlagsOutput{})
	pulumi.RegisterOutputType(ValidatorKeyPairsOutput{})
	pulumi.RegisterOutputType(ValidatorPathsOutput{})
}
