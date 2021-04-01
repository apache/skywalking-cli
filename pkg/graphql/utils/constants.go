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
	"time"

	api "skywalking.apache.org/repo/goapi/query"
)

// StepFormats is a mapping from schema.Step to its time format
var StepFormats = map[api.Step]string{
	api.StepSecond: "2006-01-02 150405",
	api.StepMinute: "2006-01-02 1504",
	api.StepHour:   "2006-01-02 15",
	api.StepDay:    "2006-01-02",
}

// StepDuration is a mapping from schema.Step to its time.Duration
var StepDuration = map[api.Step]time.Duration{
	api.StepSecond: time.Second,
	api.StepMinute: time.Minute,
	api.StepHour:   time.Hour,
	api.StepDay:    time.Hour * 24,
}

type DurationType string

const (
	BothAbsent  DurationType = "BothAbsent"
	BothPresent DurationType = "BothPresent"
	StartAbsent DurationType = "StartAbsent"
	EndAbsent   DurationType = "EndAbsent"
)

func (dt DurationType) String() string {
	return string(dt)
}
