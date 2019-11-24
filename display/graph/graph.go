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

package graph

import (
	"fmt"
	"reflect"

	"github.com/apache/skywalking-cli/display/graph/linear"
)

func Display(object interface{}) error {
	if reflect.TypeOf(object) != reflect.TypeOf(map[string]float64{}) {
		return fmt.Errorf("type of %T is not supported to be displayed as ascii graph", reflect.TypeOf(object))
	}

	kvs := object.(map[string]float64)

	return linear.Display(kvs)
}
