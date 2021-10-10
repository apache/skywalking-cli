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
)

const (
	instanceIDFlagName       = "instance-id"
	instanceNameFlagName     = "instance-name"
	destInstanceIDFlagName   = "dest-instance-id"
	destInstanceNameFlagName = "dest-instance-name"
)

// ParseInstance parses the service instance id or service instance name,
// and converts the present one to the missing one.
// See flags.InstanceFlags.
func ParseInstance(required bool) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		if err := ParseService(required)(ctx); err != nil {
			return err
		}
		return parseInstance(required, instanceIDFlagName, instanceNameFlagName, serviceIDFlagName)(ctx)
	}
}

// ParseInstanceRelation parses the source and destination service instance id or service instance name,
// and converts the present one to the missing one respectively.
// See flags.InstanceRelationFlags.
func ParseInstanceRelation(required bool) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		if err := ParseService(required)(ctx); err != nil {
			return err
		}
		if err := ParseInstance(required)(ctx); err != nil {
			return err
		}
		return parseInstance(required, destInstanceIDFlagName, destInstanceNameFlagName, destServiceIDFlagName)(ctx)
	}
}

func parseInstance(required bool, idFlagName, nameFlagName, serviceIDFlagName string) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		id := ctx.String(idFlagName)
		name := ctx.String(nameFlagName)
		serviceID := ctx.String(serviceIDFlagName)

		if id == "" && name == "" {
			if required {
				return fmt.Errorf(`either flags "--%s" or "--%s" must be given`, idFlagName, nameFlagName)
			}
			return nil
		}

		if id != "" {
			parts := strings.Split(id, "_")
			if len(parts) != 2 {
				return fmt.Errorf("invalid instance id, cannot be splitted into 2 parts. %v", id)
			}
			s, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				return err
			}
			name = string(s)
		} else if name != "" {
			if serviceID == "" {
				return fmt.Errorf(`"--%s" is specified but its related service name or id is not given`, nameFlagName)
			}
			id = serviceID + "_" + b64enc(name)
		}

		if err := ctx.Set(idFlagName, id); err != nil {
			return err
		}
		return ctx.Set(nameFlagName, name)
	}
}
