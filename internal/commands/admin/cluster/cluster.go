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

package cluster

import (
	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/pkg/admin/preflight"
	"github.com/apache/skywalking-cli/pkg/admin/status"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
)

var Command = &cli.Command{
	Name:  "cluster",
	Usage: "Inspect the OAP cluster from the admin-server `status` module",
	Subcommands: []*cli.Command{
		nodesCommand,
	},
}

var nodesCommand = &cli.Command{
	Name:  "nodes",
	Usage: "List the OAP cluster peer nodes (GET /status/cluster/nodes)",
	UsageText: `List the OAP cluster peer nodes as seen by the cluster coordinator.

Examples:
1. List cluster nodes:
$ swctl admin cluster nodes`,
	Action: func(ctx *cli.Context) error {
		nodes, err := status.ClusterNodesQuery(ctx.Context)
		if err != nil {
			return preflight.Explain(ctx.Context, err, preflight.ModuleStatus, "SW_STATUS")
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: nodes})
	},
}
