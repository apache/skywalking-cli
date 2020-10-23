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
	"testing"

	"github.com/apache/skywalking-cli/graphql/schema"
)

func TestParseScope(t *testing.T) {
	tests := []struct {
		name        string
		wantedScope schema.Scope
	}{
		{
			name:        "",
			wantedScope: schema.ScopeAll,
		},
		{
			name:        "all_percentile",
			wantedScope: schema.ScopeAll,
		},
		{
			name:        "all_heatmap",
			wantedScope: schema.ScopeAll,
		},
		{
			name:        "service_resp_time",
			wantedScope: schema.ScopeService,
		},
		{
			name:        "service_percentile",
			wantedScope: schema.ScopeService,
		},
		{
			name:        "service_relation_server_percentile ",
			wantedScope: schema.ScopeServiceRelation,
		},
		{
			name:        "service_instance_relation_client_cpm",
			wantedScope: schema.ScopeServiceInstanceRelation,
		},
		{
			name:        "service_instance_resp_time",
			wantedScope: schema.ScopeServiceInstance,
		},
		{
			name:        "endpoint_cpm",
			wantedScope: schema.ScopeEndpoint,
		},
		{
			name:        "endpoint_relation_resp_time",
			wantedScope: schema.ScopeEndpointRelation,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotScope := ParseScope(tt.name)
			if gotScope != tt.wantedScope {
				t.Errorf("ParseScope() got scope = %v, wanted scope %v", gotScope, tt.wantedScope)
			}
		})
	}
}
