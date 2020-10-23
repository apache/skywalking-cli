Apache SkyWalking CLI
===============

![](https://github.com/apache/skywalking-cli/workflows/Build/badge.svg?branch=master)
![](https://codecov.io/gh/apache/skywalking-cli/branch/master/graph/badge.svg)

<img src="http://skywalking.apache.org/assets/logo.svg" alt="Sky Walking logo" height="90px" align="right" />

The CLI (Command Line Interface) for [Apache SkyWalking](https://github.com/apache/skywalking).

SkyWalking CLI is a command interaction tool for the SkyWalking user or OPS team, as an alternative besides using browser GUI.
It is based on SkyWalking [GraphQL query protocol](https://github.com/apache/skywalking-query-protocol), same as GUI.

# Download
Go to the [download page](https://skywalking.apache.org/downloads/) to download all available binaries, including MacOS, Linux, Windows.
If you want to try the latest features, however, you can compile the latest codes yourself, as the guide below. 

# Install
As SkyWalking CLI is using `Makefile`, compiling the project is as easy as executing a command in the root directory of the project.

```shell
git clone https://github.com/apache/skywalking-cli
cd skywalking-cli
git submodule init
git submodule update
make
```

Then copy the `./bin/swctl-latest-(darwin|linux|windows)-amd64` to your `PATH` directory according to your OS,
usually `/usr/bin/` or `/usr/local/bin`, or you can copy it to any directory you like,
and add that directory to `PATH`, we recommend you to rename the `swctl-latest-(darwin|linux|windows)-amd64` to `swctl`.

# Commands
Commands in SkyWalking CLI are organized into two levels, in the form of `swctl --option <level1> --option <level2> --option`,
there're options in each level, which should follow right after the corresponding command, take the following command as example:

```shell
$ swctl --debug service list --start="2019-11-11" --end="2019-11-12"
```

where `--debug` is is an option of `swctl`, and since the `swctl` is a top-level command, `--debug` is also called global option,
and `--start` is an option of the third level command `list`, there is no option for the second level command `service`.

Generally, the second level commands are entity related, there're entities like `service`, `service instance`, `metrics` in SkyWalking,
and we have corresponding sub-command like `service`; the third level commands are operations on the entities, such as `list` command
will list all the `service`s, `service instance`s, etc.

## Common options
There're some common options that are shared by multiple commands, and they follow the same rules in different commands,

<details>

<summary>--start, --end, --timezone</summary>

`--start` and `--end` specify a time range during which the query is preformed,
they are both optional and their default values follow the rules below:

- when `start` and `end` are both absent, `start = now - 30 minutes` and `end = now`, namely past 30 minutes;
- when `start` and `end` are both present, they are aligned to the same precision by **truncating the more precise one**,
e.g. if `start = 2019-01-01 1234, end = 2019-01-01 18`, then `start` is truncated (because it's more precise) to `2019-01-01 12`,
and `end = 2019-01-01 18`;
- when `start` is absent and `end` is present, will determine the precision of `end` and then use the precision to calculate `start` (minus 30 units),
e.g. `end = 2019-11-09 1234`, the precision is `MINUTE`, so `start = end - 30 minutes = 2019-11-09 1204`,
and if `end = 2019-11-09 12`, the precision is `HOUR`, so `start = end - 30HOUR = 2019-11-08 06`;
- when `start` is present and `end` is absent, will determine the precision of `start` and then use the precision to calculate `end` (plus 30 units),
e.g. `start = 2019-11-09 1204`, the precision is `MINUTE`, so `end = start + 30 minutes = 2019-11-09 1234`,
and if `start = 2019-11-08 06`, the precision is `HOUR`, so `end = start + 30HOUR = 2019-11-09 12`;

`--timezone` specifies the timezone where `--start` `--end` are based, in the form of `+0800`:

- if `--timezone` is given in the command line option, then it's used directly;
- else if the backend support the timezone API (since 6.5.0), CLI will try to get the timezone from backend, and use it;
- otherwise, the CLI will use the current timezone in the current machine; 

</details>

## All available commands
This section covers all the available commands in SkyWalking CLI and their usages.

### `swctl`
`swctl` is the top-level command, which has some options that will take effects globally.

| option | description | default |
| :--- | :--- | :--- |
| `--config` | from where the default options values will be loaded | `~/.skywalking.yml` |
| `--debug` | enable debug mode, will print more detailed information at runtime | `false` |
| `--base-url` | base url of GraphQL backend | `http://127.0.0.1:12800/graphql` |
| `--display` | display style when printing the query result, supported styles are: `json`, `yaml`, `table`, `graph` | `json` |

Note that not all display styles (except for `json` and `yaml`) are supported in all commands due to data formats incompatibilities and the limits of
Ascii Graph, like coloring in terminal, so please use `json`  or `yaml` instead.

### `service`

<details>

<summary>service list [--start=start-time] [--end=end-time]</summary>

`service list` lists all the services in the time range of `[start, end]`.

| option | description | default |
| :--- | :--- | :--- |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

</details>

### `instance`

<details>

<summary>instance list [--start=start-time] [--end=end-time] [--service-id=service-id] [--service-name=service-name]</summary>

`instance list` lists all the instances in the time range of `[start, end]` and given `--service-id` or `--service-name`.

| option | description | default |
| :--- | :--- | :--- |
| `--service-id` | Query by service id (priority over `--service-name`)|  |
| `--service-name` | Query by service name if `--service-id` is absent |  |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

</details>

<details>

<summary>instance search [--start=start-time] [--end=end-time] [--regex=instance-name-regex] [--service-id=service-id] [--service-name=service-name]</summary>

`instance search` filter the instance in the time range of `[start, end]` and given --regex --service-id or --service-name.

| option | description | default |
| :--- | :--- | :--- |
| `--regex` | Query regex of instance name|  |
| `--service-id` | Query by service id (priority over `--service-name`)|  |
| `--service-name` | Query by service name if `service-id` is absent |  |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

</details>

### `endpoint`

<details>

<summary>endpoint list [--start=start-time] [--end=end-time] --service-id=service-id [--limit=count] [--keyword=search-keyword]</summary>

`endpoint list` lists all the endpoints of the given service id in the time range of `[start, end]`.

| option | description | default |
| :--- | :--- | :--- |
| `--service-id` | <service id> whose endpoints are to be searched | |
| `--limit` | returns at most <limit> endpoints (default: 100) | 100 |
| `--keyword` | <keyword> of the endpoint name to search for, empty to search all | "" |

</details>

### `metrics`

#### `metrics linear`

<details>

<summary>metrics linear [--start=start-time] [--end=end-time] --name=metrics-name --service=service-name</summary>

| option | description | default |
| :--- | :--- | :--- |
| `--name` | Metrics name, defined in [OAL](https://github.com/apache/skywalking/blob/master/oap-server/server-bootstrap/src/main/resources/oal/core.oal). |
| `--service` | The name of the service, when metrics' scope is `All`, no service name is required |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

</details>

#### `metrics multiple-linear`

<details>

<summary>metrics multiple-linear [--start=start-time] [--end=end-time] --name=metrics-name --service=service-name [--num=number-of-linear-metrics]</summary>

| option | description | default |
| :--- | :--- | :--- |
| `--name` | Metrics name, defined in [OAL](https://github.com/apache/skywalking/blob/master/oap-server/server-bootstrap/src/main/resources/oal/core.oal), such as `all_percentile`, etc. |
| `--service` | The name of the service, when scope is `All`, no name is required |
| `--num` | Number of the linear metrics to fetch | `5` |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

</details>

#### `metrics single`

<details>

<summary>metrics single [--start=start-time] [--end=end-time] --name=metrics-name --service=service-name</summary>

| option | description | default |
| :--- | :--- | :--- |
| `--name` | Metrics name, defined in [OAL](https://github.com/apache/skywalking/blob/master/oap-server/server-bootstrap/src/main/resources/oal/core.oal), such as `service_sla`, etc. |
| `--service` | The name of the service, when scope is `All`, no name is required |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

</details>

#### `metrics top <n>`

<details>

<summary>metrics top 5 [--start=start-time] [--end=end-time] --name=metrics-name [--order=DES]</summary>

| option | description | default |
| :--- | :--- | :--- |
| arguments | The first argument is the number of top entities | `5` |
| `--name` | Metrics name, defined in [OAL](https://github.com/apache/skywalking/blob/master/oap-server/server-bootstrap/src/main/resources/oal/core.oal), such as `service_sla`, etc. |
| `--order` | The order of metrics, `DES` or `ASC`. |`DES`|
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

</details>

#### `metrics thermodynamic`

<details>

<summary>metrics thermodynamic --name=metrics-name</summary>

| option | description | default |
| :--- | :--- | :--- |
| `--name` | Metrics name, defined in [OAL](https://github.com/apache/skywalking/blob/master/oap-server/server-bootstrap/src/main/resources/oal/core.oal), such as `all_heatmap`, etc. |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

</details>

### `trace`

<details>

<summary>trace [trace id]</summary>

`trace` displays the spans of a given trace.

| argument | description | default |
| :--- | :--- | :--- |
| `trace id` | The trace id whose spans are to displayed |  |

</details>

#### `trace ls`

<details>

<summary>trace ls</summary>

| argument | description | default |
| :--- | :--- | :--- |
| `--trace-id` | The trace id whose spans are to displayed |  |
| `--service-id` | The service id whose trace are to displayed |  |
| `--service-instance-id` | The service instance id whose trace are to displayed |  |
| `--tags` | Only tags defined in the core/default/searchableTagKeys are searchable. Check more details on the Configuration Vocabulary page | See [Configuration Vocabulary page](https://github.com/apache/skywalking/blob/master/docs/en/setup/backend/configuration-vocabulary.md) |
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

You can imitate the content of [the default template file](example/global.yml) to customize the dashboard.

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

You can imitate the content of [the default template file](example/global.yml) to customize the dashboard.

</details>

### `checkHealth`

<details>

<summary>checkHealth [--grpc=true/false] [--grpcAddr=host:port] [--grpcTLS=true/false]</summary>

| argument | description | default |
| :--- | :--- | :--- |
| `--grpc` | Enable/Disable check gRPC endpoint | `true` |
| `--grpcAddr` | The address of gRPC endpoint | `127.0.0.1:11800` |
| `--grpcTLS` | Enable/Disable TLS to access gRPC endpoint | `false` |

*Notice: Once enable gRPC TLS, checkHealth command would ignore server's cert.

</details>

# Use Cases

<details>

<summary>Query a specific service by name</summary>

```shell
# query the service named projectC
$ ./bin/swctl service ls projectC
[{"id":"4","name":"projectC"}]
```

</details>

<details>

<summary>Query instances of a specific service</summary>

If you have already got the `id` of the service:

```shell
$ ./bin/swctl instance ls --service-id=3
[{"id":"3","name":"projectD-pid:7909@skywalking-server-0001","attributes":[{"name":"os_name","value":"Linux"},{"name":"host_name","value":"skywalking-server-0001"},{"name":"process_no","value":"7909"},{"name":"ipv4s","value":"192.168.252.12"}],"language":"JAVA","instanceUUID":"ec8a79d7cb58447c978ee85846f6699a"}]
```

otherwise,

```shell
$ ./bin/swctl instance ls --service-name=projectC
[{"id":"3","name":"projectD-pid:7909@skywalking-server-0001","attributes":[{"name":"os_name","value":"Linux"},{"name":"host_name","value":"skywalking-server-0001"},{"name":"process_no","value":"7909"},{"name":"ipv4s","value":"192.168.252.12"}],"language":"JAVA","instanceUUID":"ec8a79d7cb58447c978ee85846f6699a"}]
```

</details>

<details>

<summary>Query endpoints of a specific service</summary>

If you have already got the `id` of the service:

```shell
$ ./bin/swctl endpoint ls --service-id=3
```

otherwise,

```shell
./bin/swctl service ls projectC | jq '.[].id' | xargs ./bin/swctl endpoint ls --service-id 
[{"id":"22","name":"/projectC/{value}"}]
```

</details>

<details>

<summary>Query a linear metrics graph for an instance</summary>

```shell
$ ./bin/swctl --display=graph metrics linear --name=service_instance_resp_time
┌─────────────────────────────────────────────────────────────────────────────────Press q to quit──────────────────────────────────────────────────────────────────────────────────┐
│                                                                                                                                                                                  │
│                                                                                                                                                                                  │
│         │                                                                                                                                                    ⡜⠢⡀                 │
│  1181.80│                                      ⡰⡀                                                         ⢀⡠⢢         ⡰⢣                                    ⡰⠁ ⠈⠢⡀               │
│         │                                     ⢠⠃⠱⡀              ⡀                                       ⢀⠔⠁  ⠱⡀     ⢀⠜  ⢣                        ⢀⠞⡄       ⢠⠃    ⠈⠢⡀             │
│         │                                     ⡎  ⠱⡀          ⢀⠔⠊⠱⡀                 ⢀⣀⣀⣀              ⢀⡠⠊⠁     ⠘⢄   ⢀⠎    ⢣                      ⡠⠃ ⠘⡄      ⡎       ⠈⠑⠢⢄⡀  ⢀⡠⠔⠊⠁  │
│         │          ⢀⠤⣀⡀       ⢀⡀             ⡸    ⢣        ⡠⠔⠁   ⠱⡀            ⡠⠊⠉⠉⠁   ⠉⠉⠒⠒⠤⠤⣀⣀⣀ ⢀⡠⠔⠊⠁          ⠣⡀⡠⠃      ⢣           ⢀⠔⠤⡀     ⡰⠁   ⠘⡄    ⡜            ⠈⠑⠊⠁      │
│  1043.41│⡀       ⢀⠔⠁  ⠈⠑⠒⠤⠔⠒⠊⠉⠁⠈⠒⢄          ⢀⠇     ⢣    ⢀⠤⠊       ⠱⡀         ⢀⠔⠁                ⠉⠁               ⠑⠁        ⢣         ⡠⠃  ⠈⠒⢄ ⢀⠜      ⠘⡄  ⢰⠁                      │
│         │⠈⠑⠤⣀   ⡠⠊                ⠑⠤⡀       ⡜       ⢣ ⣀⠔⠁          ⠱⡀       ⡰⠁                                              ⠣⢄⣀    ⢠⠊       ⠉⠊        ⠘⡄⢠⠃                       │
│         │    ⠑⠢⠊                    ⠈⠢⡀    ⢰⠁        ⠋              ⠱⡀  ⣀⠤⠔⠊                                                   ⠉⠒⠢⠔⠁                   ⠘⠎                        │
│         │                             ⠈⠢⡀ ⢀⠇                         ⠑⠊⠉                                                                                                         │
│      905│                               ⠈⠢⡜                                                                                                                                      │
│         └──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────  │
│          2019-12-02 2121   2019-12-02 2107   2019-12-02 2115   2019-12-02 2119   2019-12-02 2137   2019-12-02 2126   2019-12-02 2118   2019-12-02 2128   2019-12-02 2136         │
│                                                                                                                                                                                  │
│                                                                                                                                                                                  │
└──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

otherwise

```shell
$ ./bin/swctl instance ls --service-name=projectC | jq '.[] | select(.name == "projectC-pid:7895@skywalking-server-0001").id' | xargs ./bin/swctl --display=graph metrics linear --name=service_instance_resp_time --id
┌─────────────────────────────────────────────────────────────────────────────────Press q to quit──────────────────────────────────────────────────────────────────────────────────┐
│                                                                                                                                                                                  │
│                                                                                                                                                                                  │
│         │                                                                           ⡠⠒⢣                                                                                          │
│  1181.80│                          ⡠⠊⢢                                           ⣀⠔⠉   ⢣              ⡔⡄                               ⡔⡄                                        │
│         │           ⣀            ⡠⠊   ⠑⡄                                    ⣀⡠⠔⠒⠉       ⢣            ⡜ ⠈⢆                            ⢀⠎ ⠈⢢              ⡀                        │
│         │          ⡜ ⠉⠒⠤⣀   ⢀⣀⣀⡠⠊      ⠈⠢⡀               ⢀⡠⢄⣀⡀            ⡰⠉             ⢣          ⡜    ⢣                          ⡠⠃    ⠑⡄        ⢀⡠⠔⠉⠘⢄                       │
│         │        ⢀⠜      ⠉⠉⠉⠁            ⠑⢄          ⢀⡠⠔⠊⠁   ⠈⠉⠑⢢        ⡰⠁               ⢣       ⢀⠎      ⠱⡀          ⢀⠦⡀         ⢀⠜       ⠈⢢ ⢀⣀⣀⡠⠤⠒⠁     ⠣⡀                  ⡀  │
│  1043.41│       ⢀⠎                         ⠑⢄      ⢀⠔⠁           ⠱⡀     ⡰⠁                 ⢣⣀    ⢀⠎        ⠘⢄        ⢀⠎ ⠈⢢      ⢀⠤⠊          ⠉⠁            ⠘⢄               ⡠⠊   │
│         │      ⢠⠃                           ⠈⠢⡀  ⡠⠒⠁              ⠘⢄   ⡰⠁                    ⠉⠉⠉⠒⠊          ⠈⢢      ⢀⠎    ⠑⢄  ⡠⠒⠁                            ⠣⠤⣀⣀⣀       ⢀⠔⠉     │
│         │⠤⠤⠤⠤⠤⠤⠃                              ⠈⠢⠊                   ⠣⡀⡰⠁                                      ⠱⡀   ⢀⠎       ⠑⠉                                    ⠉⠉⠉⠉⠒⠒⠒⠁       │
│         │                                                            ⠑⠁                                        ⠑⡄ ⢀⠎                                                             │
│      905│                                                                                                       ⠈⢆⠎                                                              │
│         └──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────  │
│          2019-12-02 2122   2019-12-02 2137   2019-12-02 2136   2019-12-02 2128   2019-12-02 2108   2019-12-02 2130   2019-12-02 2129   2019-12-02 2115   2019-12-02 2119         │
│                                                                                                                                                                                  │
│                                                                                                                                                                                  │
└──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

</details>

<details>

<summary>Query a single metrics value for a specific endpoint</summary>

```shell
$ ./bin/swctl service ls projectC | jq '.[0].id' | xargs ./bin/swctl endpoint ls --service-id | jq '.[] | [.id] | join(",")' | xargs ./bin/swctl metrics single --name endpoint_cpm --ids
[{"id":"22","value":116}]
```

</details>

<details>

<summary>Query metrics single values for all endpoints of service of id 3</summary>

```shell
$ ./bin/swctl service ls projectC | jq '.[0].id' | xargs ./bin/swctl endpoint ls --service-id | jq '.[] | [.id] | join(",")' | xargs ./bin/swctl metrics single --name endpoint_cpm --end='2019-12-02 2137' --ids
[{"id":"3","value":116}]
```

</details>

<details>

<summary>Query multiple metrics values for all percentiles</summary>

```shell
$ ./bin/swctl --display=graph --debug metrics multiple-linear --name all_percentile

┌PRESS Q TO QUIT───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│┌───────────────────────────────#0───────────────────────────────┐┌───────────────────────────────#1───────────────────────────────┐┌─────────────────────────────────#2─────────────────────────────────┐│
││      │  ⡏⠉⠉⢹   ⢸⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⡇      ⢸⠉⠉⠉⠉⠉⠉⠉⡇  ⢸⠉⠉⠉⠉⠉⠉⠉⡇   ⡏⠉⠉⠉ ││       │     ⢸⡀                       ⢸        ⢸        ⡇       ││        │                                                  ⡠⠔⡇      ││
││960.80│ ⢀⠇  ⠘⡄  ⡜            ⢣      ⢸       ⢇  ⢸       ⡇   ⡇    ││1963.60│     ⡜⡇                       ⢸        ⢸       ⢠⡇       ││ 2600.40│                                                  ⡇ ⢣      ││
││      │ ⢸    ⡇  ⡇            ⢸      ⢸       ⢸  ⡜       ⢸  ⢸     ││       │     ⡇⢸                       ⡼⡀       ⣾       ⢸⢣       ││        │                                                 ⢸  ⢸      ││
││      │ ⢸    ⡇  ⡇            ⢸      ⡸       ⢸  ⡇       ⢸  ⢸     ││       │     ⡇⠈⡆                      ⡇⡇       ⡇⡇      ⢸⢸       ││        │                                                 ⢸  ⢸      ││
││      │ ⢸    ⢣ ⢠⠃            ⠘⡄     ⡇       ⢸  ⡇       ⢸  ⢸     ││       │    ⢰⠁ ⡇                      ⡇⡇  ⡤⢤   ⡇⡇      ⡇⢸       ││        │                                                 ⡇  ⠘⡄     ││
││824.64│ ⡇    ⢸ ⢸              ⡇     ⡇       ⠈⡆ ⡇       ⠘⡄ ⡜     ││1832.88│    ⢸  ⢣                      ⡇⡇  ⡇⢸   ⡇⡇      ⡇⢸       ││ 2486.33│                                                 ⡇   ⡇     ││
││      │ ⡇    ⢸ ⢸              ⡇     ⡇        ⡇ ⡇        ⡇ ⡇     ││       │    ⢸  ⢸                      ⡇⡇ ⢸ ⠈⡆ ⢀⠇⡇     ⢠⠃⢸       ││        │                                                ⢰⠁   ⡇     ││
││      │ ⡇    ⠈⡆⡎              ⢣     ⡇        ⡇⢸         ⡇ ⡇     ││       │    ⡎  ⢸                     ⢰⠁⡇ ⢸  ⡇ ⢸ ⡇     ⢸ ⠘⡄      ││        │                       ⡀        ⢸⠉⠲⡀  ⢀         ⢸    ⢱     ││
││      │⢰⠁     ⡇⡇              ⢸     ⡇        ⢇⢸         ⡇ ⡇     ││       │    ⡇  ⢸                     ⢸ ⢱ ⢸  ⡇ ⢸ ⢣     ⢸  ⡇      ││        │⡀                     ⢰⢱    ⢀⡄  ⡇  ⢱ ⢀⠎⡆        ⡎    ⢸  ⣀⠤ ││
││688.48│⢸      ⡇⡇              ⢸     ⡇        ⢸⢸         ⢸⢸      ││1702.16│    ⡇   ⡇                    ⢸ ⢸ ⡇  ⢣ ⢸ ⢸     ⡜  ⡇      ││ 2372.24│⠱⡀       ⡴⡀  ⢀       ⢠⠃⠈⡆  ⢀⠎⠸⡀⢠⠃   ⢣⠎ ⢸  ⣠    ⡠⠃    ⢸ ⢰⠁  ││
││      │⢸      ⢱⠁              ⠘⡄    ⡇        ⢸⢸         ⢸⢸      ││       │   ⢸    ⡇                    ⢸ ⢸ ⡇  ⢸ ⢸ ⢸     ⡇  ⡇      ││        │ ⢣      ⡜ ⠱⡀⡠⠋⡆     ⣀⠎  ⢱ ⡠⠊  ⢣⢸        ⢇⡔⠁⢣ ⣀⠔⠁     ⠈⣦⠃   ││
││      │⡜      ⠸                ⡇   ⢸         ⢸⡜         ⢸⢸      ││       │   ⢸    ⡇       ⡆     ⢀⡆     ⢸ ⢸⢀⠇  ⢸ ⡎ ⢸     ⡇  ⡇      ││        │  ⡇   ⡔⠊   ⠑⠁ ⠸⡀  ⢠⠋    ⠈⠖⠁   ⠈⠇        ⠈   ⠉         ⠏    ││
││      │⡇                       ⢣   ⢸         ⠈⡇         ⠘⡜      ││       │   ⡜    ⢱      ⢠⢣  ⢰⢄ ⡜⢸     ⡇ ⢸⢸   ⢸ ⡇ ⢸    ⢠⠃  ⢱      ││        │  ⢇   ⡇        ⢣⡀ ⡎                                        ││
││552.32│⠁                       ⠸⡀  ⢸          ⡇          ⡇      ││1571.44│   ⡇    ⢸      ⢸⢸  ⡸ ⠙ ⠘⡄    ⡇ ⠘⣼    ⡇⡇ ⢸    ⢸   ⢸      ││ 2258.16│  ⢸  ⢸          ⠈⠙                                         ││
││      │                         ⢇  ⢸                     ⠁      ││       │  ⢀⠇    ⢸      ⡜⢸  ⡇    ⢇    ⡇  ⡿    ⡇⡇  ⡇   ⢸   ⢸      ││        │  ⢸  ⢸                                                     ││
││      │                         ⢸  ⢸                            ││       │⢣ ⢸     ⠸⡀     ⡇ ⡇ ⡇    ⢸    ⡇  ⡇    ⣇⠇  ⡇   ⡜   ⢸      ││        │  ⠈⡆ ⡜                                                     ││
││      │                          ⡇ ⢸                            ││       │⠈⢆⡸      ⡇⢀   ⢠⠃ ⡇⢀⠇    ⠈⡦⠔⢇⢀⠇  ⠁    ⢹   ⡇   ⡇   ⢸      ││        │   ⡇ ⡇                                                     ││
││416.16│                          ⢱ ⢸                            ││1440.72│ ⠘⡇      ⠋⠙⡄  ⢸  ⢱⢸        ⠸⣸        ⢸   ⠱⡀  ⡇   ⠈⡆     ││2144.080│   ⡇ ⡇                                                     ││
││      │                          ⠘⡄⡎                            ││       │           ⢇  ⡎  ⢸⢸         ⢿             ⠱⡀⢠⠃    ⡇     ││        │   ⢸⢸                                                      ││
││      │                           ⡇⡇                            ││       │           ⢸ ⢰⠁  ⠸⡜         ⠈              ⠘⣼     ⠧⣀    ││        │   ⢸⢸                                                      ││
││      │                           ⢸⡇                            ││       │            ⡇⡎    ⡇                         ⠈       ⠑⢄  ││        │   ⠘⡜                                                      ││
││   280│                           ⠈⡇                            ││   1310│            ⢱⠁                                          ││    2030│    ⡇                                                      ││
││      └─────────────────────────────────────────────────────────││       └────────────────────────────────────────────────────────││        └───────────────────────────────────────────────────────────││
││       2020-03-07 0111   2020-03-07 0134   2020-03-07 0133      ││        2020-03-07 0116   2020-03-07 0121   2020-03-07 0122     ││         2020-03-07 0123   2020-03-07 0139   2020-03-07 0117        ││
│└────────────────────────────────────────────────────────────────┘└────────────────────────────────────────────────────────────────┘└────────────────────────────────────────────────────────────────────┘│
│┌────────────────────────────────────────────────#3─────────────────────────────────────────────────┐┌────────────────────────────────────────────────#4─────────────────────────────────────────────────┐│
││       │                                           ⢀⢇                                              ││        │⠤⠤⠤⠤⠤⠤⡄     ⡤⠤⢤        ⢸⠑⠒⠤⠤⠤⠤⠤⠤⠤⠤⠤⠤⠤⠤⠤⠒⠊⠉⠉⠉⠉⠉⠉⠉⠒⠢⡄     ⡤⠒⠊⡇       ⢠⠔⠒⢹           ⢠⠔⠒⠉⠑⠢⠄ ││
││       │                                           ⡸⠸⡀               ⢀⡆                            ││        │      ⡇     ⡇ ⢸        ⡸                          ⡇     ⡇  ⢇       ⢸  ⢸           ⢸       ││
││3559.60│                                          ⢀⠇ ⢇              ⢀⠎⢸                            ││54073.20│      ⢱    ⢰⠁ ⠈⡆       ⡇                          ⢱    ⢰⠁  ⢸       ⡜   ⡇          ⡎       ││
││       │           ⢀⢄                             ⡸  ⠸⡀            ⢀⠎ ⠘⡄                           ││        │      ⢸    ⢸   ⡇       ⡇                          ⢸    ⢸   ⢸       ⡇   ⡇          ⡇       ││
││       │          ⢀⠎ ⠑⢄                          ⢀⠇   ⢇           ⢀⠎   ⡇         ⣼                 ││        │      ⢸    ⢸   ⡇       ⡇                          ⢸    ⢸   ⢸       ⡇   ⡇          ⡇       ││
││       │         ⢀⠎   ⠈⢆                      ⣀  ⡸    ⠸⡀ ⣀⡀       ⡜    ⢸        ⡸⠸⡀                ││        │      ⠸⡀   ⡸   ⢇      ⢰⠁                          ⠸⡀   ⡸   ⠈⡆      ⡇   ⢣         ⢀⠇       ││
││3325.68│   ⣀⣀  ⣀⠤⠊     ⠘⡄   ⢀⣀⣀⣀⣀⡠⠤⡀      ⢀⣀⠔⠊ ⠉⠑⠃     ⠉⠉ ⠘⢄     ⡰⠁    ⠘⡄      ⢰⠁ ⡇       ⢀⣀⡠⠤⠤⠤⠄  ││43924.56│       ⡇   ⡇   ⢸      ⢸                            ⡇   ⡇    ⡇     ⢰⠁   ⢸         ⢸        ││
││       │ ⢠⠊  ⠉⠉         ⠸⡀ ⡔⠁      ⠑⢄  ⡠⠊⠉⠁                 ⠣⣀  ⢠⠃      ⡇     ⢠⠃  ⡇    ⢀⠤⠊⠁        ││        │       ⡇   ⡇   ⢸      ⢸                            ⡇   ⡇    ⡇     ⢸    ⢸         ⢸        ││
││       │⠔⠁               ⠱⠊         ⠈⠢⠊                       ⠉⠒⠎       ⠸⠤⠤⠤⠔⠊⠁   ⢇    ⢸           ││        │       ⡇   ⡇   ⢸      ⡸                            ⡇   ⡇    ⢇     ⢸    ⢸         ⢸        ││
││       │                                                                          ⢸    ⡎           ││        │       ⢱  ⢰⠁   ⠈⡆     ⡇                            ⢸  ⢸     ⢸     ⢸     ⡇        ⡎        ││
││3091.76│                                                                          ⢸    ⡇           ││33775.92│       ⢸  ⢸     ⡇     ⡇                            ⢸  ⢸     ⢸     ⡇     ⡇        ⡇        ││
││       │                                                                          ⢸   ⢀⠇           ││        │       ⢸  ⢸     ⡇     ⡇                            ⢸  ⢸     ⢸     ⡇     ⡇        ⡇        ││
││       │                                                                           ⡇  ⢸            ││        │       ⠸⡀ ⡸     ⢇    ⢰⠁                            ⠘⡄ ⡜     ⠈⡆    ⡇     ⢣       ⢠⠃        ││
││       │                                                                           ⡇  ⢸            ││        │        ⡇ ⡇     ⢸    ⢸                              ⡇ ⡇      ⡇   ⢠⠃     ⢸       ⢸         ││
││2857.84│                                                                           ⡇  ⡎            ││23627.28│        ⡇ ⡇     ⢸    ⢸                              ⡇ ⡇      ⡇   ⢸      ⢸       ⢸         ││
││       │                                                                           ⢸  ⡇            ││        │        ⡇ ⡇     ⢸    ⡸                              ⢇⢀⠇      ⢇   ⢸      ⢸       ⢸         ││
││       │                                                                           ⢸ ⢀⠇            ││        │        ⢱⢰⠁     ⠈⡆   ⡇                              ⢸⢸       ⢸   ⢸       ⡇      ⡇         ││
││       │                                                                           ⢸ ⢸             ││        │        ⢸⢸       ⡇   ⡇                              ⢸⢸       ⢸   ⡎       ⡇      ⡇         ││
││2623.92│                                                                           ⠈⡆⢸             ││13478.64│        ⢸⢸       ⡇   ⡇                              ⢸⢸       ⢸   ⡇       ⡇      ⡇         ││
││       │                                                                            ⡇⡎             ││        │        ⠸⡸       ⢇  ⢰⠁                              ⠈⡎       ⠈⡆  ⡇       ⢣     ⢠⠃         ││
││       │                                                                            ⡇⡇             ││        │         ⡇       ⢸  ⣸                                ⡇        ⡇  ⡇       ⢸     ⢸          ││
││       │                                                                            ⢱⠇             ││        │         ⠃       ⠘⠊⠉                                          ⠘⡄⢸        ⠘⠒⠊⠉⠉⠉⠉          ││
││   2390│                                                                            ⢸              ││    3330│                                                               ⠈⢾                         ││
││       └───────────────────────────────────────────────────────────────────────────────────────────││        └──────────────────────────────────────────────────────────────────────────────────────────││
││        2020-03-07 0115   2020-03-07 0139   2020-03-07 0134   2020-03-07 0136   2020-03-07 0132    ││         2020-03-07 0115   2020-03-07 0126   2020-03-07 0112   2020-03-07 0134   2020-03-07 0124   ││
│└───────────────────────────────────────────────────────────────────────────────────────────────────┘└───────────────────────────────────────────────────────────────────────────────────────────────────┘│
└──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘

```

</details>

<details>

<summary>Query the top 5 services whose sla is largest</summary>

```shell
$ ./bin/swctl metrics top 5 --name service_sla
[{"name":"load balancer1.system","id":"","value":"10000","refId":null},{"name":"load balancer2.system","id":"","value":"10000","refId":null},{"name":"projectB.business-zone","id":"","value":"10000","refId":null},{"name":"projectC.business-zone","id":"","value":"10000","refId":null},{"name":"projectD.business-zone","id":"","value":"10000","refId":null}]
```

</details>

<details>

<summary>Query the top 5 instances whose sla is largest</summary>

```shell
$ ./bin/swctl metrics top 5 --name service_instance_sla     
[{"name":"load balancer1.system - load balancer1.system","id":"","value":"10000","refId":null},{"name":"load balancer2.system - load balancer2.system","id":"","value":"10000","refId":null},{"name":"projectA.business-zone - eb38c5efeb874734a7b17de780685c55@192.168.252.12","id":"","value":"10000","refId":null},{"name":"projectB.business-zone - 4e72bad0f2c14381a5657eaaca7f33ba@192.168.252.12","id":"","value":"10000","refId":null},{"name":"projectB.business-zone - 6e0e2e1cc63145859a21fc7bf7f18d2e@192.168.252.13","id":"","value":"10000","refId":null}]
```

</details>

<details>

<summary>Query the top 5 endpoints whose sla is largest</summary>

```shell
$ ./bin/swctl metrics top 5 --name endpoint_sla  
[{"name":"load balancer1.system - /projectA/test","id":"","value":"10000","refId":null},{"name":"load balancer1.system - /","id":"","value":"10000","refId":null},{"name":"load balancer2.system - /projectA/test","id":"","value":"10000","refId":null},{"name":"load balancer2.system - /","id":"","value":"10000","refId":null},{"name":"projectA.business-zone - /projectA/{name}","id":"","value":"10000","refId":null}]
```

</details>

<details>

<summary>Query the overall heat map</summary>

```shell
$ ./bin/swctl metrics thermodynamic --name all_heatmap
{"values":[{"id":"202008290939","values":[473,3,0,0,0,0,0,0,0,0,323,0,4,0,0,0,0,0,0,0,436]},{"id":"202008290940","values":[434,0,0,0,0,0,0,0,0,0,367,0,4,0,0,0,0,0,0,0,427]},{"id":"202008290941","values":[504,0,0,0,0,0,0,0,0,0,410,0,5,0,1,0,0,0,0,0,377]},{"id":"202008290942","values":[445,0,4,0,0,0,0,0,0,0,350,0,0,0,0,0,0,0,0,0,420]},{"id":"202008290943","values":[436,0,1,0,0,0,0,0,0,0,367,0,3,0,0,0,0,0,0,0,404]},{"id":"202008290944","values":[463,0,0,0,0,0,0,0,0,0,353,0,0,0,0,0,0,0,0,0,416]},{"id":"202008290945","values":[496,0,2,3,0,0,0,0,0,0,372,0,4,0,0,0,0,0,0,0,393]},{"id":"202008290946","values":[460,0,4,0,0,0,0,0,0,0,396,0,0,0,0,0,0,0,0,0,408]},{"id":"202008290947","values":[533,0,0,0,0,0,0,0,0,0,400,0,0,0,0,0,0,0,0,0,379]},{"id":"202008290948","values":[539,0,0,0,0,0,0,0,0,0,346,0,1,0,0,0,0,0,0,0,424]},{"id":"202008290949","values":[476,0,0,0,1,0,0,0,0,0,353,0,0,0,3,0,0,0,0,0,435]},{"id":"202008290950","values":[509,0,0,0,0,0,0,0,0,0,371,0,0,0,0,0,0,0,0,0,398]},{"id":"202008290951","values":[478,0,2,0,0,0,0,0,0,0,367,0,10,0,4,0,0,0,0,0,413]},{"id":"202008290952","values":[564,0,4,0,0,0,0,0,0,0,342,0,4,0,0,0,0,0,0,0,414]},{"id":"202008290953","values":[476,0,4,0,0,0,0,0,0,0,448,0,4,0,0,0,0,0,0,0,372]},{"id":"202008290954","values":[502,0,1,0,0,0,0,0,0,0,394,0,7,0,0,0,0,0,0,0,392]},{"id":"202008290955","values":[490,0,2,0,0,0,0,0,0,0,383,0,7,0,0,0,0,0,0,0,407]},{"id":"202008290956","values":[474,0,5,0,0,0,0,0,0,0,397,0,3,0,0,0,0,0,0,0,393]},{"id":"202008290957","values":[484,0,4,0,0,0,0,0,0,0,383,0,0,0,0,0,0,0,0,0,402]},{"id":"202008290958","values":[494,0,8,0,0,0,0,0,0,0,361,0,0,0,0,0,0,0,0,0,416]},{"id":"202008290959","values":[434,0,0,0,0,0,0,0,0,0,354,0,0,0,0,0,0,0,0,0,457]},{"id":"202008291000","values":[507,0,1,0,0,0,0,0,0,0,384,0,7,0,0,0,0,0,0,0,405]},{"id":"202008291001","values":[456,0,2,0,0,0,0,0,0,0,388,0,7,0,1,0,0,0,0,0,412]},{"id":"202008291002","values":[506,0,1,0,0,0,0,0,0,0,385,0,0,0,0,0,0,0,0,0,399]},{"id":"202008291003","values":[494,0,8,0,0,0,0,0,0,0,367,0,0,0,0,0,0,0,0,0,415]},{"id":"202008291004","values":[459,0,1,0,0,0,0,0,0,0,263,0,4,0,0,0,0,0,0,0,474]},{"id":"202008291005","values":[513,0,1,0,0,0,0,0,0,0,371,0,3,0,0,0,0,0,0,0,426]},{"id":"202008291006","values":[462,0,1,0,0,0,0,0,0,0,332,0,0,0,0,0,0,0,0,0,435]},{"id":"202008291007","values":[524,0,4,0,1,0,0,0,0,0,365,0,0,0,3,0,0,0,0,0,427]},{"id":"202008291008","values":[442,0,0,0,0,0,0,0,0,0,304,0,0,0,0,0,0,0,0,0,438]},{"id":"202008291009","values":[584,0,0,0,0,0,0,0,0,0,446,0,0,0,0,0,0,0,0,0,343]}],"buckets":[{"min":"0","max":"100"},{"min":"100","max":"200"},{"min":"200","max":"300"},{"min":"300","max":"400"},{"min":"400","max":"500"},{"min":"500","max":"600"},{"min":"600","max":"700"},{"min":"700","max":"800"},{"min":"800","max":"900"},{"min":"900","max":"1000"},{"min":"1000","max":"1100"},{"min":"1100","max":"1200"},{"min":"1200","max":"1300"},{"min":"1300","max":"1400"},{"min":"1400","max":"1500"},{"min":"1500","max":"1600"},{"min":"1600","max":"1700"},{"min":"1700","max":"1800"},{"min":"1800","max":"1900"},{"min":"1900","max":"2000"},{"min":"2000","max":"infinite+"}]}
```

```shell
$ ./bin/swctl --display=graph metrics thermodynamic --name all_heatmap 
```

</details>

<details>

<summary>Query the data of global dashboard</summary>

```shell
$ ./bin/swctl db g
{"metrics":[[{"name":"projectC.business-zone","id":"","value":"210","refId":null},{"name":"projectD.business-zone","id":"","value":"198","refId":null},{"name":"projectB.business-zone","id":"","value":"198","refId":null},{"name":"load balancer1.system","id":"","value":"197","refId":null},{"name":"load balancer2.system","id":"","value":"197","refId":null},{"name":"projectA.business-zone","id":"","value":"197","refId":null}],[{"name":"projectC.business-zone","id":"","value":"4576","refId":null},{"name":"load balancer1.system","id":"","value":"2001","refId":null},{"name":"load balancer2.system","id":"","value":"2000","refId":null},{"name":"projectA.business-zone","id":"","value":"1996","refId":null},{"name":"projectB.business-zone","id":"","value":"639","refId":null},{"name":"projectD.business-zone","id":"","value":"0","refId":null}],[{"name":"load balancer2.system","id":"","value":"3086","refId":null},{"name":"projectA.business-zone","id":"","value":"3086","refId":null},{"name":"load balancer1.system","id":"","value":"3087","refId":null},{"name":"projectC.business-zone","id":"","value":"4725","refId":null},{"name":"projectB.business-zone","id":"","value":"7453","refId":null},{"name":"projectD.business-zone","id":"","value":"10000","refId":null}],[{"name":"projectC.business-zone - /projectC/{value}","id":"","value":"4576","refId":null},{"name":"projectA.business-zone - /projectA/{name}","id":"","value":"1998","refId":null},{"name":"load balancer1.system - /projectA/test","id":"","value":"1997","refId":null},{"name":"load balancer2.system - /projectA/test","id":"","value":"1997","refId":null},{"name":"projectB.business-zone - /projectB/{value}","id":"","value":"639","refId":null},{"name":"projectD.business-zone - Kafka/test-trace-topic/Consumer/test","id":"","value":"0","refId":null}]],"responseLatency":[{"2020-10-08 1911":1130,"2020-10-08 1912":1100,"2020-10-08 1913":1130,"2020-10-08 1914":1130,"2020-10-08 1915":1120,"2020-10-08 1916":1130,"2020-10-08 1917":1120,"2020-10-08 1918":1120,"2020-10-08 1919":1120,"2020-10-08 1920":1130,"2020-10-08 1921":1120,"2020-10-08 1922":1120,"2020-10-08 1923":1130,"2020-10-08 1924":1130,"2020-10-08 1925":1130,"2020-10-08 1926":1130,"2020-10-08 1927":1130,"2020-10-08 1928":1130,"2020-10-08 1929":1130,"2020-10-08 1930":1120,"2020-10-08 1931":1130,"2020-10-08 1932":1130,"2020-10-08 1933":1130,"2020-10-08 1934":1120,"2020-10-08 1935":1100,"2020-10-08 1936":1130,"2020-10-08 1937":1120,"2020-10-08 1938":1130,"2020-10-08 1939":1120,"2020-10-08 1940":1130,"2020-10-08 1941":1110},{"2020-10-08 1911":2130,"2020-10-08 1912":2130,"2020-10-08 1913":2130,"2020-10-08 1914":2130,"2020-10-08 1915":2130,"2020-10-08 1916":2130,"2020-10-08 1917":2120,"2020-10-08 1918":2130,"2020-10-08 1919":2130,"2020-10-08 1920":2130,"2020-10-08 1921":2120,"2020-10-08 1922":2120,"2020-10-08 1923":2130,"2020-10-08 1924":2120,"2020-10-08 1925":2130,"2020-10-08 1926":2130,"2020-10-08 1927":2130,"2020-10-08 1928":2170,"2020-10-08 1929":2130,"2020-10-08 1930":2120,"2020-10-08 1931":2130,"2020-10-08 1932":2130,"2020-10-08 1933":2130,"2020-10-08 1934":2020,"2020-10-08 1935":2130,"2020-10-08 1936":2130,"2020-10-08 1937":2130,"2020-10-08 1938":2130,"2020-10-08 1939":2130,"2020-10-08 1940":2130,"2020-10-08 1941":2020},{"2020-10-08 1911":3130,"2020-10-08 1912":3130,"2020-10-08 1913":3130,"2020-10-08 1914":3130,"2020-10-08 1915":3130,"2020-10-08 1916":3130,"2020-10-08 1917":3120,"2020-10-08 1918":3150,"2020-10-08 1919":3130,"2020-10-08 1920":3130,"2020-10-08 1921":3130,"2020-10-08 1922":3130,"2020-10-08 1923":3130,"2020-10-08 1924":3130,"2020-10-08 1925":3130,"2020-10-08 1926":3130,"2020-10-08 1927":3160,"2020-10-08 1928":3160,"2020-10-08 1929":3130,"2020-10-08 1930":3120,"2020-10-08 1931":3130,"2020-10-08 1932":3130,"2020-10-08 1933":3130,"2020-10-08 1934":3130,"2020-10-08 1935":3130,"2020-10-08 1936":3130,"2020-10-08 1937":3130,"2020-10-08 1938":3130,"2020-10-08 1939":3130,"2020-10-08 1940":3130,"2020-10-08 1941":3130},{"2020-10-08 1911":3190,"2020-10-08 1912":3210,"2020-10-08 1913":3290,"2020-10-08 1914":3220,"2020-10-08 1915":3260,"2020-10-08 1916":3270,"2020-10-08 1917":3230,"2020-10-08 1918":3280,"2020-10-08 1919":3260,"2020-10-08 1920":3230,"2020-10-08 1921":3180,"2020-10-08 1922":3200,"2020-10-08 1923":3210,"2020-10-08 1924":3190,"2020-10-08 1925":3440,"2020-10-08 1926":3190,"2020-10-08 1927":3240,"2020-10-08 1928":3210,"2020-10-08 1929":3260,"2020-10-08 1930":3200,"2020-10-08 1931":3240,"2020-10-08 1932":3280,"2020-10-08 1933":3200,"2020-10-08 1934":3200,"2020-10-08 1935":3260,"2020-10-08 1936":3230,"2020-10-08 1937":3360,"2020-10-08 1938":3200,"2020-10-08 1939":5190,"2020-10-08 1940":3180,"2020-10-08 1941":3170},{"2020-10-08 1911":60020,"2020-10-08 1912":8220,"2020-10-08 1913":60020,"2020-10-08 1914":60020,"2020-10-08 1915":60020,"2020-10-08 1916":60020,"2020-10-08 1917":60020,"2020-10-08 1918":60020,"2020-10-08 1919":60020,"2020-10-08 1920":60020,"2020-10-08 1921":60030,"2020-10-08 1922":60020,"2020-10-08 1923":8180,"2020-10-08 1924":8130,"2020-10-08 1925":60020,"2020-10-08 1926":8160,"2020-10-08 1927":60020,"2020-10-08 1928":8170,"2020-10-08 1929":60020,"2020-10-08 1930":60020,"2020-10-08 1931":60020,"2020-10-08 1932":8140,"2020-10-08 1933":8140,"2020-10-08 1934":8160,"2020-10-08 1935":60060,"2020-10-08 1936":60020,"2020-10-08 1937":60020,"2020-10-08 1938":60020,"2020-10-08 1939":60020,"2020-10-08 1940":8130,"2020-10-08 1941":7190}],"heatMap":{"values":[{"id":"202010081911","values":[279,171,20,0,0,0,0,1,0,0,71,247,48,11,0,0,0,0,0,1,394]},{"id":"202010081912","values":[278,226,31,0,1,0,0,0,0,0,65,223,27,0,0,0,0,0,0,0,400]},{"id":"202010081913","values":[268,156,25,1,0,0,0,0,0,0,62,257,27,9,0,0,0,0,0,0,399]},{"id":"202010081914","values":[265,143,13,6,0,0,0,0,0,1,72,263,23,4,3,0,0,0,3,3,404]},{"id":"202010081915","values":[269,166,25,0,0,0,0,1,0,1,68,232,33,6,0,0,0,0,0,0,403]},{"id":"202010081916","values":[246,141,17,3,0,0,0,0,0,0,61,255,29,8,1,0,0,0,0,0,410]},{"id":"202010081917","values":[267,182,27,0,0,0,0,0,0,0,73,248,29,4,0,0,0,0,0,0,393]},{"id":"202010081918","values":[278,192,22,5,0,0,0,0,1,0,55,241,39,1,0,0,0,0,3,0,398]},{"id":"202010081919","values":[276,165,37,4,0,0,0,0,0,1,57,268,21,3,6,0,0,0,0,0,383]},{"id":"202010081920","values":[267,181,18,0,0,0,0,0,0,0,77,259,45,3,0,0,0,0,0,0,390]},{"id":"202010081921","values":[286,169,11,12,0,0,0,0,0,0,54,264,41,14,1,0,0,0,1,3,383]},{"id":"202010081922","values":[268,178,9,0,0,3,0,0,0,0,66,268,23,3,0,1,0,0,1,1,400]},{"id":"202010081923","values":[267,169,8,0,0,0,0,0,0,0,63,284,22,0,0,0,0,0,0,0,403]},{"id":"202010081924","values":[264,174,3,0,0,0,0,0,0,0,72,292,23,0,0,0,0,0,0,0,385]},{"id":"202010081925","values":[254,189,7,0,0,0,0,0,0,0,73,238,21,0,0,0,0,0,0,0,416]},{"id":"202010081926","values":[247,136,12,0,0,0,0,0,0,0,65,257,16,5,0,0,0,0,0,0,427]},{"id":"202010081927","values":[260,135,17,3,0,0,0,1,0,0,53,244,40,3,1,3,0,3,0,0,412]},{"id":"202010081928","values":[258,116,33,3,0,0,0,0,0,0,66,253,45,4,0,0,0,0,0,0,405]},{"id":"202010081929","values":[258,190,17,0,0,0,0,0,0,0,59,214,24,4,0,0,0,0,0,0,438]},{"id":"202010081930","values":[278,197,11,4,0,0,0,0,0,0,70,254,29,0,0,0,0,0,0,0,399]},{"id":"202010081931","values":[277,184,14,0,0,0,0,0,0,0,58,253,29,3,0,0,3,1,3,0,410]},{"id":"202010081932","values":[254,169,5,0,0,0,0,0,0,0,65,217,45,10,0,0,0,5,0,0,405]},{"id":"202010081933","values":[270,152,15,1,0,0,0,0,0,0,65,284,25,0,0,0,0,0,0,0,407]},{"id":"202010081934","values":[271,188,15,0,0,0,0,0,0,0,74,308,8,9,0,0,0,0,0,0,363]},{"id":"202010081935","values":[283,174,39,7,0,0,0,0,0,0,63,254,38,2,0,0,0,0,0,0,385]},{"id":"202010081936","values":[259,150,18,3,0,0,0,0,0,0,85,260,28,0,0,0,0,0,0,0,399]},{"id":"202010081937","values":[264,191,21,1,0,0,0,1,0,0,67,226,15,1,0,0,1,3,1,3,417]},{"id":"202010081938","values":[265,175,13,0,0,0,0,0,0,0,53,259,29,3,0,0,0,0,0,0,392]},{"id":"202010081939","values":[272,191,13,0,0,0,0,0,0,0,67,267,19,3,0,0,0,0,0,0,394]},{"id":"202010081940","values":[254,136,23,0,0,0,0,0,0,0,67,274,20,2,0,0,0,0,0,0,398]},{"id":"202010081941","values":[216,141,3,0,0,0,0,0,1,0,47,205,14,4,0,0,0,0,3,1,260]}],"buckets":[{"min":"0","max":"100"},{"min":"100","max":"200"},{"min":"200","max":"300"},{"min":"300","max":"400"},{"min":"400","max":"500"},{"min":"500","max":"600"},{"min":"600","max":"700"},{"min":"700","max":"800"},{"min":"800","max":"900"},{"min":"900","max":"1000"},{"min":"1000","max":"1100"},{"min":"1100","max":"1200"},{"min":"1200","max":"1300"},{"min":"1300","max":"1400"},{"min":"1400","max":"1500"},{"min":"1500","max":"1600"},{"min":"1600","max":"1700"},{"min":"1700","max":"1800"},{"min":"1800","max":"1900"},{"min":"1900","max":"2000"},{"min":"2000","max":"infinite+"}]}}
```

</details>

<details>

<summary>Display the spans of a trace</summary>

```shell
$ ./bin/swctl --display graph trace 1585375544413.464998031.46647
```

</details>

<details>

<summary>Display the traces</summary>

```shell
$ ./bin/swctl --display graph trace ls --start='2020-08-13 1754' --end='2020-08-20 2020'  --tags='http.method=POST'
```

</details>

<details>

<summary>Display the global dashboard</summary>

```shell
$ ./bin/swctl --display graph db g
```

</details>

<details>

<summary>Display the global dashboard with a customized template</summary>

```shell
$ ./bin/swctl --display graph db g --template my-global-template.yml
```

</details>

<details>

<summary>Automatically convert to server side timezone</summary>

if your backend nodes are deployed in docker and the timezone is UTC, you may not want to convert your timezone to UTC every time you type a command, `--timezone` comes to your rescue.

```shell
$ ./bin/swctl --debug --timezone="0" service ls
```

`--timezone="+1200"` and `--timezone="-0900"` are also valid usage.

</details>

<details>

<summary>Check whether OAP server is healthy</summary>

if you want to check health status from GraphQL and the gRPC endpoint listening on 10.0.0.1:8843. 

```shell
$ ./bin/swctl checkHealth --grpcAddr=10.0.0.1:8843
```

If you only want to query GraphQL.

```shell
$ ./bin/swctl checkHealth --grpc=false
```

Once the gRPC endpoint of OAP encrypts communication by TLS.

```shell
$ ./bin/swctl checkHealth --grpcTLS=true
```

</details>

# Contributing
For developers who want to contribute to this project, see [Contribution Guide](CONTRIBUTING.md)

# License
[Apache 2.0 License.](/LICENSE)
