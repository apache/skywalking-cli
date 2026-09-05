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

package aiagent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	api "skywalking.apache.org/repo/goapi/query"

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/aiagent"
)

var filesCommand = &cli.Command{
	Name:  "files",
	Usage: "List or export the raw files of a conversation, as the OAP stores them",
	UsageText: `List every landed file and round of a conversation with its digest and size, or
export them: "--export DIR" reads each body and writes it to its id path under DIR,
which gives a storage root that "asz verify" and "asz view" read like the original.

Examples:
1. The files of a conversation:
$ swctl ai-agent files --service-name "Claude Code" --conversation 7a3c882e-0dc0-46a0-b814-6613d24b7ac2

2. Export them all:
$ swctl ai-agent files --service-name "Claude Code" --conversation 7a3c882e-0dc0-46a0-b814-6613d24b7ac2 --export ./root

3. Export two named files:
$ swctl ai-agent files --service-name "Claude Code" --conversation 7a3c882e-0dc0-46a0-b814-6613d24b7ac2 \
    --files 7a3c882e-0dc0-46a0-b814-6613d24b7ac2/streams/main/transcript-20260904T152815.774957000Z-000408.sd \
    --export ./root`,
	Flags: flags.Flags(
		flags.ServiceFlags,
		flags.InstanceFlags,
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "conversation",
				Usage:    "`id` of the conversation",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "files",
				Usage: "only these file `ids`, comma separated; without it, every file of the conversation",
			},
			&cli.StringFlag{
				Name:  "export",
				Usage: "write each file's body to its id path under this `directory`",
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseService(true),
		interceptor.ParseInstance(false),
	),
	Action: func(ctx *cli.Context) error {
		condition := &api.ConversationCondition{
			Service:      &api.ServiceCondition{ServiceName: ctx.String("service-name")},
			Conversation: ctx.String("conversation"),
			Instance:     instanceCondition(ctx),
		}
		var files []string
		if arg := strings.TrimSpace(ctx.String("files")); arg != "" {
			files = strings.Split(arg, ",")
		}
		exportDir := ctx.String("export")

		raw, err := aiagent.RawFiles(ctx.Context, condition, files, exportDir != "")
		if err != nil {
			return err
		}
		if raw.ErrorReason != nil && *raw.ErrorReason != "" {
			return fmt.Errorf("%s", *raw.ErrorReason)
		}
		if exportDir == "" {
			return display.Display(ctx.Context, &displayable.Displayable{Data: raw, Condition: condition})
		}

		written, err := export(exportDir, raw.Files)
		if err != nil {
			return err
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: written, Condition: condition})
	},
}

// Exported is one file written by "--export": its id path and size, the body left out.
type Exported struct {
	ID    string `json:"id"`
	Path  string `json:"path"`
	Bytes int    `json:"bytes"`
}

// export writes each body to its id path under dir. An id is a relative path inside the
// Sessionizer's storage root; one that would leave dir is refused.
func export(dir string, files []*api.ConversationRawFile) ([]Exported, error) {
	root, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	out := make([]Exported, 0, len(files))
	for _, f := range files {
		if f.Body == nil {
			return nil, fmt.Errorf("the OAP returned no body for %s", f.ID)
		}
		path := filepath.Join(root, filepath.FromSlash(f.ID))
		if !strings.HasPrefix(path, root+string(filepath.Separator)) {
			return nil, fmt.Errorf("refusing to write %s outside %s", f.ID, root)
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(path, []byte(*f.Body), 0o644); err != nil { // #nosec G306 -- a landed file is readable by design
			return nil, err
		}
		out = append(out, Exported{ID: f.ID, Path: path, Bytes: len(*f.Body)})
	}
	return out, nil
}
