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
	"reflect"
	"testing"
	"time"

	api "skywalking.apache.org/repo/goapi/query"
)

func TestParseDuration(t *testing.T) {
	now := time.Now()
	type args struct {
		start string
		end   string
	}
	timeFormat := "2006-01-02 1504"
	tests := []struct {
		name            string
		args            args
		wantedStartTime time.Time
		wantedEndTime   time.Time
		wantedStep      api.Step
	}{
		{
			name: "Should set current time if start is absent",
			args: args{
				start: "",
				end:   now.Format(timeFormat),
			},
			wantedStartTime: now.Add(-30 * time.Minute),
			wantedEndTime:   now,
			wantedStep:      api.StepMinute,
		},
		{
			name: "Should set current time if end is absent",
			args: args{
				start: now.Format(timeFormat),
				end:   "",
			},
			wantedStartTime: now,
			wantedEndTime:   now.Add(30 * time.Minute),
			wantedStep:      api.StepMinute,
		},
		{
			name: "Should keep both if both are present",
			args: args{
				start: now.Add(-10 * time.Minute).Format(timeFormat),
				end:   now.Add(10 * time.Minute).Format(timeFormat),
			},
			wantedStartTime: now.Add(-10 * time.Minute),
			wantedEndTime:   now.Add(10 * time.Minute),
			wantedStep:      api.StepMinute,
		},
		{
			name: "Should set both if both are absent",
			args: args{
				start: "",
				end:   "",
			},
			wantedStartTime: now.Add(-30 * time.Minute),
			wantedEndTime:   now,
			wantedStep:      api.StepMinute,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStartTime, gotEndTime, gotStep, _ := ParseDuration(tt.args.start, tt.args.end, "")
			current := gotStartTime.Truncate(time.Minute).Format(timeFormat)
			spec := tt.wantedStartTime.Truncate(time.Minute).Format(timeFormat)
			if !reflect.DeepEqual(current, spec) {
				t.Errorf(
					"ParseDuration() got start time = %v, wanted start time %v",
					gotStartTime.Truncate(time.Minute),
					tt.wantedStartTime.Truncate(time.Minute),
				)
			}
			if !reflect.DeepEqual(current, spec) {
				t.Errorf(
					"ParseDuration() got end time = %v, wanted end time %v",
					gotEndTime.Truncate(time.Minute),
					tt.wantedEndTime.Truncate(time.Minute),
				)
			}
			if gotStep != tt.wantedStep {
				t.Errorf("ParseDuration() got step = %v, wanted step %v", gotStep, tt.wantedStep)
			}
		})
	}
}

func TestAlignPrecision(t *testing.T) {
	type args struct {
		start string
		end   string
	}
	tests := []struct {
		name        string
		args        args
		wantedStart string
		wantedEnd   string
	}{
		{
			name: "Should keep both when same precision",
			args: args{
				start: "2019-01-01",
				end:   "2019-01-01",
			},
			wantedStart: "2019-01-01",
			wantedEnd:   "2019-01-01",
		},
		{
			name: "Should truncate start when it's less precise",
			args: args{
				start: "2019-01-01 1200",
				end:   "2019-01-01",
			},
			wantedStart: "2019-01-01",
			wantedEnd:   "2019-01-01",
		},
		{
			name: "Should truncate end when it's less precise",
			args: args{
				start: "2019-01-01",
				end:   "2019-01-01 1200",
			},
			wantedStart: "2019-01-01",
			wantedEnd:   "2019-01-01",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd := AlignPrecision(tt.args.start, tt.args.end)
			if gotStart != tt.wantedStart {
				t.Errorf("AlignPrecision() gotStart = %v, wantedStart %v", gotStart, tt.wantedStart)
			}
			if gotEnd != tt.wantedEnd {
				t.Errorf("AlignPrecision() gotEnd = %v, wantedStart %v", gotEnd, tt.wantedEnd)
			}
		})
	}
}
