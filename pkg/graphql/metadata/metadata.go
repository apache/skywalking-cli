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
	"context"
	"fmt"
	"regexp"
	"strconv"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/machinebox/graphql"

	"github.com/apache/skywalking-cli/assets"
	"github.com/apache/skywalking-cli/pkg/graphql/client"
	"github.com/apache/skywalking-cli/pkg/graphql/common"
)

var backendVersion = regexp.MustCompile(`^(?P<Major>\d+)\.(?P<Minor>\d+)`)

func AllServices(cliCtx context.Context, duration api.Duration) ([]api.Service, error) {
	var response map[string][]api.Service

	version, err := protocolVersion(cliCtx)
	if err != nil {
		return nil, err
	}
	request := graphql.NewRequest(assets.Read("graphqls/metadata/" + version + "/AllServices.graphql"))
	request.Var("duration", duration)

	err = client.ExecuteQuery(cliCtx, request, &response)

	return response["result"], err
}

func SearchService(ctx context.Context, serviceCode string) (service api.Service, err error) {
	var response map[string]api.Service

	majorVersion, _, err := BackendVersion(ctx)
	if err != nil {
		return api.Service{}, err
	}
	var request *graphql.Request
	if majorVersion >= 9 {
		request = graphql.NewRequest(assets.Read("graphqls/metadata/v2/FindService.graphql"))
		request.Var("serviceName", serviceCode)
	} else {
		request = graphql.NewRequest(assets.Read("graphqls/metadata/v1/SearchService.graphql"))
		request.Var("serviceCode", serviceCode)
	}

	err = client.ExecuteQuery(ctx, request, &response)

	service = response["result"]

	if service.ID == "" {
		return service, fmt.Errorf("no such service [%s]", serviceCode)
	}

	return service, err
}

func AllBrowserServices(cliCtx context.Context, duration api.Duration) ([]api.Service, error) {
	var response map[string][]api.Service

	version, err := protocolVersion(cliCtx)
	if err != nil {
		return nil, err
	}
	request := graphql.NewRequest(assets.Read("graphqls/metadata/" + version + "/AllBrowserServices.graphql"))
	request.Var("duration", duration)

	err = client.ExecuteQuery(cliCtx, request, &response)

	return response["result"], err
}

func SearchBrowserService(cliCtx context.Context, serviceCode string) (service api.Service, err error) {
	var response map[string]api.Service

	version, err := protocolVersion(cliCtx)
	if err != nil {
		return api.Service{}, err
	}
	request := graphql.NewRequest(assets.Read("graphqls/metadata/" + version + "/SearchBrowserService.graphql"))
	request.Var("serviceCode", serviceCode)

	err = client.ExecuteQuery(cliCtx, request, &response)

	service = response["result"]

	if service.ID == "" {
		return service, fmt.Errorf("no such service [%s]", serviceCode)
	}

	return service, err
}

func SearchEndpoints(cliCtx context.Context, serviceID, keyword string, limit int, duration *api.Duration) ([]api.Endpoint, error) {
	var response map[string][]api.Endpoint

	majorVersion, minorVersion, err := BackendVersion(cliCtx)
	if err != nil {
		return nil, err
	}
	var request *graphql.Request
	if majorVersion >= 10 && minorVersion >= 2 {
		request = graphql.NewRequest(assets.Read("graphqls/metadata/v2/FindEndpointsWithDuration.graphql"))
		request.Var("serviceId", serviceID)
		request.Var("keyword", keyword)
		request.Var("limit", limit)
		request.Var("duration", duration)
	} else if majorVersion >= 9 {
		request = graphql.NewRequest(assets.Read("graphqls/metadata/v2/FindEndpointsWithoutDuration.graphql"))
		request.Var("serviceId", serviceID)
		request.Var("keyword", keyword)
		request.Var("limit", limit)
	} else {
		request = graphql.NewRequest(assets.Read("graphqls/metadata/v1/SearchEndpoints.graphql"))
		request.Var("serviceId", serviceID)
		request.Var("keyword", keyword)
		request.Var("limit", limit)
	}

	err = client.ExecuteQuery(cliCtx, request, &response)
	return response["result"], err
}

