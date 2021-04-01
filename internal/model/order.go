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

// OrderEnumValue defines the values domain of --order option
type OrderEnumValue struct {
	Enum     []api.Order
	Default  api.Order
	Selected api.Order
}

// Set the --order value, from raw string to OrderEnumValue
func (s *OrderEnumValue) Set(value string) error {
	for _, enum := range s.Enum {
		if strings.EqualFold(enum.String(), value) {
			s.Selected = enum
			return nil
		}
	}
	orders := make([]string, len(api.AllOrder))
	for i, order := range api.AllOrder {
		orders[i] = order.String()
	}
	return fmt.Errorf("allowed orders are %s", strings.Join(orders, ", "))
}

// String representation of the order
func (s OrderEnumValue) String() string {
	return s.Selected.String()
}
