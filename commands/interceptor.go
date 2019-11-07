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

package commands

import (
	"github.com/apache/skywalking-cli/graphql/schema"
	"github.com/apache/skywalking-cli/logger"
	"github.com/urfave/cli"
	"time"
)

// Convenient function to chain up multiple cli.BeforeFunc
func BeforeChain(beforeFunctions []cli.BeforeFunc) cli.BeforeFunc {
	return func(ctx *cli.Context) error {
		for _, beforeFunc := range beforeFunctions {
			if err := beforeFunc(ctx); err != nil {
				return err
			}
		}
		return nil
	}
}

var StepFormats = map[schema.Step]string{
	schema.StepMonth:  "2006-01-02",
	schema.StepDay:    "2006-01-02",
	schema.StepHour:   "2006-01-02 15",
	schema.StepMinute: "2006-01-02 1504",
	schema.StepSecond: "2006-01-02 1504",
}

// Set the duration if not set, and format it according to
// the given step
func SetUpDuration(ctx *cli.Context) error {
	step := ctx.Generic("step").(*StepEnumValue).Selected
	end := ctx.String("end")
	if len(end) == 0 {
		end = time.Now().Format(StepFormats[step])
		logger.Log.Debugln("Missing --end, defaults to", end)
		if err := ctx.Set("end", end); err != nil {
			return err
		}
	}

	start := ctx.String("start")
	if len(start) == 0 {
		start = time.Now().Add(-15 * time.Minute).Format(StepFormats[step])
		logger.Log.Debugln("Missing --start, defaults to", start)
		if err := ctx.Set("start", start); err != nil {
			return err
		}
	}

	return nil
}
