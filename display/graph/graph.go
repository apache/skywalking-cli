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
	"strings"

	"github.com/urfave/cli"

	d "github.com/apache/skywalking-cli/display/displayable"
	db "github.com/apache/skywalking-cli/display/graph/dashboard"
	"github.com/apache/skywalking-cli/display/graph/gauge"
	"github.com/apache/skywalking-cli/display/graph/heatmap"
	"github.com/apache/skywalking-cli/display/graph/linear"
	"github.com/apache/skywalking-cli/display/graph/tree"
	"github.com/apache/skywalking-cli/graphql/dashboard"
	"github.com/apache/skywalking-cli/graphql/schema"
)

type (
	Thermodynamic      = schema.HeatMap
	LinearMetrics      = map[string]float64
	MultiLinearMetrics = []LinearMetrics
	Trace              = schema.Trace
	TraceBrief         = schema.TraceBrief
	GlobalMetrics      = [][]*schema.SelectedRecord
	GlobalData         = dashboard.GlobalData
)

var (
	ThermodynamicType      = reflect.TypeOf(Thermodynamic{})
	LinearMetricsType      = reflect.TypeOf(LinearMetrics{})
	MultiLinearMetricsType = reflect.TypeOf(MultiLinearMetrics{})
	TraceType              = reflect.TypeOf(Trace{})
	TraceBriefType         = reflect.TypeOf(TraceBrief{})
	GlobalMetricsType      = reflect.TypeOf(GlobalMetrics{})
	GlobalDataType         = reflect.TypeOf(&GlobalData{})
)

const multipleLinearTitles = "P50, P75, P90, P95, P99"

func Display(ctx *cli.Context, displayable *d.Displayable) error {
	data := displayable.Data

	switch reflect.TypeOf(data) {
	case ThermodynamicType:
		return heatmap.Display(displayable)

	case LinearMetricsType:
		return linear.Display(ctx, []LinearMetrics{data.(LinearMetrics)}, nil)

	case MultiLinearMetricsType:
		inputs := data.(MultiLinearMetrics)
		titles := strings.Split(multipleLinearTitles, ", ")[:len(inputs)]
		return linear.Display(ctx, inputs, titles)

	case TraceType:
		return tree.Display(tree.Adapt(data.(Trace)))

	case TraceBriefType:
		return tree.DisplayList(ctx, displayable)

	case GlobalMetricsType:
		return gauge.Display(ctx, data.(GlobalMetrics))

	case GlobalDataType:
		return db.Display(ctx, data.(*GlobalData))

	default:
		return fmt.Errorf("type of %T is not supported to be displayed as ascii graph", reflect.TypeOf(data))
	}
}
