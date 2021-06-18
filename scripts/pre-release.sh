#!/usr/bin/env sh

# Prerequisites
# 1. update change log
# 2. clear milestone issues, and create a new one if needed
# 3. export VERSION=<the version to release>

set -ex

[ -z "$VERSION" ] && echo "VERSION is not set" && exit 1

if ls skywalking-cli > /dev/null 2>&1; then
  rm -rf skywalking-cli
fi

git clone --recurse-submodules git@github.com:apache/skywalking-cli && cd skywalking-cli
git tag -a "$VERSION" -m "Release Apache SkyWalking-CLI $VERSION"
git push --tags

make clean && make release

cd ..

if ls skywalking > /dev/null 2>&1; then
  rm -rf skywalking
fi

svn co https://dist.apache.org/repos/dist/dev/skywalking/
mkdir -p skywalking/cli/"$VERSION"
cp skywalking-cli/skywalking*.tgz skywalking/cli/"$VERSION"
cp skywalking-cli/skywalking*.tgz.asc skywalking/cli/"$VERSION"
cp skywalking-cli/skywalking-cli*.tgz.sha512 skywalking/cli/"$VERSION"

cd skywalking/cli && svn add "$VERSION" && svn commit -m "Draft Apache SkyWalking-CLI release $VERSION"
cd ../..

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
   - $(cat skywalking-cli*bin*tgz.sha512) skywalking-cli-$VERSION-bin.tgz
   - $(cat skywalking-cli*src*tgz.sha512) skywalking-cli-$VERSION-src.tgz

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