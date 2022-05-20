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

package ebpf

import (
	"fmt"
	"strings"

	api "skywalking.apache.org/repo/goapi/query"
)

type ProfilingAnalyzeAggregateTypeEnumValue struct {
	Enum     []api.EBPFProfilingAnalyzeAggregateType
	Default  api.EBPFProfilingAnalyzeAggregateType
	Selected api.EBPFProfilingAnalyzeAggregateType
}

func (e *ProfilingAnalyzeAggregateTypeEnumValue) Set(value string) error {
	for _, enum := range e.Enum {
		if strings.EqualFold(enum.String(), value) {
			e.Selected = enum
			return nil
		}
	}
	orders := make([]string, len(api.AllEBPFProfilingAnalyzeAggregateType))
	for i, order := range api.AllEBPFProfilingAnalyzeAggregateType {
		orders[i] = order.String()
	}
	return fmt.Errorf("allowed analysis aggregate type are %s", strings.Join(orders, ", "))
}

func (e *ProfilingAnalyzeAggregateTypeEnumValue) String() string {
	return e.Selected.String()
}
