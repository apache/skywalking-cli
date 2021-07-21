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

package display

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	yml "gopkg.in/yaml.v2"

	d "github.com/apache/skywalking-cli/pkg/display/displayable"

	"github.com/apache/skywalking-cli/pkg/display/graph"

	"github.com/urfave/cli"

	"github.com/apache/skywalking-cli/pkg/display/json"
	"github.com/apache/skywalking-cli/pkg/display/table"
	"github.com/apache/skywalking-cli/pkg/display/yaml"
)

const (
	JSON  = "json"
	YAML  = "yaml"
	TABLE = "table"
	GRAPH = "graph"
)

type STYLE struct {
	DISPLAY [][]string `yaml:"DISPLAY"`
}

// Display the object in the style specified in flag --display
func Display(ctx *cli.Context, displayable *d.Displayable) error {
	displayStyle := ctx.GlobalString("display")
	if displayStyle == "" {
		commandFullName := ctx.Command.FullName()
		if commandFullName != "" {
			displayStyle = getDisplay(commandFullName)
		} else if ctx.Parent() != nil {
			displayStyle = getDisplay(ctx.Parent().Args()[0])
		}
	}
	if displayStyle == "" {
		displayStyle = "json"
	}
	switch strings.ToLower(displayStyle) {
	case JSON:
		return json.Display(displayable)
	case YAML:
		return yaml.Display(displayable)
	case TABLE:
		return table.Display(displayable)
	case GRAPH:
		return graph.Display(ctx, displayable)
	default:
		return fmt.Errorf("unsupported display style: %s", displayStyle)
	}
}

// getDisplay gets the default display settings
func getDisplay(fullName string) string {
	content, _ := ioutil.ReadFile("../../displaystyle.yaml")
	style := STYLE{}
	err := yml.Unmarshal(content, &style)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}
	for _, command := range style.DISPLAY {
		if fullName == command[0] {
			return command[1]
		}
	}
	return ""
}
