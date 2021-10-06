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

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/logger"
	"github.com/apache/skywalking-cli/pkg/graphql/metadata"
)

// TimezoneInterceptor sets the server timezone if the server supports the API,
// otherwise, sets to local timezone
func TimezoneInterceptor(ctx *cli.Context) error {
	// If there is timezone given by the user in command line, use it directly
	if ctx.IsSet("timezone") {
		return nil
	}

	serverTimeInfo, err := metadata.ServerTimeInfo(ctx)

	if err != nil {
		logger.Log.Debugf("Failed to get server time info: %v\n", err)
		return nil
	}

	if timezone := serverTimeInfo.Timezone; timezone != nil {
		if _, err := strconv.Atoi(*timezone); err == nil {
			for _, c := range ctx.Lineage() {
				if err := c.Set("timezone", *timezone); err == nil {
					return nil
				}
			}
			return fmt.Errorf("cannot set the timezone flag globally")
		}
	}

	return nil
}
