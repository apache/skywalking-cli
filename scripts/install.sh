#!/bin/bash

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

# Installs the latest Apache SkyWalking CLI release as /usr/local/bin/swctl on macOS and Linux.

# Treat unset variables and parameters other than the special parameters ‘@’ or ‘*’ as an error.
set -u

# Exit the script with a message.
abort() {
  printf "%s\n" "$@" >&2
  exit 1
}

# Check if there is a bash.
if [ -z "${BASH_VERSION:-}" ]; then
  abort "Bash is required to interpret this install script."
fi

# Check OS.
case "$(uname)" in
  Darwin) OS=darwin ;;
  Linux) OS=linux ;;
  *) abort "The install script is only supported on macOS and Linux." ;;
esac

# Pick the binary for this machine: the release has darwin-amd64, darwin-arm64, linux-386, linux-amd64 and linux-arm64.
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64 | amd64) ARCH=amd64 ;;
  arm64 | aarch64) ARCH=arm64 ;;
  i386 | i686) ARCH=386 ;;
  *) abort "There is no swctl binary for the $ARCH architecture, try other install methods please." ;;
esac
# A shell running under Rosetta 2 reports x86_64 on Apple silicon; install the native binary there.
if [ "$OS" = darwin ] && [ "$(sysctl -n hw.optional.arm64 2> /dev/null)" = 1 ]; then
  ARCH=arm64
fi

# Check if the commands to be used exist.
for cmd in curl tar; do
  command -v "$cmd" > /dev/null || abort "You must install $cmd before running the install script."
done
if command -v sha512sum > /dev/null; then
  verify_sha512() { sha512sum -c "$1"; }
elif command -v shasum > /dev/null; then
  verify_sha512() { shasum -a 512 -c "$1"; }
else
  abort "You must install sha512sum or shasum before running the install script."
fi

# Get the latest version of swctl: the latest GitHub release redirects to .../releases/tag/<version>.
LATEST_URL=$(curl -fsSLI -o /dev/null -w '%{url_effective}' https://github.com/apache/skywalking-cli/releases/latest) || LATEST_URL=""
VERSION=${LATEST_URL##*/}
VERSION_PATTERN='^[0-9]+\.[0-9]+\.[0-9]+$'
if ! [[ "$VERSION" =~ $VERSION_PATTERN ]]; then
  abort "Can't get the latest version. The install script may be invalid, try other install methods please."
fi
echo "The latest version of swctl is $VERSION"

TMP_DIR=$(mktemp -d) || abort "Failed to create a temporary directory."
trap 'rm -rf "$TMP_DIR"' EXIT

# Download the binary package and its checksum. downloads.apache.org has the current releases;
# archive.apache.org has all of them, but picks up a new one later.
PACKAGE="skywalking-cli-$VERSION-bin"
DOWNLOADED=false
for HOST in https://downloads.apache.org/skywalking/cli https://archive.apache.org/dist/skywalking/cli; do
  echo "Downloading $HOST/$VERSION/$PACKAGE.tgz"
  if curl -fL --progress-bar -o "$TMP_DIR/$PACKAGE.tgz" "$HOST/$VERSION/$PACKAGE.tgz" &&
    curl -fsSL -o "$TMP_DIR/$PACKAGE.tgz.sha512" "$HOST/$VERSION/$PACKAGE.tgz.sha512"; then
    DOWNLOADED=true
    break
  fi
done
[ "$DOWNLOADED" = true ] || abort "Failed to download $PACKAGE.tgz"

# Verify the integrity of the downloaded file.
(cd "$TMP_DIR" && verify_sha512 "$PACKAGE.tgz.sha512" > /dev/null) || abort "The downloaded file is incomplete."
echo "The downloaded file is complete."

BINARY="$PACKAGE/bin/swctl-$VERSION-$OS-$ARCH"
tar -xzf "$TMP_DIR/$PACKAGE.tgz" -C "$TMP_DIR" "$BINARY" || abort "Failed to extract $BINARY from $PACKAGE.tgz"

# Add swctl to the environment variable PATH, using sudo only when /usr/local/bin is not writable.
INSTALL_DIR=/usr/local/bin
SUDO=""
if [ "$(id -u)" != 0 ] && ! { [ -d "$INSTALL_DIR" ] && [ -w "$INSTALL_DIR" ]; }; then
  command -v sudo > /dev/null || abort "Writing to $INSTALL_DIR needs root, run the install script as root please."
  SUDO=sudo
fi
if ! { $SUDO mkdir -p "$INSTALL_DIR" && $SUDO install -m 0755 "$TMP_DIR/$BINARY" "$INSTALL_DIR/swctl"; }; then
  abort "Failed to install swctl to $INSTALL_DIR"
fi

echo "Installed swctl $VERSION ($OS-$ARCH) to $INSTALL_DIR/swctl. Type 'swctl --help' to get more information."
