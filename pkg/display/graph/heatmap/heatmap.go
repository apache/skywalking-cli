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
	"context"
	"fmt"
	"strings"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgetapi"

	d "github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/utils"
	"github.com/apache/skywalking-cli/pkg/heatmap"
)

const rootID = "root"

func NewHeatMapWidget(data api.HeatMap) (hp *heatmap.HeatMap, err error) {
	hp, err = heatmap.NewHeatMap()
	if err != nil {
		return hp, err
	}

	SetData(hp, data)
	return
}

func SetData(hp *heatmap.HeatMap, data api.HeatMap) {
	hpColumns, yLabels := processData(data)
	hp.SetColumns(hpColumns)
	hp.SetYLabels(yLabels)
}

// processData converts data into hpColumns and yValues for the heat map.
func processData(data api.HeatMap) (hpColumns map[string][]int64, yLabels []string) {
	hpColumns = utils.HeatMapToMap(&data)
	yLabels = utils.BucketsToStrings(data.Buckets)
	return
}

// layout controls where and how the heat map widget is placed.
// Here uses the grid layout to center the widget horizontally and vertically.
func layout(hp widgetapi.Widget) ([]container.Option, error) {
	const hpColWidthPerc = 85
	const hpRowHeightPerc = 80

	builder := grid.New()
	builder.Add(
		grid.ColWidthPerc((99-hpColWidthPerc)/2), // Use two empty cols to center the heatmap.
		grid.ColWidthPerc(hpColWidthPerc,
			grid.RowHeightPerc((99-hpRowHeightPerc)/2), // Use two empty rows to center the heatmap.
			grid.RowHeightPerc(hpRowHeightPerc, grid.Widget(hp)),
			grid.RowHeightPerc((99-hpRowHeightPerc)/2),
		),
		grid.ColWidthPerc((99-hpColWidthPerc)/2),
	)
	return builder.Build()
}

func Display(displayable *d.Displayable) error {
	t, err := termbox.New(termbox.ColorMode(terminalapi.ColorMode256))
	if err != nil {
		return err
	}
	defer t.Close()

	title := fmt.Sprintf("[%s]-PRESS Q TO QUIT", displayable.Title)
	c, err := container.New(
		t,
		container.Border(linestyle.Light),
		container.BorderTitle(title),
		container.ID(rootID))
	if err != nil {
		return err
	}

	data := displayable.Data.(api.HeatMap)
	hp, err := NewHeatMapWidget(data)
	if err != nil {
		return err
	}

	gridOpts, err := layout(hp)
	if err != nil {
		return fmt.Errorf("builder.Build => %v", err)
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
