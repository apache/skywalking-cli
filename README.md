Apache SkyWalking CLI
===============

![](https://github.com/apache/skywalking-cli/workflows/Build/badge.svg?branch=master)
![](https://codecov.io/gh/apache/skywalking-cli/branch/master/graph/badge.svg)

<img src="http://skywalking.apache.org/assets/logo.svg" alt="Sky Walking logo" height="90px" align="right" />

The CLI (Command Line Interface) for [Apache SkyWalking](https://github.com/apache/skywalking).

SkyWalking CLI is a command interaction tool for the SkyWalking user or OPS team, as an alternative besides using browser GUI.
It is based on SkyWalking [GraphQL query protocol](https://github.com/apache/skywalking-query-protocol), same as GUI.

# Install
As SkyWalking CLI is using `Makefile`, compiling the project is as easy as executing a command in the root directory of the project.

```shell
git clone https://github.com/apache/skywalking-cli
cd skywalking-cli
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

<summary>--start, --end</summary>

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
| `--display` | display style when printing the query result, supported styles are: `json`, `yaml`, `table` | `json` |

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

### `linear-metrics`

<details>

<summary>linear-metrics [--start=start-time] [--end=end-time] --name=metrics-name [--id=entity-id]</summary>

| option | description | default |
| :--- | :--- | :--- |
| `--name` | Metrics name, defined in [OAL](https://github.com/apache/skywalking/blob/master/oap-server/server-bootstrap/src/main/resources/official_analysis.oal), such as `all_p99`, etc. |
| `--id` | the related id if the metrics requires one, e.g. for metrics `service_p99`, the service `id` is required, use `--id` to specify the service id, the same for `instance`, `endpoint`, etc. |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

</details>

### `single-metrics`

<details>

<summary>single-metrics [--start=start-time] [--end=end-time] --name=metrics-name [--ids=entity-ids]</summary>

| option | description | default |
| :--- | :--- | :--- |
| `--name` | Metrics name, defined in [OAL](https://github.com/apache/skywalking/blob/master/oap-server/server-bootstrap/src/main/resources/official_analysis.oal), such as `service_sla`, etc. |
| `--ids` | IDs that are required by the metric type, such as service IDs for `service_sla` |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

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
./bin/swctl service ls projectC | jq '.[].id' | xargs ./bin/swctl-latest-darwin-amd64 endpoint ls --service-id 
[{"id":"22","name":"/projectC/{value}"}]
```

</details>

<details>

<summary>Query a linear metrics graph for an instance</summary>

If you have already got the `id` of the instance:

```shell
$ ./bin/swctl --display=graph linear-metrics --name=service_instance_resp_time --id 5
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
$ ./bin/swctl instance ls --service-name=projectC | jq '.[] | select(.name == "projectC-pid:7895@skywalking-server-0001").id' | xargs ./bin/swctl --display=graph linear-metrics --name=service_instance_resp_time --id
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
$ ./bin/swctl service ls projectC | jq '.[0].id' | xargs ./bin/swctl endpoint ls --service-id | jq '.[] | [.id] | join(",")' | xargs ./bin/swctl single-metrics --name endpoint_cpm --ids
[{"id":"22","value":116}]
```

</details>

<details>

<summary>Query metrics single values for all endpoints of service of id 3</summary>

```shell
$ ./bin/swctl service ls projectC | jq '.[0].id' | xargs ./bin/swctl endpoint ls --service-id | jq '.[] | [.id] | join(",")' | xargs ./bin/swctl single-metrics --name endpoint_cpm --end='2019-12-02 2137' --ids
[{"id":"3","value":116}]
```

</details>

# Contributing
For developers who want to contribute to this project, see [Contribution Guide](CONTRIBUTING.md)

# License
[Apache 2.0 License.](/LICENSE)
