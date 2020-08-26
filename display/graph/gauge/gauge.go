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

	"github.com/urfave/cli"

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

type MetricColumn struct {
	title  *text.Text
	gauges []*gauge.Gauge
}

func NewMetricColumn(column []*schema.SelectedRecord, config *dashboard.MetricTemplate) (*MetricColumn, error) {
	var ret MetricColumn
	var maxValue int

	t, err := text.New()
	if err != nil {
		return nil, err
	}
	if err := t.Write(config.Title, text.WriteCellOpts(cell.FgColor(cell.ColorRed))); err != nil {
		return nil, err
	}
	ret.title = t

	if len(column) == 0 {
		return nil, fmt.Errorf("the metrics data is empty, please check the GraphQL backend")
	}

	if config.Condition.Order == schema.OrderDes {
		temp, err := strconv.Atoi(*(column[0].Value))
		if err != nil {
			return nil, err
		}
		maxValue = temp
	} else if config.Condition.Order == schema.OrderAsc {
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

		if config.AggregationNum != "" {
			aggregationNum, convErr := strconv.Atoi(config.AggregationNum)
			if convErr != nil {
				return nil, convErr
			}
			strValue = fmt.Sprintf("%.4f", float64(v)/float64(aggregationNum))
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

// MetricColumnsElement is the part that separated from layout,
// which can be reused by global dashboard.
func MetricColumnsElement(columns []*MetricColumn) []grid.Element {
	var metricColumns []grid.Element
	var columnWidthPerc int

	// For the best display effect, the maximum number of columns that can be displayed
	const MaxColumnNum = 4
	// For the best display effect, the maximum number of gauges
	// that can be displayed in each column
	const MaxGaugeNum = 6
	const TitleHeight = 10

	// Number of columns to display, each column represents a global metric
	// The number should be less than or equal to MaxColumnNum
	columnNum := int(math.Min(MaxColumnNum, float64(len(columns))))

	// columnWidthPerc should be in the range (0, 100)
	if columnNum > 1 {
		columnWidthPerc = 100 / columnNum
	} else {
		columnWidthPerc = 99
	}

	for i := 0; i < columnNum; i++ {
		var column []grid.Element
		column = append(column, grid.RowHeightPerc(TitleHeight, grid.Widget(columns[i].title)))

		// Number of gauge in a column, each gauge represents a service or endpoint
		// The number should be less than or equal to MaxGaugeNum
		gaugeNum := int(math.Min(MaxGaugeNum, float64(len(columns[i].gauges))))
		gaugeHeight := int(math.Floor(float64(99-TitleHeight) / float64(gaugeNum)))

		for j := 0; j < gaugeNum; j++ {
			column = append(column, grid.RowHeightPerc(gaugeHeight, grid.Widget(columns[i].gauges[j])))
		}
		metricColumns = append(metricColumns, grid.ColWidthPerc(columnWidthPerc, column...))
	}

	return metricColumns
}

func layout(columns []grid.Element) ([]container.Option, error) {
	builder := grid.New()
	builder.Add(
		grid.RowHeightPerc(10),
		grid.RowHeightPerc(80, columns...),
	)

	gridOpts, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return gridOpts, nil
}

func Display(ctx *cli.Context, metrics [][]*schema.SelectedRecord) error {
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

	var columns []*MetricColumn

	configs, err := dashboard.LoadTemplate(ctx.String("template"))
	if err != nil {
		return nil
	}

	for i, config := range configs.Metrics {
		col, innerErr := NewMetricColumn(metrics[i], &config)
		if innerErr != nil {
			return innerErr
		}
		columns = append(columns, col)
	}

	gridOpts, err := layout(MetricColumnsElement(columns))
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

	con, cancel := context.WithCancel(context.Background())
	quitter := func(keyboard *terminalapi.Keyboard) {
		if strings.EqualFold(keyboard.Key.String(), "q") {
			cancel()
		}
	}

	err = termdash.Run(con, t, c, termdash.KeyboardSubscriber(quitter))

	return err
}
