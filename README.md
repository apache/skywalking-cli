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
make clean && make
```

and copy the `./bin/swctl` to your `PATH` directory, usually `/usr/bin/` or `/usr/local/bin`, or you can copy it to any directory you like,
and add that directory to `PATH`.

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

### `--start`, `--end`
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


## All available commands
This section covers all the available commands in SkyWalking CLI and their usages.

- [`swctl`](#swctl-top-level-command)
- [`service`](#service-second-level-command) (second level command)
  - [`list`, `ls`](#service-list---startstart-time---endend-time)
- [`instance`](#instance-second-level-command) (second level command)
  - [`list`, `ls`](#instance-list---service-idservice-id---service-nameservice-name---startstart-time---endend-time)
  - [`search`](#instance-search---regexinstance-name-regex---service-idservice-id---service-nameservice-name---startstart-time---endend-time)

### `swctl` top-level command
`swctl` is the top-level command, which has some options that will take effects globally.

| option | description | default |
| :--- | :--- | :--- |
| `--config` | from where the default options values will be loaded | `~/.skywalking.yml` |
| `--debug` | enable debug mode, will print more detailed information at runtime | `false` |
| `--base-url` | base url of GraphQL backend | `http://127.0.0.1:12800/graphql` |
| `--display` | display style when printing the query result, supported styles are: `json`, `yaml`, `table` | `json` |

### `service` second-level command
`service` second-level command is an entry for all operations related to services,
and it also has some options and third-level commands.

#### `service list [--start=<start time>] [--end=<end time>]`
`service list` lists all the services in the time range of \[`start`, `end`\].

| option | description | default |
| :--- | :--- | :--- |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

### `instance` second-level command
`instance` second-level command is an entry for all operations related to instances,
and it also has some options and third-level commands.

#### `instance list [--service-id=<service id>] [--service-name=<service name>] [--start=<start time>] [--end=<end time>]`
`instance list` lists all the instances in the time range of \[`start`, `end`\] and given --service-id or --service-name.

| option | description | default |
| :--- | :--- | :--- |
| `--service-id` | Query service id (priority over --service-name)|  |
| `--service-name` | Query service name |  |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

#### `instance search [--regex=<instance name regex>] [--service-id=<service id>] [--service-name=<service name>] [--start=<start time>] [--end=<end time>]`
`instance search` filter the instance in the time range of \[`start`, `end`\] and given --regex --service-id or --service-name.

| option | description | default |
| :--- | :--- | :--- |
| `--regex` | Query regex of instance name|  |
| `--service-id` | Query service id (priority over --service-name)|  |
| `--service-name` | Query service name |  |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

### `endpoint` second-level command
`endpoint` second-level command is an entry for all operations related to endpoints,
and it also has some options and third-level commands.

#### `endpoint list [--start=<start time>] [--end=<end time>] --service-id=<service id> [--limit=<count>] [--keyword=<search keyword>]`
`endpoint list` lists all the endpoints of the given service id in the time range of \[`start`, `end`\].

| option | description | default |
| :--- | :--- | :--- |
| `--service-id` | <service id> whose endpoints are to be searched |
| `--limit` | returns at most <limit> endpoints (default: 100) |
| `--keyword` | <keyword> of the endpoint name to search for, empty to search all |

### `linear-metrics` second-level command
`linear-metrics` second-level command is an entrance for all operations related to linear metrics,
and it also has some options.

| option | description | default |
| :--- | :--- | :--- |
| `--name` | Metrics name, defined in [OAL](https://github.com/apache/skywalking/blob/master/oap-server/server-bootstrap/src/main/resources/official_analysis.oal), such as `all_p99`, etc. |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

### `single-metrics` second-level command
`single-metrics` second-level command is an entrance for all operations related to single-value metrics,
and it also has some options.

| option | description | default |
| :--- | :--- | :--- |
| `--name` | Metrics name, defined in [OAL](https://github.com/apache/skywalking/blob/master/oap-server/server-bootstrap/src/main/resources/official_analysis.oal), such as `service_sla`, etc. |
| `--ids` | IDs that are required by the metric type, such as service IDs for `service_sla` |
| `--start` | See [Common options](#common-options) | See [Common options](#common-options) |
| `--end` | See [Common options](#common-options) | See [Common options](#common-options) |

# Developer guide

## Compiling and building
Clone the source code and simply run `make` in the source directory,
this will download all necessary dependencies and run tests, lint, and build three binary files in `./bin/`, for Windows, Linux, MacOS respectively.

```shell
make
```

## Writing a new command
All commands files locate in directory [`commands`](commands), and an individual directory for each second-level command,
an individual `go` file for each third-level command, for example, there is a directory [`service`](commands/service) for command `swctl service`, 
and a [`list.go`](commands/service/list.go) file for `swctl service list` command.

Determine what entity your command will operate on, and put your command `go` file into that directory, or create one if it doesn't exist,
for example, if you want to create a command to `list` all the `instance`s of a service, create a directory `commands/instance`,
and a `go` file `commands/instance/list.go`.

## Reusing common options
There're some [common options](#common-options) that can be shared by multiple commands, check [`commands/flags`](commands/flags)
to get all the shared options, and reuse them when possible, an example shares the options is [`commands/service/list.go`](commands/service/list.go#L35)

## Linting your codes
We have some rules for the code style and please lint your codes locally before opening a pull request

```shell
make lint
```

if you found some errors in the output of the above command, try `make fix` to fix some obvious style issues, as for the complicated errors, please fix them manually.

## Checking license
The Apache Software Foundation requires every source file to contain a license header, run `make license` to check that there is license header in every source file.

```shell
make license
``` 

## Running tests
Before submitting a pull request, add some test code to test the added/modified codes,
and run the tests locally, make sure all tests passed.

```shell
make test
```

# License
[Apache 2.0 License.](/LICENSE)
