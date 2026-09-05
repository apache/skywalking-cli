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

// Package view fetches the asz.view document of an AI agent conversation from the OAP's
// GET /ai-agent/conversations/{conversation}/v1/view route. The route lives on the query
// host beside /graphql, not on the admin host, because the document is what the UI reads;
// it is streamed, since a long conversation renders to tens of megabytes, so the body is
// copied through and never held whole.
package view

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"

	"github.com/apache/skywalking-cli/pkg/contextkey"
	"github.com/apache/skywalking-cli/pkg/transport"
)

const (
	// MediaTypeJSON names the document, its version a parameter: the route's default body.
	MediaTypeJSON = "application/vnd.skywalking.asz.view+json"
	// MediaTypeYAML is the same document as YAML, chosen by Accept.
	MediaTypeYAML = "application/vnd.skywalking.asz.view+yaml"
	// problemType is the route's error body, RFC 9457.
	problemType = "application/problem+json"

	defaultBaseURL = "http://127.0.0.1:12800/graphql"
)

// Problem is the route's error: an RFC 9457 problem document carrying the status.
type Problem struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
	URL    string `json:"-"`
}

func (p *Problem) Error() string {
	if p.Detail != "" {
		return fmt.Sprintf("%d %s: %s (%s)", p.Status, p.Title, p.Detail, p.URL)
	}
	return fmt.Sprintf("%d %s (%s)", p.Status, p.Title, p.URL)
}

// CoreURL is the root of the query host the route lives on, derived from the GraphQL
// base URL by dropping its path: http://host:12800/graphql becomes http://host:12800.
// A base URL that does not parse is returned trimmed, so the error surfaces on the call.
func CoreURL(baseURL string) string {
	trimmed := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	u, err := url.Parse(trimmed)
	if err != nil || u.Host == "" {
		return trimmed
	}
	return u.Scheme + "://" + u.Host
}

// Path is the route of one conversation.
func Path(conversation string) string {
	return "/ai-agent/conversations/" + url.PathEscape(conversation) + "/v1/view"
}

// Fetch streams the document of the conversation to out and returns the Content-Type
// it came with. serviceName is required; instanceName narrows the read to one sender.
// A non-2xx answer is returned as a *Problem when the OAP sent one.
func Fetch(ctx context.Context, conversation, serviceName, instanceName string, yaml bool, out io.Writer) (string, error) {
	query := url.Values{"service": {serviceName}}
	if instanceName != "" {
		query.Set("instance", instanceName)
	}
	full := CoreURL(transport.GetValue(ctx, contextkey.BaseURL{}, defaultBaseURL)) + Path(conversation) + "?" + query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, full, http.NoBody)
	if err != nil {
		return "", err
	}
	if yaml {
		req.Header.Set("Accept", MediaTypeYAML)
	} else {
		req.Header.Set("Accept", MediaTypeJSON)
	}
	if authorization := transport.AuthHeader(ctx); authorization != "" {
		req.Header.Set("Authorization", authorization)
	}

	resp, err := transport.HTTPClient(ctx).Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return contentType, readError(resp, full)
	}
	_, err = io.Copy(out, resp.Body)
	return contentType, err
}

// readError turns a non-2xx response into an error: the problem document when the OAP
// sent one, otherwise the status and whatever the body says.
func readError(resp *http.Response, full string) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	mediaType, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if mediaType == problemType {
		problem := &Problem{URL: full}
		if json.Unmarshal(body, problem) == nil && problem.Status != 0 {
			return problem
		}
	}
	return fmt.Errorf("%s: %s: %s", full, resp.Status, strings.TrimSpace(string(body)))
}