func Instances(cliCtx context.Context, serviceID string, duration api.Duration) ([]api.ServiceInstance, error) {
	var response map[string][]api.ServiceInstance

	version, err := protocolVersion(cliCtx)
	if err != nil {
		return nil, err
	}
	request := graphql.NewRequest(assets.Read("graphqls/metadata/" + version + "/Instances.graphql"))
	request.Var("serviceId", serviceID)
	request.Var("duration", duration)

	err = client.ExecuteQuery(cliCtx, request, &response)

	return response["result"], err
}

func GetInstance(cliCtx context.Context, instanceID string) (api.ServiceInstance, error) {
	var response map[string]api.ServiceInstance

	request := graphql.NewRequest(assets.Read("graphqls/metadata/v2/GetInstance.graphql"))
	request.Var("instanceId", instanceID)

	err := client.ExecuteQuery(cliCtx, request, &response)

	return response["result"], err
}

func GetEndpointInfo(cliCtx context.Context, endpointID string) (api.EndpointInfo, error) {
	var response map[string]api.EndpointInfo

	request := graphql.NewRequest(assets.Read("graphqls/metadata/v2/GetEndpointInfo.graphql"))
	request.Var("endpointId", endpointID)

	err := client.ExecuteQuery(cliCtx, request, &response)

	return response["result"], err
}

func Processes(cliCtx context.Context, instanceID string, duration api.Duration) ([]api.Process, error) {
	var response map[string][]api.Process

	request := graphql.NewRequest(assets.Read("graphqls/metadata/v2/Processes.graphql"))
	request.Var("instanceId", instanceID)
	request.Var("duration", duration)

	err := client.ExecuteQuery(cliCtx, request, &response)

	return response["result"], err
}

func GetProcess(cliCtx context.Context, processID string) (api.Process, error) {
	var response map[string]api.Process

	request := graphql.NewRequest(assets.Read("graphqls/metadata/v2/GetProcess.graphql"))
	request.Var("processId", processID)

	err := client.ExecuteQuery(cliCtx, request, &response)

	return response["result"], err
}

func EstimateProcessScale(cliCtx context.Context, serviceID string, labels []string) (int64, error) {
	var response map[string]int64

	request := graphql.NewRequest(assets.Read("graphqls/metadata/v2/EstimateProcessScale.graphql"))
	request.Var("serviceId", serviceID)
	request.Var("labels", labels)

	err := client.ExecuteQuery(cliCtx, request, &response)

	return response["result"], err
}

func ServerTimeInfo(cliCtx context.Context) (api.TimeInfo, error) {
	var response map[string]api.TimeInfo

	request := graphql.NewRequest(assets.Read("graphqls/metadata/v2/ServerTimeInfo.graphql"))

	if err := client.ExecuteQuery(cliCtx, request, &response); err != nil {
		return api.TimeInfo{}, err
	}

	return response["result"], nil
}

func ListLayers(cliCtx context.Context) ([]string, error) {
	var response map[string][]string

	request := graphql.NewRequest(assets.Read("graphqls/metadata/v2/ListLayers.graphql"))

	if err := client.ExecuteQuery(cliCtx, request, &response); err != nil {
		return make([]string, 0), err
	}

	return response["result"], nil
}

func ListLayerService(cliCtx context.Context, layer string) ([]api.Service, error) {
	var response map[string][]api.Service

	request := graphql.NewRequest(assets.Read("graphqls/metadata/v2/ListService.graphql"))
	request.Var("layer", layer)

	err := client.ExecuteQuery(cliCtx, request, &response)

	return response["result"], err
}

func BackendVersion(ctx context.Context) (major, minor int, err error) {
	version, err := common.Version(ctx)
	if err != nil {
		return 0, 0, err
	}
	if version == "" {
		return 0, 0, fmt.Errorf("failed to detect OAP version")
	}

	versions := backendVersion.FindStringSubmatch(version)
	if len(versions) != 3 {
		return 0, 0, fmt.Errorf("parsing OAP version failure: %s", version)
	}
	major, err = strconv.Atoi(versions[1])
	if err != nil {
		return 0, 0, fmt.Errorf("parse major failure: %s", version)
	}
	minor, err = strconv.Atoi(versions[2])
	if err != nil {
		return 0, 0, fmt.Errorf("parse minor failure: %s", version)
	}
	return major, minor, nil
}

func protocolVersion(ctx context.Context) (string, error) {
	if majorVersion, _, err := BackendVersion(ctx); err != nil {
		return "", err
	} else if majorVersion >= 9 {
		return "v2", nil
	}
	return "v1", nil
}
