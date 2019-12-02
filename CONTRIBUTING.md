# Contributing to Apache SkyWalking CLI

Firstly, thanks for your interest in contributing! We hope that this will be a
pleasant first experience for you, and that you will return to continue
contributing.

## Code of Conduct

This project and everyone participating in it is governed by the Apache
software Foundation's [Code of Conduct](http://www.apache.org/foundation/policies/conduct.html).
By participating, you are expected to adhere to this code. If you are aware of unacceptable behavior, please visit the
[Reporting Guidelines page](http://www.apache.org/foundation/policies/conduct.html#reporting-guidelines)
and follow the instructions there.

## How to contribute?

Most of the contributions that we receive are code contributions, but you can
also contribute to the documentation or simply report solid bugs
for us to fix.

## How to report a bug?

* **Ensure the bug was not already reported** by searching on GitHub under [Issues](https://github.com/apache/skywalking/issues).

* If you're unable to find an open issue addressing the problem, [open a new one](https://github.com/apache/skywalking/issues/new).
Be sure to include a **title and clear description**, as much relevant information as possible,
and a **code sample** or an **executable test case** demonstrating the expected behavior that is not occurring.

## How to add a new feature or change an existing one

_Before making any significant changes, please [open an issue](https://github.com/apache/skywalking/issues)._
Discussing your proposed changes ahead of time will make the contribution process smooth for everyone.

Once we've discussed your changes and you've got your code ready, make sure that tests are passing and open your pull request. Your PR is most likely to be accepted if it:

* Update the README.md with details of changes to the interface.
* Includes tests for new functionality.
* References the original issue in description, e.g. "Resolves #123".
* Has a [good commit message](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html).

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