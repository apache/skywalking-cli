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

	"github.com/guptarohit/asciigraph"
)

func Display(object interface{}) error {
	var unifiedValues []float64

	switch reflect.TypeOf(object) {
	case reflect.TypeOf([]float64{}):
	case reflect.TypeOf([]float32{}):
		for _, v := range object.([]float32) {
			unifiedValues = append(unifiedValues, float64(v))
		}
	case reflect.TypeOf([]int64{}):
		for _, v := range object.([]int64) {
			unifiedValues = append(unifiedValues, float64(v))
		}
	case reflect.TypeOf([]int32{}):
		for _, v := range object.([]int32) {
			unifiedValues = append(unifiedValues, float64(v))
		}
	case reflect.TypeOf([]int16{}):
		for _, v := range object.([]int16) {
			unifiedValues = append(unifiedValues, float64(v))
		}
	case reflect.TypeOf([]int8{}):
		for _, v := range object.([]int8) {
			unifiedValues = append(unifiedValues, float64(v))
		}
	case reflect.TypeOf([]int{}):
		for _, v := range object.([]int) {
			unifiedValues = append(unifiedValues, float64(v))
		}
	default:
		return fmt.Errorf("type of %T is not supported to display as graph", object)
	}

	graph := asciigraph.Plot(unifiedValues)
	fmt.Println(graph)

	return nil
}
