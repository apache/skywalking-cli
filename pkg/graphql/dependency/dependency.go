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

package dependency

import (
	"context"

	"github.com/machinebox/graphql"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/apache/skywalking-cli/assets"
	"github.com/apache/skywalking-cli/pkg/graphql/client"
)

func EndpointDependency(ctx context.Context, endpointID string, duration api.Duration) (api.EndpointTopology, error) {
	var response map[string]api.EndpointTopology

	request := graphql.NewRequest(assets.Read("graphqls/dependency/EndpointDependency.graphql"))
	request.Var("endpointId", endpointID)
	request.Var("duration", duration)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func GlobalTopology(ctx context.Context, layer string, duration api.Duration) (api.Topology, error) {
	var response map[string]api.Topology

	request := graphql.NewRequest(assets.Read("graphqls/dependency/GlobalTopology.graphql"))
	request.Var("layer", layer)
	request.Var("duration", duration)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func GlobalTopologyWithoutLayer(ctx context.Context, duration api.Duration) (api.Topology, error) {
	var response map[string]api.Topology

	request := graphql.NewRequest(assets.Read("graphqls/dependency/GlobalTopologyWithoutLayer.graphql"))
	request.Var("duration", duration)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func ServiceTopology(ctx context.Context, serviceID string, duration api.Duration) (api.Topology, error) {
	var response map[string]api.Topology

	request := graphql.NewRequest(assets.Read("graphqls/dependency/ServiceTopology.graphql"))
	request.Var("serviceId", serviceID)
	request.Var("duration", duration)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func InstanceTopology(ctx context.Context, clientServiceID, serverServiceID string, duration api.Duration) (api.ServiceInstanceTopology, error) {
	var response map[string]api.ServiceInstanceTopology

	request := graphql.NewRequest(assets.Read("graphqls/dependency/InstanceTopology.graphql"))
	request.Var("clientServiceId", clientServiceID)
	request.Var("serverServiceId", serverServiceID)
	request.Var("duration", duration)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func ProcessTopology(ctx context.Context, instanceID string, duration api.Duration) (api.ProcessTopology, error) {
	var response map[string]api.ProcessTopology

	request := graphql.NewRequest(assets.Read("graphqls/dependency/ProcessTopology.graphql"))
	request.Var("serviceInstanceId", instanceID)
	request.Var("duration", duration)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}
