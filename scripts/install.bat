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

@REM Installation (Note: you need to start cmd or powershell in administrator mode.)
@echo off
setlocal ENABLEDELAYEDEXPANSION

@REM  Get the latest version of swctl.
set FLAG="FALSE"
set VERSION= UNKNOW
curl -LO "https://raw.githubusercontent.com/apache/skywalking-website/master/data/releases.yml"
if EXIST "releases.yml" (
    for /F "tokens=1,2,*" %%i in ('FINDSTR "name version" "./releases.yml"') do (
        if !FLAG! EQU "TRUE" (
            set FLAG="FALSE"
            set VERSION=%%k
        )
        if "%%k" == "SkyWalking CLI" (set FLAG="TRUE")
    )
)
set VERSION=%VERSION:~1%
@echo The latest version of swctl is %VERSION%

if "%VERSION%" NEQ "UNKNOW" (

    @REM Download the binary package.
    curl -LO "https://apache.website-solution.net/skywalking/cli/%VERSION%/skywalking-cli-%VERSION%-bin.tgz"
    if EXIST "skywalking-cli-%VERSION%-bin.tgz" (
        tar -zxvf ".\skywalking-cli-%VERSION%-bin.tgz"

        @REM Verify the integrity of the downloaded file.
        curl -LO "https://archive.apache.org/dist/skywalking/cli/%VERSION%/skywalking-cli-%VERSION%-bin.tgz.sha512"
        CertUtil -hashfile skywalking-cli-%VERSION%-bin.tgz sha512 | findstr /X "[0-9a-zA-Z]*" > verify.txt
        for /F "tokens=*" %%i in ( 'type ".\verify.txt"' ) do ( set VERIFY1="%%i  skywalking-cli-%VERSION%-bin.tgz" )
        for /F "tokens=*" %%i in ( 'type ".\skywalking-cli-%VERSION%-bin.tgz.sha512"' ) do ( set VERIFY2="%%i" )
        if "!VERIFY1!" EQU "!VERIFY2!" (
            @echo Through verification, the file is complete.
            mkdir "C:\Program Files\swctl-cli"

            @REM Add swctl to the environment variable PATH.
            copy ".\skywalking-cli-%VERSION%-bin\bin\swctl-%VERSION%-windows-amd64" "C:\Program Files\swctl-cli\swctl.exe"
            setx "Path" "C:\Program Files\swctl-cli\;%path%" /m

            @REM Delete unnecessary files.
            del ".\skywalking-cli-%VERSION%-bin.tgz" ".\verify.txt" 
            del ".\skywalking-cli-%VERSION%-bin.tgz.sha512" ".\releases.yml"
            rd /S /Q ".\skywalking-cli-%VERSION%-bin"
            
            @echo Reopen the terminal and type 'swctl --help' to get more information.
        ) else (
            @echo The file is incomplete.
        )
    ) else (
        @echo Failed to download skywalking-cli-%VERSION%-bin.tgz
    )
) else (
    @echo Can't get the latest version. The install script may be invalid, try other install methods please.
)
