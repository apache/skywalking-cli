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

package inspect

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/admin/inspect"
	"github.com/apache/skywalking-cli/pkg/admin/preflight"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
)

var Command = &cli.Command{
	Name:  "inspect",
	Usage: "Browse the metric catalog and entities from the admin-server `inspect` module",
	Subcommands: []*cli.Command{
		metricsCommand,
		entitiesCommand,
	},
}

var metricsCommand = &cli.Command{
	Name:  "metrics",
	Usage: "List the registered metric catalog (GET /inspect/metrics)",
	UsageText: `List the registered metrics with their type, scope and supported downsamplings.

Examples:
1. List every metric:
$ swctl admin inspect metrics

2. List service metrics that /inspect/entities accepts:
$ swctl admin inspect metrics --catalog SERVICE --mqe-queryable`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "regex",
			Usage: "filter metric names by a Java `regex` (default matches all)",
		},
		&cli.StringSliceFlag{
			Name:  "type",
			Usage: "filter by metric `type` (REGULAR_VALUE / LABELED_VALUE / HEATMAP / SAMPLED_RECORD); repeatable",
		},
		&cli.StringSliceFlag{
			Name:  "catalog",
			Usage: "filter by `catalog` (SERVICE / SERVICE_INSTANCE / ENDPOINT / *_RELATION); repeatable",
		},
		&cli.BoolFlag{
			Name:  "mqe-queryable",
			Usage: "return only metrics that /inspect/entities accepts (REGULAR_VALUE, LABELED_VALUE)",
		},
	},
	Action: func(ctx *cli.Context) error {
		metrics, err := inspect.ListMetrics(ctx.Context, inspect.MetricsOptions{
			Regex:        ctx.String("regex"),
			Types:        ctx.StringSlice("type"),
			Catalogs:     ctx.StringSlice("catalog"),
			MQEQueryable: ctx.Bool("mqe-queryable"),
		})
		if err != nil {
			return preflight.Explain(ctx.Context, err, preflight.ModuleInspect, "SW_INSPECT")
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: metrics})
	},
}

var entitiesCommand = &cli.Command{
	Name:  "entities",
	Usage: "Enumerate the entities holding values for a metric (GET /inspect/entities)",
	UsageText: `Enumerate the entities currently holding values for a metric over a time range.
Each row carries an MQE-ready entity to paste into a follow-up "swctl metrics exec"
(execExpression) query. Only REGULAR_VALUE / LABELED_VALUE metrics are accepted.

Examples:
1. Entities reporting service_cpm in the last 30 minutes:
$ swctl admin inspect entities --metric service_cpm`,
	Flags: flags.Flags(
		flags.DurationFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "metric",
				Usage:    "the `metric` name to enumerate entities for",
				Required: true,
			},
			&cli.IntFlag{
				Name:  "limit",
				Usage: fmt.Sprintf("max rows scanned at the storage layer (1-%d, server default 300)", inspect.MaxLimit),
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.DurationInterceptor,
	),
	Action: func(ctx *cli.Context) error {
		limit := ctx.Int("limit")
		if limit < 0 || limit > inspect.MaxLimit {
			return fmt.Errorf("--limit must be between 1 and %d", inspect.MaxLimit)
		}
		step := ctx.Generic("step").(*model.StepEnumValue).Selected

		entities, err := inspect.ListEntities(ctx.Context, inspect.EntitiesOptions{
			Metric: ctx.String("metric"),
			Start:  ctx.String("start"),
			End:    ctx.String("end"),
			Step:   string(step),
			Limit:  limit,
		})
		if err != nil {
			return preflight.Explain(ctx.Context, err, preflight.ModuleInspect, "SW_INSPECT")
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: entities})
	},
}
