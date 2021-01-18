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
	"strings"

	"github.com/apache/skywalking-cli/internal/logger"

	"github.com/urfave/cli"
)

// ParseParameters parses parameters of the event message from args.
func ParseParameters(paras cli.Args) map[string]string {
	ret := make(map[string]string, len(paras))

	for _, para := range paras {
		sepIndex := strings.Index(para, "=")
		// To make sure that len(k) > 0 && len(v) > 0
		if len(para) >= 3 && sepIndex >= 1 && sepIndex < len(para)-1 {
			k := para[:sepIndex]
			v := para[sepIndex+1:]
			ret[k] = v
		} else {
			logger.Log.Warnf("%s is not a vaild parameter, should like `key=value`\n", para)
		}
	}

	return ret
}
