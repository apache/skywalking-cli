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

type MetricsName string

const (
	GlobalP50  MetricsName = "all_p50"
	GlobalP75  MetricsName = "all_p75"
	GlobalP90  MetricsName = "all_p90"
	GlobalP95  MetricsName = "all_p95"
	GlobalP99  MetricsName = "all_p99"
	ServiceP50 MetricsName = "service_p50"
	ServiceP75 MetricsName = "service_p75"
	ServiceP90 MetricsName = "service_p90"
	ServiceP95 MetricsName = "service_p95"
	ServiceP99 MetricsName = "service_p99"
)

func (e MetricsName) String() string {
	return string(e)
}
