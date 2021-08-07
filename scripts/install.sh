#!/usr/bin/env bash

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

# Get the latest version number.
VERSION=$(curl "https://raw.githubusercontent.com/apache/skywalking-website/master/data/releases.yml" | grep --after-context=7 "name: SkyWalking CLI" | grep "version" | grep -o "[0-9].[0-9].[0-9]")
if [ "$VERSION" != "" ]; then
    echo "Latest version: $VERSION"
    # Download the package.
    curl -LO "https://www.apache.org/dyn/closer.cgi/skywalking/cli/$VERSION/skywalking-cli-$VERSION-bin.tgz"
    if [ -f "skywalking-cli-$VERSION-bin.tgz" ]; then
        # Verify the integrity.
        curl -LO "https://downloads.apache.org/skywalking/cli/$VERSION/skywalking-cli-$VERSION-bin.tgz.sha512"
        VERIFY=$(sha512sum --check "skywalking-cli-$VERSION-bin.tgz.sha512")
        VERIFY="${VERIFY#* }"
        if [ "$VERIFY" = "OK" ]; then
            echo "Through verification, the file is complete."
            tar -zxvf skywalking-cli-$VERSION-bin.tgz
            # Add swctl to the environment variable PATH.
            if [ "$(uname -s)" = "Darwin" ]; then
                sudo cp skywalking-cli-$VERSION-bin/bin/swctl-$VERSION-darwin-amd64 /usr/local/bin/swctl
            else 
                sudo cp skywalking-cli-$VERSION-bin/bin/swctl-$VERSION-linux-amd64 /usr/local/bin/swctl
            fi
            # Delete unnecessary files.
            sudo rm -rf "./skywalking-cli-$VERSION-bin.tgz.sha512"
            sudo rm -rf "./skywalking-cli-$VERSION-bin.tgz"
            sudo rm -rf "./skywalking-cli-$VERSION-bin"
            echo "Type 'swctl --help' to get more information."
        else
            echo "The file is incomplete."
        fi
    else
        echo "Could not found skywalking-cli-$VERSION-bin.tgz"
    fi
else
    echo "Can't get the latest version."
fi