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

package flamegraph

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/apache/skywalking-cli/internal/logger"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/urfave/cli/v2"
)

func DisplayByTrace(ctx *cli.Context, analysis api.ProfileAnalyzation) error {
	trees := make([]*ProfilingDataTree, 0)
	for _, tree := range analysis.Trees {
		elements := make([]ProfilingDataStackElement, 0)
		for _, e := range tree.Elements {
			elements = append(elements, &traceStackElementAdapter{e})
		}
		trees = append(trees, &ProfilingDataTree{Elements: elements})
	}

	return display(ctx, trees)
}

func DisplayByEBPF(ctx *cli.Context, analysis *api.EBPFProfilingAnalyzation) error {
	trees := make([]*ProfilingDataTree, 0)
	for _, tree := range analysis.Trees {
		elements := make([]ProfilingDataStackElement, 0)
		for _, e := range tree.Elements {
			elements = append(elements, &eBPFStackElementAdapter{e})
		}
		trees = append(trees, &ProfilingDataTree{Elements: elements})
	}

	return display(ctx, trees)
}

func display(_ *cli.Context, trees []*ProfilingDataTree) error {
	if len(trees) == 0 {
		return fmt.Errorf("could not find the analysis data")
	}

	// generate flame graph file path for write
	flameGraphPath, err := generateFlameGraphFile()
	if err != nil {
		return err
	}

	// build data
	data := make(map[string]interface{})
	elements, maxDepth := buildFlameGraphElements(trees)
	data["elements"] = elements
	data["maxDepth"] = maxDepth + 1
	data["canvasHeight"] = maxDepth*16 + 30

	// render template
	return renderFlameGraphAndWrite(flameGraphPath, data)
}

func renderFlameGraphAndWrite(path string, data map[string]interface{}) error {
	// render template
	var b bytes.Buffer
	tmpl, err := template.New("flameGraphTemplate").Parse(flameGraphHTML)
	if err != nil {
		return fmt.Errorf("failed to parse flame graph template: %v", err)
	}
	if err = tmpl.Execute(&b, data); err != nil {
		return fmt.Errorf("failed to render flame graph: %v", err)
	}

	// write to file
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create the flame graph file, %v", err)
	}
	_, err = file.Write(b.Bytes())
	if err != nil {
		return fmt.Errorf("could not write the flame graph to file, %v", err)
	}
	logger.Log.Infof("success write the flame graph to: %s", path)
	return nil
}

func generateFlameGraphFile() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get current directory: %v", err)
	}
	flameGraphPath := filepath.Join(wd, buildFileName())
	return flameGraphPath, nil
}

func buildFileName() string {
	return fmt.Sprintf("flame_graph_%d.html", time.Now().Unix())
}

func buildFlameGraphElements(trees []*ProfilingDataTree) (result []*StackElement, maxDepth int64) {
	// adding the root element
	root := &StackElement{Symbol: "all"}
	result = append(result, root)

	// build elements
	var left int64
	for _, t := range trees {
		result, left = buildFlameGraphChildElements(result, t.Elements, nil, left)
	}

	// calculate the root total count
	for _, r := range result {
		if r.ParentID == "0" {
			root.Count += r.Count
		}
		if maxDepth < r.Depth {
			maxDepth = r.Depth
		}
	}
	return result, maxDepth
}

func buildFlameGraphChildElements(renderElements []*StackElement, dataElements []ProfilingDataStackElement,
	parent *StackElement, relativeLeft int64) (result []*StackElement, left int64) {
	parentID := "0"
	var depth int64 = 1
	left = relativeLeft
	if parent != nil {
		parentID = parent.ID
		depth = parent.Depth + 1
		left += parent.Left
	}
	es := findProfilingElementByParentID(dataElements, parentID)
	for _, element := range es {
		current := GenerateStackElementByProfilingData(element, depth, left)
		renderElements = append(renderElements, current)

		left += element.DumpCount()
		renderElements, _ = buildFlameGraphChildElements(renderElements, dataElements, current, relativeLeft)
	}

	return renderElements, left
}

func findProfilingElementByParentID(elements []ProfilingDataStackElement, parentID string) []ProfilingDataStackElement {
	children := make([]ProfilingDataStackElement, 0)
	for _, e := range elements {
		if e.ParentID() == parentID {
			children = append(children, e)
		}
	}
	return children
}
