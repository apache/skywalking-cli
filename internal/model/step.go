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

// StepEnumValue defines the values domain of --step option
type StepEnumValue struct {
	Enum     []api.Step
	Default  api.Step
	Selected api.Step
}

// Set the --step value, from raw string to StepEnumValue
func (s *StepEnumValue) Set(value string) error {
	for _, enum := range s.Enum {
		if strings.EqualFold(enum.String(), value) {
			s.Selected = enum
			return nil
		}
	}
	steps := make([]string, len(api.AllStep))
	for i, step := range api.AllStep {
		steps[i] = step.String()
	}
	return fmt.Errorf("allowed steps are %s", strings.Join(steps, ", "))
}

// String representation of the step
func (s StepEnumValue) String() string {
	return s.Selected.String()
}
