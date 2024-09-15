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

// urn is a helper function to build an urn for running integration tests.
func urn(typ string) resource.URN {
	return resource.NewURN("stack", "proj", "",
		tokens.Type("svm:svm:"+typ), "name")
}

// Create a test server.
func provider() integration.Server {
	return integration.NewServer(svm.Name, semver.MustParse("0.0.1"), svm.Provider())

}
