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

package heatmap

import (
	"fmt"
	"math"
	"time"

	"github.com/apache/skywalking-cli/graphql/utils"
	"github.com/apache/skywalking-cli/util"

	ui "github.com/gizak/termui/v3"

	d "github.com/apache/skywalking-cli/display/displayable"
	"github.com/apache/skywalking-cli/graphql/schema"
	"github.com/apache/skywalking-cli/lib"
)

func Display(displayable *d.Displayable) error {
	data := displayable.Data.(schema.Thermodynamic)

	nodes := data.Nodes
	duration := displayable.Duration

	rows, cols, min, max := statistics(nodes)

	if err := ui.Init(); err != nil {
		return err
	}
	defer ui.Close()

	termW, _ := ui.TerminalDimensions()

	hm := lib.NewHeatMap()
	hm.Title = fmt.Sprintf(" %s ", displayable.Title)
	hm.XLabels = make([]string, rows)
	hm.YLabels = make([]string, cols)
	for i := 0; i < rows; i++ {
		step := utils.StepDuration[duration.Step]
		format := utils.StepFormats[duration.Step]
		startTime, err := time.Parse(format, duration.Start)

		if err != nil {
			return err
		}

		hm.XLabels[i] = startTime.Add(time.Duration(i) * step).Format("15:04")
	}
	for i := 0; i < cols; i++ {
		hm.YLabels[i] = fmt.Sprintf("%4d", i*data.AxisYStep)
	}

	hm.Data = make([][]float64, rows)
	hm.CellColors = make([][]ui.Color, rows)
	hm.NumStyles = make([][]ui.Style, rows)
	for row := 0; row < rows; row++ {
		hm.Data[row] = make([]float64, cols)
		hm.CellColors[row] = make([]ui.Color, cols)
		hm.NumStyles[row] = make([]ui.Style, cols)
	}

	scale := max - min
	for _, node := range nodes {
		color := ui.Color(255 - (float64(*node[2])/scale)*23)
		hm.Data[*node[0]][*node[1]] = float64(*node[2])
		hm.CellColors[*node[0]][*node[1]] = color
		hm.NumStyles[*node[0]][*node[1]] = ui.Style{Fg: ui.ColorMagenta}
	}

	hm.Formatter = nil
	hm.XLabelStyles = []ui.Style{{Fg: ui.ColorWhite}}
	hm.CellGap = 0
	hm.CellWidth = int(float64(termW) / float64(rows))
	realWidth := (hm.CellWidth+hm.CellGap)*(rows+1) - hm.CellGap + 5
	hm.SetRect(int(float64(termW-realWidth)/2), 2, realWidth, cols+5)

	ui.Render(hm)

	events := ui.PollEvents()
	for e := <-events; e.ID != "q" && e.ID != "<C-c>"; e = <-events {
	}
	return nil
}

func statistics(nodes [][]*int) (rows, cols int, min, max float64) {
	min = math.MaxFloat64

	for _, node := range nodes {
		rows = util.MaxInt(rows, *node[0])
		cols = util.MaxInt(cols, *node[1])
		max = math.Max(max, float64(*node[2]))
		min = math.Min(min, float64(*node[2]))
	}

	rows++
	cols++
	return
}
