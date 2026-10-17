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
# 1. update change log
# 2. clear milestone issues, and create a new one if needed
# 3. export VERSION=<the version to release>
# 4. build with the Go version CI uses, e.g. export GOTOOLCHAIN=go1.26.<latest patch>
# 5. run this script from an empty directory, it clones skywalking-cli and checks out the svn dev area there

set -ex

[ -z "$VERSION" ] && echo "VERSION is not set" && exit 1

# Both directories are recreated below, so ask before deleting anything that already exists here.
# Under `set -e`, a failing command that is not the last one of an `a && b` chain does not stop the script, so every
# step below that must succeed sits on its own line.
existing=""
for dir in skywalking-cli skywalking; do
  if [ -e "$dir" ] || [ -L "$dir" ]; then
    existing="$existing $dir"
  fi
done
if [ -n "$existing" ]; then
  printf 'Already in %s:%s. Remove and continue? [y/N] ' "$(pwd)" "$existing"
  read -r answer
  # $existing is unquoted on purpose: it only holds the fixed names listed in the prompt.
  # shellcheck disable=SC2086
  case "$answer" in
    [yY] | [yY][eE][sS]) rm -rf $existing ;;
    *) echo "Aborted, nothing was removed" && exit 1 ;;
  esac
fi

git clone --recurse-submodules git@github.com:apache/skywalking-cli
cd skywalking-cli
git tag -a "$VERSION" -m "Release Apache SkyWalking-CLI $VERSION"

make clean
make release

# Push the tag only after the packages are built and signed, so a failed build leaves no published tag.
git push origin "$VERSION"

cd ..

svn co https://dist.apache.org/repos/dist/dev/skywalking/
mkdir -p skywalking/cli/"$VERSION"
cp skywalking-cli/skywalking*.tgz skywalking/cli/"$VERSION"
cp skywalking-cli/skywalking*.tgz.asc skywalking/cli/"$VERSION"
cp skywalking-cli/skywalking-cli*.tgz.sha512 skywalking/cli/"$VERSION"

cd skywalking/cli
svn add "$VERSION"
svn commit -m "Draft Apache SkyWalking-CLI release $VERSION"
cd ../..

# The svn working copy is only needed for the upload above, the packages stay in ./skywalking-cli.
rm -rf skywalking

cd skywalking-cli

cat << EOF
=========================================================================
Subject: [VOTE] Release Apache SkyWalking CLI version $VERSION

Content:

Hi The SkyWalking Community:

This is a call for vote to release Apache SkyWalking CLI version $VERSION.

Release notes:

 * https://github.com/apache/skywalking-cli/blob/$VERSION/CHANGES.md

Release Candidate:

 * https://dist.apache.org/repos/dist/dev/skywalking/cli/$VERSION/

 * sha512 checksums
   - $(cat skywalking-cli*bin*tgz.sha512)
   - $(cat skywalking-cli*src*tgz.sha512)

Release Tag :

 * (Git Tag) $VERSION

Release CommitID :

 * https://github.com/apache/skywalking-cli/tree/$(git rev-list -n 1 "$VERSION")

Keys to verify the Release Candidate :

 * https://dist.apache.org/repos/dist/release/skywalking/KEYS

Guide to build the release from source :

 * https://github.com/apache/skywalking-cli/blob/$VERSION/README.md#install

Voting will start now and will remain open for at least 72 hours, all PMC members are required to give their votes.

[ ] +1 Release this package.
[ ] +0 No opinion.
[ ] -1 Do not release this package because....

Thanks.

[1] https://github.com/apache/skywalking-cli/blob/master/docs/How-to-release.md#vote-check

EOF
