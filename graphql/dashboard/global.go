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

package dashboard

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/machinebox/graphql"
	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/assets"
	"github.com/apache/skywalking-cli/graphql/client"
	"github.com/apache/skywalking-cli/graphql/schema"
)

type MetricConfig struct {
	Condition      schema.TopNCondition `json:"condition"`
	Title          string               `json:"title"`
	Aggregation    string               `json:"aggregation"`
	AggregationNum string               `json:"aggregationNum"`
}

type ChartConfig struct {
	Condition schema.MetricsCondition `json:"condition"`
	Title     string                  `json:"title"`
	Unit      string                  `json:"unit"`
	Labels    string                  `json:"labels"`
}

type GlobalConfig struct {
	Metrics         []MetricConfig `json:"metrics"`
	ResponseLatency ChartConfig    `json:"responseLatency"`
	HeatMap         ChartConfig    `json:"heatMap"`
}

type GlobalData struct {
	Metrics         [][]*schema.SelectedRecord `json:"metrics"`
	ResponseLatency []*schema.MetricsValues    `json:"responseLatency"`
	HeatMap         schema.HeatMap             `json:"heatMap"`
}

func LoadConfig(filename string) (*GlobalConfig, error) {
	var config GlobalConfig

	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err := json.Unmarshal(byteValue, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func Metrics(ctx *cli.Context, duration schema.Duration) [][]*schema.SelectedRecord {
	var ret [][]*schema.SelectedRecord
	configs, err := LoadConfig("assets/config/Dashboard.Global.json")
	if err != nil {
		return nil
	}

	for _, m := range configs.Metrics {
		var response map[string][]*schema.SelectedRecord
		request := graphql.NewRequest(assets.Read("graphqls/dashboard/SortMetrics.graphql"))
		request.Var("condition", m.Condition)
		request.Var("duration", duration)

		client.ExecuteQueryOrFail(ctx, request, &response)
		ret = append(ret, response["result"])
	}

	return ret
}

func responseLatency(ctx *cli.Context, duration schema.Duration) []*schema.MetricsValues {
	var response map[string][]*schema.MetricsValues

	request := graphql.NewRequest(assets.Read("graphqls/dashboard/LabeledMetricsValues.graphql"))
	request.Var("duration", duration)

	client.ExecuteQueryOrFail(ctx, request, &response)

	return response["result"]
}

func heatMap(ctx *cli.Context, duration schema.Duration) schema.HeatMap {
	var response map[string]schema.HeatMap

	request := graphql.NewRequest(assets.Read("graphqls/dashboard/HeatMap.graphql"))
	request.Var("duration", duration)

	client.ExecuteQueryOrFail(ctx, request, &response)

	return response["result"]
}

func Global(ctx *cli.Context, duration schema.Duration) *GlobalData {
	var globalData GlobalData

	globalData.Metrics = Metrics(ctx, duration)
	globalData.ResponseLatency = responseLatency(ctx, duration)
	globalData.HeatMap = heatMap(ctx, duration)

	return &globalData
}
