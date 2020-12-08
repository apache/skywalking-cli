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
	"strings"

	"github.com/apache/skywalking-cli/graphql/schema"
)

// ParseScope defines the scope according to name's prefix.
func ParseScope(name string) schema.Scope {
	ret := schema.ScopeAll

	if strings.HasPrefix(name, "service_relation") {
		ret = schema.ScopeServiceRelation
	} else if strings.HasPrefix(name, "service_instance_relation") {
		ret = schema.ScopeServiceInstanceRelation
	} else if strings.HasPrefix(name, "service_instance") || strings.HasPrefix(name, "instance_") {
		ret = schema.ScopeServiceInstance
	} else if strings.HasPrefix(name, "service_") || strings.HasPrefix(name, "database_") {
		ret = schema.ScopeService
	} else if strings.HasPrefix(name, "endpoint_relation") {
		ret = schema.ScopeEndpointRelation
	} else if strings.HasPrefix(name, "endpoint_") {
		ret = schema.ScopeEndpoint
	}

	return ret
}
