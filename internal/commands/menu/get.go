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

package menu

import (
	"fmt"
	"strings"

	"github.com/apache/skywalking-cli/pkg/display"
	"github.com/apache/skywalking-cli/pkg/display/displayable"
	"github.com/apache/skywalking-cli/pkg/graphql/menu"

	"github.com/urfave/cli/v2"
)

var Get = &cli.Command{
	Name:  "get",
	Usage: "Get the UI menu items",
	UsageText: `Get the UI menu items.

Examples:
1. Get the UI menu items:
$swctl menu get`,
	Action: func(ctx *cli.Context) error {
		menuItems, err := menu.GetItems(ctx.Context)
		if err != nil {
			if isMenuUnsupported(err) {
				return fmt.Errorf("this OAP version no longer serves the UI menu: the `getMenuItems` " +
					"GraphQL query was retired in SkyWalking 11.0.0, where the OAP backend stopped storing " +
					"and serving the sidebar menu. The menu is now owned client-side by Horizon UI " +
					"(https://github.com/apache/skywalking-horizon-ui). Upgrade/downgrade swctl to match your " +
					"OAP version, or stop using `swctl menu get` against 11.0.0+ backends")
			}
			return err
		}
		return display.Display(ctx.Context, &displayable.Displayable{Data: menuItems})
	},
}

// isMenuUnsupported reports whether err is the GraphQL schema-validation error raised
// by an OAP backend that no longer defines the retired `getMenuItems` query (11.0.0+),
// as opposed to a transport or other runtime error. graphql-java phrases this as a
// "FieldUndefined" validation error; we match defensively on the field name plus a
// validation marker so wording changes do not break detection.
func isMenuUnsupported(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	if !strings.Contains(msg, "getMenuItems") {
		return false
	}
	lower := strings.ToLower(msg)
	return strings.Contains(lower, "fieldundefined") ||
		strings.Contains(lower, "undefined") ||
		strings.Contains(lower, "cannot query field") ||
		strings.Contains(lower, "validation error")
}
