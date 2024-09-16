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
	"github.com/abklabs/pulumi-svm/provider/pkg/svm"
	"github.com/abklabs/pulumi-svm/provider/pkg/svm/genesis"
	"github.com/abklabs/pulumi-svm/provider/pkg/svm/validator"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

const Name string = "svm"

func Provider() p.Provider {
	// We tell the provider what resources it needs to support.
	// In this case, a single custom resource.
	return infer.Provider(infer.Options{
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