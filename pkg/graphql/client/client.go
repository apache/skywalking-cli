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
	"crypto/tls"
	"encoding/base64"
	"net/http"

	"github.com/machinebox/graphql"

	"github.com/apache/skywalking-cli/pkg/contextkey"
	"github.com/apache/skywalking-cli/pkg/logger"
)

// newClient creates a new GraphQL client with configuration from context.
//
// CLI has its own defaults, but for calls from other servers,
// we also add validation and certain default values here.
func newClient(ctx context.Context) *graphql.Client {
	options := []graphql.ClientOption{}

	var insecure bool
	if val := ctx.Value(contextkey.Insecure{}); val != nil {
		if flag, ok := val.(bool); ok {
			insecure = flag
		}
	}
	if insecure {
		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: insecure} // #nosec G402
		httpClient := &http.Client{Transport: customTransport}
		options = append(options, graphql.WithHTTPClient(httpClient))
	}

	var baseURL string
	if val := ctx.Value(contextkey.BaseURL{}); val != nil {
		if str, ok := val.(string); ok {
			baseURL = str
		}
	}
	if baseURL == "" {
		baseURL = "http://127.0.0.1:12800/graphql"
	}

	client := graphql.NewClient(baseURL, options...)
	client.Log = func(msg string) {
		logger.Log.Debugln(msg)
	}
	return client
}

// ExecuteQuery executes the `request` and parse to the `response`, returning `error` if there is any.
func ExecuteQuery(ctx context.Context, request *graphql.Request, response any) error {
	var username, password, authorization string
	if val := ctx.Value(contextkey.Username{}); val != nil {
		if str, ok := val.(string); ok {
			username = str
		}
	}
	if val := ctx.Value(contextkey.Password{}); val != nil {
		if str, ok := val.(string); ok {
			password = str
		}
	}
	if val := ctx.Value(contextkey.Authorization{}); val != nil {
		if str, ok := val.(string); ok {
			authorization = str
		}
	}

	if authorization == "" && username != "" && password != "" {
		authorization = "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
	}
	if authorization != "" {
		request.Header.Set("Authorization", authorization)
	}

	client := newClient(ctx)
	err := client.Run(ctx, request, response)
	return err
}
