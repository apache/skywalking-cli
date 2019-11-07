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

package client

import (
	"context"
	"github.com/apache/skywalking-cli/config"
	"github.com/apache/skywalking-cli/graphql/schema"
	"github.com/apache/skywalking-cli/logger"
	"github.com/machinebox/graphql"
)

func Services(duration schema.Duration) []schema.Service {
	ctx := context.Background()
	client := graphql.NewClient(config.Config.Global.BaseUrl)
	client.Log = func(msg string) {
		logger.Log.Debugln(msg)
	}

	var response map[string][]schema.Service
	request := graphql.NewRequest(`
		query ($duration: Duration!) {
			services: getAllServices(duration: $duration) {
				id name
			}
		}
	`)
	request.Var("duration", duration)
	if err := client.Run(ctx, request, &response); err != nil {
		logger.Log.Fatalln(err)
		panic(err)
	}

	return response["services"]
}
