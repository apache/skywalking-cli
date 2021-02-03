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

package healthcheck

import (
	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/pkg/healthcheck"

	"github.com/apache/skywalking-cli/internal/logger"
	hc "github.com/apache/skywalking-cli/pkg/graphql/healthcheck"
)

var Command = cli.Command{
	Name:    "checkHealth",
	Aliases: []string{"ch"},
	Usage:   "Check the health status of OAP server",
	Flags: []cli.Flag{
		cli.BoolTFlag{
			Name:     "grpc",
			Usage:    "Check gRPC by HealthCheck service",
			Required: false,
		},
		cli.BoolFlag{
			Name:     "grpcTLS",
			Usage:    "use TLS for gRPC",
			Required: false,
		},
	},
	Action: func(ctx *cli.Context) error {
		healthStatus, err := hc.CheckHealth(ctx)

		if err != nil {
			logger.Log.Fatalln(err)
		}

		if healthStatus.Score != 0 {
			return cli.NewExitError(healthStatus.Details, healthStatus.Score)
		}
		logger.Log.Println("OAP modules are healthy")
		if !ctx.BoolT("grpc") {
			return nil
		}
		retCode := healthcheck.HealthCheck(ctx.GlobalString("grpcAddr"), ctx.Bool("grpcTLS"))
		if retCode != 0 {
			return cli.NewExitError("gRPC: failed to check health", retCode)
		}
		logger.Log.Println("OAP gRPC is healthy")
		return nil
	},
}
