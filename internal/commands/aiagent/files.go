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
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/apache/skywalking-cli/internal/commands/interceptor"
	"github.com/apache/skywalking-cli/internal/flags"
	"github.com/apache/skywalking-cli/pkg/aiagent/files"
	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
)

var filesCommand = &cli.Command{
	Name:  "files",
	Usage: "List or export chosen stored files of a conversation, as the OAP stores them",
	UsageText: `Read chosen Session Data files of a conversation's session from the OAP's route
GET /ai-agent/conversations/{conversation}/v1/files, on the "--base-url" host, and list each
with its digest and size, or export them: "--export DIR" writes each file to its name under DIR.
A file is chosen by "--session" and "--seqs", its landed seq; the asz.view document's files
list gives both. Each file's bytes are checked against the digest the OAP names.

Examples:
1. Two files of a session, listed:
$ swctl ai-agent files --service-name "Claude Code" --instance-name laptop --conversation 7a3c882e-0dc0-46a0-b814-6613d24b7ac2 \
    --session 7a3c882e-0dc0-46a0-b814-6613d24b7ac2 --seqs 408,409

2. The same files, exported:
$ swctl ai-agent files --service-name "Claude Code" --instance-name laptop --conversation 7a3c882e-0dc0-46a0-b814-6613d24b7ac2 \
    --session 7a3c882e-0dc0-46a0-b814-6613d24b7ac2 --seqs 408,409 --export ./root`,
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
				Name:     "session",
				Usage:    "the `session` the files belong to",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "seqs",
				Usage:    "the landed `seqs` of the files, comma separated",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "export",
				Usage: "write each file to its name under this `directory`",
			},
		},
	),
	Before: interceptor.BeforeChain(
		interceptor.ParseService(true),
		interceptor.ParseInstance(true),
	),
	Action: func(ctx *cli.Context) error {
		seqs, err := numbers(ctx.String("seqs"))
		if err != nil {
			return err
		}
		if len(seqs) == 0 {
			return fmt.Errorf("--seqs needs at least one number")
		}
		var root *os.Root
		if exportDir := ctx.String("export"); exportDir != "" {
			if mkErr := os.MkdirAll(exportDir, 0o755); mkErr != nil {
				return mkErr
			}
			if root, err = os.OpenRoot(exportDir); err != nil {
				return err
			}
			defer root.Close()
		}

		out := List{Files: []files.File{}}
		var written []Exported
		err = files.Read(ctx.Context, ctx.String("conversation"), ctx.String("service-name"), ctx.String("instance-name"),
			ctx.String("session"), seqs, func(f files.File, content []byte) error {
				if root == nil {
					out.Files = append(out.Files, f)
					return nil
				}
				path, exportErr := export(root, f.ID, content)
				if exportErr != nil {
					return exportErr
				}
				written = append(written, Exported{ID: f.ID, Path: path, Bytes: len(content)})
				return nil
			})
		if err != nil {
			return err
		}
		if root == nil {
			return display.Display(ctx.Context, &displayable.Displayable{Data: out})
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: written})
	},
}

// List is what "files" prints without "--export": each stored file's naming line.
type List struct {
	Files []files.File `json:"files"`
}

// Exported is one file written by "--export": its name, its path and its size.
type Exported struct {
	ID    string `json:"file"`
	Path  string `json:"path"`
	Bytes int    `json:"bytes"`
}

// numbers reads a comma separated list of positive whole numbers.
func numbers(arg string) ([]int64, error) {
	var out []int64
	for _, part := range strings.Split(arg, ",") {
		if part = strings.TrimSpace(part); part == "" {
			continue
		}
		n, err := strconv.ParseInt(part, 10, 64)
		if err != nil || n <= 0 {
			return nil, fmt.Errorf("%q is not a positive whole number", part)
		}
		out = append(out, n)
	}
	return out, nil
}

// export writes one file to its name under root. A name is a relative path inside the Sessionizer's
// storage root. The root refuses a name that would leave it, through ".." or through a symbolic
// link, so a file the OAP names can only land inside the export directory.
func export(root *os.Root, name string, content []byte) (string, error) {
	rel := filepath.FromSlash(name)
	if !filepath.IsLocal(rel) {
		return "", fmt.Errorf("refusing to write %s outside %s", name, root.Name())
	}
	if dir := filepath.Dir(rel); dir != "." {
		if err := root.MkdirAll(dir, 0o755); err != nil {
			return "", err
		}
	}
	if err := root.WriteFile(rel, content, 0o644); err != nil { // #nosec G306 -- a landed file is readable by design
		return "", err
	}
	return filepath.Join(root.Name(), rel), nil
}
