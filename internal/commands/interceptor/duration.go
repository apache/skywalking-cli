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
	"fmt"
	"strconv"
	"time"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/apache/skywalking-cli/pkg/graphql/utils"

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/logger"
	"github.com/apache/skywalking-cli/internal/model"
)

func TryParseTime(unparsed string, userStep api.Step) (api.Step, time.Time, error) {
	var possibleError error
	for step, layout := range utils.StepFormats {
		t, err := time.Parse(layout, unparsed)
		if err == nil {
			return step, t, nil
		}
		possibleError = err
	}
	duration, err := time.ParseDuration(unparsed)
	if err == nil {
		return userStep, time.Now().Add(duration), nil
	}
	return userStep, time.Time{}, fmt.Errorf("the given time %v is neither absolute time nor relative time: %+v %+v", unparsed, possibleError, err)
}

// DurationInterceptor sets the duration if absent, and formats it accordingly,
// see ParseDuration
func DurationInterceptor(ctx *cli.Context) error {
	start := ctx.String("start")
	end := ctx.String("end")
	userStep := ctx.Generic("step")
	if timezone := ctx.String("timezone"); timezone != "" {
		if offset, err := strconv.Atoi(timezone); err == nil {
			// `offset` is in form of "+1300", while `time.FixedZone` takes offset in seconds
			time.Local = time.FixedZone("", offset/100*60*60)
		}
	}

	var s api.Step
	if userStep != nil {
		s = userStep.(*model.StepEnumValue).Selected
	}

	startTime, endTime, step, dt := ParseDuration(start, end, s)

	if err := ctx.Set("start", startTime.Format(utils.StepFormats[step])); err != nil {
		return err
	} else if err := ctx.Set("end", endTime.Format(utils.StepFormats[step])); err != nil {
		return err
	} else if err := ctx.Set("step", step.String()); err != nil {
		return err
	} else if err := ctx.Set("duration-type", dt.String()); err != nil {
		return err
	}
	return nil
}

// ParseDuration parses the `start` and `end` to a triplet, (startTime, endTime, step),
// based on the given `timezone`, however, if the given `timezone` is empty, UTC becomes the default timezone.
// if --start and --end are both absent,
// then: start := now - 30min; end := now
// if --start is given, --end is absent,
// then: end := now + 30 units, where unit is the precision of `start`, (hours, minutes, etc.)
// if --start is absent, --end is given,
// then: start := end - 30 units, where unit is the precision of `end`, (hours, minutes, etc.)
func ParseDuration(start, end string, userStep api.Step) (startTime, endTime time.Time, step api.Step, dt utils.DurationType) {
	logger.Log.Debugln("Start time:", start, "end time:", end, "timezone:", time.Local)

	now := time.Now()

	// both are absent
	if start == "" && end == "" {
		return now.Add(-30 * time.Minute), now, api.StepMinute, utils.BothAbsent
	}

	var err error

	// both are present
	if len(start) > 0 && len(end) > 0 {
		if userStep, startTime, err = TryParseTime(start, userStep); err != nil {
			logger.Log.Fatalln("Unsupported time format:", start, err)
		}
		if step, endTime, err = TryParseTime(end, userStep); err != nil {
			logger.Log.Fatalln("Unsupported time format:", end, err)
		}

		return startTime, endTime, step, utils.BothPresent
	} else if end == "" { // end is absent
		if step, startTime, err = TryParseTime(start, userStep); err != nil {
			logger.Log.Fatalln("Unsupported time format:", start, err)
		}
		return startTime, startTime.Add(30 * utils.StepDuration[step]), step, utils.EndAbsent
	} else { // start is absent
		if step, endTime, err = TryParseTime(end, userStep); err != nil {
			logger.Log.Fatalln("Unsupported time format:", end, err)
		}
		return endTime.Add(-30 * utils.StepDuration[step]), endTime, step, utils.StartAbsent
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
