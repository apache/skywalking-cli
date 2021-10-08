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

import (
	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/model"

	api "skywalking.apache.org/repo/goapi/query"
)

// MetricsFlags can be reused in metrics commands.
var MetricsFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "name",
		Usage:    "`metrics` name, which should be defined in OAL files",
		Required: true,
	},
	&cli.GenericFlag{
		Name:  "scope",
		Usage: "the `scope` of the metrics entity, see https://skywalking.apache.org/docs/main/latest/en/concepts-and-designs/scope-definitions/",
		Value: &model.ScopeEnumValue{
			Enum:     api.AllScope,
			Default:  api.ScopeService,
			Selected: api.ScopeService,
		},
	},
}
