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

package flags

import (
	"github.com/urfave/cli/v2"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/apache/skywalking-cli/internal/model"
)

var startEndUsage = `"start" and "end" specify a time range during which the query is preformed,
		they can be absolute time like "2019-01-01 12", "2019-01-01 1213", or relative time (to the
		current time) like "-30m", "30m". They are both optional and their default values follow the rules below: 
		1. when "start" and "end" are both absent, "start = now - 30 minutes" and "end = now", 
		namely past 30 minutes; 
		2. when "start" and "end" are both present, they are aligned to the same precision by 
		truncating the more precise one, e.g. if "start = 2019-01-01 1234, end = 2019-01-01 18", 
		then "start" is truncated (because it's more precise) to "2019-01-01 12", and "end = 2019-01-01 18"; 
		3. when "start" is absent and "end" is present, will determine the precision of "end" 
		and then use the precision to calculate "start" (minus 30 units), e.g. "end = 2019-11-09 1234", 
		the precision is "MINUTE",  so "start = end - 30 minutes = 2019-11-09 1204", 
		and if "end = 2019-11-09 12", the precision is "HOUR", so "start = end - 30HOUR = 2019-11-08 06"; 
		4. when "start" is present and "end" is absent, will determine the precision of "start" 
		and then use the precision to calculate "end" (plus 30 units), e.g. "start = 2019-11-09 1204", 
		the precision is "MINUTE", so "end = start + 30 minutes = 2019-11-09 1234", 
		and if "start = 2019-11-08 06", the precision is "HOUR", so "end = start + 30HOUR = 2019-11-09 12".
		Examples:
		1. Query the metrics from 20 minutes ago to 10 minutes ago
		$ swctl metrics linear --name=service_resp_time --service-name business-zone::projectB --start "-20m" --end "-10m"
		2. Query the metrics from 1 hour ago to 10 minutes ago
		$ swctl metrics linear --name=service_resp_time --service-name business-zone::projectB --start "-1h" --end "-10m"
		3. Query the metrics from 1 hour ago to now
		$ swctl metrics linear --name=service_resp_time --service-name business-zone::projectB --start "-1h" --end "0m"
		4. Query the metrics from "2021-10-26 1047" to "2021-10-26 1127"
		$ swctl metrics linear --name=service_resp_time --service-name business-zone::projectB --start "2021-10-26 1047" --end "2021-10-26 1127"`

// DurationFlags are healthcheck flags that involves a duration, composed
// by a start time, an end time, and a step, which is commonly used
// in most of the commands
var DurationFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "start",
		Usage: startEndUsage,
	},
	&cli.StringFlag{
		Name:  "end",
		Usage: `end time of the query duration. Check the usage of "start"`,
	},
	&cli.GenericFlag{
		Name:  "step",
		Usage: `time step between start time and end time, should be one of SECOND, MINUTE, HOUR, DAY`,
		Value: &model.StepEnumValue{
			Enum:     api.AllStep,
			Default:  api.StepMinute,
			Selected: api.StepMinute,
		},
	},
	&cli.StringFlag{
		Name:   "duration-type",
		Usage:  "the type of duration",
		Hidden: true,
	},
}
