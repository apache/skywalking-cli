Apache SkyWalking CLI
===============

![](https://github.com/apache/skywalking-cli/workflows/Build/badge.svg?branch=master)
![](https://codecov.io/gh/apache/skywalking-cli/branch/master/graph/badge.svg)

<img src="http://skywalking.apache.org/assets/logo.svg" alt="Sky Walking logo" height="90px" align="right" />

The CLI (Command Line Interface) for [Apache SkyWalking](https://github.com/apache/skywalking).

SkyWalking CLI is a command interaction tool for the SkyWalking user or OPS team, as an alternative besides using
browser GUI. It is based on SkyWalking [GraphQL query protocol](https://github.com/apache/skywalking-query-protocol),
same as GUI.

## Install

### Quick install

#### Linux or macOS

Install the latest version with the following command:

```shell
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/apache/skywalking-cli/tree/master/scripts/install.sh)"
```

#### Windows

Note: you need to start cmd or powershell in administrator mode.

```shell
curl -LO "https://raw.githubusercontent.com/apache/skywalking-cli/tree/master/scripts/install.bat" && .\install.bat
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

After reloading your shell, swctl auto-completion should be working.

### `metrics`

#### `metrics linear`

<details>

<summary>metrics linear [--start=start-time] [--end=end-time] --name=metrics-name [--service=service-name] [--instance=instance-name] [--endpoint=endpoint-name] [--destService=dest-service-name] [--destInstance=dest-instance-name] [--destEndpoint=dest-endpoint-name] [--isDestNormal=true/false]</summary>

| option | description | default |
| :--- | :--- | :--- |
| `--name` | Metrics name, defined in [OAL](https://github.com/apache/skywalking/blob/master/oap-server/server-bootstrap/src/main/resources/oal/core.oal). |
| `--service` | The name of the service. | "" |
| `--instance` | The name of the service instance. | "" |
| `--endpoint` | The name of the endpoint. | "" |
| `--destService` | The name of the destination service. | "" |
| `--destInstance` | The name of the destination instance. | "" |
| `--destEndpoint` | The name of the destination endpoint. | "" |
| `--isDestNormal` | Set the destination service to normal or unnormal. | `true` |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

</details>

#### `metrics multiple-linear`

<details>

<summary>metrics multiple-linear [--start=start-time] [--end=end-time] --name=metrics-name [--service=service-name] [--num=number-of-linear-metrics] [--instance=instance-name] [--endpoint=endpoint-name] [--isNormal=true/false] [--destService=dest-service-name] [--destInstance=dest-instance-name] [--destEndpoint=dest-endpoint-name] [--isDestNormal=true/false]</summary>

| option | description | default |
| :--- | :--- | :--- |
| `--name` | Metrics name that ends with `_percentile`, defined in [OAL](https://github.com/apache/skywalking/blob/master/oap-server/server-bootstrap/src/main/resources/oal/core.oal), such as `all_percentile`, etc. |
| `--service` | The name of the service, when scope is `All`, no name is required. | "" |
| `--labels` | The labels you need to query | `0,1,2,3,4` |
| `--instance` | The name of the service instance. | "" |
| `--endpoint` | The name of the endpoint. | "" |
| `--isNormal` | Set the service to normal or unnormal. | `true` |
| `--destService` | The name of the destination service. | "" |
| `--destInstance` | The name of the destination instance. | "" |
| `--destEndpoint` | The name of the destination endpoint. | "" |
| `--isDestNormal` | Set the destination service to normal or unnormal. | `true` |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

</details>

#### `metrics single`

<details>

<summary>metrics single [--start=start-time] [--end=end-time] --name=metrics-name --service=service-name [--instance=instance-name] [--endpoint=endpoint-name] [--isNormal=true/false] [--destService=dest-service-name] [--destInstance=dest-instance-name] [--destEndpoint=dest-endpoint-name] [--isDestNormal=true/false]</summary>

| option | description | default |
| :--- | :--- | :--- |
| `--name` | Metrics name, defined in [OAL](https://github.com/apache/skywalking/blob/master/oap-server/server-bootstrap/src/main/resources/oal/core.oal), such as `service_sla`, etc. |
| `--service` | The name of the service. | "" |
| `--instance` | The name of the service instance. | "" |
| `--endpoint` | The name of the endpoint. | "" |
| `--isNormal` | Set the service to normal or unnormal. | `true` |
| `--destService` | The name of the destination service. | "" |
| `--destInstance` | The name of the destination instance. | "" |
| `--destEndpoint` | The name of the destination endpoint. | "" |
| `--isDestNormal` | Set the destination service to normal or unnormal. | `true` |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

</details>

### `dashboard`

#### `dashboard global-metrics`

<details>

<summary>dashboard global-metrics [--template=template]</summary>

`dashboard global-metrics` displays global metrics in the form of a dashboard.

| argument | description | default |
| :--- | :--- | :--- |
| `--template` | The template file to customize how to display information | `templates/Dashboard.Global.json` |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

You can imitate the content of [the default template file](examples/global.yml) to customize the dashboard.

</details>

#### `dashboard global`

<details>

<summary>dashboard global [--template=template]</summary>

`dashboard global` displays global metrics, global response latency and global heat map in the form of a dashboard.

| argument | description | default |
| :--- | :--- | :--- |
| `--template` | The template file to customize how to display information | `templates/dashboard/global.yml` |
| `--refresh` | The interval of auto-refresh (s). When `start` and `end` are both present, auto-refresh is disabled. | `6` |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

You can imitate the content of [the default template file](examples/global.yml) to customize the dashboard.

</details>

### `install`

#### `manifest`

<details>

<summary>install manifest oap|ui [--name=string] [--namespace=string] [-f=filepath/-]</summary>

| argument | description | default |
| :--- | :--- | :--- |
| `--name` | The name of prefix of generated resources  | `skywalking` |
| `--namespace` |  The namespace where resource will be deployed | `skywalking-system` |
| `-f` | The custom resource file describing custom resources defined by swck | |

</details>

### `event`

#### `report`

<details>

<summary>event report [--uuid=uuid] [--service=service] [--name=name] [--message=message] [--startTime=startTime] [--endTime=endTime] [--instance=instance] [--endpoint=endpoint] [--type=type] [parameters...]</summary>

`event report` reports an event to OAP server via gRPC.

| argument | description | default |
| :--- | :--- | :--- |
| `uuid` | The unique ID of the event. |  |
| `service` | The service of the event occurred on. |  |
| `instance` | The service instance of the event occurred on. |  |
| `endpoint` | The endpoint of the event occurred on. |  |
| `name` | The name of the event. For example, 'Reboot' and 'Upgrade' etc. |  |
| `type` | The type of the event, could be `Normal` or `Error`. | `Normal` |
| `message` | The detail of the event. This should be a one-line message that briefly describes why the event is reported. |  |
| `startTime` | The start time (in milliseconds) of the event, measured between the current time and midnight, January 1, 1970 UTC. |  |
| `endTime` | The end time (in milliseconds) of the event, measured between the current time and midnight, January 1, 1970 UTC. |  |

</details>

#### `list`

<details>

<summary>event list [--service=service] [--instance=instance] [--endpoint=endpoint] [--name=name] [--type=type] [--start=start-time] [--end=end-time] </summary>

| argument | description | default |
| :--- | :--- | :--- |
| `service` | The service name whose events are to displayed. |  |
| `instance` | The service instance name whose events are to displayed. |  |
| `endpoint` | The endpoint name whose logs are to displayed. |  |
| `name` | The name of the event. |  |
| `type` | The type of the event, could be `Normal` or `Error`. | `Normal` |

</details>

### `logs`

#### `list`

<details>

<summary>logs list [--service-id=service-id] [--instance-id=instance-id] [--endpoint-id=endpoint-id] [--trace-id=trace-id] [--tags=tags] [--start=start-time] [--end=end-time] </summary>

| argument | description | default |
| :--- | :--- | :--- |
| `service-id` | The service id whose logs are to displayed. |  |
| `instance-id` | The service instance id whose logs are to displayed. |  |
| `endpoint-id` | The service endpoint id whose logs are to displayed. |  |
| `trace-id` | The trace id whose logs are to displayed. |  |
| `tags` | Only tags defined in the core/default/searchableLogsTags are searchable. Check more details on the Configuration Vocabulary page | See [Configuration Vocabulary page](https://github.com/apache/skywalking/blob/master/docs/en/setup/backend/configuration-vocabulary.md) |

</details>

### `profile`

#### `create`

<details>

<summary>profile create [--service-id=service-id] [--service-name=service-name] [--endpoint=endpoint] [--start-time=start-time] [--duration=duration] [--min-duration-threshold=min-duration-threshold] [--dump-period=dump-period] [--max-sampling-count=max-sampling-count] </summary>

| argument | description | default |
| :--- | :--- | :--- |
| `service-id` | <service-id> whose endpoints are to be profile. |  |
| `service-name` | <service-name> whose endpoints are to be profile. |  |
| `endpoint` | which endpoint should profile. |  |
| `start-time` | profile task start time(millisecond). |  |
| `duration` | profile task continuous time(minute). | |
| `min-duration-threshold` | profiled endpoint must greater duration(millisecond). | |
| `dump-period` | profiled endpoint dump period(millisecond). | |
| `max-sampling-count` | profile task max sampling count. | |

</details>

#### `list`

<details>

<summary>profile list [--service-id=service-id] [--service-name=service-name] [--endpoint=endpoint] </summary>

| argument | description | default |
| :--- | :--- | :--- |
| `service-id` | `<service id>` whose profile task are to be searched. |  |
| `service-name` | `<service name>` whose profile task are to be searched. |  |
| `endpoint` | `<endpoint>` whose profile task are to be searched |  |

</details>

#### `segment-list`

<details>

<summary>profile segment-list [--task-id=task-id] </summary>

| argument | description | default |
| :--- | :--- | :--- |
| `task-id` | `<task id>` whose profiled segment are to be searched. |  |

</details>

#### `profiled-segment`

<details>

<summary>profile profiled-segment [--segment-id=segment-id] </summary>

| argument | description | default |
| :--- | :--- | :--- |
| `segment-id` | profiled segment id. |  |

</details>

#### `profiled-analyze`

<details>

<summary>profile profiled-analyze [--segment-id=segment-id] [--time-ranges=time-ranges] </summary>

| argument | description | default |
| :--- | :--- | :--- |
| `segment-id` | profiled segment id. |  |
| `time-ranges` | need to analyze time ranges in the segment: start-end,start-end. |  |

</details>

### `dependency`

#### `service`

<details>

<summary>dependency service [service-id] [--start=start-time] [--end=end-time]</summary>

`dependency service` shows all the dependencies of given `[service-id]` in the time range of `[start, end]`.

| argument | description | default |
| :--- | :--- | :--- |
| `service-id` | The service id whose dependencies are to displayed. |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

</details>

#### `endpoint`

<details>

<summary>dependency endpoint [endpoint-id] [--start=start-time] [--end=end-time]</summary>

`dependency endpoint` shows all the dependencies of given `[endpoint-id]` in the time range of `[start, end]`.

| argument | description | default |
| :--- | :--- | :--- |
| `endpoint-id` | The service endpoint id whose dependencies are to displayed. |  |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

</details>

#### `instance`

<details>

<summary>dependency instance [clientService-id] [serverService-id] [--start=start-time] [--end=end-time]</summary>

`dependency instance` shows the instance topology of given `[clientService-id]` and `[serverService-id]` in the time
range of `[start, end]`.

| argument | description | default |
| :--- | :--- | :--- |
| `clientService-id` | The service id of the client. |  |
| `serverService-id` | The service id of the server. |  |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

</details>

# Use Cases

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

      # your package / deployment steps... 
```

</details>

# Contributing

For developers who want to contribute to this project, see [Contribution Guide](CONTRIBUTING.md)

# License

[Apache 2.0 License.](/LICENSE)
