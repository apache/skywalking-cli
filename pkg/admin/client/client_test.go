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

import "testing"

func TestDeriveAdminURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		adminURL string
		want     string
	}{
		{
			name:    "derive from default graphql base-url",
			baseURL: "http://127.0.0.1:12800/graphql",
			want:    "http://127.0.0.1:17128",
		},
		{
			name:    "derive keeps the base host",
			baseURL: "http://1.2.3.4:12800/graphql",
			want:    "http://1.2.3.4:17128",
		},
		{
			name:    "derive preserves https scheme",
			baseURL: "https://oap.example.com:12800/graphql",
			want:    "https://oap.example.com:17128",
		},
		{
			name:     "explicit admin-url wins and trailing slash trimmed",
			baseURL:  "http://1.2.3.4:12800/graphql",
			adminURL: "http://admin.example.com:17128/",
			want:     "http://admin.example.com:17128",
		},
		{
			name:     "explicit admin-url with whitespace",
			baseURL:  "http://1.2.3.4:12800/graphql",
			adminURL: "  http://admin:17128  ",
			want:     "http://admin:17128",
		},
		{
			name:    "empty base-url falls back to default admin url",
			baseURL: "",
			want:    DefaultAdminURL,
		},
		{
			name:    "host without scheme defaults to http",
			baseURL: "//demo.skywalking.apache.org:12800/graphql",
			want:    "http://demo.skywalking.apache.org:17128",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeriveAdminURL(tt.baseURL, tt.adminURL); got != tt.want {
				t.Errorf("DeriveAdminURL(%q, %q) = %q, want %q", tt.baseURL, tt.adminURL, got, tt.want)
			}
		})
	}
}
