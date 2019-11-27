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

package schema

import "time"

// StepFormats is a mapping from schema.Step to its time format
var StepFormats = map[Step]string{
	StepSecond: "2006-01-02 150400",
	StepMinute: "2006-01-02 1504",
	StepHour:   "2006-01-02 15",
	StepDay:    "2006-01-02",
	StepMonth:  "2006-01",
}

// StepDuration is a mapping from schema.Step to its time.Duration
var StepDuration = map[Step]time.Duration{
	StepSecond: time.Second,
	StepMinute: time.Minute,
	StepHour:   time.Hour,
	StepDay:    time.Hour * 24,
	StepMonth:  time.Hour * 24 * 30,
}
