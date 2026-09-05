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

package view

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apache/skywalking-cli/pkg/contextkey"
)

func TestCoreURL(t *testing.T) {
	cases := map[string]string{
		"http://127.0.0.1:12800/graphql":   "http://127.0.0.1:12800",
		"https://oap.example.com/graphql/": "https://oap.example.com",
		"http://[::1]:12800/graphql":       "http://[::1]:12800",
		"http://oap:12800":                 "http://oap:12800",
		"not a url":                        "not a url",
	}
	for in, want := range cases {
		if got := CoreURL(in); got != want {
			t.Errorf("CoreURL(%q) = %q, want %q", in, got, want)
		}
	}
}

// server answers the route the way the OAP does: the document as JSON in several flushes or
// as YAML by Accept, a problem document for anything else, and 418 for a wrong request.
func server(t *testing.T, document []byte) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ai-agent/conversations/c 1/v1/view" ||
			r.URL.Query().Get("service") != "agent" || r.URL.Query().Get("instance") != "sender" {
			w.WriteHeader(http.StatusTeapot)
			return
		}
		if r.Header.Get("Authorization") != "Basic dTpw" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		switch r.Header.Get("Accept") {
		case MediaTypeYAML:
			w.Header().Set("Content-Type", MediaTypeYAML+"; version=1.0")
			_, _ = w.Write([]byte("format: asz.view\n"))
		case MediaTypeJSON:
			w.Header().Set("Content-Type", MediaTypeJSON+"; version=1.0")
			for i := 0; i < len(document); i += 8192 {
				end := i + 8192
				if end > len(document) {
					end = len(document)
				}
				_, _ = w.Write(document[i:end])
				w.(http.Flusher).Flush()
			}
		default:
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"type":"about:blank","title":"Not Found","status":404,"detail":"no round"}`))
		}
	}))
}

func testContext(serverURL string) context.Context {
	ctx := context.WithValue(context.Background(), contextkey.BaseURL{}, serverURL+"/graphql")
	ctx = context.WithValue(ctx, contextkey.Username{}, "u")
	return context.WithValue(ctx, contextkey.Password{}, "p")
}

func TestFetchStreamsTheDocument(t *testing.T) {
	document := bytes.Repeat([]byte("{\"format\":\"asz.view\"}\n"), 4096)
	srv := server(t, document)
	defer srv.Close()

	var out bytes.Buffer
	contentType, err := Fetch(testContext(srv.URL), "c 1", "agent", "sender", false, &out)
	if err != nil {
		t.Fatal(err)
	}
	if contentType != MediaTypeJSON+"; version=1.0" || !bytes.Equal(out.Bytes(), document) {
		t.Fatalf("content type %q, %d bytes", contentType, out.Len())
	}
}

func TestFetchAsksForYAML(t *testing.T) {
	srv := server(t, nil)
	defer srv.Close()

	var out bytes.Buffer
	contentType, err := Fetch(testContext(srv.URL), "c 1", "agent", "sender", true, &out)
	if err != nil || contentType != MediaTypeYAML+"; version=1.0" || out.String() != "format: asz.view\n" {
		t.Fatalf("%v, content type %q, body %q", err, contentType, out.String())
	}
}

func TestAProblemDocumentIsTheError(t *testing.T) {
	srv := server(t, nil)
	defer srv.Close()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, srv.URL+Path("c 1")+"?service=agent&instance=sender", http.NoBody)
	req.Header.Set("Authorization", "Basic dTpw")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var problem *Problem
	if err := readError(resp, req.URL.String()); !errors.As(err, &problem) ||
		problem.Status != 404 || problem.Detail != "no round" || problem.Title != "Not Found" {
		t.Fatalf("problem: %v", err)
	}
}
