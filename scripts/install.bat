@REM Licensed to the Apache Software Foundation (ASF) under one or more
@REM contributor license agreements.  See the NOTICE file distributed with
@REM this work for additional information regarding copyright ownership.
@REM The ASF licenses this file to You under the Apache License, Version 2.0
@REM (the "License"); you may not use this file except in compliance with
@REM the License.  You may obtain a copy of the License at
@REM
@REM     http://www.apache.org/licenses/LICENSE-2.0
@REM
@REM Unless required by applicable law or agreed to in writing, software
@REM distributed under the License is distributed on an "AS IS" BASIS,
@REM WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
@REM See the License for the specific language governing permissions and
@REM limitations under the License.

@REM Prerequisites
@REM 1. update change log
@REM 2. clear milestone issues, and create a new one if needed
@REM 3. export VERSION=<the version to release>

@echo off
setlocal ENABLEDELAYEDEXPANSION
set FLAG="FALSE"
set VERSION="UNKNOW"
curl -LO "https://endpoint.fastgit.org/https://github.com/apache/skywalking-website/blob/5da4b1082da44c0548b968417005b8f4821c1712/data/releases.yml"
@REM Get the latest version number
for /F "tokens=1,2,*" %%i in ('FINDSTR "name version" "./releases.yml"') do (
    if !FLAG! EQU "TRUE" (
        set FLAG="FALSE"
        set VERSION=%%k
    )
    if "%%k" == "SkyWalking CLI" (set FLAG="TRUE")
)
del "./releases.yml"
set VERSION=%VERSION:~1%
if VERSION NEQ "UNKNOW" (
    @echo Latest version:%VERSION%
    @REM Download the package
    curl -LO "https://apache.osuosl.org/skywalking/cli/%VERSION%/skywalking-cli-%VERSION%-bin.tgz"
    if EXIST "skywalking-cli-%VERSION%-bin.tgz" (
        @REM Installation (this requires you to be in privileged mode)
        tar -zxvf ".\skywalking-cli-%VERSION%-bin.tgz"
        mkdir "C:\Program Files\swctl-cli"
        @REM Add swctl to the environment variable PATH
        copy ".\skywalking-cli-%VERSION%-bin\bin\swctl-%VERSION%-windows-amd64" "C:\Program Files\swctl-cli\swctl.exe"
        setx "Path" "C:\Program Files\swctl-cli\;%path%" /m
		del ".\skywalking-cli-%VERSION%-bin.tgz"
        rd /S /Q ".\skywalking-cli-%VERSION%-bin"
        @echo Reopen the terminal and type "swctl --help" to get more information.
    ) else (
        @echo Could not found "skywalking-cli-%VERSION%-bin.tgz"
    )
) else (
    @echo Can't get the latest version.
)