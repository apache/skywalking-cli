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

package gauge

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/apache/skywalking-cli/graphql/dashboard"
	"github.com/apache/skywalking-cli/graphql/schema"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/gauge"
	"github.com/mum4k/termdash/widgets/text"
)

const RootID = "root"

type metricColumn struct {
	title  *text.Text
	gauges []*gauge.Gauge
}

func newMetricColumn(title string, column []*schema.SelectedRecord, isDec bool) (*metricColumn, error) {
	var ret metricColumn
	var maxValue int

	t, err := text.New()
	if err != nil {
		return nil, err
	}
	if err := t.Write(title, text.WriteCellOpts(cell.FgColor(cell.ColorRed))); err != nil {
		return nil, err
	}
	ret.title = t

	if isDec {
		temp, err := strconv.Atoi(*(column[0].Value))
		if err != nil {
			return nil, err
		}
		maxValue = temp
	} else {
		temp, err := strconv.Atoi(*(column[len(column)-1].Value))
		if err != nil {
			return nil, err
		}
		maxValue = temp
	}

	for _, item := range column {
		strValue := *(item.Value)
		v, err := strconv.Atoi(strValue)
		if err != nil {
			return nil, err
		}

		if !isDec {
			strValue = fmt.Sprintf("%.4f", float64(v)/10000)
		}

		g, err := gauge.New(
			gauge.Height(1),
			gauge.Border(linestyle.Light),
			gauge.Color(cell.ColorMagenta),
			gauge.BorderTitle("["+strValue+"]"),
			gauge.HideTextProgress(),
			gauge.TextLabel(item.Name),
		)
		if err != nil {
			return nil, err
		}

		if err := g.Absolute(v, maxValue); err != nil {
			return nil, err
		}
		ret.gauges = append(ret.gauges, g)
	}

	return &ret, nil
}

func layout(columns ...*metricColumn) ([]container.Option, error) {
	const GaugeNum = 6
	const TitleHeight = 10

	gaugeHeight := math.Floor((100 - TitleHeight) / GaugeNum)
	var metricColumns []grid.Element

	for _, c := range columns {
		var column []grid.Element
		column = append(column, grid.RowHeightPerc(TitleHeight, grid.Widget(c.title)))

		for i := 0; i < GaugeNum; i++ {
			column = append(column, grid.RowHeightPerc(int(gaugeHeight), grid.Widget(c.gauges[i])))
		}
		metricColumns = append(metricColumns, grid.ColWidthPerc(25, column...))
	}

	builder := grid.New()
	builder.Add(
		grid.RowHeightPerc(10),
		grid.RowHeightPerc(70, metricColumns...),
	)

	gridOpts, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return gridOpts, nil
}

func Display(metrics *dashboard.GlobalMetrics) error {
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

	var columns []*metricColumn

	col, err := newMetricColumn(" Service Load (calls/min) ", metrics.ServiceLoad, true)
	if err != nil {
		return err
	}
	columns = append(columns, col)

	col, err = newMetricColumn("    Slow Services (ms)    ", metrics.SlowServices, true)
	if err != nil {
		return err
	}
	columns = append(columns, col)

	col, err = newMetricColumn("Un-Health Services (Apdex)", metrics.UnhealthyServices, false)
	if err != nil {
		return err
	}
	columns = append(columns, col)

	col, err = newMetricColumn("    Slow Endpoints (ms)   ", metrics.SlowEndpoints, true)
	if err != nil {
		return err
	}
	columns = append(columns, col)

	gridOpts, err := layout(columns...)
	if err != nil {
		return err
	}

	err = c.Update(RootID, append(
		gridOpts,
		container.Border(linestyle.Light),
		container.BorderTitle("[Global Metrics]-PRESS Q TO QUIT"))...,
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
