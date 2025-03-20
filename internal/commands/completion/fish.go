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

var fishCommand = &cli.Command{
	Name:      "fish",
	Aliases:   []string{"f"},
	Usage:     "Output shell completion code for fish",
	ArgsUsage: "[parameters...]",
	Action: func(_ *cli.Context) error {
		fmt.Print(fishScript)
		return nil
	},
}

const fishScript = `
function __fish_swctl_complete
    set -l command (commandline -cp)
    set -l last_token (commandline -ct)

    if test "$last_token" = ""
        return
    end

    # Get completions using the auto-complete flag
    set -l completions (eval "$command --auto_complete" 2> /dev/null)

    for completion in $completions
        echo -e "$completion\t(swctl suggestion)"
    end
end

complete -c swctl -f -a "(__fish_swctl_complete)"
`
