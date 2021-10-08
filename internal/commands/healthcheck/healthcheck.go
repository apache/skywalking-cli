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
	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/pkg/healthcheck"

	"github.com/apache/skywalking-cli/internal/logger"
	hc "github.com/apache/skywalking-cli/pkg/graphql/healthcheck"
)

var Command = &cli.Command{
	Name:  "health",
	Usage: "Checks whether OAP server is healthy",
	UsageText: `Checks whether OAP server is healthy.
Before using this, please make sure the OAP enables the health checker,
refer to https://skywalking.apache.org/docs/main/latest/en/setup/backend/backend-health-check/

Note: once enable gRPC TLS, checkHealth command would ignore server's cert.

Examples:
1. Check health status from GraphQL and the gRPC endpoint listening on 10.0.0.1:8843
$ swctl health --grpc-addr=10.0.0.1:8843

2. Once the gRPC endpoint of OAP encrypts communication by TLS:
$ swctl health --grpcTLS=true

3. Check health status from GraphQL service only:
$ swctl health --grpc=false
`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:     "grpc",
			Usage:    "check OAP gRPC by HealthCheck service",
			Required: false,
			Value:    true,
		},
		&cli.BoolFlag{
			Name:     "grpcTLS",
			Usage:    "use TLS for gRPC",
			Required: false,
		},
	},
	Action: func(ctx *cli.Context) error {
		healthStatus, err := hc.CheckHealth(ctx)

		if err != nil {
			return err
		}

		if healthStatus.Score != 0 {
			return cli.Exit(healthStatus.Details, healthStatus.Score)
		}
		logger.Log.Println("OAP modules are healthy")
		if !ctx.Bool("grpc") {
			return nil
		}
		retCode := healthcheck.HealthCheck(ctx.String("grpc-addr"), ctx.Bool("grpcTLS"))
		if retCode != 0 {
			return cli.Exit("gRPC: failed to check health", retCode)
		}
		logger.Log.Println("OAP gRPC is healthy")
		return nil
	},
}
