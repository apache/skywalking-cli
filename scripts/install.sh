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

# Prerequisites
# 1. update change log
# 2. clear milestone issues, and create a new one if needed
# 3. export VERSION=<the version to release>

# Get the latest version number.
VERSION=$(curl "https://endpoint.fastgit.org/https://github.com/apache/skywalking-website/blob/5da4b1082da44c0548b968417005b8f4821c1712/data/releases.yml" | grep --after-context=7 "name: SkyWalking CLI" | grep "version" | grep -o "[0-9].[0-9].[0-9]")
if [ $VERSION != "" ]; then
    echo Latest version:$VERSION
    # Download the package.
    curl -LO "https://mirrors.advancedhosters.com/apache/skywalking/cli/$VERSION/skywalking-cli-$VERSION-bin.tgz"
    if [ -f "skywalking-cli-$VERSION-bin.tgz" ]; then        
        # Installation
        tar -zxvf skywalking-cli-$VERSION-bin.tgz 
        sudo cp skywalking-cli-$VERSION-bin/bin/swctl-$VERSION-linux-amd64 /usr/local/bin/swctl
        sudo rm -rf "./skywalking-cli-$VERSION-bin.tgz"
        sudo rm -rf "./skywalking-cli-$VERSION-bin"
        echo "Type 'swctl --help' to get more information."
    else
        echo Could not found skywalking-cli-$VERSION-bin.tgz
    fi
else
    echo "Can't get the latest version."
fi