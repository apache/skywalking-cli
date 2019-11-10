/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package client

import (
	"context"
	"github.com/apache/skywalking-cli/graphql/schema"
	"github.com/apache/skywalking-cli/logger"
	"github.com/machinebox/graphql"
	"github.com/urfave/cli"
)

func newGraphqlClient(cliCtx *cli.Context) (client *graphql.Client) {
	client = graphql.NewClient(cliCtx.GlobalString("base-url"))
	client.Log = func(msg string) {
		logger.Log.Debugln(msg)
	}
	return
}

func Services(cliCtx *cli.Context, duration schema.Duration) []schema.Service {
	client := newGraphqlClient(cliCtx)

	var response map[string][]schema.Service
	request := graphql.NewRequest(`
		query ($duration: Duration!) {
			services: getAllServices(duration: $duration) {
				id name
			}
		}
	`)
	request.Var("duration", duration)

	ctx := context.Background()
	if err := client.Run(ctx, request, &response); err != nil {
		logger.Log.Fatalln(err)
		panic(err)
	}

	return response["services"]
}

func Instances(cliCtx *cli.Context, serviceId string, duration schema.Duration) []schema.ServiceInstance {
	client := newGraphqlClient(cliCtx)

	var response map[string][]schema.ServiceInstance
	request := graphql.NewRequest(`
		query ($serviceId: ID!, $duration: Duration!) {
    		instances: getServiceInstances(duration: $duration, serviceId: $serviceId) {
      			id
      			name
      			language
				instanceUUID
      			attributes {
        			name
        			value
      			}
			}
		}
	`)
	request.Var("serviceId", serviceId)
	request.Var("duration", duration)

	ctx := context.Background()
	if err := client.Run(ctx, request, &response); err != nil {
		logger.Log.Fatalln(err)
		panic(err)
	}

	return response["instances"]
}
