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
	"os"
	"strings"
	"sync"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	api "skywalking.apache.org/repo/goapi/query"

	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"

	"gopkg.in/yaml.v2"

	"github.com/apache/skywalking-cli/assets"
	"github.com/apache/skywalking-cli/internal/logger"
	"github.com/apache/skywalking-cli/pkg/graphql/metrics"
	"github.com/apache/skywalking-cli/pkg/graphql/utils"
)

type ButtonTemplate struct {
	Texts    []string `mapstructure:"texts"`
	ColorNum int      `mapstructure:"colorNumber"`
	Height   int      `mapstructure:"height"`
}

type MetricTemplate struct {
	Condition      api.TopNCondition `mapstructure:"condition"`
	Title          string            `mapstructure:"title"`
	Aggregation    string            `mapstructure:"aggregation"`
	AggregationNum string            `mapstructure:"aggregationNum"`
}

type ChartTemplate struct {
	Condition api.MetricsCondition `mapstructure:"condition"`
	Title     string               `mapstructure:"title"`
	Unit      string               `mapstructure:"unit"`
	Relabels  string               `mapstructure:"relabels"`
	Labels    string               `mapstructure:"labels"`
}

type GlobalTemplate struct {
	Buttons         ButtonTemplate   `mapstructure:"buttons"`
	Metrics         []MetricTemplate `mapstructure:"metrics"`
	ResponseLatency ChartTemplate    `mapstructure:"responseLatency"`
	HeatMap         ChartTemplate    `mapstructure:"heatMap"`
}

type GlobalData struct {
	Metrics         [][]*api.SelectedRecord       `json:"metrics"`
	ResponseLatency map[string]map[string]float64 `json:"responseLatency"`
	HeatMap         api.HeatMap                   `json:"heatMap"`
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
		byteValue, err = os.ReadFile(filename)
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
			ret = append(ret, cases.Title(language.Und).String(s))
		}
	}
	return ret, nil
}

func Metrics(ctx *cli.Context, duration api.Duration) ([][]*api.SelectedRecord, error) {
	var ret [][]*api.SelectedRecord

	template, err := LoadTemplate(ctx.String("template"))
	if err != nil {
		return nil, nil
	}

	// Check if there is a template of metrics.
	if template.Metrics == nil {
		return nil, nil
	}

	for _, m := range template.Metrics {
		sortMetrics, err := metrics.SortMetrics(ctx, m.Condition, duration)
		if err != nil {
			return nil, err
		}
		ret = append(ret, sortMetrics)
	}

	return ret, nil
}

func responseLatency(ctx *cli.Context, duration api.Duration) map[string]map[string]float64 {
	template, err := LoadTemplate(ctx.String("template"))
	if err != nil {
		return nil
	}

	// Check if there is a template of response latency.
	if template.ResponseLatency == (ChartTemplate{}) {
		return nil
	}

	// Labels in the template file is string type, like "0, 1, 2",
	// need use ", " to split into string array for graphql query.
	labels := strings.Split(template.ResponseLatency.Labels, ",")
	relabels := strings.Split(template.ResponseLatency.Relabels, ",")

	responseLatency, err := metrics.MultipleLinearIntValues(ctx, template.ResponseLatency.Condition, labels, duration)

	if err != nil {
		logger.Log.Fatalln(err)
	}

	mapping := make(map[string]string, len(labels))
	for i := 0; i < len(labels); i++ {
		mapping[labels[i]] = relabels[i]
	}

	// Convert metrics values to map type data.
	return utils.MetricsValuesArrayToMap(duration, responseLatency, mapping)
}

func heatMap(ctx *cli.Context, duration api.Duration) (api.HeatMap, error) {
	template, err := LoadTemplate(ctx.String("template"))
	if err != nil {
		return api.HeatMap{}, nil
	}

	// Check if there is a template of heat map.
	if template.HeatMap == (ChartTemplate{}) {
		return api.HeatMap{}, nil
	}

	return metrics.Thermodynamic(ctx, template.HeatMap.Condition, duration)
}

func Global(ctx *cli.Context, duration api.Duration) (*GlobalData, error) {
	// Load template to `globalTemplate`, so the subsequent three calls can ues it directly.
	_, err := LoadTemplate(ctx.String("template"))
	if err != nil {
		return nil, nil
	}

	errors := make(chan error)
	done := make(chan bool)

	// Use three goroutines to enable concurrent execution of three graphql queries.
	var wg sync.WaitGroup
	wg.Add(3)
	var m [][]*api.SelectedRecord
	go func() {
		m, err = Metrics(ctx, duration)
		if err != nil {
			errors <- err
		}
		wg.Done()
	}()
	var rl map[string]map[string]float64
	go func() {
		rl = responseLatency(ctx, duration)
		wg.Done()
	}()
	var hm api.HeatMap
	go func() {
		hm, err = heatMap(ctx, duration)
		if err != nil {
			errors <- err
		}
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		break
	case err := <-errors:
		close(errors)
		return nil, err
	}

	var globalData GlobalData
	globalData.Metrics = m
	globalData.ResponseLatency = rl
	globalData.HeatMap = hm

	return &globalData, nil
}
