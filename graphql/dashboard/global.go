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
	"bytes"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/machinebox/graphql"
	"github.com/spf13/viper"
	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/assets"
	"github.com/apache/skywalking-cli/graphql/client"
	"github.com/apache/skywalking-cli/graphql/schema"
	"github.com/apache/skywalking-cli/graphql/utils"
	"github.com/apache/skywalking-cli/logger"

	"gopkg.in/yaml.v2"
)

type ButtonTemplate struct {
	Texts    []string `mapstructure:"texts"`
	ColorNum int      `mapstructure:"colorNumber"`
	Height   int      `mapstructure:"height"`
}

type MetricTemplate struct {
	Condition      schema.TopNCondition `mapstructure:"condition"`
	Title          string               `mapstructure:"title"`
	Aggregation    string               `mapstructure:"aggregation"`
	AggregationNum string               `mapstructure:"aggregationNum"`
}

type ChartTemplate struct {
	Condition   schema.MetricsCondition `mapstructure:"condition"`
	Title       string                  `mapstructure:"title"`
	Unit        string                  `mapstructure:"unit"`
	Labels      string                  `mapstructure:"labels"`
	LabelsIndex string                  `mapstructure:"labelsIndex"`
}

type GlobalTemplate struct {
	Buttons         ButtonTemplate   `mapstructure:"buttons"`
	Metrics         []MetricTemplate `mapstructure:"metrics"`
	ResponseLatency ChartTemplate    `mapstructure:"responseLatency"`
	HeatMap         ChartTemplate    `mapstructure:"heatMap"`
}

type GlobalData struct {
	Metrics         [][]*schema.SelectedRecord `json:"metrics"`
	ResponseLatency []map[string]float64       `json:"responseLatency"`
	HeatMap         schema.HeatMap             `json:"heatMap"`
}

// Use singleton pattern to make sure to load template only once.
var globalTemplate *GlobalTemplate

const templateType = "yml"
const DefaultTemplatePath = "templates/dashboard/global.yml"

// newGlobalTemplate create a new GlobalTemplate and set default values for buttons' template.
func newGlobalTemplate() GlobalTemplate {
	return GlobalTemplate{
		Buttons: ButtonTemplate{
			ColorNum: 220,
			Height:   1,
		},
	}
}

// LoadTemplate reads UI template from yaml file.
func LoadTemplate(filename string) (*GlobalTemplate, error) {
	if globalTemplate != nil {
		return globalTemplate, nil
	}

	gt := newGlobalTemplate()
	viper.SetConfigType(templateType)

	var err error
	var byteValue []byte
	if filename == DefaultTemplatePath {
		byteValue = []byte(assets.Read(filename))
	} else {
		byteValue, err = ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
	}

	gt.Buttons.Texts, err = getButtonTexts(byteValue)
	if err != nil {
		return nil, err
	}

	if err := viper.ReadConfig(bytes.NewReader(byteValue)); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&gt); err != nil {
		return nil, err
	}

	globalTemplate = &gt
	return globalTemplate, nil
}

// getButtonTexts get keys in the template file,
// which will be set as texts of buttons in the dashboard.
func getButtonTexts(byteValue []byte) ([]string, error) {
	var ret []string

	c := make(map[string]interface{})
	if err := yaml.Unmarshal(byteValue, &c); err != nil {
		return nil, err
	}

	for s := range c {
		if s != "style" {
			ret = append(ret, strings.Title(s))
		}
	}
	return ret, nil
}

func Metrics(ctx *cli.Context, duration schema.Duration) [][]*schema.SelectedRecord {
	var ret [][]*schema.SelectedRecord

	template, err := LoadTemplate(ctx.String("template"))
	if err != nil {
		return nil
	}

	// Check if there is a template of metrics.
	if template.Metrics == nil {
		return nil
	}

	for _, m := range template.Metrics {
		var response map[string][]*schema.SelectedRecord
		request := graphql.NewRequest(assets.Read("graphqls/dashboard/SortMetrics.graphql"))
		request.Var("condition", m.Condition)
		request.Var("duration", duration)

		client.ExecuteQueryOrFail(ctx, request, &response)
		ret = append(ret, response["result"])
	}

	return ret
}

func responseLatency(ctx *cli.Context, duration schema.Duration) []map[string]float64 {
	var response map[string][]*schema.MetricsValues

	template, err := LoadTemplate(ctx.String("template"))
	if err != nil {
		return nil
	}

	// Check if there is a template of response latency.
	if template.ResponseLatency == (ChartTemplate{}) {
		return nil
	}

	// LabelsIndex in the template file is string type, like "0, 1, 2",
	// need use ", " to split into string array for graphql query.
	labelsIndex := strings.Split(template.ResponseLatency.LabelsIndex, ", ")

	request := graphql.NewRequest(assets.Read("graphqls/dashboard/LabeledMetricsValues.graphql"))
	request.Var("duration", duration)
	request.Var("condition", template.ResponseLatency.Condition)
	request.Var("labels", labelsIndex)

	client.ExecuteQueryOrFail(ctx, request, &response)

	// Convert metrics values to map type data.
	responseLatency := response["result"]
	reshaped := make([]map[string]float64, len(responseLatency))
	for _, mvs := range responseLatency {
		index, err := strconv.Atoi(strings.TrimSpace(*mvs.Label))
		if err != nil {
			logger.Log.Fatalln(err)
			return nil
		}
		reshaped[index] = utils.MetricsToMap(duration, *mvs.Values)
	}

	return reshaped
}

func heatMap(ctx *cli.Context, duration schema.Duration) schema.HeatMap {
	var response map[string]schema.HeatMap

	template, err := LoadTemplate(ctx.String("template"))
	if err != nil {
		return schema.HeatMap{}
	}

	// Check if there is a template of heat map.
	if template.HeatMap == (ChartTemplate{}) {
		return schema.HeatMap{}
	}

	request := graphql.NewRequest(assets.Read("graphqls/dashboard/HeatMap.graphql"))
	request.Var("duration", duration)
	request.Var("condition", template.HeatMap.Condition)

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
