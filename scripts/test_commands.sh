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

swctl="bin/swctl-${VERSION}-${os}-amd64 --base-url=http://localhost:12800/graphql"

retries=1
max_retries=10
# Check whether OAP server is healthy.
while ! ${swctl} ch > /dev/null 2>&1; do
  if [[ $retries -ge $max_retries ]]; then
    echo "OAP server is not healthy after $retries retires, will exit now"
    exit 1
  fi
  echo "OAP server is not healthy, retrying [$retries/$max_retries] ..."
  sleep 3
  retries=$(($retries+1))
done;

${swctl} --display=json metrics ls > /dev/null 2>&1

${swctl} --display=json service ls > /dev/null 2>&1

${swctl} --display=json endpoint ls --service-id="test" > /dev/null 2>&1

SERVICE_SCOPE_METRICS=(
  service_resp_time
  service_sla
  service_cpm
  service_apdex
)

for metrics in "${SERVICE_SCOPE_METRICS[@]}"; do
  ${swctl} --display=json metrics linear --name="$metrics" --service="test" > /dev/null 2>&1

  ${swctl} --display=json metrics single --name="$metrics" --service="test" > /dev/null 2>&1

  ${swctl} --display=json metrics top 3 --name="$metrics" > /dev/null 2>&1
done

${swctl} --display=json metrics multiple-linear --name="all_percentile" > /dev/null 2>&1

# Test `metrics thermodynamic`
${swctl} --display=json metrics hp --name="all_heatmap" >/dev/null 2>&1

${swctl} --display=json trace ls >/dev/null 2>&1

# Test `dashboard global`
${swctl} --display=json db g >/dev/null 2>&1

# Test `dependency`
${swctl} --display=json dependency service "test" > /dev/null 2>&1

${swctl} --display=json dependency endpoint "test" > /dev/null 2>&1