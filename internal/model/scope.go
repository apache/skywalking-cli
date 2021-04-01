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

package model

import (
	"fmt"
	"strings"

	api "skywalking.apache.org/repo/goapi/query"
)

// ScopeEnumValue defines the values domain of --scope option
type ScopeEnumValue struct {
	Enum     []api.Scope
	Default  api.Scope
	Selected api.Scope
}

// Set the --scope value, from raw string to ScopeEnumValue
func (s *ScopeEnumValue) Set(value string) error {
	for _, enum := range s.Enum {
		if strings.EqualFold(enum.String(), value) {
			s.Selected = enum
			return nil
		}
	}
	scopes := make([]string, len(api.AllScope))
	for i, scope := range api.AllScope {
		scopes[i] = scope.String()
	}
	return fmt.Errorf("allowed scopes are %s", strings.Join(scopes, ", "))
}

// String representation of the scope
func (s ScopeEnumValue) String() string {
	return s.Selected.String()
}
