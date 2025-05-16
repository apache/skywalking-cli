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

package completion

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var zshCommand = &cli.Command{
	Name:      "zsh",
	Aliases:   []string{"z"},
	Usage:     "Output shell completion code for zsh",
	ArgsUsage: "[parameters...]",
	Action: func(_ *cli.Context) error {
		fmt.Print(zshScript)
		return nil
	},
}

const zshScript = `
#compdef swctl

_cli_zsh_autocomplete() {
    local -a completions
    local word

    word="${words[CURRENT]}"
    completions=("${(@f)$( ${words[1,CURRENT-1]} --auto_complete )}")

    _describe -t commands 'swctl commands' completions
}

compdef _cli_zsh_autocomplete swctl
`
