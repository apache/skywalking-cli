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
	"time"

	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/graphql/schema"
	"github.com/apache/skywalking-cli/logger"
)

func tryParseTime(unparsed string) (schema.Step, time.Time, error) {
	var possibleError error = nil
	for step, layout := range schema.StepFormats {
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

	if err := ctx.Set("start", startTime.Format(schema.StepFormats[step])); err != nil {
		return err
	} else if err := ctx.Set("end", endTime.Format(schema.StepFormats[step])); err != nil {
		return err
	} else if err := ctx.Set("step", step.String()); err != nil {
		return err
	}
	return nil
}

// ParseDuration parses the `start` and `end` to a triplet, (startTime, endTime, step)
// if --start and --end are both absent,
//   then: start := now - 30min; end := now
// if --start is given, --end is absent,
//   then: end := now + 30 units, where unit is the precision of `start`, (hours, minutes, etc.)
// if --start is absent, --end is given,
//   then: start := end - 30 unis, where unit is the precision of `end`, (hours, minutes, etc.)
// NOTE that when either(both) `start` or `end` is(are) given, there is no timezone info
// in the format, (e.g. 2019-11-09 1001), so they'll be considered as UTC-based,
// and generate the missing `start`(`end`) based on the same timezone, UTC
func ParseDuration(start, end string) (startTime, endTime time.Time, step schema.Step) {
	logger.Log.Debugln("Start time:", start, "end time:", end)

	now := time.Now()

	// both are absent
	if start == "" && end == "" {
		return now.Add(-30 * time.Minute), now, schema.StepMinute
	}

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
	} else if end == "" { // end is absent
		if step, startTime, err = tryParseTime(start); err != nil {
			logger.Log.Fatalln("Unsupported time format:", start, err)
		}
		return startTime, startTime.Add(30 * schema.StepDuration[step]), step
	} else { // start is absent
		if step, endTime, err = tryParseTime(end); err != nil {
			logger.Log.Fatalln("Unsupported time format:", end, err)
		}
		return endTime.Add(-30 * schema.StepDuration[step]), endTime, step
	}
}

// AlignPrecision aligns the two time strings to same precision
// by truncating the more precise one
func AlignPrecision(start, end string) (_, _ string) {
	if len(start) < len(end) {
		return start, end[0:len(start)]
	}
	if len(start) > len(end) {
		return start[0:len(end)], end
	}
	return start, end
}
