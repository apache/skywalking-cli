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

package aggregation

import (
	"github.com/machinebox/graphql"
	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/graphql/client"
	"github.com/apache/skywalking-cli/graphql/schema"
)

func ServiceTopN(ctx *cli.Context, name string, topN int, duration schema.Duration, order schema.Order) []schema.TopNEntity {
	var response map[string][]schema.TopNEntity

	request := graphql.NewRequest(`
		query ($name: String!, $topN: Int!, $duration: Duration!, $order: Order!) {
			result: getServiceTopN(
				duration: $duration,
				name: $name,
				topN: $topN,
				order: $order
			) {
				id name value
			}
		}
	`)
	request.Var("name", name)
	request.Var("topN", topN)
	request.Var("duration", duration)
	request.Var("order", order)

	client.ExecuteQueryOrFail(ctx, request, &response)

	return response["result"]
}

func AllServiceInstanceTopN(ctx *cli.Context, name string, topN int, duration schema.Duration, order schema.Order) []schema.TopNEntity {
	var response map[string][]schema.TopNEntity

	request := graphql.NewRequest(`
		query ($name: String!, $topN: Int!, $duration: Duration!, $order: Order!) {
			result: getAllServiceInstanceTopN(
				duration: $duration,
				name: $name,
				topN: $topN,
				order: $order
			) {
				id name value
			}
		}
	`)
	request.Var("name", name)
	request.Var("topN", topN)
	request.Var("duration", duration)
	request.Var("order", order)

	client.ExecuteQueryOrFail(ctx, request, &response)

	return response["result"]
}

func ServiceInstanceTopN(ctx *cli.Context, serviceID, name string, topN int, duration schema.Duration, order schema.Order) []schema.TopNEntity {
	var response map[string][]schema.TopNEntity

	request := graphql.NewRequest(`
		query ($serviceId: ID!, $name: String!, $topN: Int!, $duration: Duration!, $order: Order!) {
			result: getServiceInstanceTopN(
				serviceId: $serviceId,
				duration: $duration,
				name: $name,
				topN: $topN,
				order: $order
			) {
				id name value
			}
		}
	`)
	request.Var("serviceId", serviceID)
	request.Var("name", name)
	request.Var("topN", topN)
	request.Var("duration", duration)
	request.Var("order", order)

	client.ExecuteQueryOrFail(ctx, request, &response)

	return response["result"]
}

func AllEndpointTopN(ctx *cli.Context, name string, topN int, duration schema.Duration, order schema.Order) []schema.TopNEntity {
	var response map[string][]schema.TopNEntity

	request := graphql.NewRequest(`
		query ($name: String!, $topN: Int!, $duration: Duration!, $order: Order!) {
			result: getAllEndpointTopN(
				duration: $duration,
				name: $name,
				topN: $topN,
				order: $order
			) {
				id name value
			}
		}
	`)
	request.Var("name", name)
	request.Var("topN", topN)
	request.Var("duration", duration)
	request.Var("order", order)

	client.ExecuteQueryOrFail(ctx, request, &response)

	return response["result"]
}

func EndpointTopN(ctx *cli.Context, serviceID, name string, topN int, duration schema.Duration, order schema.Order) []schema.TopNEntity {
	var response map[string][]schema.TopNEntity

	request := graphql.NewRequest(`
		query ($serviceId: ID!, $name: String!, $topN: Int!, $duration: Duration!, $order: Order!) {
			result: getEndpointTopN(
				serviceId: $serviceId,
				duration: $duration,
				name: $name,
				topN: $topN,
				order: $order
			) {
				id name value
			}
		}
	`)
	request.Var("serviceId", serviceID)
	request.Var("name", name)
	request.Var("topN", topN)
	request.Var("duration", duration)
	request.Var("order", order)

	client.ExecuteQueryOrFail(ctx, request, &response)

	return response["result"]
}
