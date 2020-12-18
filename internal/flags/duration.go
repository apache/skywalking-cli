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
	"github.com/apache/skywalking-cli/api"
	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/internal/model"
)

// DurationFlags are healthcheck flags that involves a duration, composed
// by a start time, an end time, and a step, which is commonly used
// in most of the commands
var DurationFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "start",
		Usage: "query start `TIME`",
	},
	cli.StringFlag{
		Name:  "end",
		Usage: "query end `TIME`",
	},
	cli.GenericFlag{
		Name:   "step",
		Hidden: true,
		Value: &model.StepEnumValue{
			Enum:     api.AllStep,
			Default:  api.StepMinute,
			Selected: api.StepMinute,
		},
	},
	cli.StringFlag{
		Name:  "durationType",
		Usage: "the type of duration",
	},
}
