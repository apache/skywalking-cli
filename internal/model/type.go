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

	event "skywalking.apache.org/repo/goapi/collect/event/v3"
)

// EventTypeEnumValue defines the values domain of --type option.
type EventTypeEnumValue struct {
	Enum     []event.Type
	Default  event.Type
	Selected event.Type
}

// Set the --type value, from raw string to EventTypeEnumValue.
func (s *EventTypeEnumValue) Set(value string) error {
	for _, enum := range s.Enum {
		if strings.EqualFold(enum.String(), value) {
			s.Selected = enum
			return nil
		}
	}
	types := make([]string, len(event.Type_name))
	for index := range event.Type_name {
		types[index] = event.Type_name[index]
	}
	return fmt.Errorf("allowed types are %s", strings.Join(types, ", "))
}

// String representation of the event type.
func (s EventTypeEnumValue) String() string {
	return s.Selected.String()
}
