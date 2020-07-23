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

	"github.com/apache/skywalking-cli/display/graph/gauge"

	"github.com/apache/skywalking-cli/display/graph/tree"

	"github.com/apache/skywalking-cli/display/graph/heatmap"
	"github.com/apache/skywalking-cli/graphql/schema"

	d "github.com/apache/skywalking-cli/display/displayable"
	"github.com/apache/skywalking-cli/display/graph/linear"
)

type (
	Thermodynamic      = schema.Thermodynamic
	LinearMetrics      = map[string]float64
	MultiLinearMetrics = []LinearMetrics
	Trace              = schema.Trace
	GlobalMetrics      = [][]*schema.SelectedRecord
)

var (
	ThermodynamicType      = reflect.TypeOf(Thermodynamic{})
	LinearMetricsType      = reflect.TypeOf(LinearMetrics{})
	MultiLinearMetricsType = reflect.TypeOf(MultiLinearMetrics{})
	TraceType              = reflect.TypeOf(Trace{})
	GlobalMetricsType      = reflect.TypeOf(GlobalMetrics{})
)

func Display(displayable *d.Displayable) error {
	data := displayable.Data

	switch reflect.TypeOf(data) {
	case ThermodynamicType:
		return heatmap.Display(displayable)

	case LinearMetricsType:
		return linear.Display([]LinearMetrics{data.(LinearMetrics)})

	case MultiLinearMetricsType:
		return linear.Display(data.(MultiLinearMetrics))

	case TraceType:
		return tree.Display(tree.Adapt(data.(Trace)))

	case GlobalMetricsType:
		return gauge.Display(data.(GlobalMetrics))

	default:
		return fmt.Errorf("type of %T is not supported to be displayed as ascii graph", reflect.TypeOf(data))
	}
}
