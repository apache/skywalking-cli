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

package dashboard

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/apache/skywalking-cli/graphql/utils"
	"github.com/apache/skywalking-cli/lib/heatmap"

	"github.com/mattn/go-runewidth"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/display/graph/gauge"
	"github.com/apache/skywalking-cli/display/graph/linear"
	"github.com/apache/skywalking-cli/graphql/dashboard"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/widgets/button"
	"github.com/mum4k/termdash/widgets/linechart"
)

// rootID is the ID assigned to the root container.
const rootID = "root"

type layoutType int

const (
	// layoutMetrics displays all the widgets.
	layoutMetrics layoutType = iota

	// layoutLineChart focuses onto the line chart.
	layoutLineChart

	// layoutHeatMap focuses onto the heat map.
	layoutHeatMap
)

// strToLayoutType ensures the order of buttons is fixed.
var strToLayoutType = map[string]layoutType{
	"Metrics":         layoutMetrics,
	"ResponseLatency": layoutLineChart,
	"HeatMap":         layoutHeatMap,
}

// widgets holds the widgets used by the dashboard.
type widgets struct {
	gauges  []*gauge.MetricColumn
	linears []*linechart.LineChart
	heatmap *heatmap.HeatMap

	// buttons are used to change the layout.
	buttons []*button.Button
}

// linearTitles are titles of each line chart, load from the template file.
var linearTitles []string

// setLayout sets the specified layout.
func setLayout(c *container.Container, w *widgets, lt layoutType) error {
	gridOpts, err := gridLayout(w, lt)
	if err != nil {
		return err
	}
	return c.Update(rootID, gridOpts...)
}

// newLayoutButtons returns buttons that dynamically switch the layouts.
func newLayoutButtons(c *container.Container, w *widgets, template *dashboard.ButtonTemplate) ([]*button.Button, error) {
	var buttons []*button.Button

	opts := []button.Option{
		button.WidthFor(longestString(template.Texts)),
		button.FillColor(cell.ColorNumber(template.ColorNum)),
		button.Height(template.Height),
	}

	for _, text := range template.Texts {
		// declare a local variable lt to avoid closure.
		lt, ok := strToLayoutType[text]
		if !ok {
			return nil, fmt.Errorf("the %s is not supposed to be the button's text", text)
		}

		b, err := button.New(text, func() error {
			return setLayout(c, w, lt)
		}, opts...)
		if err != nil {
			return nil, err
		}
		buttons = append(buttons, b)
	}

	return buttons, nil
}

// gridLayout prepares container options that represent the desired screen layout.
func gridLayout(w *widgets, lt layoutType) ([]container.Option, error) {
	const buttonRowHeight = 15

	buttonColWidthPerc := 100 / len(w.buttons)
	var buttonCols []grid.Element

	for _, b := range w.buttons {
		buttonCols = append(buttonCols, grid.ColWidthPerc(buttonColWidthPerc, grid.Widget(b)))
	}

	rows := []grid.Element{
		grid.RowHeightPerc(buttonRowHeight, buttonCols...),
	}

	switch lt {
	case layoutMetrics:
		rows = append(rows,
			grid.RowHeightPerc(70, gauge.MetricColumnsElement(w.gauges)...),
		)

	case layoutLineChart:
		lcElements := linear.LineChartElements(w.linears, linearTitles)
		percentage := int(math.Min(99, float64((100-buttonRowHeight)/len(lcElements))))

		for _, e := range lcElements {
			rows = append(rows,
				grid.RowHeightPerc(percentage, e...),
			)
		}

	case layoutHeatMap:
		const heatmapColWidth = 85

		rows = append(rows,
			grid.RowHeightPerc(
				99-buttonRowHeight,
				grid.ColWidthPerc((99-heatmapColWidth)/2), // Use two empty cols to center the heatmap.
				grid.ColWidthPerc(heatmapColWidth, grid.Widget(w.heatmap)),
				grid.ColWidthPerc((99-heatmapColWidth)/2),
			),
		)
	}

	builder := grid.New()
	builder.Add(
		grid.RowHeightPerc(99, rows...),
	)
	gridOpts, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return gridOpts, nil
}

// newWidgets creates all widgets used by the dashboard.
func newWidgets(data *dashboard.GlobalData, template *dashboard.GlobalTemplate) (*widgets, error) {
	var columns []*gauge.MetricColumn
	var linears []*linechart.LineChart

	// Create gauges to display global metrics.
	for i, t := range template.Metrics {
		col, err := gauge.NewMetricColumn(data.Metrics[i], &t)
		if err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}

	// Create line charts to display global response latency.
	for _, input := range data.ResponseLatency {
		l, err := linear.NewLineChart(input)
		if err != nil {
			return nil, err
		}
		linears = append(linears, l)
	}

	// Create a heat map.
	hp, err := heatmap.NewHeatMap()
	if err != nil {
		return nil, err
	}
	hpColumns := utils.HeatMapToMap(&data.HeatMap)
	yLabels := utils.BucketsToStrings(data.HeatMap.Buckets)
	hp.SetColumns(hpColumns)
	hp.SetYLabels(yLabels)

	return &widgets{
		gauges:  columns,
		linears: linears,
		heatmap: hp,
	}, nil
}

func Display(ctx *cli.Context, data *dashboard.GlobalData) error {
	t, err := termbox.New(termbox.ColorMode(terminalapi.ColorMode256))
	if err != nil {
		return err
	}
	defer t.Close()

	c, err := container.New(
		t,
		container.Border(linestyle.Light),
		container.BorderTitle("[Global Dashboard]-PRESS Q TO QUIT"),
		container.ID(rootID))
	if err != nil {
		return err
	}

	template, err := dashboard.LoadTemplate(ctx.String("template"))
	if err != nil {
		return err
	}
	linearTitles = strings.Split(template.ResponseLatency.Labels, ", ")

	w, err := newWidgets(data, template)
	if err != nil {
		return err
	}
	lb, err := newLayoutButtons(c, w, &template.Buttons)
	if err != nil {
		return err
	}
	w.buttons = lb

	gridOpts, err := gridLayout(w, layoutMetrics)
	if err != nil {
		return err
	}

	if e := c.Update(rootID, gridOpts...); e != nil {
		return e
	}

	con, cancel := context.WithCancel(context.Background())
	quitter := func(keyboard *terminalapi.Keyboard) {
		if strings.EqualFold(keyboard.Key.String(), "q") {
			cancel()
		}
	}

	err = termdash.Run(con, t, c, termdash.KeyboardSubscriber(quitter))

	return err
}

// longestString returns the longest string in the string array.
func longestString(strs []string) (ret string) {
	maxLen := 0
	for _, s := range strs {
		if l := runewidth.StringWidth(s); l > maxLen {
			ret = s
			maxLen = l
		}
	}
	return
}
