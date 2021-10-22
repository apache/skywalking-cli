// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package interceptor

import "github.com/urfave/cli/v2"

const (
	versionIDFlagName   = "version-id"
	versionNameFlagName = "version-name"
)

// ParseVersion parses the service instance id or service instance name,
// and converts the present one to the missing one.
// See flags.VersionFlags.
func ParseVersion(required bool) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		if err := ParseBrowserService(required)(ctx); err != nil {
			return err
		}
		return parseInstance(required, versionIDFlagName, versionNameFlagName, serviceIDFlagName)(ctx)
	}
}
