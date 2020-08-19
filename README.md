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

and copy the `./bin/swctl-latest-(darwin|linux|windows)-amd64` to your `PATH` directory according to your OS,
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

<summary>metrics linear [--start=start-time] [--end=end-time] --name=metrics-name [--id=entity-id]</summary>

| option | description | default |
| :--- | :--- | :--- |
| `--name` | Metrics name, defined in [OAL](https://github.com/apache/skywalking/blob/master/oap-server/server-bootstrap/src/main/resources/official_analysis.oal), such as `all_p99`, etc. |
| `--id` | the related id if the metrics requires one, e.g. for metrics `service_p99`, the service `id` is required, use `--id` to specify the service id, the same for `instance`, `endpoint`, etc. |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

</details>

#### `metrics multiple-linear`

<details>

<summary>metrics multiple-linear [--start=start-time] [--end=end-time] --name=metrics-name [--id=entity-id] [--num=number-of-linear-metrics]</summary>

| option | description | default |
| :--- | :--- | :--- |
| `--name` | Metrics name, defined in [OAL](https://github.com/apache/skywalking/blob/master/oap-server/server-bootstrap/src/main/resources/official_analysis.oal), such as `all_p99`, etc. |
| `--id` | the related id if the metrics requires one, e.g. for metrics `service_p99`, the service `id` is required, use `--id` to specify the service id, the same for `instance`, `endpoint`, etc. |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--num` | Number of the linear metrics to fetch | `5` |

</details>

#### `metrics single`

<details>

<summary>metrics single [--start=start-time] [--end=end-time] --name=metrics-name [--ids=entity-ids]</summary>

| option | description | default |
| :--- | :--- | :--- |
| `--name` | Metrics name, defined in [OAL](https://github.com/apache/skywalking/blob/master/oap-server/server-bootstrap/src/main/resources/official_analysis.oal), such as `service_sla`, etc. |
| `--ids` | IDs that are required by the metric type, such as service IDs for `service_sla` |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

</details>

#### `metrics top <n>`

<details>

<summary>metrics top 3 [--start=start-time] [--end=end-time] --name endpoint_sla [--service-id 3]</summary>

| option | description | default |
| :--- | :--- | :--- |
| `--name` | Metrics name, defined in [OAL](https://github.com/apache/skywalking/blob/master/oap-server/server-bootstrap/src/main/resources/official_analysis.oal), such as `service_sla`, etc. |
| `--service-id` | service ID that are required by the metric type, such as service IDs for `service_sla` |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |
| arguments | the first argument is the number of top entities | `3` |

</details>

#### `metrics thermodynamic`

<details>

<summary>metrics thermodynamic --name=thermodynamic name --ids=entity-ids</summary>

| option | description | default |
| :--- | :--- | :--- |
| `--name` | Metrics name, defined in [OAL](https://github.com/apache/skywalking/blob/master/oap-server/server-bootstrap/src/main/resources/official_analysis.oal), such as `service_sla`, etc. |
| `--ids` | IDs that are required by the metric type, such as service IDs for `service_heatmap` |
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

You can imitate the content of [the default template file](example/Dashboard.Global.json) to customize the dashboard.

</details>

#### `dashboard global`

<details>

<summary>dashboard global [--template=template]</summary>

`dashboard global` displays global metrics, global response latency and global heat map in the form of a dashboard.

| argument | description | default |
| :--- | :--- | :--- |
| `--template` | The template file to customize how to display information | `templates/Dashboard.Global.json` |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

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

If you have already got the `id` of the instance:

```shell
$ ./bin/swctl --display=graph metrics linear --name=service_instance_resp_time --id 5
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
[{"name":"projectB","id":"2","value":10000},{"name":"projectC","id":"3","value":10000},{"name":"projectA","id":"4","value":10000},{"name":"projectD","id":"5","value":10000}]
```

</details>

<details>

<summary>Query the top 5 instances whose sla is largest, of service (id = 3)</summary>

```shell
$ ./bin/swctl metrics top 5 --name service_instance_sla --service-id 3        
[{"name":"projectC-pid:30335@skywalking-server-0002","id":"13","value":10000},{"name":"projectC-pid:22037@skywalking-server-0001","id":"2","value":10000}]
```

</details>

<details>

<summary>Query the top 5 endpoints whose sla is largest, of service (id = 3)</summary>

```shell
$ ./bin/swctl metrics top 5 --name endpoint_sla --service-id 3        
[{"name":"/projectC/{value}","id":"4","value":10000}]
```

</details>

<details>

<summary>Query the overall heat map</summary>

```shell
$ ./bin/swctl metrics thermodynamic --name all_heatmap
{"nodes":[[0,0,238],[0,1,1],[0,2,39],[0,3,31],[0,4,12],[0,5,13],[0,6,4],[0,7,3],[0,8,3],[0,9,0],[0,10,48],[0,11,3],[0,12,49],[0,13,54],[0,14,11],[0,15,9],[0,16,2],[0,17,4],[0,18,0],[0,19,1],[0,20,186],[1,0,264],[1,1,3],[1,2,51],[1,3,38],[1,4,16],[1,5,14],[1,6,3],[1,7,2],[1,8,1],[1,9,2],[1,10,51],[1,11,1],[1,12,41],[1,13,56],[1,14,16],[1,15,15],[1,16,7],[1,17,7],[1,18,3],[1,19,1],[1,20,174],[2,0,231],[2,1,3],[2,2,42],[2,3,41],[2,4,18],[2,5,4],[2,6,2],[2,7,1],[2,8,2],[2,9,0],[2,10,54],[2,11,4],[2,12,55],[2,13,48],[2,14,14],[2,15,4],[2,16,3],[2,17,2],[2,18,4],[2,19,4],[2,20,187],[3,0,231],[3,1,3],[3,2,55],[3,3,38],[3,4,18],[3,5,9],[3,6,1],[3,7,1],[3,8,1],[3,9,1],[3,10,56],[3,11,6],[3,12,38],[3,13,50],[3,14,16],[3,15,12],[3,16,4],[3,17,4],[3,18,2],[3,19,2],[3,20,183],[4,0,238],[4,1,2],[4,2,47],[4,3,49],[4,4,11],[4,5,7],[4,6,0],[4,7,0],[4,8,2],[4,9,2],[4,10,55],[4,11,3],[4,12,41],[4,13,47],[4,14,12],[4,15,7],[4,16,3],[4,17,2],[4,18,10],[4,19,0],[4,20,190],[5,0,238],[5,1,3],[5,2,42],[5,3,28],[5,4,18],[5,5,4],[5,6,2],[5,7,4],[5,8,4],[5,9,1],[5,10,54],[5,11,2],[5,12,65],[5,13,56],[5,14,17],[5,15,9],[5,16,2],[5,17,3],[5,18,0],[5,19,2],[5,20,179],[6,0,218],[6,1,1],[6,2,34],[6,3,37],[6,4,10],[6,5,5],[6,6,1],[6,7,1],[6,8,0],[6,9,3],[6,10,49],[6,11,7],[6,12,47],[6,13,43],[6,14,19],[6,15,15],[6,16,1],[6,17,4],[6,18,2],[6,19,3],[6,20,183],[7,0,242],[7,1,0],[7,2,41],[7,3,34],[7,4,21],[7,5,4],[7,6,3],[7,7,4],[7,8,1],[7,9,0],[7,10,71],[7,11,4],[7,12,47],[7,13,50],[7,14,19],[7,15,8],[7,16,6],[7,17,3],[7,18,2],[7,19,4],[7,20,174],[8,0,220],[8,1,3],[8,2,40],[8,3,36],[8,4,6],[8,5,8],[8,6,1],[8,7,5],[8,8,0],[8,9,1],[8,10,61],[8,11,2],[8,12,43],[8,13,50],[8,14,17],[8,15,11],[8,16,4],[8,17,5],[8,18,1],[8,19,1],[8,20,183],[9,0,239],[9,1,1],[9,2,48],[9,3,37],[9,4,8],[9,5,12],[9,6,2],[9,7,0],[9,8,0],[9,9,0],[9,10,74],[9,11,1],[9,12,58],[9,13,53],[9,14,17],[9,15,13],[9,16,5],[9,17,2],[9,18,2],[9,19,0],[9,20,178],[10,0,249],[10,1,2],[10,2,40],[10,3,49],[10,4,12],[10,5,8],[10,6,0],[10,7,1],[10,8,0],[10,9,0],[10,10,58],[10,11,1],[10,12,54],[10,13,47],[10,14,21],[10,15,12],[10,16,6],[10,17,4],[10,18,3],[10,19,2],[10,20,165],[11,0,240],[11,1,1],[11,2,50],[11,3,47],[11,4,10],[11,5,2],[11,6,1],[11,7,1],[11,8,2],[11,9,1],[11,10,52],[11,11,4],[11,12,41],[11,13,51],[11,14,17],[11,15,6],[11,16,1],[11,17,6],[11,18,1],[11,19,0],[11,20,199],[12,0,240],[12,1,3],[12,2,40],[12,3,41],[12,4,17],[12,5,10],[12,6,5],[12,7,2],[12,8,2],[12,9,0],[12,10,86],[12,11,1],[12,12,56],[12,13,49],[12,14,16],[12,15,7],[12,16,4],[12,17,8],[12,18,4],[12,19,3],[12,20,157],[13,0,234],[13,1,1],[13,2,53],[13,3,38],[13,4,12],[13,5,4],[13,6,0],[13,7,2],[13,8,0],[13,9,0],[13,10,59],[13,11,2],[13,12,53],[13,13,48],[13,14,18],[13,15,8],[13,16,3],[13,17,8],[13,18,1],[13,19,1],[13,20,187],[14,0,269],[14,1,0],[14,2,66],[14,3,47],[14,4,17],[14,5,4],[14,6,1],[14,7,0],[14,8,0],[14,9,0],[14,10,55],[14,11,1],[14,12,53],[14,13,48],[14,14,18],[14,15,8],[14,16,3],[14,17,3],[14,18,4],[14,19,0],[14,20,179],[15,0,254],[15,1,0],[15,2,57],[15,3,45],[15,4,8],[15,5,9],[15,6,9],[15,7,4],[15,8,3],[15,9,0],[15,10,68],[15,11,1],[15,12,52],[15,13,51],[15,14,19],[15,15,7],[15,16,4],[15,17,0],[15,18,0],[15,19,1],[15,20,177],[16,0,257],[16,1,1],[16,2,65],[16,3,50],[16,4,16],[16,5,3],[16,6,1],[16,7,0],[16,8,0],[16,9,0],[16,10,61],[16,11,3],[16,12,63],[16,13,59],[16,14,14],[16,15,9],[16,16,5],[16,17,2],[16,18,0],[16,19,0],[16,20,174],[17,0,243],[17,1,1],[17,2,63],[17,3,44],[17,4,5],[17,5,3],[17,6,0],[17,7,3],[17,8,0],[17,9,0],[17,10,66],[17,11,4],[17,12,56],[17,13,38],[17,14,11],[17,15,10],[17,16,4],[17,17,2],[17,18,3],[17,19,0],[17,20,181],[18,0,236],[18,1,3],[18,2,38],[18,3,49],[18,4,16],[18,5,5],[18,6,3],[18,7,3],[18,8,1],[18,9,0],[18,10,41],[18,11,4],[18,12,59],[18,13,49],[18,14,13],[18,15,9],[18,16,4],[18,17,1],[18,18,2],[18,19,0],[18,20,192],[19,0,238],[19,1,2],[19,2,49],[19,3,37],[19,4,15],[19,5,2],[19,6,1],[19,7,1],[19,8,3],[19,9,0],[19,10,60],[19,11,3],[19,12,58],[19,13,53],[19,14,17],[19,15,4],[19,16,2],[19,17,2],[19,18,2],[19,19,0],[19,20,185],[20,0,242],[20,1,0],[20,2,55],[20,3,36],[20,4,10],[20,5,6],[20,6,1],[20,7,1],[20,8,1],[20,9,0],[20,10,57],[20,11,4],[20,12,46],[20,13,58],[20,14,15],[20,15,11],[20,16,3],[20,17,2],[20,18,7],[20,19,0],[20,20,188],[21,0,231],[21,1,3],[21,2,50],[21,3,43],[21,4,13],[21,5,1],[21,6,0],[21,7,1],[21,8,0],[21,9,0],[21,10,57],[21,11,3],[21,12,51],[21,13,36],[21,14,15],[21,15,8],[21,16,7],[21,17,2],[21,18,3],[21,19,1],[21,20,188],[22,0,241],[22,1,2],[22,2,60],[22,3,42],[22,4,11],[22,5,8],[22,6,0],[22,7,0],[22,8,0],[22,9,0],[22,10,56],[22,11,4],[22,12,57],[22,13,46],[22,14,20],[22,15,8],[22,16,6],[22,17,1],[22,18,1],[22,19,0],[22,20,191],[23,0,240],[23,1,0],[23,2,46],[23,3,44],[23,4,20],[23,5,3],[23,6,3],[23,7,4],[23,8,1],[23,9,1],[23,10,62],[23,11,4],[23,12,64],[23,13,44],[23,14,15],[23,15,3],[23,16,4],[23,17,2],[23,18,3],[23,19,1],[23,20,181],[24,0,255],[24,1,0],[24,2,61],[24,3,41],[24,4,17],[24,5,7],[24,6,0],[24,7,1],[24,8,0],[24,9,0],[24,10,60],[24,11,3],[24,12,62],[24,13,49],[24,14,17],[24,15,10],[24,16,3],[24,17,2],[24,18,3],[24,19,2],[24,20,177],[25,0,244],[25,1,1],[25,2,56],[25,3,35],[25,4,12],[25,5,12],[25,6,2],[25,7,1],[25,8,0],[25,9,0],[25,10,66],[25,11,3],[25,12,53],[25,13,55],[25,14,20],[25,15,13],[25,16,3],[25,17,1],[25,18,3],[25,19,2],[25,20,173],[26,0,234],[26,1,1],[26,2,45],[26,3,34],[26,4,9],[26,5,6],[26,6,0],[26,7,3],[26,8,0],[26,9,1],[26,10,54],[26,11,6],[26,12,59],[26,13,48],[26,14,20],[26,15,10],[26,16,1],[26,17,2],[26,18,2],[26,19,0],[26,20,182],[27,0,228],[27,1,1],[27,2,46],[27,3,35],[27,4,5],[27,5,7],[27,6,2],[27,7,3],[27,8,2],[27,9,3],[27,10,61],[27,11,2],[27,12,61],[27,13,43],[27,14,15],[27,15,7],[27,16,3],[27,17,1],[27,18,3],[27,19,1],[27,20,187],[28,0,248],[28,1,4],[28,2,60],[28,3,45],[28,4,11],[28,5,9],[28,6,5],[28,7,1],[28,8,1],[28,9,1],[28,10,58],[28,11,2],[28,12,53],[28,13,38],[28,14,20],[28,15,10],[28,16,4],[28,17,6],[28,18,1],[28,19,2],[28,20,178],[29,0,241],[29,1,2],[29,2,46],[29,3,28],[29,4,16],[29,5,8],[29,6,4],[29,7,2],[29,8,1],[29,9,0],[29,10,66],[29,11,3],[29,12,51],[29,13,51],[29,14,28],[29,15,9],[29,16,3],[29,17,4],[29,18,3],[29,19,4],[29,20,153],[30,0,151],[30,1,1],[30,2,26],[30,3,26],[30,4,8],[30,5,4],[30,6,2],[30,7,2],[30,8,3],[30,9,1],[30,10,32],[30,11,3],[30,12,33],[30,13,25],[30,14,10],[30,15,3],[30,16,1],[30,17,3],[30,18,2],[30,19,0],[30,20,82]],"axisYStep":0}
```

```shell
$ ./bin/swctl --display=graph metrics thermodynamic --name all_heatmap 
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
