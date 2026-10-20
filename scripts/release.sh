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
# 1. the vote passed after at least 72 hours, and the [RESULT][VOTE] email is sent
# 2. a PMC member runs this script, only PMC members can write to dist/release
# 3. run it from the skywalking-cli checkout of the release tag, with `gh` logged in, and `docker login` done
# Every step checks whether an earlier run already did it, so rerun the script to finish a skipped or failed step.

set -e

[ -z "$VERSION" ] && echo "VERSION is not set" && exit 1

DIST=https://dist.apache.org/repos/dist
RELEASE_DIR="$DIST/release/skywalking/cli"
IMAGE="${HUB:-docker.io/apache}/skywalking-cli:$VERSION"

# ask <question>: true when the answer is y or yes. POSIX `read` needs a variable name, dash fails without one.
ask() {
  printf '\n%s [y/N] ' "$1"
  read -r answer || answer=""
  case "$answer" in
    [yY] | [yY][eE][sS]) return 0 ;;
    *) return 1 ;;
  esac
}

# run <command>...: prints the command, then runs it.
run() {
  echo "+ $*"
  "$@"
}

# 1. Promote the packages from dist/dev to dist/release.
if svn info "$RELEASE_DIR/$VERSION" > /dev/null 2>&1; then
  dist_status="already done"
  echo "$RELEASE_DIR/$VERSION already exists, skipping svn mv"
else
  # Not in a condition: if the packages cannot be promoted, stop before the other steps.
  run svn mv "$DIST/dev/skywalking/cli/$VERSION" "$RELEASE_DIR/$VERSION" -m"Release Apache SkyWalking-CLI $VERSION"
  dist_status="done"
fi

# 2. The GitHub release, with the $VERSION section of CHANGES.md at the release tag as its notes.
tmp=${TMPDIR:-/tmp}
notes="${tmp%/}/skywalking-cli-$VERSION-release-notes.md"
git show "$VERSION:CHANGES.md" 2> /dev/null |
  awk -v v="$VERSION" '$0 == v {f = 1; next} f && /^[0-9]+\.[0-9]+\.[0-9]+$/ {exit} f && !/^-+$/ && (NF || p) {print; p = 1}' > "$notes"
if gh release view "$VERSION" -R apache/skywalking-cli > /dev/null 2>&1; then
  github_status="already done"
elif [ ! -s "$notes" ]; then
  github_status="skipped, no $VERSION section found in CHANGES.md at tag $VERSION"
elif ask "Publish the GitHub release $VERSION, with the $VERSION section of CHANGES.md as its notes?"; then
  if run gh release create "$VERSION" -R apache/skywalking-cli --title "$VERSION" --latest --notes-file "$notes"; then
    github_status="done"
  else
    github_status="failed"
  fi
else
  github_status="skipped"
fi

# 3. The Docker images. `make docker.push` builds the working tree, so it must be the release tag.
if docker buildx imagetools inspect "$IMAGE" > /dev/null 2>&1; then
  docker_status="already done"
elif [ "$(git describe --tags --exact-match HEAD 2> /dev/null)" != "$VERSION" ]; then
  docker_status="skipped, $(pwd) is not a checkout of tag $VERSION"
elif ask "Build and push $IMAGE and :latest for linux/386, linux/amd64 and linux/arm64?"; then
  if run make docker.push VERSION="$VERSION"; then
    docker_status="done"
  else
    docker_status="failed"
  fi
else
  docker_status="skipped"
fi

case "$github_status $docker_status" in
  *skipped* | *failed*) summary="Not all release steps for SkyWalking CLI $VERSION are done, rerun this script to finish them:" ;;
  *) summary="All release steps for SkyWalking CLI $VERSION are done:" ;;
esac

cat << EOF

=========================================================================
$summary
- packages in $RELEASE_DIR/$VERSION: $dist_status
- GitHub release: $github_status
- Docker images: $docker_status
=========================================================================

Remember to:
1. Update the events, download links and docs menu on the website, like
   https://github.com/apache/skywalking-website/pull/751.
2. Once https://skywalking.apache.org/downloads/ lists $VERSION, send the announcement below from your apache.org
   address to dev@skywalking.apache.org and announce@apache.org.
3. Add the release at https://reporter.apache.org/addrelease.html?skywalking

Subject: [ANNOUNCEMENT] Apache SkyWalking CLI $VERSION Released

Hi the SkyWalking Community

On behalf of the SkyWalking CLI Team, I’m glad to announce that SkyWalking CLI $VERSION is now released.

SkyWalking CLI: A command line interface for SkyWalking.

SkyWalking: APM (application performance monitor) tool for distributed systems, especially designed for microservices, cloud native and container-based (Docker, Kubernetes, Mesos) architectures.

Download Links: http://skywalking.apache.org/downloads/

Release Notes : https://github.com/apache/skywalking-cli/blob/$VERSION/CHANGES.md

Website: http://skywalking.apache.org/

SkyWalking CLI Resources:
- Issue: https://github.com/apache/skywalking/issues
- Mailing list: dev@skywalking.apache.org
- Documents: https://github.com/apache/skywalking-cli/blob/$VERSION/README.md

The Apache SkyWalking Team
=========================================================================
EOF

# 4. The older releases. archive.apache.org keeps them, so dist/release only needs the new one.
if ! listing=$(svn ls "$RELEASE_DIR/" 2> /dev/null); then
  echo "Could not list $RELEASE_DIR/ to find the older releases, check it with svn ls."
  exit 0
fi
old_versions=""
old_urls=""
for v in $(echo "$listing" | sed -n 's|/$||p' | grep -vxF "$VERSION" || true); do
  old_versions="$old_versions $v"
  old_urls="$old_urls $RELEASE_DIR/$v"
done
if [ -z "$old_versions" ]; then
  echo "No older release in dist/release to remove."
elif ask "Remove the older releases$old_versions from dist/release? The downloads page links to them until the website is updated."; then
  # $old_urls is unquoted on purpose: it holds one URL per older version, with no spaces inside.
  # shellcheck disable=SC2086
  run svn rm $old_urls -m "Remove old SkyWalking CLI releases"
else
  echo "Kept the older releases$old_versions in dist/release."
fi
