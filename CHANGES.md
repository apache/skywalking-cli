Changes by Version
==================
Release Notes.

0.11.0
------------------

### Features

- Update the `profiling ebpf create` task relate to the service level.(#141)
- Update the sub-command `process list/get` to add the `labels` field.(#141)
- Add the sub-command `profiling ebpf create prepare` to query data for prepare creating task.(#141)
- Add the sub-command `process estimate scale` to estimate the process scale.(#142)
- Update the `process list` to must be provided the instance and duration.(#144)
- Remove `layer` field in the `instance` and `process` commands for 9.1.0 GraphQL query protocol's breaking changes (#149)
- Remove `duration` flag in `profiling ebpf schedules`.(#150)
- [Breaking Change] Remove `total` field in `trace list` and `logs list` commands.(#152)
- [Breaking Change] Remove `total` field in `event list`, `browser logs list`, `alarm list` commands.(#153)

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
