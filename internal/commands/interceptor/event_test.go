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

package interceptor

import (
	"reflect"
	"testing"

	"github.com/urfave/cli/v2"
)

type args []string

func (a *args) Get(n int) string {
	if len(*a) > n {
		return (*a)[n]
	}
	return ""
}

func (a *args) First() string {
	return a.Get(0)
}

func (a *args) Tail() []string {
	if a.Len() >= 2 {
		tail := []string((*a)[1:])
		ret := make([]string, len(tail))
		copy(ret, tail)
		return ret
	}
	return []string{}
}

func (a *args) Len() int {
	return len(*a)
}

func (a *args) Present() bool {
	return a.Len() != 0
}

func (a *args) Slice() []string {
	ret := make([]string, len(*a))
	copy(ret, *a)
	return ret
}

func TestParseParameters(t *testing.T) {
	tests := []struct {
		name    string
		args    cli.Args
		want    map[string]string
		wantErr bool
	}{
		{
			name:    "no parameters",
			args:    &args{},
			want:    map[string]string{},
			wantErr: false,
		},
		{
			name: "all parameters are invalid",
			args: &args{
				"key",
				"key=",
				"=value",
				"=",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "all parameters are valid",
			args: &args{
				"key=value",
				"k=v===",
				"kk=====",
			},
			want: map[string]string{
				"key": "value",
				"k":   "v===",
				"kk":  "====",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseParameters(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseParameters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseParameters() got = %v, want %v", got, tt.want)
			}
		})
	}
}
