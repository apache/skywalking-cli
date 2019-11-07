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

package commands

import (
	"fmt"
	"github.com/apache/skywalking-cli/graphql/schema"
	"strings"
)

type StepEnumValue struct {
	Enum     []schema.Step
	Default  schema.Step
	Selected schema.Step
}

func (s *StepEnumValue) Set(value string) error {
	for _, enum := range s.Enum {
		if enum.String() == value {
			s.Selected = enum
			return nil
		}
	}
	steps := make([]string, len(schema.AllStep))
	for i, step := range schema.AllStep {
		steps[i] = step.String()
	}
	return fmt.Errorf("allowed steps are %s", strings.Join(steps, ", "))
}

func (s StepEnumValue) String() string {
	return s.Selected.String()
}
