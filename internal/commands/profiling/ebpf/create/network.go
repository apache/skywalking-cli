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

package create

import (
	"fmt"
	"os"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/profiling"

	"gopkg.in/yaml.v2"

	"github.com/urfave/cli/v2"

	api "skywalking.apache.org/repo/goapi/query"
)

type SamplingConfig struct {
	Samplings []*SamplingRule `yaml:"samplings"`
}

type SamplingRule struct {
	URIPattern  *string          `yaml:"uri_pattern"`
	MinDuration *int             `yaml:"min_duration"`
	When4xx     *bool            `yaml:"when_4xx"`
	When5xx     *bool            `yaml:"when_5xx"`
	Setting     *SamplingSetting `yaml:"setting"`
}

type SamplingSetting struct {
	RequireRequest  *bool `yaml:"require_request"`
	MaxRequestSize  *int  `yaml:"max_request_size"`
	RequireResponse *bool `yaml:"require_response"`
	MaxResponseSize *int  `yaml:"max_response_size"`
}

var NetworkCreateCommand = &cli.Command{
	Name:    "network",
	Aliases: []string{"net"},
	Usage:   "Create a new ebpf network profiling task",
	UsageText: `Create a new ebpf network profiling task

Examples:
1. Create ebpf network profiling task
$ swctl profiling ebpf create network --service-instance-id=abc`,
	Flags: flags.Flags(
		flags.ServiceFlags,
		flags.InstanceFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "sampling-config",
				Usage:    "the `sampling-config` file define how to sampling the network data, if not then then ignore data sampling",
				Required: false,
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseInstance(true),
	),
	Action: func(ctx *cli.Context) error {
		instanceID := ctx.String("instance-id")

		samplingConfigFile := ctx.String("sampling-config")

		// convert the sampling rule
		var samplings = make([]*api.EBPFNetworkSamplingRule, 0)
		if samplingConfigFile != "" {
			config, err := os.ReadFile(samplingConfigFile)
			if err != nil {
				return err
			}
			r := &SamplingConfig{}
			if e := yaml.Unmarshal(config, r); e != nil {
				return e
			}

			samplings, err = parsingNetworkSampling(r)
			if err != nil {
				return err
			}
		}

		request := &api.EBPFProfilingNetworkTaskRequest{
			InstanceID: instanceID,
			Samplings:  samplings,
		}

		task, err := profiling.CreateEBPFNetworkProfilingTask(ctx, request)

		if err != nil {
			return err
		}

		return display.Display(ctx, &displayable.Displayable{Data: task, Condition: request})
	},
}

func parsingNetworkSampling(config *SamplingConfig) ([]*api.EBPFNetworkSamplingRule, error) {
	rules := make([]*api.EBPFNetworkSamplingRule, 0)
	if config == nil || len(config.Samplings) == 0 {
		return rules, nil
	}
	for _, conf := range config.Samplings {
		rule := &api.EBPFNetworkSamplingRule{}
		rule.URIRegex = conf.URIPattern
		rule.MinDuration = conf.MinDuration
		if conf.When4xx == nil {
			return nil, fmt.Errorf("the when_4xx is required")
		}
		rule.When4xx = *conf.When4xx
		if conf.When5xx == nil {
			return nil, fmt.Errorf("the when_5xx is required")
		}
		rule.When5xx = *conf.When5xx

		confSetting := conf.Setting
		if confSetting == nil {
			return nil, fmt.Errorf("the sampling settings is required")
		}
		setting := &api.EBPFNetworkDataCollectingSettings{}
		rule.Settings = setting
		if confSetting.RequireRequest == nil {
			return nil, fmt.Errorf("the sampling request is required")
		}
		setting.RequireCompleteRequest = *confSetting.RequireRequest
		setting.MaxRequestSize = confSetting.MaxRequestSize
		if confSetting.RequireResponse == nil {
			return nil, fmt.Errorf("the sampling response is required")
		}
		setting.RequireCompleteResponse = *confSetting.RequireResponse
		setting.MaxResponseSize = confSetting.MaxResponseSize

		rules = append(rules, rule)
	}

	return rules, nil
}
