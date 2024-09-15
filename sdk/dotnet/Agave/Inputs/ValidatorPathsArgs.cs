// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Svm.Agave.Inputs
{

    public sealed class ValidatorPathsArgs : global::Pulumi.ResourceArgs
    {
        [Input("accounts", required: true)]
        public Input<string> Accounts { get; set; } = null!;

        [Input("ledger", required: true)]
        public Input<string> Ledger { get; set; } = null!;

        [Input("log", required: true)]
        public Input<string> Log { get; set; } = null!;

        public ValidatorPathsArgs()
        {
        }
        public static new ValidatorPathsArgs Empty => new ValidatorPathsArgs();
    }
}
