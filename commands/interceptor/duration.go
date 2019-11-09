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
	"github.com/apache/skywalking-cli/logger"
	"github.com/urfave/cli"
	"time"
)

var stepFormats = map[schema.Step]string{
	schema.StepSecond: "2006-01-02 1504",
	schema.StepMinute: "2006-01-02 1504",
	schema.StepHour:   "2006-01-02 15",
	schema.StepDay:    "2006-01-02",
	schema.StepMonth:  "2006-01-02",
}

var supportedTimeLayouts = []string{
	"2006-01-02 150400",
	"2006-01-02 1504",
	"2006-01-02 15",
	"2006-01-02",
	"2006-01",
}

func tryParseTime(unparsed string, parsed *time.Time) error {
	var possibleError error = nil
	for _, layout := range supportedTimeLayouts {
		t, err := time.Parse(layout, unparsed)
		if err == nil {
			*parsed = t
			return nil
		}
		possibleError = err
	}
	return possibleError
}

// DurationInterceptor sets the duration if absent, and formats it accordingly,
// see ParseDuration
func DurationInterceptor(ctx *cli.Context) error {
	start := ctx.String("start")
	end := ctx.String("end")

	startTime, endTime, step := ParseDuration(start, end)

	if err := ctx.Set("start", startTime.Format(stepFormats[step])); err != nil {
		return err
	} else if err := ctx.Set("end", endTime.Format(stepFormats[step])); err != nil {
		return err
	} else if err := ctx.Set("step", step.String()); err != nil {
		return err
	}
	return nil
}

// ParseDuration parses the `start` and `end` to a triplet, (startTime, endTime, step)
// if --start and --end are both absent, then: start := now - 30min; end := now
// if --start is given, --end is absent, then: end := now
// if --start is absent, --end is given, then: start := end - 30min
// NOTE that when either(both) `start` or `end` is(are) given, there is no timezone info
// in the format, (e.g. 2019-11-09 1001), so they'll be considered as UTC-based,
// and generate the missing `start`(`end`) based on the same timezone, UTC
func ParseDuration(start string, end string) (time.Time, time.Time, schema.Step) {
	now := time.Now().UTC()

	startTime := now
	endTime := now
	logger.Log.Debugln("Start time:", start, "end time:", end)
	if len(start) == 0 && len(end) == 0 { // both absent
		startTime = now.Add(-30 * time.Minute)
		endTime = now
	} else if len(end) == 0 { // start is present
		if err := tryParseTime(start, &startTime); err != nil {
			logger.Log.Fatalln("Unsupported time format:", start, err)
		}
	} else if len(start) == 0 { // end is present
		if err := tryParseTime(end, &endTime); err != nil {
			logger.Log.Fatalln("Unsupported time format:", end, err)
		}
	} else { // both are present
		if err := tryParseTime(start, &startTime); err != nil {
			logger.Log.Fatalln("Unsupported time format:", start, err)
		}
		if err := tryParseTime(end, &endTime); err != nil {
			logger.Log.Fatalln("Unsupported time format:", end, err)
		}
	}
	duration := endTime.Sub(startTime)
	step := schema.StepSecond
	if duration.Hours() >= 24*30 { // time range > 1 month
		step = schema.StepMonth
	} else if duration.Hours() > 24 { // time range > 1 day
		step = schema.StepDay
	} else if duration.Minutes() > 60 { // time range > 1 hour
		step = schema.StepHour
	} else if duration.Seconds() > 60 { // time range > 1 minute
		step = schema.StepMinute
	} else if duration.Seconds() <= 0 { // illegal
		logger.Log.Fatalln("end time must be later than start time, end time:", endTime, ", start time:", startTime)
	}
	return startTime, endTime, step
}
