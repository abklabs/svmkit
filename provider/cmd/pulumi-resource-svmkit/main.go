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

package main

import (
	"fmt"
	"os"
	"strings"

	p "github.com/pulumi/pulumi-go-provider"

	svmkit "github.com/abklabs/svmkit/provider"
	"github.com/abklabs/svmkit/provider/pkg/version"
)

// Serve the provider against Pulumi's Provider protocol.
func main() {
	version := strings.TrimPrefix(version.Version, "v")

	err := p.RunProvider(svmkit.Name, version, svmkit.Provider())

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}
}
