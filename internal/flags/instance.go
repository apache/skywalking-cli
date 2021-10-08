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
)

// InstanceFlags take either service instance id or service instance name as input,
// and transform to the other one.
var InstanceFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "instance-id",
		Usage:    "`instance id`, if you don't have instance id, use `--instance-name` instead",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "instance-name",
		Usage:    "`instance name`, if you already have instance id, prefer to use `--instance-id`",
		Required: false,
	},
}

// InstanceRelationFlags take either destination instance id or destination instance name as input,
// and transform to the other one.
var InstanceRelationFlags = append(
	InstanceFlags,

	&cli.StringFlag{
		Name:     "dest-instance-id",
		Usage:    "`destination` instance id, if you don't have instance id, use `--dest-instance-name` instead",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "dest-instance-name",
		Usage:    "`destination` instance name, if you already have instance id, prefer to use `--dest-instance-id`",
		Required: false,
	},
)
