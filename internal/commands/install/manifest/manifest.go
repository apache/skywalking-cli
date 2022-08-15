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

package manifest

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/apache/skywalking-swck/operator/pkg/kubernetes"
	repo "github.com/apache/skywalking-swck/operator/pkg/operator/manifests"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime"
)

var Command = &cli.Command{
	Name:    "manifest",
	Aliases: []string{"mf"},
	Usage:   "Output the Kubernetes manifest for installing OAP server and UI to stdout",
	Subcommands: []*cli.Command{
		oapCmd,
		uiCmd,
	},
}

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:     "name",
		Usage:    "The name of prefix of generated resources",
		Required: false,
		Value:    "skywalking",
	},
	&cli.StringFlag{
		Name:     "namespace",
		Usage:    "The namespace where resource will be deployed",
		Required: false,
		Value:    "skywalking-system",
	},
	&cli.StringFlag{
		Name:     "f",
		Usage:    "The custom resource file describing custom resources defined by swck",
		Required: false,
	},
}

func usage(command, exampleOverlays string) string {
	return fmt.Sprintf(`
Examples:

%s

1. Output manifest with default custom resource
$ swctl install manifest %s

2. Load overlay custom resource from flag
$ swctl install manifest %s -f %s-cr.yaml

3. Load overlay custom resource from stdin
$ cat %s-cr.yaml | swctl install manifest %s -f=-

4. Apply directly to Kubernetes
$ swctl install manifest %s -f %s-cr.yaml | kubectl apply -f-
`, exampleOverlays, command, command, command, command, command, command, command)
}

func loadOverlay(file string, in io.Reader, out interface{}) error {
	if file == "" {
		return nil
	}
	if file == "-" {
		scanner := bufio.NewScanner(in)
		ll := make([]string, 0)
		for scanner.Scan() {
			ll = append(ll, scanner.Text())
		}
		if len(ll) > 0 {
			if err := yaml.Unmarshal([]byte(strings.Join(ll, "\n")), out); err != nil {
				return err
			}
		}
	} else {
		b, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(b, out); err != nil {
			return err
		}
	}
	return nil
}

func render(repoName string, ctx *cli.Context, base, overlay runtime.Object) error {
	if ctx.String("f") != "" {
		if err := applyOverlay(ctx, base, overlay); err != nil {
			return err
		}
	}
	r := repo.NewRepo(repoName)
	ff, err := r.GetFilesRecursive("templates")
	if err != nil {
		return fmt.Errorf("failed to load resource templates: %v", err)
	}

	for _, f := range ff {
		manifests, err := r.ReadFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read template file content: %v", err)
			continue
		}
		bb, err := kubernetes.GenerateManifests(string(manifests), base, nil)
		if err != nil && err != kubernetes.ErrNothingLoaded {
			fmt.Fprintf(os.Stderr, "failed to generate manifest: %v", err)
			continue
		}
		fmt.Println(string(bb))
		fmt.Println("---")
	}
	return nil
}

func applyOverlay(ctx *cli.Context, base, overlay runtime.Object) error {
	if err := loadOverlay(ctx.String("f"), os.Stdin, overlay); err != nil {
		return fmt.Errorf("falied to load overlay from flag -f: %v", err)
	}
	if err := kubernetes.ApplyOverlay(base, overlay); err != nil {
		return fmt.Errorf("failed to apply overlay to base: %v", err)
	}
	return nil
}
