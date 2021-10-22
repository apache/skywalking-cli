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

package interceptor

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	api "skywalking.apache.org/repo/goapi/query"

	"github.com/apache/skywalking-cli/internal/logger"
	"github.com/apache/skywalking-cli/pkg/graphql/metadata"
)

type nodeType int

const (
	serviceIDFlagName       = "service-id"
	serviceNameFlagName     = "service-name"
	destServiceIDFlagName   = "dest-service-id"
	destServiceNameFlagName = "dest-service-name"

	normal  nodeType = iota
	browser nodeType = iota
)

// ParseService parses the service id or service name,
// and converts the present one to the missing one.
// See flags.ServiceFlags.
func ParseService(required bool) func(*cli.Context) error {
	return parseService(required, serviceIDFlagName, serviceNameFlagName, normal)
}

// ParseBrowserService parses the service id or service name,
// and converts the present one to the missing one.
// See flags.ServiceFlags.
func ParseBrowserService(required bool) func(*cli.Context) error {
	return parseService(required, serviceIDFlagName, serviceNameFlagName, browser)
}

// ParseServiceRelation parses the source and destination service id or service name,
// and converts the present one to the missing one respectively.
// See flags.ServiceRelationFlags.
func ParseServiceRelation(required bool) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		if err := ParseService(required)(ctx); err != nil {
			return err
		}
		return parseService(required, destServiceIDFlagName, destServiceNameFlagName, normal)(ctx)
	}
}

func parseService(required bool, idFlagName, nameFlagName string, nodeType nodeType) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		id := ctx.String(idFlagName)
		name := ctx.String(nameFlagName)

		if id == "" && name == "" {
			if required {
				return fmt.Errorf(`either flags "--%s" or "--%s" must be given`, idFlagName, nameFlagName)
			}
			return nil
		}

		if id != "" {
			parts := strings.Split(id, ".")
			if len(parts) != 2 {
				return fmt.Errorf("invalid service id, cannot be splitted into 2 parts. %v", id)
			}
			s, err := base64.StdEncoding.DecodeString(parts[0])
			if err != nil {
				return err
			}
			name = string(s)
		} else if name != "" {
			var service api.Service
			var err error
			switch nodeType {
			case normal:
				service, err = metadata.SearchService(ctx, name)
			case browser:
				service, err = metadata.SearchBrowserService(ctx, name)
			}
			if err != nil {
				return err
			}
			id = service.ID
		}

		if err := ctx.Set(idFlagName, id); err != nil {
			return err
		}
		if err := ctx.Set(nameFlagName, name); err != nil {
			return err
		}

		logger.Log.Debugf("%v=%v, %v=%v", idFlagName, id, nameFlagName, name)

		return nil
	}
}
