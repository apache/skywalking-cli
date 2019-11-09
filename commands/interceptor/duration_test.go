/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package interceptor

import (
	"github.com/apache/skywalking-cli/graphql/schema"
	"reflect"
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	now := time.Now().UTC()

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
		wantedStep      schema.Step
	}{
		{
			name: "Should set current time if start is absent",
			args: args{
				start: "",
				end:   now.Add(10 * time.Minute).Format(timeFormat),
			},
			wantedStartTime: now,
			wantedEndTime:   now.Add(10 * time.Minute),
			wantedStep:      schema.StepMinute,
		},
		{
			name: "Should set current time if end is absent",
			args: args{
				start: now.Add(-10 * time.Minute).Format(timeFormat),
				end:   "",
			},
			wantedStartTime: now.Add(-10 * time.Minute),
			wantedEndTime:   now,
			wantedStep:      schema.StepMinute,
		},
		{
			name: "Should keep both if both are present",
			args: args{
				start: now.Add(-10 * time.Minute).Format(timeFormat),
				end:   now.Add(10 * time.Minute).Format(timeFormat),
			},
			wantedStartTime: now.Add(-10 * time.Minute),
			wantedEndTime:   now.Add(10 * time.Minute),
			wantedStep:      schema.StepMinute,
		},
		{
			name: "Should set both if both are absent",
			args: args{
				start: "",
				end:   "",
			},
			wantedStartTime: now.Add(-30 * time.Minute),
			wantedEndTime:   now,
			wantedStep:      schema.StepMinute,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStartTime, gotEndTime, gotStep := ParseDuration(tt.args.start, tt.args.end)
			if !reflect.DeepEqual(gotStartTime.Truncate(time.Minute), tt.wantedStartTime.Truncate(time.Minute)) {
				t.Errorf("ParseDuration() got start time = %v, wanted start time %v", gotStartTime.Truncate(time.Minute), tt.wantedStartTime.Truncate(time.Minute))
			}
			if !reflect.DeepEqual(gotEndTime.Truncate(time.Minute), tt.wantedEndTime.Truncate(time.Minute)) {
				t.Errorf("ParseDuration() got end time = %v, wanted end time %v", gotEndTime.Truncate(time.Minute), tt.wantedEndTime.Truncate(time.Minute))
			}
			if gotStep != tt.wantedStep {
				t.Errorf("ParseDuration() got step = %v, wanted step %v", gotStep, tt.wantedStep)
			}
		})
	}
}
