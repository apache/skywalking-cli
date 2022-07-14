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
	"crypto/sha256"
	"fmt"

	"github.com/urfave/cli/v2"
)

const (
	processIDFlagName       = "process-id"
	processNameFlagName     = "process-name"
	destProcessNameFlagName = "dest-process-name"
)

// ParseProcess parses the process id or process name,
// and converts the present one to the missing one.
// See flags.InstanceFlags.
func ParseProcess(required bool) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		if err := ParseService(required)(ctx); err != nil {
			return err
		}
		if err := ParseInstance(required)(ctx); err != nil {
			return err
		}
		return parseProcess(required, processIDFlagName, processNameFlagName, instanceIDFlagName)(ctx)
	}
}

// ParseProcessRelation parses the source and destination service instance id or service instance name,
// and converts the present one to the missing one respectively.
// See flags.InstanceRelationFlags.
func ParseProcessRelation(required bool) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		if err := ParseService(required)(ctx); err != nil {
			return err
		}
		if err := ParseInstance(required)(ctx); err != nil {
			return err
		}
		if err := ParseProcess(required)(ctx); err != nil {
			return err
		}
		if ctx.String(destProcessNameFlagName) == "" && required {
			return fmt.Errorf(`flag "--%s" must given`, destProcessNameFlagName)
		}
		return nil
	}
}

func parseProcess(required bool, idFlagName, nameFlagName, instanceIDFlagName string) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		id := ctx.String(idFlagName)
		name := ctx.String(nameFlagName)
		instanceID := ctx.String(instanceIDFlagName)

		if id == "" && name == "" {
			if required {
				return fmt.Errorf(`either flags "--%s" or "--%s" must be given`, idFlagName, nameFlagName)
			}
			return nil
		}

		if name != "" {
			if instanceID == "" {
				return fmt.Errorf(`"--%s" is specified but its related service name or id is not given`, nameFlagName)
			}
			id = fmt.Sprintf("%x", sha256.New().Sum([]byte(fmt.Sprintf("%s_%s", instanceID, name))))
		}

		return ctx.Set(idFlagName, id)
	}
}
