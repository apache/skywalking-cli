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
	"time"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/pkg/graphql/utils"
	lib "github.com/apache/skywalking-cli/pkg/heatmap"

	"github.com/mattn/go-runewidth"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/display/graph/gauge"
	"github.com/apache/skywalking-cli/pkg/display/graph/heatmap"
	"github.com/apache/skywalking-cli/pkg/display/graph/linear"
	"github.com/apache/skywalking-cli/pkg/graphql/dashboard"

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
	linears map[string]*linechart.LineChart
	heatmap *lib.HeatMap

	// buttons are used to change the layout.
	buttons []*button.Button
}

// template determines how the global dashboard is displayed.
var template *dashboard.GlobalTemplate

var allWidgets *widgets

var initStartStr string
var initStep = api.StepMinute
var initEndStr string

var curStartTime time.Time
var curEndTime time.Time

// setLayout sets the specified layout.
func setLayout(c *container.Container, lt layoutType) error {
	gridOpts, err := gridLayout(lt)
	if err != nil {
		return err
	}
	return c.Update(rootID, gridOpts...)
}

// newLayoutButtons returns buttons that dynamically switch the layouts.
func newLayoutButtons(c *container.Container) ([]*button.Button, error) {
	buttons := make([]*button.Button, len(strToLayoutType))

	ls := longestString(template.Buttons.Texts)
	if ls == "" {
		return nil, fmt.Errorf("failed to parse texts of buttons")
	}

	opts := []button.Option{
		button.WidthFor(ls),
		button.FillColor(cell.ColorNumber(template.Buttons.ColorNum)),
		button.Height(template.Buttons.Height),
	}

	for _, text := range template.Buttons.Texts {
		// declare a local variable lt to avoid closure.
		lt, ok := strToLayoutType[text]
		if !ok {
			return nil, fmt.Errorf("the '%s' is not supposed to be the button's text", text)
		}

		b, err := button.New(text, func() error {
			return setLayout(c, lt)
		}, opts...)
		if err != nil {
			return nil, err
		}

		buttons[lt] = b
	}

	return buttons, nil
}

// gridLayout prepares container options that represent the desired screen layout.
func gridLayout(lt layoutType) ([]container.Option, error) {
	const buttonRowHeight = 15

	buttonColWidthPerc := 99 / len(allWidgets.buttons)
	var buttonCols []grid.Element

	for _, b := range allWidgets.buttons {
		if b != nil {
			buttonCols = append(buttonCols, grid.ColWidthPerc(buttonColWidthPerc, grid.Widget(b)))
		}
	}

	rows := []grid.Element{
		grid.RowHeightPerc(buttonRowHeight, buttonCols...),
	}

	switch lt {
	case layoutMetrics:
		rows = append(rows,
			grid.RowHeightPerc(70, gauge.MetricColumnsElement(allWidgets.gauges)...),
		)

	case layoutLineChart:
		lcElements := linear.LineChartElements(allWidgets.linears)
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
				grid.ColWidthPerc(heatmapColWidth, grid.Widget(allWidgets.heatmap)),
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
func newWidgets(data *dashboard.GlobalData) error {
	var columns []*gauge.MetricColumn
	linears := make(map[string]*linechart.LineChart)

	// Create gauges to display global metrics.
	for i := range template.Metrics {
		col, err := gauge.NewMetricColumn(data.Metrics[i], &template.Metrics[i])
		if err != nil {
			return err
		}
		columns = append(columns, col)
	}

	// Create line charts to display global response latency.
	for label, input := range data.ResponseLatency {
		l, err := linear.NewLineChart(input)
		if err != nil {
			return err
		}
		linears[label] = l
	}

	// Create a heat map.
	hp, err := heatmap.NewHeatMapWidget(data.HeatMap)
	if err != nil {
		return err
	}

	allWidgets.gauges = columns
	allWidgets.linears = linears
	allWidgets.heatmap = hp
	return nil
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

	te, err := dashboard.LoadTemplate(ctx.String("template"))
	if err != nil {
		return err
	}
	template = te

	// Initialization
	allWidgets = &widgets{
		gauges:  nil,
		linears: nil,
		heatmap: nil,
		buttons: nil,
	}
	err = newWidgets(data)
	if err != nil {
		return err
	}
	lb, err := newLayoutButtons(c)
	if err != nil {
		return err
	}
	allWidgets.buttons = lb

	gridOpts, err := gridLayout(layoutMetrics)
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

	refreshInterval := time.Duration(ctx.Int("refresh")) * time.Second
	dt := utils.DurationType(ctx.String("duration-type"))

	// Only when users use the relative time, the duration will be adjusted to refresh.
	if dt != utils.BothPresent {
		go refresh(con, ctx, refreshInterval)
	}

	err = termdash.Run(con, t, c, termdash.KeyboardSubscriber(quitter), termdash.RedrawInterval(refreshInterval))

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

// refresh updates the duration and query the new data to update all of widgets, once every delay.
func refresh(con context.Context, ctx *cli.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	initStartStr = ctx.String("start")
	initEndStr = ctx.String("end")

	if s := ctx.Generic("step"); s != nil {
		initStep = s.(*model.StepEnumValue).Selected
	}

	_, start, err := interceptor.TryParseTime(initStartStr, initStep)
	if err != nil {
		return
	}
	_, end, err := interceptor.TryParseTime(initEndStr, initStep)
	if err != nil {
		return
	}

	curStartTime = start
	curEndTime = end

	for {
		select {
		case <-ticker.C:
			d, err := updateDuration(interval)
			if err != nil {
				continue
			}

			data, err := dashboard.Global(ctx, d)
			if err != nil {
				continue
			}

			if err := updateAllWidgets(data); err != nil {
				continue
			}
		case <-con.Done():
			return
		}
	}
}

// updateDuration will check if the duration changes after adding the interval.
// If the duration doesn't change, an error will be returned, and the dashboard will not refresh.
// Otherwise, a new duration will be returned, which is used to get the latest global data.
func updateDuration(interval time.Duration) (api.Duration, error) {
	step, _, err := interceptor.TryParseTime(initStartStr, initStep)
	if err != nil {
		return api.Duration{}, err
	}

	curStartTime = curStartTime.Add(interval)
	curEndTime = curEndTime.Add(interval)

	curStartStr := curStartTime.Format(utils.StepFormats[step])
	curEndStr := curEndTime.Format(utils.StepFormats[step])

	if curStartStr == initStartStr && curEndStr == initEndStr {
		return api.Duration{}, fmt.Errorf("the duration does not update")
	}

	initStartStr = curStartStr
	initEndStr = curEndStr
	return api.Duration{
		Start: curStartStr,
		End:   curEndStr,
		Step:  step,
	}, nil
}

// updateAllWidgets will update all of widgets' data to be displayed.
func updateAllWidgets(data *dashboard.GlobalData) error {
	// Update gauges
	for i, mcData := range data.Metrics {
		if err := allWidgets.gauges[i].Update(mcData); err != nil {
			return err
		}
	}

	// Update line charts.
	for i, inputs := range data.ResponseLatency {
		if err := linear.SetLineChartSeries(allWidgets.linears[i], inputs); err != nil {
			return err
		}
	}

	// Update the heat map.
	heatmap.SetData(allWidgets.heatmap, data.HeatMap)

	return nil
}
