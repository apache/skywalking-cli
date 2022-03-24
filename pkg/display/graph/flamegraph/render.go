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
	_ "embed"
)

type StackType int

var (
	StackTypeUserSpace   StackType
	StackTypeKernelSpace StackType = 2
)

type ProfilingDataTree struct {
	Elements []ProfilingDataStackElement
}

type ProfilingDataStackElement interface {
	ID() string
	ParentID() string
	DumpCount() int64
	Symbol() string
	StackType() StackType
}

type StackElement struct {
	ID       string
	ParentID string
	Depth    int64
	Left     int64
	Count    int64
	Type     StackType
	Symbol   string
}

func GenerateStackElementByProfilingData(element ProfilingDataStackElement, depth, left int64) *StackElement {
	return &StackElement{
		ID:       element.ID(),
		ParentID: element.ParentID(),
		Depth:    depth,
		Left:     left,
		Count:    element.DumpCount(),
		Type:     element.StackType(),
		Symbol:   element.Symbol(),
	}
}

//go:embed flamegraph.html
var flameGraphHTML string
