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

	"github.com/apache/skywalking-cli/internal/logger"
)

// MetricsValuesArrayToMap converts Array of MetricsValues into a map that uses time as key.
func MetricsValuesArrayToMap(duration api.Duration, mvArray []api.MetricsValues, labelsMap map[string]string) map[string]map[string]float64 {
	ret := make(map[string]map[string]float64, len(mvArray))
	for _, mvs := range mvArray {
		label := *mvs.Label
		if l, ok := labelsMap[label]; ok {
			label = l
		}
		ret[label] = MetricsValuesToMap(duration, mvs)
	}
	return ret
}

// MetricsValuesToMap converts MetricsValues into a map that uses time as key.
func MetricsValuesToMap(duration api.Duration, metricsValues api.MetricsValues) map[string]float64 {
	kvInts := metricsValues.Values.Values
	ret := map[string]float64{}
	format := StepFormats[duration.Step]
	startTime, err := time.Parse(format, duration.Start)

	if err != nil {
		logger.Log.Fatalln(err)
	}

	step := StepDuration[duration.Step]
	for idx, value := range kvInts {
		ret[startTime.Add(time.Duration(idx)*step).Format(format)] = float64(value.Value)
	}

	return ret
}

// HeatMapToMap converts a HeatMap into a map that uses time as key.
func HeatMapToMap(hp *api.HeatMap) map[string][]int64 {
	ret := make(map[string][]int64)
	for _, col := range hp.Values {
		// col.id is a string represents date, like "202007292131",
		// extracts its time part as key.
		t := col.ID[8:10] + ":" + col.ID[10:12]

		// Reverse the array.
		for i, j := 0, len(col.Values)-1; i < j; i, j = i+1, j-1 {
			col.Values[i], col.Values[j] = col.Values[j], col.Values[i]
		}

		ret[t] = col.Values
	}
	return ret
}

// BucketsToStrings extracts strings from buckets as a chart's labels.
func BucketsToStrings(buckets []*api.Bucket) []string {
	var ret []string
	for _, b := range buckets {
		ret = append(ret, b.Min)
	}
	return ret
}
