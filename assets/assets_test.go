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

package assets

import "testing"

func TestRead(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Trim license header",
			args: args{
				filename: "graphqls/dependency/EndpointDependency.graphql",
			},
			want: `
query ($endpointId:ID!, $duration: Duration!) {
    result: getEndpointDependencies(duration: $duration, endpointId: $endpointId) {
        nodes {
            id
            name
            serviceId
            serviceName
            type
            isReal
        }
        calls {
            id
            source
            target
            detectPoints
        }
    }
}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Read(tt.args.filename); got != tt.want {
				t.Errorf("Read() = %v, want %v", got, tt.want)
			}
		})
	}
}
