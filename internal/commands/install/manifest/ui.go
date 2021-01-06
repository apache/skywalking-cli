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
	"fmt"

	operatorv1alpha1 "github.com/apache/skywalking-swck/apis/operator/v1alpha1"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	controllerruntime "sigs.k8s.io/controller-runtime"

	"github.com/apache/skywalking-cli/assets"
)

var uiCmd = cli.Command{
	Name:      "ui",
	ShortName: "u",
	Usage:     "Output the Kubernetes manifest for installing UI to stdout",
	UsageText: usage("ui"),
	Flags:     flags,
	Action: func(ctx *cli.Context) error {
		base := &operatorv1alpha1.UI{
			TypeMeta: controllerruntime.TypeMeta{
				Kind: "UI",
			},
			ObjectMeta: controllerruntime.ObjectMeta{
				Name:      ctx.String("name"),
				Namespace: ctx.String("namespace"),
			},
		}
		err := yaml.Unmarshal([]byte(assets.Read("cr/ui.yaml")), base)
		if err != nil {
			return fmt.Errorf("failed to convert yaml to UI: %v", err)
		}
		base.Default()
		if err := base.ValidateCreate(); err != nil {
			return fmt.Errorf("failed to validate UI: %v", err)
		}
		if err := render("ui", ctx, base, &operatorv1alpha1.UI{}); err != nil {
			return err
		}
		return nil
	},
}
