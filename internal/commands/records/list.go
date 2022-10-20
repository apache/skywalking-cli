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

package records

import (
	"github.com/apache/skywalking-cli/internal/commands/metrics/aggregation"

	"github.com/urfave/cli/v2"
)

var ListCommand = &cli.Command{
	Name:      "list",
	Aliases:   []string{"ls"},
	Usage:     "List records according to the specified options",
	ArgsUsage: "<n>",
	UsageText: `List the top <n> records according to the specified options.

Examples:
1. Query the top 5 database statements whose execute duration are largest:
$ swctl records list --name top_n_database_statement 5`,
	Flags:  aggregation.SampledRecords.Flags,
	Before: aggregation.SampledRecords.Before,
	Action: aggregation.SampledRecords.Action,
}
