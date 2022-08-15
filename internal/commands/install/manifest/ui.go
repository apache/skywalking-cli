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

	operatorv1alpha1 "github.com/apache/skywalking-swck/operator/apis/operator/v1alpha1"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	controllerruntime "sigs.k8s.io/controller-runtime"

	"github.com/apache/skywalking-cli/assets"
)

var uiCmd = &cli.Command{
	Name:    "ui",
	Aliases: []string{"u"},
	Usage:   "Output the Kubernetes manifest for installing UI to stdout",
	UsageText: usage("ui", `Some examples of customized resource overlay files (ui-cr.yaml)

1. Set OAP server address to 'oap.test'', use an ingress to expose UI

	spec:
	  OAPServerAddress: oap.test
	  service:
		ingress:
		  host: ui.skywalking.test

2. Use a LoadBalancer to expose UI

	spec:
	  service:
		serviceSpec:
		  type: LoadBalancer
		  ports:
			- name: page
			  port: 80
			  targetPort: 8080`),
	Flags: flags,
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
		return render("ui", ctx, base, &operatorv1alpha1.UI{})
	},
}
