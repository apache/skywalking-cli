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

package table

import (
	"testing"

	"github.com/apache/skywalking-cli/graphql/schema"
)

func TestTableDisplay(t *testing.T) {
	var result []schema.Service
	display(t, result)
	result = make([]schema.Service, 0)
	display(t, result)
	result = append(result, schema.Service{
		ID:   "1",
		Name: "table",
	})
	display(t, result)
}

func display(t *testing.T, result []schema.Service) {
	if err := Display(result); err != nil {
		t.Error(err)
	}
}
