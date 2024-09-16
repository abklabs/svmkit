// Copyright 2016-2023, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"github.com/abklabs/svmkit/provider/pkg/genesis"
	"github.com/abklabs/svmkit/provider/pkg/svm"
	"github.com/abklabs/svmkit/provider/pkg/validator"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-go-provider/middleware/schema"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

const Name string = "svmkit"

func Provider() p.Provider {
	// We tell the provider what resources it needs to support.
	// In this case, a single custom resource.
	return infer.Provider(infer.Options{
		Metadata: schema.Metadata{
			DisplayName: "Svmkit",
			Description: "The Pulumi Command Provider enables you to execute commands and scripts either locally or remotely as part of the Pulumi resource model.",
			Keywords: []string{
				"pulumi",
				"svmkit",
				"solana",
				"blockchain",
			},
			Homepage:   "https://abklabs.com",
			License:    "Apache-3.0",
			Repository: "https://github.com/abklabs/svmkit",
			Publisher:  "ABK Labs",
		},
		Resources: []infer.InferredResource{
			infer.Resource[svm.KeyPair, svm.KeyPairArgs, svm.KeyPairState](),
			infer.Resource[validator.Agave, validator.AgaveArgs, validator.AgaveState](),
			infer.Resource[genesis.Solana, genesis.SolanaArgs, genesis.SolanaState](),
		},
		ModuleMap: map[tokens.ModuleName]tokens.ModuleName{
			"svm": "index",
		},
	})
}
