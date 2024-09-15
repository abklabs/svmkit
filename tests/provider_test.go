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

package tests

import (
	"testing"

	"github.com/blang/semver"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	svm "github.com/abklabs/pulumi-svm/provider"
)

func TestKeyPairCreate(t *testing.T) {
	prov := provider()

	response, err := prov.Create(p.CreateRequest{
		Urn:        urn("KeyPair"),
		Properties: resource.PropertyMap{},
		Preview:    false,
	})

	require.NoError(t, err)

	publicKey := response.Properties["publicKey"].StringValue()
	privateKey := response.Properties["privateKey"].ArrayValue()
	jsonKey := response.Properties["json"].StringValue()

	assert.IsType(t, "", publicKey)
	assert.IsType(t, "", jsonKey)
	assert.IsType(t, []resource.PropertyValue{}, privateKey)
}

func TestValidatorPreview(t *testing.T) {
	prov := provider()

	_, err := prov.Create(p.CreateRequest{
		Urn: urn("Validator"),
		Properties: resource.PropertyMap{
			"keyPairs": resource.NewObjectProperty(resource.PropertyMap{
				"identity":    resource.NewStringProperty("dummy-identity-key"),
				"voteAccount": resource.NewStringProperty("dummy-vote-account-key"),
			}),
			"flags": resource.NewObjectProperty(resource.PropertyMap{
				"entryPoint": resource.NewArrayProperty([]resource.PropertyValue{
					resource.NewStringProperty("entrypoint.testnet.solana.com:8001"),
					resource.NewStringProperty("entrypoint2.testnet.solana.com:8001"),
					resource.NewStringProperty("entrypoint3.testnet.solana.com:8001"),
				}),
				"knownValidator": resource.NewArrayProperty([]resource.PropertyValue{
					resource.NewStringProperty("5D1fNXzvv5NjV1ysLjirC4WY92RNsVH18vjmcszZd8on"),
					resource.NewStringProperty("7XSY3MrYnK8vq693Rju17bbPkCN3Z7KvvfvJx4kdrsSY"),
					resource.NewStringProperty("Ft5fbkqNa76vnsjYNwjDZUXoTWpP7VYm3mtsaQckQADN"),
					resource.NewStringProperty("9QxCLckBiJc783jnMvXZubK4wH86Eqqvashtrwvcsgkv"),
				}),
				"expectedGenesisHash":          resource.NewStringProperty("4uhcVJyU9pJkvQyS88uRDiswHXSCkY3zQawwpjk2NsNY"),
				"useSnapshotArchivesAtStartup": resource.NewStringProperty("when-newest"),
				"rpcPort":                      resource.NewNumberProperty(8899),
				"privateRPC":                   resource.NewBoolProperty(true),
				"onlyKnownRPC":                 resource.NewBoolProperty(true),
				"dynamicPortRange":             resource.NewStringProperty("8002-8020"),
				"gossipPort":                   resource.NewNumberProperty(8001),
				"rpcBindAddress":               resource.NewStringProperty("0.0.0.0"),
				"walRecoveryMode":              resource.NewStringProperty("skip_any_corrupted_record"),
				"limitLedgerSize":              resource.NewNumberProperty(50000000),
				"blockProductionMethod":        resource.NewStringProperty("central-scheduler"),
				"noWaitForVoteToStartLeader":   resource.NewBoolProperty(false),
				"fullSnapshotIntervalSlots":    resource.NewNumberProperty(1000),
				"paths": resource.NewObjectProperty(resource.PropertyMap{
					"accounts": resource.NewStringProperty("dummy-accounts-path"),
					"ledger":   resource.NewStringProperty("dummy-ledger-path"),
					"log":      resource.NewStringProperty("dummy-log-path"),
				}),
			}),
		},
		Preview: true,
	})

	require.NoError(t, err)
}

// urn is a helper function to build an urn for running integration tests.
func urn(typ string) resource.URN {
	return resource.NewURN("stack", "proj", "",
		tokens.Type("svm:svm:"+typ), "name")
}

// Create a test server.
func provider() integration.Server {
	return integration.NewServer(svm.Name, semver.MustParse("0.0.1"), svm.Provider())

}
