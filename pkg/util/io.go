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

package util

import (
	"os/user"
	"strings"

	"github.com/apache/skywalking-cli/internal/logger"
)

// UserHomeDir returns the current user's home directory absolute path,
// which is usually represented as `~` in most shells
func UserHomeDir() string {
	if currentUser, err := user.Current(); err != nil {
		logger.Log.Warnln("Cannot obtain user home directory")
	} else {
		return currentUser.HomeDir
	}
	return ""
}

// ExpandFilePath expands the leading `~` to absolute path
func ExpandFilePath(path string) string {
	if strings.HasPrefix(path, "~") {
		return strings.Replace(path, "~", UserHomeDir(), 1)
	}
	return path
}
