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
	"context"

	"github.com/apache/skywalking-cli/assets"
	"github.com/apache/skywalking-cli/pkg/graphql/client"

	"github.com/machinebox/graphql"
)

// Item is one entry of the UI menu the OAP served before 11.0.0. The query protocol
// retired getItems, so goapi no longer generates the type; the command stays for the
// older backends and carries the shape itself.
type Item struct {
	Title        string  `json:"title"`
	Icon         *string `json:"icon,omitempty"`
	Layer        string  `json:"layer"`
	Activate     bool    `json:"activate"`
	SubItems     []*Item `json:"subItems"`
	Description  *string `json:"description,omitempty"`
	DocumentLink *string `json:"documentLink,omitempty"`
	I18nKey      *string `json:"i18nKey,omitempty"`
}

func GetItems(ctx context.Context) ([]*Item, error) {
	var response map[string][]*Item

	request := graphql.NewRequest(assets.Read("graphqls/menu/GetItems.graphql"))

	err := client.ExecuteQuery(ctx, request, &response)
	return response["result"], err
}
