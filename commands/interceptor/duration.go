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
	schema.StepSecond: "2006-01-02 150400",
	schema.StepMinute: "2006-01-02 1504",
	schema.StepHour:   "2006-01-02 15",
	schema.StepDay:    "2006-01-02",
	schema.StepMonth:  "2006-01",
}

var stepDuration = map[schema.Step]time.Duration{
	schema.StepSecond: time.Second,
	schema.StepMinute: time.Minute,
	schema.StepHour:   time.Hour,
	schema.StepDay:    time.Hour * 24,
	schema.StepMonth:  time.Hour * 24 * 30,
}

func tryParseTime(unparsed string) (schema.Step, time.Time, error) {
	var possibleError error = nil
	for step, layout := range stepFormats {
		t, err := time.Parse(layout, unparsed)
		if err == nil {
			return step, t, nil
		}
		possibleError = err
	}
	return schema.StepSecond, time.Time{}, possibleError
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
// if --start is given, --end is absent, then: end := now + 30 units, where unit is the precision of `start`, (hours, minutes, etc.)
// if --start is absent, --end is given, then: start := end - 30 unis, where unit is the precision of `end`, (hours, minutes, etc.)
// NOTE that when either(both) `start` or `end` is(are) given, there is no timezone info
// in the format, (e.g. 2019-11-09 1001), so they'll be considered as UTC-based,
// and generate the missing `start`(`end`) based on the same timezone, UTC
func ParseDuration(start string, end string) (time.Time, time.Time, schema.Step) {
	logger.Log.Debugln("Start time:", start, "end time:", end)

	now := time.Now().UTC()

	// both are absent
	if len(start) == 0 && len(end) == 0 {
		return now.Add(-30 * time.Minute), now, schema.StepMinute
	}

	var startTime, endTime time.Time
	var step schema.Step
	var err error

	// both are present
	if len(start) > 0 && len(end) > 0 {
		start, end = AlignPrecision(start, end)

		if _, startTime, err = tryParseTime(start); err != nil {
			logger.Log.Fatalln("Unsupported time format:", start, err)
		}
		if step, endTime, err = tryParseTime(end); err != nil {
			logger.Log.Fatalln("Unsupported time format:", end, err)
		}

		return startTime, endTime, step
	}

	// end is absent
	if len(end) == 0 {
		if step, startTime, err = tryParseTime(start); err != nil {
			logger.Log.Fatalln("Unsupported time format:", start, err)
		}
		return startTime, startTime.Add(30 * stepDuration[step]), step
	}

	// start is present
	if len(start) == 0 {
		if step, endTime, err = tryParseTime(end); err != nil {
			logger.Log.Fatalln("Unsupported time format:", end, err)
		}
		return endTime.Add(-30 * stepDuration[step]), endTime, step
	}

	logger.Log.Fatalln("Should never happen")

	return startTime, endTime, step
}

// AlignPrecision aligns the two time strings to same precision
// by truncating the more precise one
func AlignPrecision(start string, end string) (string, string) {
	if len(start) < len(end) {
		return start, end[0:len(start)]
	}
	if len(start) > len(end) {
		return start[0:len(end)], end
	}
	return start, end
}
