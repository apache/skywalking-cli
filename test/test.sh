# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -ex

if [ "$(uname)" == "Darwin" ]; then
  os="darwin"
else
  os="linux"
fi

swctl="bin/swctl-latest-${os}-amd64"

# Check whether OAP server is healthy.
if ! ${swctl} ch --grpc=false > /dev/null 2>&1; then
  echo "OAP server is not healthy"
  exit 1
fi

if ! ${swctl} metrics ls > /dev/null 2>&1; then
  exit 1
fi

if ! ${swctl} service ls > /dev/null 2>&1; then
  exit 1
fi

if ! ${swctl} metrics linear --name="database_access_resp_time" --service="test" > /dev/null 2>&1; then
  exit 1
fi

if ! ${swctl} metrics multiple-linear --name="all_percentile" > /dev/null 2>&1; then
  exit 1
fi

# Test `metrics thermodynamic`
if ! ${swctl} metrics hp --name="all_heatmap" > /dev/null 2>&1; then
  exit 1
fi

if ! ${swctl} metrics single --name="service_resp_time" --service="test" > /dev/null 2>&1; then
  exit 1
fi

if ! ${swctl} metrics top 5 --name="service_resp_time" > /dev/null 2>&1; then
  exit 1
fi

if ! ${swctl} trace ls > /dev/null 2>&1; then
  exit 1
fi

# Test `dashboard global`
if ! ${swctl} db g > /dev/null 2>&1; then
  exit 1
fi
