# Contributing to Apache SkyWalking CLI

Firstly, thanks for your interest in contributing! We hope that this will be a pleasant first experience for you, and
that you will return to continue contributing.

## Code of Conduct

This project and everyone participating in it is governed by the Apache software
Foundation's [Code of Conduct](http://www.apache.org/foundation/policies/conduct.html). By participating, you are
expected to adhere to this code. If you are aware of unacceptable behavior, please visit the
[Reporting Guidelines page](http://www.apache.org/foundation/policies/conduct.html#reporting-guidelines)
and follow the instructions there.

## How to contribute?

Most of the contributions that we receive are code contributions, but you can also contribute to the documentation or
simply report solid bugs for us to fix.

## How to report a bug?

* **Ensure the bug was not already reported** by searching on GitHub
  under [Issues](https://github.com/apache/skywalking/issues).

* If you're unable to find an open issue addressing the
  problem, [open a new one](https://github.com/apache/skywalking/issues/new). Be sure to include a **title and clear
  description**, as much relevant information as possible, and a **code sample** or an **executable test case**
  demonstrating the expected behavior that is not occurring.

## How to add a new feature or change an existing one

_Before making any significant changes, please [open an issue](https://github.com/apache/skywalking/issues)._
Discussing your proposed changes ahead of time will make the contribution process smooth for everyone.

Once we've discussed your changes and you've got your code ready, make sure that tests are passing and open your pull
request. Your PR is most likely to be accepted if it:

* Update the README.md with details of changes to the interface.
* Includes tests for new functionality.
* References the original issue in description, e.g. "Resolves #123".
* Has a [good commit message](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html).

## Compiling and building

Building needs Go 1.26 or later, and `make` also needs `license-eye` and `golangci-lint`. Clone the source code and
simply run `make` in the source directory, this will download all necessary dependencies, check the license headers,
lint, run tests, and build the binaries in `./bin/`, one per platform: `darwin` (`amd64`, `arm64`), `linux` (`386`,
`amd64`, `arm64`) and `windows` (`386`, `amd64`). `make build` only builds them.

```shell
make
```

## Writing a new command

All commands files locate in directory [`commands`](internal/commands), and an individual directory for each
second-level command, an individual `go` file for each third-level command, for example, there is a
directory [`service`](internal/commands/service) for command `swctl service`, and
a [`list.go`](internal/commands/service/list.go) file for `swctl service list` command.

Determine what entity your command will operate on, and put your command `go` file into that directory, or create one if
it doesn't exist, for example, if you want to create a command to `list` all the `instance`s of a service, create a
directory `commands/instance`, and a `go` file `commands/instance/list.go`.

## Reusing common options

There're some [common options](#common-options) that can be shared by multiple commands,
check [`commands/flags`](internal/flags)
to get all the shared options, and reuse them when possible, an example shares the options
is [`commands/service/list.go`](internal/commands/service/list.go#L35)

## Linting your codes

We have some rules for the code style and please lint your codes locally before opening a pull request

```shell
make lint
```

if you found some errors in the output of the above command, try `make fix` to fix some obvious style issues, as for the
complicated errors, please fix them manually.

## Checking license

The Apache Software Foundation requires every source file to contain a license header, run `make license` to check that
there is license header in every source file.

```shell
make license
``` 

## Running tests

Before submitting a pull request, add some test code to test the added/modified codes, and run the tests locally, make
sure all tests passed.

```shell
make test
```

## How to release

This section guides committers and PMC members to release SkyWalking CLI in Apache Way.

### Prerequisites

- [x] [GNU Make](https://www.gnu.org/software/make/manual/make.html) is installed
- [x] [GPG tool](https://gpgtools.org) is installed
- [x] [Add your GPG key](docs/How-to-release.md#add-your-gpg-public-key-to-apache-svn)

### Release steps

- Export the version that is to be released, `export VERSION=1.0.1`, and the Go version CI uses,
  `export GOTOOLCHAIN=go1.26.<latest patch>`
- [Build and sign the source code package](docs/How-to-release.md#build-and-sign-the-source-code-package): tag the
  release commit with the annotated tag `${VERSION}` (no `v` prefix), run `make release` to verify licenses, build and
  sign the distribution packages, and push the tag only after it succeeds
- [Upload the packages to SVN repository](docs/How-to-release.md#upload-to-apache-svn)
- `scripts/pre-release.sh` runs the two steps above and prints the vote email
- [Send internal announcement](docs/How-to-release.md#make-the-internal-announcement)
- [Wait at least 48 hours for test responses](docs/How-to-release.md#wait-at-least-48-hours-for-test-responses)
- [Call for vote](docs/How-to-release.md#call-for-vote-in-dev-mailing-list)
- [Publish release](docs/How-to-release.md#publish-release), `scripts/release.sh` runs its steps
