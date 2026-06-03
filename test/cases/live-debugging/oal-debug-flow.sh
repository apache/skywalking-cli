#!/usr/bin/env bash
#
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

# Drives swctl's DSL live-debugger commands end-to-end against the OAL pipeline:
#   admin dsl-debug status  → injection enabled
#   admin dsl-debug session start (catalog=oal, name=core.oal, ruleName=<metric>)
#   admin dsl-debug session get   → poll until the live trace flow is captured
#   admin dsl-debug session stop
#
# The consumer→provider demo traffic fires the OAL source events that the session
# captures. Mirrors apache/skywalking test/e2e-v2/cases/dsl-debugging/oal but exercises
# swctl instead of curl. All diagnostics go to stderr; the only stdout is the final
# "ok" the e2e verify matches.

set -euo pipefail

log() { echo "[oal-debug] $*" >&2; }
fail() { log "FAIL: $*"; exit 1; }

OAP_HOST="${OAP_HOST:-127.0.0.1}"
OAP_ADMIN_PORT="${OAP_ADMIN_PORT:-17128}"
SETTLE_SECONDS="${SETTLE_SECONDS:-300}"

# OAL debug is per-metric. service_cpm is a shipped core.oal rule on the Service source
# that every instrumented service fires on each inbound request, so the demo traffic
# drives it reliably.
CATALOG="oal"
NAME="core.oal"
METRIC="${METRIC:-service_cpm}"
# OAL source class the metric is derived from — surfaces as the input sample's
# payload.type (service_cpm = from(Service.*).cpm() → Service).
SOURCE_TYPE="${SOURCE_TYPE:-Service}"
CLIENT_ID="e2e-oal-debug-$$"

ADMIN=(--display=json "--admin-url=http://${OAP_HOST}:${OAP_ADMIN_PORT}")

# --- Phase 0: module status -----------------------------------------------------------
log "=== Phase 0: dsl-debug status ==="
swctl "${ADMIN[@]}" admin dsl-debug status | yq e '.injectionEnabled' - | grep -qx true \
    || fail "injectionEnabled is not true"

# --- Phase 1: start session -----------------------------------------------------------
log "=== Phase 1: start OAL session (catalog=${CATALOG}, name=${NAME}, ruleName=${METRIC}) ==="
start_out="$(swctl "${ADMIN[@]}" admin dsl-debug session start \
    --catalog "${CATALOG}" --name "${NAME}" --rule-name "${METRIC}" --client-id "${CLIENT_ID}")"
log "  start → ${start_out}"
SESSION_ID="$(echo "${start_out}" | yq e '.sessionId // ""' -)"
[ -n "${SESSION_ID}" ] || fail "session start did not return a sessionId"
log "✓ session started: ${SESSION_ID}"

# --- Phase 2: poll until the live OAL pipeline is captured -----------------------------
log "=== Phase 2: poll for captured records (budget ${SETTLE_SECONDS}s) ==="
deadline=$(( $(date +%s) + SETTLE_SECONDS ))
records=0
body=""
while [ "$(date +%s)" -lt "${deadline}" ]; do
    body="$(swctl "${ADMIN[@]}" admin dsl-debug session get "${SESSION_ID}")"
    records="$(echo "${body}" | yq e '[.nodes[].records[]] | length' -)"
    [ "${records}" -gt 0 ] && break
    sleep 5
done
[ "${records}" -gt 0 ] || fail "no records captured within ${SETTLE_SECONDS}s"
log "✓ captured ${records} record(s)"

# --- Phase 3: verify the capture is EXACTLY the bound metric --------------------------
log "=== Phase 3: assert the captured pipeline is exactly ${METRIC} ==="

samples="$(echo "${body}" | yq e '[.nodes[].records[].samples[]] | length' -)"
[ "${samples}" -gt 0 ] || fail "captured records carry no samples"

# 3a. Per-metric gate isolation: every captured record is bound to ${METRIC}; no sibling
#     rule on the same OAL dispatcher leaked into this per-metric session.
foreign="$(echo "${body}" | yq e "[.nodes[].records[] | select(.rule.ruleName != \"${METRIC}\")] | length" -)"
[ "${foreign}" = "0" ] || fail "${foreign} record(s) bound to a rule other than ${METRIC}"

# 3b. Each record carries the verbatim core.oal source of ${METRIC}.
dsl_hits="$(echo "${body}" | yq e "[.nodes[].records[] | select(.dsl | contains(\"${METRIC}\"))] | length" -)"
[ "${dsl_hits}" -gt 0 ] || fail "no record's .dsl carries the ${METRIC} OAL source"

# 3c. Source stage: an input sample drawn from the ${SOURCE_TYPE} source.
src="$(echo "${body}" | yq e "[.nodes[].records[].samples[] | select(.type == \"input\" and .payload.type == \"${SOURCE_TYPE}\")] | length" -)"
[ "${src}" -gt 0 ] || fail "no input sample from the ${SOURCE_TYPE} source"

# 3d. Aggregation stage: the verbatim cpm() function from the rule.
agg="$(echo "${body}" | yq e '[.nodes[].records[].samples[] | select(.type == "aggregation" and .sourceText == "cpm()")] | length' -)"
[ "${agg}" -gt 0 ] || fail "no cpm() aggregation sample for ${METRIC}"

# 3e. Output stage: the materialised ${METRIC} metric.
out="$(echo "${body}" | yq e "[.nodes[].records[].samples[] | select(.type == \"output\" and .sourceText == \"${METRIC}\")] | length" -)"
[ "${out}" -gt 0 ] || fail "no output sample for metric ${METRIC}"

# 3f. Gate isolation on output: no OTHER metric's output leaked through this session.
leak="$(echo "${body}" | yq e "[.nodes[].records[].samples[] | select(.type == \"output\" and .sourceText != \"${METRIC}\")] | length" -)"
[ "${leak}" = "0" ] || fail "${leak} output sample(s) for a metric other than ${METRIC} leaked into the session"

log "✓ capture is exactly ${METRIC}: ${samples} samples — ${SOURCE_TYPE} source → cpm() → ${METRIC} output, no foreign rules"

# --- Phase 4: stop session ------------------------------------------------------------
log "=== Phase 4: stop session ==="
swctl "${ADMIN[@]}" admin dsl-debug session stop "${SESSION_ID}" | yq e '.localStopped' - | grep -qx true \
    || fail "localStopped is not true"

log "=== OAL LIVE-DEBUG FLOW PASSED (${records} records, ${samples} samples) ==="
echo ok
