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

// Package status wraps the OAP admin-server `status` feature module: cluster
// membership, effective configuration / TTL, and alarm runtime status. These
// endpoints were served on the public REST port before OAP 11.0.0 and now live on
// the admin host (default 17128); only /status/config/ttl is also mirrored on 12800.
package status

import (
	"context"
	"net/url"

	"github.com/apache/skywalking-cli/pkg/admin/client"
)

// ClusterNode is a single OAP node as seen by the cluster coordinator.
type ClusterNode struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	Self bool   `json:"self"`
}

// ClusterNodes is the response of GET /status/cluster/nodes.
type ClusterNodes struct {
	Nodes []ClusterNode `json:"nodes"`
}

// ClusterAlarmStatus is the per-node envelope returned by every /status/alarm/* call.
// OAP fans the query out across cluster members; non-evaluating nodes return a stub
// status and errorMsg is omitted on success.
type ClusterAlarmStatus struct {
	OapInstances []OapInstanceStatus `json:"oapInstances"`
}

// OapInstanceStatus is one cluster member's slice of an alarm status response.
type OapInstanceStatus struct {
	Address  string `json:"address"`
	ErrorMsg string `json:"errorMsg,omitempty"`
	Status   any    `json:"status"`
}

// ClusterNodesQuery returns the OAP cluster peer list (GET /status/cluster/nodes).
// The self flag is normalized from either `self` or `isSelf` on the wire.
func ClusterNodesQuery(ctx context.Context) (*ClusterNodes, error) {
	var raw struct {
		Nodes []struct {
			Host   string `json:"host"`
			Port   int    `json:"port"`
			Self   bool   `json:"self"`
			IsSelf bool   `json:"isSelf"`
		} `json:"nodes"`
	}
	if err := client.GetJSON(ctx, "/status/cluster/nodes", nil, &raw); err != nil {
		return nil, err
	}
	out := &ClusterNodes{Nodes: make([]ClusterNode, 0, len(raw.Nodes))}
	for _, n := range raw.Nodes {
		out.Nodes = append(out.Nodes, ClusterNode{Host: n.Host, Port: n.Port, Self: n.Self || n.IsSelf})
	}
	return out, nil
}

// ConfigTTL returns the effective TTL configuration (GET /status/config/ttl).
func ConfigTTL(ctx context.Context) (any, error) {
	var out any
	err := client.GetJSON(ctx, "/status/config/ttl", nil, &out)
	return out, err
}

// ConfigDump returns the effective, secrets-redacted configuration as a flat map of
// `<module>.<provider>.<property>` keys (GET /debugging/config/dump).
func ConfigDump(ctx context.Context) (any, error) {
	var out any
	err := client.GetJSON(ctx, "/debugging/config/dump", nil, &out)
	return out, err
}

// AlarmRules returns the loaded alarm rules per OAP node (GET /status/alarm/rules).
func AlarmRules(ctx context.Context) (*ClusterAlarmStatus, error) {
	var out ClusterAlarmStatus
	err := client.GetJSON(ctx, "/status/alarm/rules", nil, &out)
	return &out, err
}

// AlarmRule returns the definition + running state of a single alarm rule
// (GET /status/alarm/{ruleId}).
func AlarmRule(ctx context.Context, ruleID string) (*ClusterAlarmStatus, error) {
	var out ClusterAlarmStatus
	err := client.GetJSON(ctx, "/status/alarm/"+url.PathEscape(ruleID), nil, &out)
	return &out, err
}

// AlarmRuleEntity returns the per-entity alarm window/evaluation state for a rule
// (GET /status/alarm/{ruleId}/{entityName}).
func AlarmRuleEntity(ctx context.Context, ruleID, entityName string) (*ClusterAlarmStatus, error) {
	var out ClusterAlarmStatus
	path := "/status/alarm/" + url.PathEscape(ruleID) + "/" + url.PathEscape(entityName)
	err := client.GetJSON(ctx, path, nil, &out)
	return &out, err
}
