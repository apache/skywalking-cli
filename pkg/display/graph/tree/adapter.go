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

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/apache/skywalking-cli/pkg/util"
)

func Adapt(trace api.Trace) (roots []*Node, serviceNames []string) {
	all := make(map[string]*Node)
	set := make(map[string]bool)
	var traceID string
	for _, span := range trace.Spans {
		if !set[span.ServiceCode] {
			serviceNames = append(serviceNames, span.ServiceCode)
			set[span.ServiceCode] = true
		}
		all[id(span)] = node(span)
		if traceID == "" {
			traceID = span.TraceID
		}
	}

	seen := make(map[string]bool)

	for _, span := range trace.Spans {
		if isRoot(span) {
			roots = append(roots, all[id(span)])
			seen[id(span)] = true
		}
		for _, ref := range span.Refs {
			if all[id0(ref)] == nil {
				for i := 0; i <= ref.ParentSpanID; i++ {
					if traceID != ref.TraceID {
						continue
					}
					virtualSpan := virtualSpan(i, *ref)
					if all[id(virtualSpan)] != nil {
						continue
					}
					all[id(virtualSpan)] = node(virtualSpan)
					if i == 0 {
						roots = append(roots, all[id(virtualSpan)])
						seen[id(virtualSpan)] = true
					} else if all[id1(ref)] != nil {
						all[id1(ref)].Children = append(all[id1(ref)].Children, all[id(virtualSpan)])
						seen[id(virtualSpan)] = true
					}
				}
			}
		}
	}
	buildTree(all, seen, trace)
	return roots, serviceNames
}

func buildTree(all map[string]*Node, seen map[string]bool, trace api.Trace) {
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
				refData := all[id0(ref)]
				if refData != nil {
					refData.Children = append(refData.Children, all[id(span)])
					seen[id(span)] = true
				}
			}
		}
	}
}

func virtualSpan(spanID int, ref api.Ref) *api.Span {
	endpointName := fmt.Sprintf("VNode: %s", ref.ParentSegmentID)
	component := fmt.Sprintf("VirtualNode: #%d", spanID)
	peer := "No Peer"
	fail := true
	layer := "Broken"
	span := api.Span{
		TraceID:      ref.TraceID,
		SegmentID:    ref.ParentSegmentID,
		SpanID:       spanID,
		ParentSpanID: spanID - 1,
		EndpointName: &endpointName,
		ServiceCode:  "VirtualNode",
		Type:         fmt.Sprintf("[Broken] %s", ref.Type),
		Peer:         &peer,
		Component:    &component,
		IsError:      &fail,
		Layer:        &layer,
		Tags:         nil,
		Logs:         nil,
	}
	return &span
}

func isRoot(span *api.Span) bool {
	return span.SpanID == 0 && span.ParentSpanID == -1 && len(span.Refs) == 0
}

func id(span *api.Span) string {
	return fmt.Sprintf("%s:%s:%d", span.TraceID, span.SegmentID, span.SpanID)
}

func pid(span *api.Span) string {
	return fmt.Sprintf("%s:%s:%d", span.TraceID, span.SegmentID, span.ParentSpanID)
}

func id0(ref *api.Ref) string {
	return fmt.Sprintf("%s:%s:%d", ref.TraceID, ref.ParentSegmentID, ref.ParentSpanID)
}

func id1(ref *api.Ref) string {
	return fmt.Sprintf("%s:%s:%d", ref.TraceID, ref.ParentSegmentID, ref.ParentSpanID-1)
}

func node(span *api.Span) *Node {
	return &Node{
		Children: []*Node{},
		Value:    util.Stringify{Str: value(span)},
		Detail:   detail(span),
	}
}

func value(span *api.Span) string {
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

func detail(span *api.Span) string {
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
