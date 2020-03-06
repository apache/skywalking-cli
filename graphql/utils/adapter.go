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

	"github.com/apache/skywalking-cli/graphql/schema"
	"github.com/apache/skywalking-cli/logger"
)

type IntValues schema.IntValues

func MetricsToMap(duration schema.Duration, intValues schema.IntValues) map[string]float64 {
	kvInts := intValues.Values
	values := map[string]float64{}
	format := StepFormats[duration.Step]
	startTime, err := time.Parse(format, duration.Start)

	if err != nil {
		logger.Log.Fatalln(err)
	}

	step := StepDuration[duration.Step]
	for idx, value := range kvInts {
		values[startTime.Add(time.Duration(idx)*step).Format(format)] = float64(value.Value)
	}

	return values
}
