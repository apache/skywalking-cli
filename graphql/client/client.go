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

package client

import (
	"context"
	"fmt"
	"time"

	"github.com/machinebox/graphql"
	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/graphql/schema"
	"github.com/apache/skywalking-cli/logger"
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

func SearchEndpoints(cliCtx *cli.Context, serviceID, keyword string, limit int) []schema.Endpoint {
	var response map[string][]schema.Endpoint
	request := graphql.NewRequest(`
		query ($keyword: String!, $serviceId: ID!, $limit: Int!) {
			endpoints: searchEndpoint(keyword: $keyword, serviceId: $serviceId, limit: $limit) {
				id name
			}
		}
	`)
	request.Var("serviceId", serviceID)
	request.Var("keyword", keyword)
	request.Var("limit", limit)

	executeQuery(cliCtx, request, &response)
	return response["endpoints"]
}

func GetEndpointInfo(cliCtx *cli.Context, endpointID string) schema.Endpoint {
	var response map[string]schema.Endpoint
	request := graphql.NewRequest(`
		query ($endpointId: ID!) {
			endpoint: getEndpointInfo(endpointId: $endpointId) {
				id name
			}
		}
	`)
	request.Var("endpointId", endpointID)

	executeQuery(cliCtx, request, &response)
	return response["endpoint"]
}

func Instances(cliCtx *cli.Context, serviceID string, duration schema.Duration) []schema.ServiceInstance {
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
	request.Var("serviceId", serviceID)
	request.Var("duration", duration)

	executeQuery(cliCtx, request, &response)
	return response["instances"]
}

func SearchService(cliCtx *cli.Context, serviceCode string) (service schema.Service, err error) {
	var response map[string]schema.Service
	request := graphql.NewRequest(`
		query searchService($serviceCode: String!) {
			service: searchService(serviceCode: $serviceCode) {
				id name
			}
		}
	`)
	request.Var("serviceCode", serviceCode)

	executeQuery(cliCtx, request, &response)
	service = response["service"]
	if service.ID == "" {
		return service, fmt.Errorf("no such service [%s]", serviceCode)
	}
	return service, nil
}

func LinearIntValues(ctx *cli.Context, condition schema.MetricCondition, duration schema.Duration) map[string]float64 {
	var response map[string]schema.IntValues

	request := graphql.NewRequest(`
		query ($metric: MetricCondition!, $duration: Duration!) {
			metrics: getLinearIntValues(metric: $metric, duration: $duration) {
				values { value }
			}
		}
	`)
	request.Var("metric", condition)
	request.Var("duration", duration)

	executeQuery(ctx, request, &response)

	values := metricsToMap(duration, response["metrics"].Values)

	return values
}

func IntValues(ctx *cli.Context, condition schema.BatchMetricConditions, duration schema.Duration) []*schema.KVInt {
	var response map[string]schema.IntValues

	request := graphql.NewRequest(`
		query ($metric: BatchMetricConditions!, $duration: Duration!) {
			metrics: getValues(metric: $metric, duration: $duration) {
				values { id value }
			}
		}
	`)
	request.Var("metric", condition)
	request.Var("duration", duration)

	executeQuery(ctx, request, &response)

	return response["metrics"].Values
}

func metricsToMap(duration schema.Duration, kvInts []*schema.KVInt) map[string]float64 {
	values := map[string]float64{}
	format := schema.StepFormats[duration.Step]
	startTime, err := time.Parse(format, duration.Start)
	if err != nil {
		logger.Log.Fatalln(err)
	}
	step := schema.StepDuration[duration.Step]
	for idx, value := range kvInts {
		values[startTime.Add(time.Duration(idx)*step).Format(format)] = float64(value.Value)
	}
	return values
}
