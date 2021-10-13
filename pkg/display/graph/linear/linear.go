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

package linear

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/linechart"

	"github.com/urfave/cli/v2"
)

const RootID = "root"

const defaultSeriesLabel = "linear"

func NewLineChart(inputs map[string]float64) (lineChart *linechart.LineChart, err error) {
	if lineChart, err = linechart.New(linechart.YAxisAdaptive()); err != nil {
		return
	}
	if err = SetLineChartSeries(lineChart, inputs); err != nil {
		return
	}
	return lineChart, err
}

func SetLineChartSeries(lc *linechart.LineChart, inputs map[string]float64) error {
	xLabels, yValues := processInputs(inputs)
	return lc.Series(defaultSeriesLabel, yValues, linechart.SeriesXLabels(xLabels))
}

// processInputs converts inputs into xLabels and yValues for line charts.
func processInputs(inputs map[string]float64) (xLabels map[int]string, yValues []float64) {
	index := 0

	xLabels = map[int]string{}
	yValues = make([]float64, len(inputs))

	// The iteration order of map is uncertain, so the keys must be sorted explicitly.
	var names []string
	for name := range inputs {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		xLabels[index] = name
		yValues[index] = inputs[name]
		index++
	}
	return
}

// LineChartElements is the part that separated from layout,
// which can be reused by global dashboard.
func LineChartElements(lineCharts map[string]*linechart.LineChart) [][]grid.Element {
	cols := maxSqrt(len(lineCharts))

	rows := make([][]grid.Element, int(math.Ceil(float64(len(lineCharts))/float64(cols))))

	var charts []*linechart.LineChart
	var titles []string
	for t := range lineCharts {
		titles = append(titles, t)
	}
	sort.Strings(titles)
	for _, title := range titles {
		charts = append(charts, lineCharts[title])
	}

	for r := 0; r < len(rows); r++ {
		var row []grid.Element
		for c := 0; c < cols && r*cols+c < len(lineCharts); c++ {
			percentage := int(math.Floor(float64(100) / float64(cols)))
			if r == len(rows)-1 {
				percentage = int(math.Floor(float64(100) / float64(len(lineCharts)-r*cols)))
			}

			title := titles[r*cols+c]
			chart := charts[r*cols+c]

			row = append(row, grid.ColWidthPerc(
				int(math.Min(99, float64(percentage))),
				grid.Widget(
					chart,
					container.Border(linestyle.Light),
					container.BorderTitleAlignCenter(),
					container.BorderTitle(title),
				),
			))
		}
		rows[r] = row
	}

	return rows
}

func layout(rows [][]grid.Element) ([]container.Option, error) {
	builder := grid.New()

	for _, row := range rows {
		percentage := int(math.Min(99, float64(100/len(rows))))
		builder.Add(grid.RowHeightPerc(percentage, row...))
	}

	return builder.Build()
}

func Display(cliCtx *cli.Context, inputs map[string]map[string]float64) error {
	t, err := termbox.New()
	if err != nil {
		return err
	}
	defer t.Close()

	c, err := container.New(
		t,
		container.ID(RootID),
	)
	if err != nil {
		return err
	}

	elements := make(map[string]*linechart.LineChart)

	for title, input := range inputs {
		w, e := NewLineChart(input)
		if e != nil {
			return e
		}
		elements[title] = w
	}

	gridOpts, err := layout(LineChartElements(elements))
	if err != nil {
		return err
	}

	err = c.Update(RootID, append(
		gridOpts,
		container.Border(linestyle.Light),
		container.BorderTitle(fmt.Sprintf("[%s]-PRESS Q TO QUIT", cliCtx.String("name"))),
		container.BorderTitleAlignLeft())...,
	)

	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	quitter := func(keyboard *terminalapi.Keyboard) {
		if strings.EqualFold(keyboard.Key.String(), "q") {
			cancel()
		}
	}

	err = termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quitter))

	return err
}

func maxSqrt(num int) int {
	return int(math.Ceil(math.Sqrt(float64(num))))
}
