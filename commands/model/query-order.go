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

	"github.com/apache/skywalking-cli/graphql/schema"
)

// QueryOrderEnumValue defines the values domain of --query-order option
type QueryOrderEnumValue struct {
	Enum      []schema.QueryOrder
	StartTime schema.QueryOrder
	Duration  schema.QueryOrder
}

// Set the --order value, from raw string to QueryOrderEnumValue
func (s *QueryOrderEnumValue) Set(value string) error {
	for _, enum := range s.Enum {
		if strings.EqualFold(enum.String(), value) {
			s.Duration = enum
			return nil
		}
	}
	orders := make([]string, len(schema.AllQueryOrder))
	for i, order := range schema.AllQueryOrder {
		orders[i] = order.String()
	}
	return fmt.Errorf("allowed query orders are %s", strings.Join(orders, ", "))
}

// String representation of the query order
func (s QueryOrderEnumValue) String() string {
	return s.Duration.String()
}
