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
	"github.com/apache/skywalking-cli/graphql/schema"

	"github.com/urfave/cli"
)

func ParseEntity(ctx *cli.Context) *schema.Entity {
	service := ctx.String("service")
	normal := ctx.BoolT("isNormal")
	instance := ctx.String("instance")
	endpoint := ctx.String("endpoint")

	destService := ctx.String("destService")
	destNormal := ctx.BoolT("isDestNormal")
	destInstance := ctx.String("destServiceInstance")
	destEndpoint := ctx.String("destEndpoint")

	return &schema.Entity{
		Scope:                   parseScope(ctx),
		ServiceName:             &service,
		Normal:                  &normal,
		ServiceInstanceName:     &instance,
		EndpointName:            &endpoint,
		DestServiceName:         &destService,
		DestNormal:              &destNormal,
		DestServiceInstanceName: &destInstance,
		DestEndpointName:        &destEndpoint,
	}
}

// parseScope defines the scope based on the input parameters.
func parseScope(ctx *cli.Context) schema.Scope {
	ret := schema.ScopeAll

	if ctx.String("destEndpoint") != "" {
		ret = schema.ScopeEndpointRelation
	} else if ctx.String("destInstance") != "" {
		ret = schema.ScopeServiceInstanceRelation
	} else if ctx.String("destService") != "" {
		ret = schema.ScopeServiceRelation
	} else if ctx.String("endpointName") != "" {
		ret = schema.ScopeEndpoint
	} else if ctx.String("instanceName") != "" {
		ret = schema.ScopeServiceInstance
	} else if ctx.String("serviceName") != "" {
		ret = schema.ScopeService
	}

	return ret
}
