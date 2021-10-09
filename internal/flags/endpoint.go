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

// EndpointFlags take either endpoint id or endpoint name as input,
// and transform to the other one.
var EndpointFlags = append(
	ServiceFlags, // endpoint level requires service level by default

	&cli.StringFlag{
		Name:     "endpoint-id",
		Usage:    "`endpoint id`, if you don't have endpoint id, use `--endpoint-name` instead",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "endpoint-name",
		Usage:    "`endpoint name`, if you already have endpoint id, prefer to use `--endpoint-id`",
		Required: false,
	},
)

// EndpointRelationFlags take either destination endpoint id or destination endpoint name as input,
// and transform to the other one.
var EndpointRelationFlags = append(
	append(EndpointFlags[len(ServiceFlags):], ServiceRelationFlags...),

	&cli.StringFlag{
		Name:     "dest-endpoint-id",
		Usage:    "`destination` service endpoint id, if you don't have service endpoint id, use `--dest-endpoint-name` instead",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "dest-endpoint-name",
		Usage:    "`destination` service endpoint name, if you already have endpoint id, prefer to use `--dest-endpoint-id`",
		Required: false,
	},
)
