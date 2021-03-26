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

package event

import (
	common "skywalking/network/common/v3"
	event "skywalking/network/event/v3"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/logger"
	"github.com/apache/skywalking-cli/internal/model"
	"github.com/apache/skywalking-cli/pkg/grpc"

	"github.com/urfave/cli"
)

func Report(ctx *cli.Context) (*common.Commands, error) {
	parameters, err := interceptor.ParseParameters(ctx.Args())
	if err != nil {
		return nil, err
	}

	e := event.Event{
		Uuid: ctx.String("uuid"),
		Source: &event.Source{
			Service:         ctx.String("service"),
			ServiceInstance: ctx.String("instance"),
			Endpoint:        ctx.String("endpoint"),
		},
		Name:       ctx.String("name"),
		Type:       ctx.Generic("type").(*model.EventTypeEnumValue).Selected,
		Message:    ctx.String("message"),
		Parameters: parameters,
		StartTime:  ctx.Int64("startTime"),
		EndTime:    ctx.Int64("endTime"),
	}

	reply, err := grpc.ReportEvent(ctx.GlobalString("grpcAddr"), &e)
	if err != nil {
		logger.Log.Fatalln(err)
		return nil, err
	}

	return reply, nil
}
