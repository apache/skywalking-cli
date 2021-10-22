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

package flags

import "github.com/urfave/cli/v2"

// VersionFlags take either browser service version id or browser service version name as input,
// and transform to the other one.
var VersionFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "version-id",
		Usage:    "`version id`, if you don't have version id, use `--version-name` instead",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "version-name",
		Usage:    "`version-name`, if you already have version id, prefer to use `--version-id`",
		Required: false,
	},
}
