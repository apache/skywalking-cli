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

package hierarchy

import (
	"context"

	"github.com/apache/skywalking-cli/assets"
	"github.com/apache/skywalking-cli/pkg/graphql/client"

	"github.com/machinebox/graphql"

	api "skywalking.apache.org/repo/goapi/query"
)

func ServiceHierarchy(ctx context.Context, serviceID, layer string) (api.ServiceHierarchy, error) {
	var response map[string]api.ServiceHierarchy

	request := graphql.NewRequest(assets.Read("graphqls/hierarchy/ServiceHierarchy.graphql"))
	request.Var("serviceId", serviceID)
	request.Var("layer", layer)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func InstanceHierarchy(ctx context.Context, instanceID, layer string) (api.InstanceHierarchy, error) {
	var response map[string]api.InstanceHierarchy

	request := graphql.NewRequest(assets.Read("graphqls/hierarchy/InstanceHierarchy.graphql"))
	request.Var("instanceId", instanceID)
	request.Var("layer", layer)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func ListLayerLevels(ctx context.Context) ([]api.LayerLevel, error) {
	var response map[string][]api.LayerLevel

	request := graphql.NewRequest(assets.Read("graphqls/hierarchy/ListLayerLevels.graphql"))

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}
