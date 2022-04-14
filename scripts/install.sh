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

# Treat unset variables and parameters other than the special parameters ‘@’ or ‘*’ as an error.
set -u

# Exit the script with a message.
abort() {
  printf "%s\n" "$@"
  exit 1
}

# Check if there is a bash.
if [ -z "${BASH_VERSION:-}" ]; then
  abort "Bash is required to interpret this install script."
fi

# Check OS.
OS="$(uname)"
if [[ "$OS" != "Darwin" && "$OS" != "Linux" ]]; then
  abort "The install script is only supported on macOS and Linux."
fi

check_cmd() {
    if ! command -v "$@" &> /dev/null 
    then
        abort "You must install "$@" before running the install script."
    fi
}

# Check if the commands to be used exist.
for cmd in shasum curl tar awk; do
  check_cmd $cmd
done

# Convert the string to lower case.
OS=$(echo $OS | awk '{print tolower($0)}')

# Get the latest version of swctl.
VERSION=$(curl "https://raw.githubusercontent.com/apache/skywalking-website/master/data/releases.yml" | grep --after-context=7 "name: SkyWalking CLI" | grep "version" | grep -Eo "[0-9]+.[0-9]+.[0-9]+")
if [ "$VERSION" != "" ]; then
    echo "The latest version of swctl is $VERSION"
    
    # Download the binary package.
    curl -sSLO "https://archive.apache.org/dist/skywalking/cli/$VERSION/skywalking-cli-$VERSION-bin.tgz" > /dev/null
    if [ -f "skywalking-cli-$VERSION-bin.tgz" ]; then
        # Verify the integrity of the downloaded file.
        curl -sSLO "https://archive.apache.org/dist/skywalking/cli/$VERSION/skywalking-cli-$VERSION-bin.tgz.sha512" > /dev/null
        VERIFY=$(shasum -a512 -c "skywalking-cli-$VERSION-bin.tgz.sha512")
        if [ "${VERIFY#* }" = "OK" ]; then
            echo "The downloaded file is complete."
            tar -zxvf skywalking-cli-$VERSION-bin.tgz
            
            # Add swctl to the environment variable PATH.
            sudo cp skywalking-cli-$VERSION-bin/bin/swctl-$VERSION-$OS-amd64 /usr/local/bin/swctl
            
            # Delete unnecessary files.
            sudo rm -rf "./skywalking-cli-$VERSION-bin.tgz.sha512" "./skywalking-cli-$VERSION-bin.tgz" "./skywalking-cli-$VERSION-bin"
 
            echo "Type 'swctl --help' to get more information."
        else
            abort "The downloaded file is incomplete."
        fi
    else
        abort "Failed to download skywalking-cli-$VERSION-bin.tgz"
    fi
else
    echo $VERSION
    abort "Can't get the latest version. The install script may be invalid, try other install methods please."
fi
