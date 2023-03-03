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

package continuous

import (
	"fmt"
	"os"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"

	"github.com/urfave/cli/v2"

	"gopkg.in/yaml.v2"

	api "skywalking.apache.org/repo/goapi/query"
)

type PolicyConfig struct {
	Policies []*PolicyTarget `yaml:"policy"`
}

type PolicyTarget struct {
	Type     string        `yaml:"type"`
	Checkers []*PolicyItem `yaml:"checkers"`
}

type PolicyItem struct {
	Type      string   `yaml:"type"`
	Threshold string   `yaml:"threshold"`
	Period    int      `yaml:"period"`
	Count     int      `yaml:"count"`
	URIList   []string `yaml:"uriList"`
	URIRegex  string   `yaml:"uriRegex"`
}

var SetPolicyCommand = &cli.Command{
	Name:    "set",
	Aliases: []string{"s"},
	Usage:   "Set the continuous profiling policy to service",
	UsageText: `Set the continuous profiling policy to service

Examples:
1. Set the service continuous profiling policy
$ swctl profiling continuous set --service-name=abc --config=/path/to/config.yaml`,
	Flags: flags.Flags(
		flags.ServiceFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Usage:    "the `config` file define the service policy configuration, if not provide then make the service policy is empty",
				Required: false,
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseService(true),
	),
	Action: func(ctx *cli.Context) error {
		serviceID := ctx.String("service-id")

		configFile := ctx.String("config")

		targets := make([]*api.ContinuousProfilingPolicyTargetCreation, 0)
		if configFile != "" {
			config, err := os.ReadFile(configFile)
			if err != nil {
				return err
			}
			r := &PolicyConfig{}
			if e := yaml.Unmarshal(config, r); e != nil {
				return err
			}

			targets, err = parsingPolicyConfig(r)
			if err != nil {
				return err
			}
		}

		request := &api.ContinuousProfilingPolicyCreation{
			ServiceID: serviceID,
			Targets:   targets,
		}

		result, err := profiling.SetContinuousProfilingPolicy(ctx, request)
		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: result, Condition: request})
	},
}

func parsingPolicyConfig(conf *PolicyConfig) ([]*api.ContinuousProfilingPolicyTargetCreation, error) {
	result := make([]*api.ContinuousProfilingPolicyTargetCreation, 0)
	if len(conf.Policies) == 0 {
		return result, nil
	}
	for _, t := range conf.Policies {
		var realTarget api.ContinuousProfilingTargetType
		for _, targetType := range api.AllContinuousProfilingTargetType {
			if t.Type == targetType.String() {
				realTarget = targetType
				break
			}
		}
		if realTarget == "" {
			return nil, fmt.Errorf("cannot found the target: %s", t.Type)
		}

		target := &api.ContinuousProfilingPolicyTargetCreation{
			TargetType: realTarget,
		}
		for _, c := range t.Checkers {
			var realMonitorType api.ContinuousProfilingMonitorType
			for _, monitorType := range api.AllContinuousProfilingMonitorType {
				if c.Type == monitorType.String() {
					realMonitorType = monitorType
					break
				}
			}
			if realMonitorType == "" {
				return nil, fmt.Errorf("cannot fount the monitor type: %s", c.Type)
			}

			item := &api.ContinuousProfilingPolicyItemCreation{
				Type:      realMonitorType,
				Threshold: c.Threshold,
				Period:    c.Period,
				Count:     c.Count,
				URIList:   c.URIList,
			}

			if c.URIRegex != "" {
				item.URIRegex = &c.URIRegex
			}

			target.CheckItems = append(target.CheckItems, item)
		}

		result = append(result, target)
	}
	return result, nil
}
