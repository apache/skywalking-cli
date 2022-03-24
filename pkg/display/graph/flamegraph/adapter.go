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

import api "skywalking.apache.org/repo/goapi/query"

type eBPFStackElementAdapter struct {
	*api.EBPFProfilingStackElement
}

func (e *eBPFStackElementAdapter) ID() string {
	return e.EBPFProfilingStackElement.ID
}

func (e *eBPFStackElementAdapter) ParentID() string {
	return e.EBPFProfilingStackElement.ParentID
}

func (e *eBPFStackElementAdapter) DumpCount() int64 {
	return e.EBPFProfilingStackElement.DumpCount
}

func (e *eBPFStackElementAdapter) Symbol() string {
	return e.EBPFProfilingStackElement.Symbol
}

func (e *eBPFStackElementAdapter) StackType() StackType {
	var stackType StackType
	switch e.EBPFProfilingStackElement.StackType {
	case api.EBPFProfilingStackTypeKernelSpace:
		stackType = StackTypeKernelSpace
	case api.EBPFProfilingStackTypeUserSpace:
		stackType = StackTypeUserSpace
	}
	return stackType
}

type traceStackElementAdapter struct {
	*api.ProfileStackElement
}

func (e *traceStackElementAdapter) ID() string {
	return e.ProfileStackElement.ID
}

func (e *traceStackElementAdapter) ParentID() string {
	return e.ProfileStackElement.ParentID
}

func (e *traceStackElementAdapter) DumpCount() int64 {
	return int64(e.ProfileStackElement.Count)
}

func (e *traceStackElementAdapter) Symbol() string {
	return e.ProfileStackElement.CodeSignature
}

func (e *traceStackElementAdapter) StackType() StackType {
	return StackTypeUserSpace
}
