// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Svm.Agave.Outputs
{

    [OutputType]
    public sealed class ValidatorFlags
    {
        public readonly string BlockProductionMethod;
        public readonly string DynamicPortRange;
        public readonly ImmutableArray<string> EntryPoint;
        public readonly string ExpectedGenesisHash;
        public readonly bool FullRpcAPI;
        public readonly int FullSnapshotIntervalSlots;
        public readonly int GossipPort;
        public readonly ImmutableArray<string> KnownValidator;
        public readonly int LimitLedgerSize;
        public readonly bool NoVoting;
        public readonly bool NoWaitForVoteToStartLeader;
        public readonly bool OnlyKnownRPC;
        public readonly Outputs.ValidatorPaths Paths;
        public readonly bool PrivateRPC;
        public readonly string RpcBindAddress;
        public readonly int RpcPort;
        public readonly int TvuReceiveThreads;
        public readonly string UseSnapshotArchivesAtStartup;
        public readonly string WalRecoveryMode;

        [OutputConstructor]
        private ValidatorFlags(
            string blockProductionMethod,

            string dynamicPortRange,

            ImmutableArray<string> entryPoint,

            string expectedGenesisHash,

            bool fullRpcAPI,

            int fullSnapshotIntervalSlots,

            int gossipPort,

            ImmutableArray<string> knownValidator,

            int limitLedgerSize,

            bool noVoting,

            bool noWaitForVoteToStartLeader,

            bool onlyKnownRPC,

            Outputs.ValidatorPaths paths,

            bool privateRPC,

            string rpcBindAddress,

            int rpcPort,

            int tvuReceiveThreads,

            string useSnapshotArchivesAtStartup,

            string walRecoveryMode)
        {
            BlockProductionMethod = blockProductionMethod;
            DynamicPortRange = dynamicPortRange;
            EntryPoint = entryPoint;
            ExpectedGenesisHash = expectedGenesisHash;
            FullRpcAPI = fullRpcAPI;
            FullSnapshotIntervalSlots = fullSnapshotIntervalSlots;
            GossipPort = gossipPort;
            KnownValidator = knownValidator;
            LimitLedgerSize = limitLedgerSize;
            NoVoting = noVoting;
            NoWaitForVoteToStartLeader = noWaitForVoteToStartLeader;
            OnlyKnownRPC = onlyKnownRPC;
            Paths = paths;
            PrivateRPC = privateRPC;
            RpcBindAddress = rpcBindAddress;
            RpcPort = rpcPort;
            TvuReceiveThreads = tvuReceiveThreads;
            UseSnapshotArchivesAtStartup = useSnapshotArchivesAtStartup;
            WalRecoveryMode = walRecoveryMode;
        }
    }
}
