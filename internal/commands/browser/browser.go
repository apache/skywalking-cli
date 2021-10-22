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

package browser

import (
	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/commands/browser/logs"
	"github.com/apache/skywalking-cli/internal/commands/browser/page"
	"github.com/apache/skywalking-cli/internal/commands/browser/service"
	"github.com/apache/skywalking-cli/internal/commands/browser/version"
)

var Command = &cli.Command{
	Name:    "browser",
	Aliases: []string{"b"},
	Usage:   "Browser related sub-command",
	Subcommands: cli.Commands{
		service.Command,
		version.Command,
		page.Command,
		logs.Command,
	},
}
