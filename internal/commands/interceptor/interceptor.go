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
	"context"

	"github.com/apache/skywalking-cli/pkg/contextkey"

	"github.com/urfave/cli/v2"
)

// BeforeChain is a convenient function to chain up multiple cli.BeforeFunc
func BeforeChain(beforeFunctions ...cli.BeforeFunc) cli.BeforeFunc {
	return func(cliCtx *cli.Context) error {
		ctx := cliCtx.Context
		ctx = context.WithValue(ctx, contextkey.BaseURL{}, cliCtx.String("base-url"))
		ctx = context.WithValue(ctx, contextkey.Insecure{}, cliCtx.Bool("insecure"))
		ctx = context.WithValue(ctx, contextkey.Username{}, cliCtx.String("username"))
		ctx = context.WithValue(ctx, contextkey.Password{}, cliCtx.String("password"))
		ctx = context.WithValue(ctx, contextkey.Authorization{}, cliCtx.String("authorization"))
		ctx = context.WithValue(ctx, contextkey.Display{}, cliCtx.String("display"))
		cliCtx.Context = ctx

		// --timezone is global option, it should be applied always.
		if err := TimezoneInterceptor(cliCtx); err != nil {
			return err
		}
		for _, beforeFunc := range beforeFunctions {
			if err := beforeFunc(cliCtx); err != nil {
				return err
			}
		}
		return nil
	}
}
