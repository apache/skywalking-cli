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

package utils

import (
	"testing"

	api "skywalking.apache.org/repo/goapi/query"
)

//nolint:funlen // disable function length check for the test case count
func TestParseScope(t *testing.T) {
	empty := ""
	nonEmpty := "test"
	tests := []struct {
		name string
		args *api.Entity
		want api.Scope
	}{
		{
			name: "all of names are empty",
			args: &api.Entity{
				ServiceName:             &empty,
				ServiceInstanceName:     &empty,
				EndpointName:            &empty,
				ProcessName:             &empty,
				DestServiceName:         &empty,
				DestServiceInstanceName: &empty,
				DestEndpointName:        &empty,
				DestProcessName:         &empty,
			},
			want: api.ScopeAll,
		},
		{
			name: "all of names are not empty",
			args: &api.Entity{
				ServiceName:             &nonEmpty,
				ServiceInstanceName:     &nonEmpty,
				EndpointName:            &nonEmpty,
				ProcessName:             &nonEmpty,
				DestServiceName:         &nonEmpty,
				DestServiceInstanceName: &nonEmpty,
				DestEndpointName:        &nonEmpty,
				DestProcessName:         &nonEmpty,
			},
			want: api.ScopeProcessRelation,
		},
		{
			name: "only serviceName is not empty",
			args: &api.Entity{
				ServiceName:             &nonEmpty,
				ServiceInstanceName:     &empty,
				EndpointName:            &empty,
				ProcessName:             &empty,
				DestServiceName:         &empty,
				DestServiceInstanceName: &empty,
				DestEndpointName:        &empty,
				DestProcessName:         &empty,
			},
			want: api.ScopeService,
		},
		{
			name: "instanceName is not empty",
			args: &api.Entity{
				ServiceName:             &nonEmpty,
				ServiceInstanceName:     &nonEmpty,
				EndpointName:            &empty,
				ProcessName:             &empty,
				DestServiceName:         &empty,
				DestServiceInstanceName: &empty,
				DestEndpointName:        &empty,
				DestProcessName:         &empty,
			},
			want: api.ScopeServiceInstance,
		},
		{
			name: "endpointName is not empty",
			args: &api.Entity{
				ServiceName:             &nonEmpty,
				ServiceInstanceName:     &empty,
				EndpointName:            &nonEmpty,
				ProcessName:             &empty,
				DestServiceName:         &empty,
				DestServiceInstanceName: &empty,
				DestEndpointName:        &empty,
				DestProcessName:         &empty,
			},
			want: api.ScopeEndpoint,
		},
		{
			name: "destService is not empty",
			args: &api.Entity{
				ServiceName:             &nonEmpty,
				ServiceInstanceName:     &empty,
				EndpointName:            &empty,
				ProcessName:             &empty,
				DestServiceName:         &nonEmpty,
				DestServiceInstanceName: &empty,
				DestEndpointName:        &empty,
				DestProcessName:         &empty,
			},
			want: api.ScopeServiceRelation,
		},
		{
			name: "destInstance is not empty",
			args: &api.Entity{
				ServiceName:             &nonEmpty,
				ServiceInstanceName:     &nonEmpty,
				EndpointName:            &empty,
				ProcessName:             &empty,
				DestServiceName:         &nonEmpty,
				DestServiceInstanceName: &nonEmpty,
				DestEndpointName:        &empty,
				DestProcessName:         &empty,
			},
			want: api.ScopeServiceInstanceRelation,
		},
		{
			name: "destProcess is not empty",
			args: &api.Entity{
				ServiceName:             &nonEmpty,
				ServiceInstanceName:     &nonEmpty,
				EndpointName:            &empty,
				ProcessName:             &nonEmpty,
				DestServiceName:         &nonEmpty,
				DestServiceInstanceName: &nonEmpty,
				DestEndpointName:        &empty,
				DestProcessName:         &nonEmpty,
			},
			want: api.ScopeProcessRelation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseScope(tt.args); got != tt.want {
				t.Errorf("ParseScope() = %v, want %v", got, tt.want)
			}
		})
	}
}
