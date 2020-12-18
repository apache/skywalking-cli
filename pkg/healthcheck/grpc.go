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

package healthcheck

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/apache/skywalking-cli/internal/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

const (
	Healthy           = 0
	ConnectionFailure = 1
	RPCFailure        = 2
	Unhealthy         = 3
)

func HealthCheck(addr string, enableTLS bool) int {
	ctx := context.Background()

	opts := []grpc.DialOption{
		grpc.WithUserAgent("swctl_health_probe"),
		grpc.WithBlock()}
	if enableTLS {
		// #nosec
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	dialCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(dialCtx, addr, opts...)
	if err != nil {
		if err == context.DeadlineExceeded {
			logger.Log.Printf("timeout: failed to connect service %q within 1 second", addr)
		} else {
			logger.Log.Printf("error: failed to connect service at %q: %+v", addr, err)
		}
		return ConnectionFailure
	}
	defer conn.Close()

	rpcCtx, rpcCancel := context.WithTimeout(ctx, time.Second)
	defer rpcCancel()
	resp, err := healthpb.NewHealthClient(conn).Check(rpcCtx, &healthpb.HealthCheckRequest{Service: ""})
	if err != nil {
		if stat, ok := status.FromError(err); ok && stat.Code() == codes.DeadlineExceeded {
			logger.Log.Printf("timeout: health request did not complete within 1 second")
		} else {
			logger.Log.Printf("error: health request failed: %+v", err)
		}
		return RPCFailure
	}

	if resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		logger.Log.Printf("OAP gRPC service is unhealthy %q", resp.GetStatus().String())
		return Unhealthy
	}
	return Healthy
}
