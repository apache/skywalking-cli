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

package asyncprofiler

import (
	"fmt"
	api "skywalking.apache.org/repo/goapi/query"
	"strings"
)

type AsyncProfilerEventTypeEnumValue struct {
	Enum     []api.AsyncProfilerEventType
	Default  []api.AsyncProfilerEventType
	Selected []api.AsyncProfilerEventType
}

func (e *AsyncProfilerEventTypeEnumValue) Set(value string) error {
	values := strings.Split(value, ",")
	types := make([]api.AsyncProfilerEventType, 0)
	for _, v := range values {
		for _, enum := range e.Enum {
			if strings.EqualFold(enum.String(), v) {
				types = append(types, enum)
				break
			}
		}
	}

	if len(types) != 0 {
		e.Selected = types
		return nil
	}

	orders := make([]string, len(api.AllAsyncProfilerEventType))
	for i, order := range api.AllAsyncProfilerEventType {
		orders[i] = order.String()
	}
	return fmt.Errorf("allowed analysis aggregate type are %s", strings.Join(orders, ", "))
}

func (e *AsyncProfilerEventTypeEnumValue) String() string {
	selected := make([]string, len(e.Selected))
	for i, item := range e.Selected {
		selected[i] = item.String()
	}
	return strings.Join(selected, ",")
}
