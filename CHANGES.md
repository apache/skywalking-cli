Changes by Version
==================
Release Notes.

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

#### Features
- Add command `swctl service` to list services
- Add command `swctl instance` and `swctl search` to list and search instances of service.
- Add command `swctl endpoint` to list endpoints of service.
- Add command `swctl linear-metrics` to query linear metrics and plot the metrics in Ascii Graph mode.
- Add command `swctl single-metrics` to query single-value metrics.

#### Chores
- Set up GitHub actions to check code styles, licenses, and tests.
