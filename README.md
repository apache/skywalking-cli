Apache SkyWalking CLI
===============

![](https://github.com/apache/skywalking-cli/workflows/Build/badge.svg?branch=master)
![](https://codecov.io/gh/apache/skywalking-cli/branch/master/graph/badge.svg)

<img src="https://skywalking.apache.org/assets/logo.svg" alt="Sky Walking logo" height="90px" align="right" />

The CLI (Command Line Interface) for [Apache SkyWalking](https://github.com/apache/skywalking).

SkyWalking CLI is a command interaction tool for the SkyWalking user or OPS team, as an alternative besides using
browser GUI. It is based on SkyWalking [GraphQL query protocol](https://github.com/apache/skywalking-query-protocol),
same as GUI.

## Install

### Quick install

#### Linux or macOS

Install the latest version with the following command:

```shell
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/apache/skywalking-cli/master/scripts/install.sh)"
```

#### Windows

Note: you need to start cmd or powershell in administrator mode.

```shell
curl -LO "https://raw.githubusercontent.com/apache/skywalking-cli/master/scripts/install.bat" && .\install.bat
```

### Install by available binaries

Go to the [download page](https://skywalking.apache.org/downloads/#SkyWalkingCLI) to download all available binaries,
including macOS, Linux, Windows.

### Build from source

If you want to try the latest features, you can compile the latest source code and build `swctl` by yourself. Since
SkyWalking CLI is using `Makefile`, compiling the project is as easy as executing a command in the root directory of the
project.

```shell
git clone https://github.com/apache/skywalking-cli
cd skywalking-cli
make
```

Then copy the `./bin/swctl-latest-(darwin|linux|windows)-amd64` to your `PATH` directory according to your OS,
usually `/usr/bin/` or `/usr/local/bin`.

You can also copy it to any directory you like, then add that directory to `PATH`. **We recommend you to rename
the `swctl-latest-(darwin|linux|windows)-amd64` to `swctl`.**

### Run from Docker image

If you prefer to use Docker, skywalking-cli also provides Docker images for convenient usages since 0.9.0. We also push
the snapshot Docker images to GitHub registry for developers who want to test the latest features, note that this is not
Apache releases, and it's for test only, **DO NOT** use it in your production environment.

```shell
docker run -it --rm apache/skywalking-cli service ls

# Or to use the snapshot Docker image

docker run -it --rm ghcr.io/apache/skywalking-cli/skywalking-cli  service ls
```

## Autocompletion

`swctl` provides auto-completion support for bash and powershell, which can save you a lot of typing.

### Bash

The swctl completion script for bash can be generated with the command `swctl completion bash`. Sourcing the completion
script in your shell enables swctl auto-completion:

```shell
swctl completion bash > bash_autocomplete &&
    sudo cp ./bash_autocomplete /etc/bash_completion.d/swctl &&
    echo >> ~/.bashrc &&
    echo "export PROG=swctl" >> ~/.bashrc
```

After reloading your shell, swctl auto-completion should be working.

### powershell

Similarly, run the following command in your powershell terminal to enable auto-completion:

```shell 
set-executionpolicy remotesigned -Scope CurrentUser
swctl completion powershell >> $profile
```

If you get an error like `OpenError: (:) [Out-File], DirectoryNotFoundException`, then you need to run the following
command to create `$profile` file:

```shell
New-Item -Type file -Force $profile
```

After reloading your shell, `swctl` auto-completion should be working.

## Help / Manual

`swctl` builds the help manual inside the command itself with some useful example commands, please type `swctl help`
after installation. If you want to look up detail manual for a specific sub-command, insert `help` before the last
sub-command, for example, `swctl service help list` shows the manual for command `swctl service list`,
and `swctl install manifest help oap` shows the manual for `swctl install manifest oap`.

# More Use Cases

<details>
<summary>Report events in CD workflows - GitHub Actions</summary>

Integrate skywalking-cli into your CD workflows to report events, this is an implementation of GitHub Actions, but we
welcome you to contribute plugins of other CD platforms, like Jenkins, GitLab, etc.

The usage of integration for GitHub Actions is as follows.

```yaml
# ...

jobs:
  deploy:
    strategy:
      matrix:
        instance:
          - asia-southeast
          - asia-northeast
    name: Deploy Product Service
    runs-on: ubuntu-latest
    steps:
      # other steps such as checkout ...

      - name: Wrap the deployment steps with skywalking-cli
        uses: apache/skywalking-cli@main # we always suggest using a revision instead of `main`
        with:
          oap-url: ${{ secrets.OAP_URL }}                       # Required. Set the OAP backend URL, such as example.com:11800
          auth-token: ${{ secrets.OAP_AUTH_TOKEN }}             # Optional. OAP auth token if you enable authentication in OAP
          service: product                                      # Required. Name of the service to be deployed
          instance: ${{ matrix.instance }}                      # Required. Name of the instance to be deployed
          endpoint: ""                                          # Optional. Endpoint of the service, if any
          message: "Upgrade from {fromVersion} to {toVersion}"  # Optional. The message of the event
          parameters: ""                                        # Optional. The parameters in the message, if any
          layer: "GENERAL"                                      # Required. Name of the layer to which the event belongs (case-insensitive)

      # your package / deployment steps... 
```

</details>

# Compatibility

SkyWalking CLI and SkyWalking OAP communicate with different query version, here is a summary of the compatible version of both.

| SkyWalking CLI | OAP Server Version |
|----------------|---------------|
| \> = 0.11.0    | \> = 9.2.0    |
| \> = 0.12.0    | \> = 9.3.0    |

# Contributing

For developers who want to contribute to this project, see [Contribution Guide](CONTRIBUTING.md)

# License

[Apache 2.0 License.](/LICENSE)
