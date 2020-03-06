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

package client

import (
	"context"

	"github.com/machinebox/graphql"
	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/logger"
)

func newClient(cliCtx *cli.Context) (client *graphql.Client) {
	client = graphql.NewClient(cliCtx.GlobalString("base-url"))
	client.Log = func(msg string) {
		logger.Log.Debugln(msg)
	}
	return
}

// ExecuteQuery executes the `request` and parse to the `response`, returning `error` if there is any.
func ExecuteQuery(cliCtx *cli.Context, request *graphql.Request, response interface{}) error {
	client := newClient(cliCtx)
	ctx := context.Background()
	err := client.Run(ctx, request, response)
	return err
}

// ExecuteQuery executes the `request` and parse to the `response`, panic if there is any `error`.
func ExecuteQueryOrFail(cliCtx *cli.Context, request *graphql.Request, response interface{}) {
	if err := ExecuteQuery(cliCtx, request, response); err != nil {
		logger.Log.Fatalln(err)
	}
}
