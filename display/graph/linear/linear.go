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

	"github.com/mum4k/termdash/linestyle"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/linechart"
)

const RootID = "root"

func NewLineChart(inputs map[string]float64) (lineChart *linechart.LineChart, err error) {
	index := 0

	xLabels := map[int]string{}
	yValues := make([]float64, len(inputs))

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

	if lineChart, err = linechart.New(linechart.YAxisAdaptive()); err != nil {
		return
	}

	err = lineChart.Series("graph-linear", yValues, linechart.SeriesXLabels(xLabels))

	return lineChart, err
}

// LineChartElements is the part that separated from layout,
// which can be reused by global dashboard.
func LineChartElements(lineCharts []*linechart.LineChart, titles []string) [][]grid.Element {
	cols := maxSqrt(len(lineCharts))

	rows := make([][]grid.Element, int(math.Ceil(float64(len(lineCharts))/float64(cols))))

	for r := 0; r < len(rows); r++ {
		var row []grid.Element
		for c := 0; c < cols && r*cols+c < len(lineCharts); c++ {
			percentage := int(math.Floor(float64(100) / float64(cols)))
			if r == len(rows)-1 {
				percentage = int(math.Floor(float64(100) / float64(len(lineCharts)-r*cols)))
			}

			var title string
			if titles == nil {
				title = fmt.Sprintf("#%v", r*cols+c)
			} else {
				title = titles[r*cols+c]
			}

			row = append(row, grid.ColWidthPerc(
				int(math.Min(99, float64(percentage))),
				grid.Widget(
					lineCharts[r*cols+c],
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

func Display(inputs []map[string]float64) error {
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

	var elements []*linechart.LineChart

	for _, input := range inputs {
		w, e := NewLineChart(input)
		if e != nil {
			return e
		}
		elements = append(elements, w)
	}

	gridOpts, err := layout(LineChartElements(elements, nil))
	if err != nil {
		return err
	}

	err = c.Update(RootID, append(
		gridOpts,
		container.Border(linestyle.Light),
		container.BorderTitle("PRESS Q TO QUIT"))...,
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
