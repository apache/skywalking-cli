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
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/linechart"
	"strings"
)

const RootID = "root"

func newWidgets(inputs map[string]float64) (lineChart *linechart.LineChart, err error) {

	index := 0

	xLabels := map[int]string{}
	var yValues []float64
	for xLabel, yValue := range inputs {
		xLabels[index] = xLabel
		index++
		yValues = append(yValues, yValue)
	}

	if lineChart, err = linechart.New(
		linechart.YAxisAdaptive(),
	); err != nil {
		return
	}

	err = lineChart.Series("graph-linear", yValues, linechart.SeriesXLabels(xLabels))

	return lineChart, err
}

func gridLayout(lineChart *linechart.LineChart) ([]container.Option, error) {
	widget := grid.Widget(lineChart,
		container.Border(linestyle.Light),
		container.BorderTitleAlignCenter(),
		container.BorderTitle("Press q to quit"),
	)

	builder := grid.New()
	builder.Add(widget)

	return builder.Build()
}

func Display(inputs map[string]float64) error {
	t, err := termbox.New()
	if err != nil {
		return err
	}
	defer t.Close()

	c, err := container.New(
		t,
		container.ID(RootID),
		container.PaddingTop(2),
		container.PaddingRight(2),
		container.PaddingBottom(2),
		container.PaddingLeft(2),
	)
	if err != nil {
		return err
	}

	w, err := newWidgets(inputs)
	if err != nil {
		return err
	}

	gridOpts, err := gridLayout(w)
	if err != nil {
		return err
	}

	if err := c.Update(RootID, gridOpts...); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	quitter := func(keyboard *terminalapi.Keyboard) {
		if strings.ToLower(keyboard.Key.String()) == "q" {
			cancel()
		}
	}

	err = termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quitter))

	return err
}
