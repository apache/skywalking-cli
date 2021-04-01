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
	"reflect"
	"testing"

	api "skywalking.apache.org/repo/goapi/query"
)

func TestMetricsToMap(t *testing.T) {
	type args struct {
		duration      api.Duration
		metricsValues api.MetricsValues
	}
	tests := []struct {
		name string
		args args
		want map[string]float64
	}{
		{
			name: "Should convert to map",
			args: args{
				duration: api.Duration{
					Start: "2020-01-01 0000",
					End:   "2020-01-01 0007",
					Step:  api.StepMinute,
				},
				metricsValues: api.MetricsValues{
					Values: &api.IntValues{
						Values: []*api.KVInt{
							{Value: 0},
							{Value: 1},
							{Value: 2},
							{Value: 3},
							{Value: 4},
							{Value: 5},
							{Value: 6},
							{Value: 7},
						},
					},
				},
			},
			want: map[string]float64{
				"2020-01-01 0000": 0,
				"2020-01-01 0001": 1,
				"2020-01-01 0002": 2,
				"2020-01-01 0003": 3,
				"2020-01-01 0004": 4,
				"2020-01-01 0005": 5,
				"2020-01-01 0006": 6,
				"2020-01-01 0007": 7,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MetricsValuesToMap(tt.args.duration, tt.args.metricsValues); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetricsValuesToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
