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
	"errors"
	"fmt"
	"github.com/apache/skywalking-cli/graphql/schema"
	"github.com/apache/skywalking-cli/logger"
	"github.com/machinebox/graphql"
	"github.com/urfave/cli"
)

func newClient(cliCtx *cli.Context) (client *graphql.Client) {
	client = graphql.NewClient(cliCtx.GlobalString("base-url"))
	client.Log = func(msg string) {
		logger.Log.Debugln(msg)
	}
	return
}

func executeQuery(cliCtx *cli.Context, request *graphql.Request, response interface{}) {
	client := newClient(cliCtx)
	ctx := context.Background()
	if err := client.Run(ctx, request, response); err != nil {
		logger.Log.Fatalln(err)
	}
}

func Services(cliCtx *cli.Context, duration schema.Duration) []schema.Service {
	var response map[string][]schema.Service
	request := graphql.NewRequest(`
		query ($duration: Duration!) {
			services: getAllServices(duration: $duration) {
				id name
			}
		}
	`)
	request.Var("duration", duration)

	executeQuery(cliCtx, request, &response)
	return response["services"]
}

func Instances(cliCtx *cli.Context, nameOrID string, duration schema.Duration) []schema.ServiceInstance {
	var serviceId string
	service, err := searchServices(cliCtx, nameOrID, duration)
	if err != nil {
		serviceId = nameOrID
	} else {
		serviceId = service.ID
	}

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

	executeQuery(cliCtx, request, &response)
	return response["instances"]
}

func searchServices(cliCtx *cli.Context, serviceName string, duration schema.Duration) (service schema.Service, err error) {
	var response map[string][]schema.Service
	request := graphql.NewRequest(`
		query searchServices($keyword: String!, $duration: Duration!) {
    		service: searchServices(duration: $duration, keyword: $keyword) {
      				id name
			}
		}
	`)
	request.Var("keyword", serviceName)
	request.Var("duration", duration)

	executeQuery(cliCtx, request, &response)
	services := response["service"]
	if services == nil || len(services) < 1 {
		return service, errors.New(fmt.Sprintf("no such service [%s]", serviceName))
	}
	return services[0], nil
}
