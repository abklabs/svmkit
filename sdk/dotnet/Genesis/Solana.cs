// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Svm.Genesis
{
    [SvmResourceType("svm:genesis:Solana")]
    public partial class Solana : global::Pulumi.CustomResource
    {
        [Output("connection")]
        public Output<Pulumi.Svm.Ssh.Outputs.Connection> Connection { get; private set; } = null!;

        [Output("flags")]
        public Output<Pulumi.Svm.Solana.Outputs.GenesisFlags> Flags { get; private set; } = null!;

        [Output("genesisHash")]
        public Output<string> GenesisHash { get; private set; } = null!;

        [Output("primordial")]
        public Output<ImmutableArray<Outputs.PrimorialEntry>> Primordial { get; private set; } = null!;


        /// <summary>
        /// Create a Solana resource with the given unique name, arguments, and options.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resource</param>
        /// <param name="args">The arguments used to populate this resource's properties</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public Solana(string name, SolanaArgs args, CustomResourceOptions? options = null)
            : base("svm:genesis:Solana", name, args ?? new SolanaArgs(), MakeResourceOptions(options, ""))
        {
        }

        private Solana(string name, Input<string> id, CustomResourceOptions? options = null)
            : base("svm:genesis:Solana", name, null, MakeResourceOptions(options, id))
        {
        }

        private static CustomResourceOptions MakeResourceOptions(CustomResourceOptions? options, Input<string>? id)
        {
            var defaultOptions = new CustomResourceOptions
            {
                Version = Utilities.Version,
            };
            var merged = CustomResourceOptions.Merge(defaultOptions, options);
            // Override the ID if one was specified for consistency with other language SDKs.
            merged.Id = id ?? merged.Id;
            return merged;
        }
        /// <summary>
        /// Get an existing Solana resource's state with the given name, ID, and optional extra
        /// properties used to qualify the lookup.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resulting resource.</param>
        /// <param name="id">The unique provider ID of the resource to lookup.</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public static Solana Get(string name, Input<string> id, CustomResourceOptions? options = null)
        {
            return new Solana(name, id, options);
        }
    }

    public sealed class SolanaArgs : global::Pulumi.ResourceArgs
    {
        [Input("connection", required: true)]
        public Input<Pulumi.Svm.Ssh.Inputs.ConnectionArgs> Connection { get; set; } = null!;

        [Input("flags", required: true)]
        public Input<Pulumi.Svm.Solana.Inputs.GenesisFlagsArgs> Flags { get; set; } = null!;

        [Input("primordial", required: true)]
        private InputList<Inputs.PrimorialEntryArgs>? _primordial;
        public InputList<Inputs.PrimorialEntryArgs> Primordial
        {
            get => _primordial ?? (_primordial = new InputList<Inputs.PrimorialEntryArgs>());
            set => _primordial = value;
        }

        public SolanaArgs()
        {
        }
        public static new SolanaArgs Empty => new SolanaArgs();
    }
}