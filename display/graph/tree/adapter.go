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

package tree

import (
	"fmt"
	"strings"

	"github.com/apache/skywalking-cli/graphql/schema"
	"github.com/apache/skywalking-cli/util"
)

func Adapt(trace schema.Trace) []*Node {
	all := make(map[string]*Node)

	for _, span := range trace.Spans {
		all[id(span)] = node(span)
	}

	seen := make(map[string]bool)

	var roots []*Node

	for _, span := range trace.Spans {
		if isRoot(span) {
			roots = append(roots, all[id(span)])
			seen[id(span)] = true
		}
	}

	for len(seen) < len(trace.Spans) {
		for _, span := range trace.Spans {
			if seen[id(span)] {
				continue
			}

			if all[pid(span)] != nil {
				all[pid(span)].Children = append(all[pid(span)].Children, all[id(span)])
				seen[id(span)] = true
			}

			for _, ref := range span.Refs {
				if all[id0(ref)] != nil {
					all[id0(ref)].Children = append(all[id0(ref)].Children, all[id(span)])
					seen[id(span)] = true
				}
			}
		}
	}

	return roots
}

func isRoot(span *schema.Span) bool {
	return span.SpanID == 0 && span.ParentSpanID == -1 && len(span.Refs) == 0
}

func id(span *schema.Span) string {
	return fmt.Sprintf("%s:%s:%d", span.TraceID, span.SegmentID, span.SpanID)
}

func pid(span *schema.Span) string {
	return fmt.Sprintf("%s:%s:%d", span.TraceID, span.SegmentID, span.ParentSpanID)
}

func id0(ref *schema.Ref) string {
	return fmt.Sprintf("%s:%s:%d", ref.TraceID, ref.ParentSegmentID, ref.ParentSpanID)
}

func node(span *schema.Span) *Node {
	return &Node{
		Children: []*Node{},
		Value:    util.Stringify{Str: value(span)},
		Detail:   detail(span),
	}
}

func value(span *schema.Span) string {
	if *span.IsError {
		return fmt.Sprintf(
			"[|%s| %s [%s/%s]](mod:bold,fg:white,bg:red)",
			span.Type, *span.EndpointName, *span.Component, *span.Layer,
		)
	}

	return fmt.Sprintf("[|%s|](fg:bold,fg:green) %s [[%s/%s]](mod:bold,fg:green)",
		span.Type, *span.EndpointName, *span.Component, *span.Layer,
	)
}

func detail(span *schema.Span) string {
	var lines []string

	lines = append(lines,
		fmt.Sprintf("[Endpoint    :](mod:bold,fg:red) %s", *span.EndpointName),
		fmt.Sprintf("[Span Type   :](mod:bold,fg:red) %s", span.Type),
		fmt.Sprintf("[Component   :](mod:bold,fg:red) %s", *span.Component),
		fmt.Sprintf("[Peer        :](mod:bold,fg:red) %s", *span.Peer),
		fmt.Sprintf("[Error       :](mod:bold,fg:red) %t", *span.IsError),
	)

	for _, tag := range span.Tags {
		lines = append(lines, fmt.Sprintf("[%-12s:](mod:bold,fg:red) %s", tag.Key, *tag.Value))
	}

	for _, log := range span.Logs {
		for _, datum := range log.Data {
			lines = append(lines, fmt.Sprintf("%-12s: %s", datum.Key, *datum.Value))
		}
	}

	return strings.Join(lines, "\n")
}
