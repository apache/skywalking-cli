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

@REM Installs the latest Apache SkyWalking CLI release as "C:\Program Files\swctl-cli\swctl.exe" and adds that directory
@REM to the system PATH. (Note: you need to start cmd or powershell in administrator mode.)
@echo off
setlocal

set "INSTALL_DIR=C:\Program Files\swctl-cli"

@REM Get the latest version of swctl: the latest GitHub release redirects to .../releases/tag/<version>.
set "LATEST_URL=none"
for /F "delims=" %%i in ('curl -fsSLI -o NUL -w "%%{url_effective}" https://github.com/apache/skywalking-cli/releases/latest') do set "LATEST_URL=%%i"
set "VERSION=%LATEST_URL:*/tag/=%"
echo %VERSION%| findstr /R /X "[0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*" >NUL
if errorlevel 1 (
    echo Can't get the latest version. The install script may be invalid, try other install methods please.
    exit /b 1
)
echo The latest version of swctl is %VERSION%

@REM Pick the binary for this machine: the release has windows-386 and windows-amd64, which Windows on ARM runs too.
set "ARCH=amd64"
if /I "%PROCESSOR_ARCHITECTURE%"=="x86" if not defined PROCESSOR_ARCHITEW6432 set "ARCH=386"

set "PACKAGE=skywalking-cli-%VERSION%-bin"
set "BINARY=%PACKAGE%/bin/swctl-%VERSION%-windows-%ARCH%"
set "WORK_DIR=%TEMP%\swctl-install-%RANDOM%%RANDOM%"
set "PACKAGE_FILE=%WORK_DIR%\%PACKAGE%.tgz"
set "RESULT=1"
mkdir "%WORK_DIR%" || exit /b 1

@REM Download the binary package and its checksum. downloads.apache.org has the current releases;
@REM archive.apache.org has all of them, but picks up a new one later.
call :download https://downloads.apache.org/skywalking/cli
if errorlevel 1 call :download https://archive.apache.org/dist/skywalking/cli
if errorlevel 1 (
    echo Failed to download %PACKAGE%.tgz
    goto :cleanup
)

@REM Verify the integrity of the downloaded file.
powershell -NoProfile -NonInteractive -Command "$expected = (Get-Content -Raw -LiteralPath ($env:PACKAGE_FILE + '.sha512')).Split()[0]; if ((Get-FileHash -Algorithm SHA512 -LiteralPath $env:PACKAGE_FILE).Hash -ne $expected) { exit 1 }"
if errorlevel 1 (
    echo The downloaded file is incomplete.
    goto :cleanup
)
echo The downloaded file is complete.

@REM Extract with a relative path: a GNU tar from Git for Windows reads "C:\..." as a remote host.
pushd "%WORK_DIR%"
tar -xzf "%PACKAGE%.tgz" "%BINARY%"
set "TAR_RESULT=%ERRORLEVEL%"
popd
if not "%TAR_RESULT%"=="0" (
    echo Failed to extract %BINARY% from %PACKAGE%.tgz
    goto :cleanup
)

if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"
copy /Y "%WORK_DIR%\%BINARY:/=\%" "%INSTALL_DIR%\swctl.exe" >NUL
if errorlevel 1 (
    echo Failed to copy swctl.exe to %INSTALL_DIR%, start cmd or powershell in administrator mode please.
    goto :cleanup
)

@REM Add swctl to the system PATH, once. setx would cut a PATH longer than 1024 characters.
powershell -NoProfile -NonInteractive -Command "$path = [Environment]::GetEnvironmentVariable('Path', 'Machine'); if (($path -split ';' | ForEach-Object { $_.TrimEnd('\') }) -notcontains $env:INSTALL_DIR) { [Environment]::SetEnvironmentVariable('Path', $path.TrimEnd(';') + ';' + $env:INSTALL_DIR, 'Machine') }"
if errorlevel 1 (
    echo Failed to add %INSTALL_DIR% to the system PATH.
    goto :cleanup
)

echo Installed swctl %VERSION% (windows-%ARCH%) to %INSTALL_DIR%\swctl.exe.
echo Reopen the terminal and type 'swctl --help' to get more information.
set "RESULT=0"

:cleanup
rd /S /Q "%WORK_DIR%" 2>NUL
exit /b %RESULT%

@REM Downloads the package and its checksum from the given host: call :download <host>
:download
echo Downloading %~1/%VERSION%/%PACKAGE%.tgz
curl -fL --progress-bar -o "%PACKAGE_FILE%" "%~1/%VERSION%/%PACKAGE%.tgz" || exit /b 1
curl -fsSL -o "%PACKAGE_FILE%.sha512" "%~1/%VERSION%/%PACKAGE%.tgz.sha512" || exit /b 1
exit /b 0
