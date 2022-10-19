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

package aggregation

import (
	"fmt"
	"strconv"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/graphql/utils"

	"github.com/urfave/cli/v2"

	api "skywalking.apache.org/repo/goapi/query"
)

// buildSortedCondition from context of cli, the first argument must be the count of top N
func buildSortedCondition(ctx *cli.Context, parseScope bool) (*api.TopNCondition, *api.Duration, error) {
	start := ctx.String("start")
	end := ctx.String("end")
	step := ctx.Generic("step").(*model.StepEnumValue).Selected

	metricsName := ctx.String("name")
	var scope *api.Scope
	if parseScope {
		tmp := utils.ParseScopeInTop(metricsName)
		scope = &tmp
	}
	order := ctx.Generic("order").(*model.OrderEnumValue).Selected
	topN := 5
	parentServiceID := ctx.String("service-id")
	parentService, normal, err := interceptor.ParseServiceID(parentServiceID)
	if err != nil {
		return nil, nil, err
	}

	if ctx.NArg() > 0 {
		nn, err2 := strconv.Atoi(ctx.Args().First())
		if err2 != nil {
			return nil, nil, fmt.Errorf("the 1st argument must be a number: %v", err2)
		}
		topN = nn
	}

	return &api.TopNCondition{
			Name:          metricsName,
			ParentService: &parentService,
			Normal:        &normal,
			Scope:         scope,
			TopN:          topN,
			Order:         order,
		}, &api.Duration{
			Start: start,
			End:   end,
			Step:  step,
		}, nil
}

func buildReadRecordsCondition(ctx *cli.Context) (*api.RecordCondition, *api.Duration, error) {
	start := ctx.String("start")
	end := ctx.String("end")
	step := ctx.Generic("step").(*model.StepEnumValue).Selected

	metricsName := ctx.String("name")
	order := ctx.Generic("order").(*model.OrderEnumValue).Selected
	topN := 5
	entity, err := interceptor.ParseEntity(ctx)
	if err != nil {
		return nil, nil, err
	}

	if ctx.NArg() > 0 {
		nn, err2 := strconv.Atoi(ctx.Args().First())
		if err2 != nil {
			return nil, nil, fmt.Errorf("the 1st argument must be a number: %v", err2)
		}
		topN = nn
	}

	return &api.RecordCondition{
			Name:         metricsName,
			ParentEntity: entity,
			TopN:         topN,
			Order:        order,
		}, &api.Duration{
			Start: start,
			End:   end,
			Step:  step,
		}, nil
}
