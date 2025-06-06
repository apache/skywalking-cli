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
	"context"
	"fmt"
	"strings"

	"github.com/apache/skywalking-cli/pkg/contextkey"
	d "github.com/apache/skywalking-cli/pkg/display/displayable"

	"github.com/apache/skywalking-cli/pkg/display/graph"

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

// Display the object in the style specified in flag --display
func Display(ctx context.Context, displayable *d.Displayable) error {
	displayStyle := ctx.Value(contextkey.Display{}).(string)
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
