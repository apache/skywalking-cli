#!/usr/bin/env sh

# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Prerequisites
# 1. vote passed
# 2. update GitHub release and the website

set -ex

[ -z "$VERSION" ] && echo "VERSION is not set" && exit 1

 if ls skywalking > /dev/null 2>&1; then
   rm -rf skywalking
 fi

 svn mv https://dist.apache.org/repos/dist/dev/skywalking/cli/"$VERSION" https://dist.apache.org/repos/dist/release/skywalking/cli/"$VERSION" -m"Release Apache SkyWalking-CLI $VERSION" || true

 echo "Make sure you have released the GitHub release, updated the events, download links and menus?" && read -r

cat << EOF
=========================================================================
Hi the SkyWalking Community

On behalf of the SkyWalking CLI Team, I’m glad to announce that SkyWalking CLI $VERSION is now released.

SkyWalking CLI: A command line interface for SkyWalking.

SkyWalking: APM (application performance monitor) tool for distributed systems, especially designed for microservices, cloud native and container-based (Docker, Kubernetes, Mesos) architectures.

Download Links: http://skywalking.apache.org/downloads/

Release Notes : https://github.com/apache/skywalking-cli/blob/$VERSION/CHANGES.md

Website: http://skywalking.apache.org/

SkyWalking CLI Resources:
- Issue: https://github.com/apache/skywalking/issues
- Mailing list: dev@skywalkiing.apache.org
- Documents: https://github.com/apache/skywalking-cli/blob/$VERSION/README.md

The Apache SkyWalking Team

EOF
