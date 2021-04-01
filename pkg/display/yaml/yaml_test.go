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

package yaml

import (
	"testing"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/apache/skywalking-cli/pkg/display/displayable"
)

func TestYamlDisplay(t *testing.T) {
	var result []api.Service
	display(t, result)
	result = make([]api.Service, 0)
	display(t, result)
	result = append(result, api.Service{
		ID:   "1",
		Name: "yaml",
	})
	display(t, result)
}

func display(t *testing.T, result []api.Service) {
	if err := Display(&displayable.Displayable{Data: result}); err != nil {
		t.Error(err)
	}
}
