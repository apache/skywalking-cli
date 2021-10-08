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

package metadata

import (
	"fmt"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/apache/skywalking-cli/assets"

	"github.com/machinebox/graphql"
	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/pkg/graphql/client"
)

func AllServices(cliCtx *cli.Context, duration api.Duration) ([]api.Service, error) {
	var response map[string][]api.Service

	request := graphql.NewRequest(assets.Read("graphqls/metadata/AllServices.graphql"))
	request.Var("duration", duration)

	err := client.ExecuteQuery(cliCtx, request, &response)

	return response["result"], err
}

func SearchService(cliCtx *cli.Context, serviceCode string) (service api.Service, err error) {
	var response map[string]api.Service

	request := graphql.NewRequest(assets.Read("graphqls/metadata/SearchService.graphql"))
	request.Var("serviceCode", serviceCode)

	err = client.ExecuteQuery(cliCtx, request, &response)

	service = response["result"]

	if service.ID == "" {
		return service, fmt.Errorf("no such service [%s]", serviceCode)
	}

	return service, err
}

func SearchEndpoints(cliCtx *cli.Context, serviceID, keyword string, limit int) ([]api.Endpoint, error) {
	var response map[string][]api.Endpoint

	request := graphql.NewRequest(assets.Read("graphqls/metadata/SearchEndpoints.graphql"))
	request.Var("serviceId", serviceID)
	request.Var("keyword", keyword)
	request.Var("limit", limit)

	err := client.ExecuteQuery(cliCtx, request, &response)

	return response["result"], err
}

func Instances(cliCtx *cli.Context, serviceID string, duration api.Duration) ([]api.ServiceInstance, error) {
	var response map[string][]api.ServiceInstance

	request := graphql.NewRequest(assets.Read("graphqls/metadata/Instances.graphql"))
	request.Var("serviceId", serviceID)
	request.Var("duration", duration)

	err := client.ExecuteQuery(cliCtx, request, &response)

	return response["result"], err
}

func ServerTimeInfo(cliCtx *cli.Context) (api.TimeInfo, error) {
	var response map[string]api.TimeInfo

	request := graphql.NewRequest(assets.Read("graphqls/metadata/ServerTimeInfo.graphql"))

	if err := client.ExecuteQuery(cliCtx, request, &response); err != nil {
		return api.TimeInfo{}, err
	}

	return response["result"], nil
}
