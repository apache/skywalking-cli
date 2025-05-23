Changes by Version
==================
Release Notes.

0.15.0
------------------

### Features

* Add the sub-command `profiling async` for async-profiler query API by @zhengziyi0117 in https://github.com/apache/skywalking-cli/pull/203
* Support the owner in MQE response by using [10.2 MQE query protocol](https://github.com/apache/skywalking-query-protocol/pull/141) by @zhengziyi0117 in https://github.com/apache/skywalking-cli/pull/203
* Add the sub-command `alarm autocomplete-keys` and `alarm auto-complete-values` for alarm query API by @mrproliu in https://github.com/apache/skywalking-cli/pull/210
* Adapt the alarm message query API by @mrproliu in https://github.com/apache/skywalking-cli/pull/210
* Add the owner field in the `metrics exec` query API by @mrproliu in https://github.com/apache/skywalking-cli/pull/210

### Bug Fixes

* Fix wrong process id format by @mrproliu in https://github.com/apache/skywalking-cli/pull/215

0.14.0
------------------

### Features

* Add the sub-command `dependency global` for adapt the global dependency query API by @mrproliu in https://github.com/apache/skywalking-cli/pull/198
* Upgrade crypto lib to fix cve by @mrproliu in https://github.com/apache/skywalking-cli/pull/199
* Add the **hierarchy** related commands `hierarchy service`, `hierarchy instance` and `hierarchy layer-levels` by @mrproliu in https://github.com/apache/skywalking-cli/pull/200
* Add the `layers` field to `nodes` in the `dependency service` command by @mrproliu in https://github.com/apache/skywalking-cli/pull/200
* Add the duration related flags in the `endpoint list` command by @mrproliu in https://github.com/apache/skywalking-cli/pull/201

### Bug Fixes

0.13.0
------------------

### Features

### Bug Fixes

* Fix the record list query does not support new OAP versions (with major version number > 9).

## What's Changed

* Add the sub-command `menu get` for get the ui menu items by @mrproliu in https://github.com/apache/skywalking-cli/pull/187

0.12.0
------------------

### Features
## What's Changed

* Add the sub-command `records list` for adapt the new record query API by @mrproliu in https://github.com/apache/skywalking-cli/pull/167
* Add the attached events fields into the `trace` sub-command by @mrproliu in https://github.com/apache/skywalking-cli/pull/169
* Add the sampling config file into the `profiling ebpf create network` sub-command by @mrproliu in https://github.com/apache/skywalking-cli/pull/171
* Add the sub-command `profiling continuous` for adapt the new continuous profiling API by @mrproliu in https://github.com/apache/skywalking-cli/pull/173
* Adapt the sub-command `metrics` for deprecate scope fron entity by @mrproliu in https://github.com/apache/skywalking-cli/pull/173
* Add components in topology related sub-commands. @mrproliu in https://github.com/apache/skywalking-cli/pull/175
* Add the sub-command `metrics nullable` for query the nullable metrics value. @mrproliu in https://github.com/apache/skywalking-cli/pull/176
* Adapt the sub-command `profiling trace` for adapt the new trace profiling protocol. @mrproliu in https://github.com/apache/skywalking-cli/pull/177
* Add `isEmptyValue` field in metrics related sub-commands. @mrproliu in https://github.com/apache/skywalking-cli/pull/180
* Add the sub-command `metrics execute` for execute the metrics query. @mrproliu in https://github.com/apache/skywalking-cli/pull/182
* Add the sub-command `profiling continuous monitoring` for query all continuous profiling monitoring instances. @mrproliu in https://github.com/apache/skywalking-cli/pull/182
* Add `continuousProfilingCauses.message` field in the `profiling ebpf list` comamnds by @mrproliu in https://github.com/apache/skywalking-cli/pull/184

0.11.0
------------------

### Features
## What's Changed

* Add `.github/scripts` to release source tarball by @kezhenxu94 in https://github.com/apache/skywalking-cli/pull/140
* Let the eBPF profiling could performs by service level by @mrproliu in https://github.com/apache/skywalking-cli/pull/141
* Add the sub-command for estimate the process scale by @mrproliu in https://github.com/apache/skywalking-cli/pull/142
* feature: update install.sh version regex by @Alexxxing in https://github.com/apache/skywalking-cli/pull/143
* Update the commands relate to the process by @mrproliu in https://github.com/apache/skywalking-cli/pull/144
* Add layer to event related commands by @fgksgf in https://github.com/apache/skywalking-cli/pull/145
* Add layer to events.graphql by @fgksgf in https://github.com/apache/skywalking-cli/pull/146
* Add layer field to alarms.graphql by @fgksgf in https://github.com/apache/skywalking-cli/pull/147
* Upgrade crypto lib to fix cve by @kezhenxu94 in https://github.com/apache/skywalking-cli/pull/148
* Remove `layer` field in the `instance` and `process` commands by @mrproliu in https://github.com/apache/skywalking-cli/pull/149
* Remove `duration` flag in `profiling ebpf schedules` by @mrproliu in https://github.com/apache/skywalking-cli/pull/150
* Remove `total` field in `trace list` and `logs list` commands by @mrproliu in https://github.com/apache/skywalking-cli/pull/152
* Remove `total` field in `event list`, `browser logs`, `alarm list` commands. by @mrproliu in https://github.com/apache/skywalking-cli/pull/153
* Add `aggregate` flag in `profiling ebpf analysis` commands by @mrproliu in https://github.com/apache/skywalking-cli/pull/154
* event: fix event query should query all types by default by @kezhenxu94 in https://github.com/apache/skywalking-cli/pull/155
* Fix a possible lint error and update CI lint version by @JarvisG495 in https://github.com/apache/skywalking-cli/pull/156
* Add commands for support network profiling by @mrproliu in https://github.com/apache/skywalking-cli/pull/158
* Add the components field in the process relation by @mrproliu in https://github.com/apache/skywalking-cli/pull/159
* Trim license headers in query string by @kezhenxu94 in https://github.com/apache/skywalking-cli/pull/160
* Bump up dependency swck version to fix CVE by @kezhenxu94 in https://github.com/apache/skywalking-cli/pull/161
* Bump up swck dependency for transitive dep upgrade by @kezhenxu94 in https://github.com/apache/skywalking-cli/pull/162
* Add the sub-commands for query sorted metrics/records by @mrproliu in https://github.com/apache/skywalking-cli/pull/163
* Add compatibility documentation by @mrproliu in https://github.com/apache/skywalking-cli/pull/164

0.10.0
------------------

### Features

- Allow setting `start` and `end` with relative time (#128)
- Add some commands for the browser (#126)
- Add the sub-command `service layer` to query services according to layer (#133)
- Add the sub-command `layer list` to query layer list (#133)
- Add the sub-command `instance get` to query single instance (#134)
- Add the sub-command `endpoint get` to query single endpoint info (#134)
- Change the GraphQL method to the v9 version according to the server version (#134)
- Add `normal` field to Service entity (#136)
- Add the command `process` for query Process metadata (#137)
- Add the command `profiling ebpf` for process ebpf profiling (#138)
- Support `getprofiletasklogs` query (#125)
- Support query list alarms (#127)
- [Breaking Change] Update the command `profile` as a sub-command `profiling trace`, and update `profiled-analyze` command to `analysis` (#138)
- `profiling ebpf/trace analysis` generates the profiling graph HTML on default and saves it to the current work directory (#138)

### Bug Fixes

- Fix quick install (#131)
- Set correct go version in publishing snapshot docker image (#124)
- Stop build kit container after finishing (#130)

### Chores

- Add cross platform build targets (#129)
- Update download host (#132)

0.9.0
------------------

### Features

- Add the sub-command `dependency instance` to query instance relationships (#117)

### Bug Fixes

- fix: `multiple-linear` command's `labels` type can be string type (#122)
- Add missing `dest-service-id` `dest-service-name` to `metrics linear` command (#121)
- Fix the wrong name when getting `destInstance` flag (#118)

### Chores

- Upgrade Go version to 1.16 (#120)
- Migrate tests to infra-e2e, overhaul the flags names (#119)
- Publish Docker snapshot images to ghcr (#116)
- Remove dist directory when build release source tar (#115)

0.8.0
------------------

### Features

- Add profile command
- Add logs command
- Add dependency command
- Support query events protocol
- Support auto-completion for bash and powershell

### Bug Fixes

- Fix missing service instance name in trace command

### Chores

- Optimize output by adding color to help information
- Set display style explicitly for commands in the test script
- Set different default display style for different commands
- Add scripts for quick install
- Update release doc and add scripts for release
- split into multiple workflows to speed up CI

0.7.0
------------------

### Features

- Add GitHub Action for integration of event reporter

### Bug Fixes

- Fix `metrics top` can't infer the scope automatically

### Chores

- Upgrade dependency crypto
- Refactor project to use goapi
- Move `parseScope` to pkg
- Update release doc

0.6.0
------------------

### Features

- Support authorization when connecting to the OAP
- Add `install` command and `manifest` sub-command
- Add `event` command and `report` sub-command

### Bug Fixes

- Fix the bug that can't query JVM instance metrics

### Chores

- Set up a simple test with GitHub Actions
- Reorganize the project layout
- Update year in NOTICE
- Add missing license of swck
- Use license-eye to check license header

0.5.0
------------------

### Features

- Use template files in yaml format instead
- Refactor `metrics` command to adopt metrics-v2 protocol
- Use goroutine to speed up `dashboard global` command
- Add `metrics list` command

### Bug Fixes

- Add flags of instance, endpoint and normal for `metrics` command
- Fix the problem of unable to query database metrics

### Chores

- Update release guide doc
- Add screenshots for use cases in `README.md`
- Introduce generated codes into codebase

0.4.0
------------------

### Features

- Add `dashboard global` command with auto-refresh
- Add `dashboard global-metrics` command
- Add traces search
- Refactor `metrics thermodynamic` command to adopt the new query protocol

### Bug Fixes

- Fix wrong golang standard time

0.3.0
------------------

### Features

- Add health check command
- Add `trace` command

### Bug Fixes

- Fix wrong metrics graphql path

### Chores

- Move tools setup into Makefile to easy the setup work locally

0.2.0
------------------

### Features

- Support visualization of heat map
- Support top N entities, `swctl metrics top 5 --name service_sla`
- Support thermodynamic metrics, `swctl metrics thermodynamic --name all_heatmap`
- Support multiple linear metrics, `swctl --display=graph --debug metrics multiple-linear --name all_percentile`
- Automatically make use of server timezone API when possible

### Chores

- Generate GraphQL codes dynamically
- Update merge buttons to only allow squash and commit
- Add release guide doc
- Update NOTICE year

0.1.0
------------------

### Features

- Add command `swctl service` to list services
- Add command `swctl instance` and `swctl search` to list and search instances of service.
- Add command `swctl endpoint` to list endpoints of service.
- Add command `swctl linear-metrics` to query linear metrics and plot the metrics in Ascii Graph mode.
- Add command `swctl single-metrics` to query single-value metrics.

### Chores

- Set up GitHub actions to check code styles, licenses, and tests.
