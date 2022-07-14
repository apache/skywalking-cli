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

package interceptor

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/logger"
	"github.com/apache/skywalking-cli/pkg/graphql/utils"

	api "skywalking.apache.org/repo/goapi/query"
)

func ParseEntity(ctx *cli.Context) (*api.Entity, error) {
	serviceID := ctx.String("service-id")
	instance := ctx.String("instance-name")
	endpoint := ctx.String("endpoint-name")
	process := ctx.String("process-name")

	destServiceID := ctx.String("dest-service-id")
	destInstance := ctx.String("dest-instance-name")
	destEndpoint := ctx.String("dest-endpoint-name")
	destProcess := ctx.String("dest-process-name")

	serviceName, isNormal, err := ParseServiceID(serviceID)
	if err != nil {
		return nil, err
	}

	destServiceName, destIsNormal, err := ParseServiceID(destServiceID)
	if err != nil {
		return nil, err
	}

	entity := &api.Entity{
		ServiceName:             &serviceName,
		Normal:                  &isNormal,
		ServiceInstanceName:     &instance,
		EndpointName:            &endpoint,
		ProcessName:             &process,
		DestServiceName:         &destServiceName,
		DestNormal:              &destIsNormal,
		DestServiceInstanceName: &destInstance,
		DestEndpointName:        &destEndpoint,
		DestProcessName:         &destProcess,
	}
	entity.Scope = utils.ParseScope(entity)

	// adapt for the old version of backend
	if *entity.ProcessName == "" {
		entity.ProcessName = nil
	}
	if *entity.DestProcessName == "" {
		entity.DestProcessName = nil
	}

	if logger.Log.GetLevel() <= logrus.DebugLevel {
		s, _ := json.Marshal(&entity)
		logger.Log.Debugf("entity: %+v", string(s))
	}

	return entity, nil
}

func ParseServiceID(id string) (name string, isNormal bool, err error) {
	if id == "" {
		return "", false, nil
	}
	parts := strings.Split(id, ".")
	if len(parts) != 2 {
		return "", false, fmt.Errorf("invalid service id, cannot be splitted into 2 parts. %v", id)
	}
	nameBytes, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return "", false, err
	}
	name = string(nameBytes)
	isNormal = parts[1] == "1"

	return name, isNormal, nil
}
