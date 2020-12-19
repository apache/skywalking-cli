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
	"github.com/urfave/cli"
)

// EntityFlags are attributes of Entity in the metrics v2 protocol.
var EntityFlags = []cli.Flag{
	cli.StringFlag{
		Name:     "instance",
		Usage:    "the name of the service instance",
		Value:    "",
		Required: false,
	},
	cli.StringFlag{
		Name:     "endpoint",
		Usage:    "the name of the endpoint",
		Value:    "",
		Required: false,
	},
	cli.StringFlag{
		Name:     "destService",
		Usage:    "the name of the destination endpoint",
		Value:    "",
		Required: false,
	},
	cli.BoolTFlag{
		Name:     "isDestNormal",
		Usage:    "set the destination service to normal or unnormal",
		Required: false,
	},
	cli.StringFlag{
		Name:     "destInstance",
		Usage:    "the name of the destination endpoint",
		Value:    "",
		Required: false,
	},
	cli.StringFlag{
		Name:     "destEndpoint",
		Usage:    "the name of the destination endpoint",
		Value:    "",
		Required: false,
	},
}
