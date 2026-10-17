#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

VERSION ?= dev-$(shell git rev-parse --short HEAD)
APP_NAME = skywalking-cli
OUT_DIR = bin
BINARY = swctl
DESTDIR ?= /usr/local/bin

HUB ?= docker.io/apache

RELEASE_BIN = skywalking-cli-$(VERSION)-bin
RELEASE_SRC = skywalking-cli-$(VERSION)-src

# COPYFILE_DISABLE stops macOS bsdtar from adding AppleDouble `._*` files, and --no-xattrs (bsdtar and GNU tar)
# keeps extended attributes such as com.apple.provenance out of the tar headers.
TAR = COPYFILE_DISABLE=1 tar --no-xattrs

GO = go
GO_PATH = $$($(GO) env GOPATH)
GO_BUILD = $(GO) build
GO_GET = $(GO) get
GO_INSTALL = $(GO) install
GO_TEST = $(GO) test
GO_BUILD_FLAGS = -v
GO_BUILD_LDFLAGS = -X main.version=$(VERSION)

GO_LINT = golangci-lint
LICENSE_EYE = license-eye

SHELL = /bin/bash

BUILDS := darwin-amd64 darwin-arm64 linux-386 linux-amd64 linux-arm64 windows-386 windows-amd64
BUILD_RULE = GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO_BUILD) $(GO_BUILD_FLAGS) -ldflags "$(GO_BUILD_LDFLAGS)" -o $(OUT_DIR)/$(BINARY)-$(VERSION)-$(GOOS)-$(GOARCH) cmd/swctl/main.go

all: clean license deps lint test build

.PHONY: $(BUILDS)
$(BUILDS): GOOS = $(word 1,$(subst -, ,$@))
$(BUILDS): GOARCH = $(word 2,$(subst -, ,$@))
$(BUILDS):
	$(BUILD_RULE)

.PHONY: build
build: $(BUILDS)

.PHONY: deps
deps:
	@$(GO_GET) -v -t ./...

$(GO_LINT):
	@$(GO_LINT) version > /dev/null 2>&1 || go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.4
$(LICENSE_EYE):
	@$(LICENSE_EYE) --version > /dev/null 2>&1 || go install github.com/apache/skywalking-eyes/cmd/license-eye@d38fe05

.PHONY: lint
lint: $(GO_LINT)
	$(GO_LINT) run -v --timeout 5m ./...
.PHONY: fix-lint
fix-lint: $(GO_LINT)
	$(GO_LINT) run -v --fix ./...

# `all` and `verify` depend on `license`. Without this target, case-insensitive macOS filesystems let the LICENSE file
# satisfy it (so no license check ran), and Linux failed with "No rule to make target 'license'".
.PHONY: license
license: license-header dependency-license

.PHONY: license-header
license-header: clean $(LICENSE_EYE)
	@$(LICENSE_EYE) header check

.PHONY: fix-license-header
fix-license-header: clean $(LICENSE_EYE)
	@$(LICENSE_EYE) header fix

.PHONY: dependency-license
dependency-license: clean $(LICENSE_EYE)
	@$(LICENSE_EYE) dependency resolve --summary ./dist/LICENSE.tpl --output ./dist/licenses || exit 1
	@if [ ! -z "`git diff -U0 ./dist`" ]; then \
		echo "LICENSE file is not updated correctly"; \
		git diff -U0 ./dist; \
		exit 1; \
	fi

.PHONY: fix-dependency-license
fix-dependency-license: clean $(LICENSE_EYE)
	@$(LICENSE_EYE) dependency resolve --summary ./dist/LICENSE.tpl --output ./dist/licenses

.PHONY: fix-license
fix-license: fix-license-header fix-dependency-license

.PHONY: fix
fix: fix-lint fix-license

.PHONY: test
test: clean
	$(GO_TEST) ./... -coverprofile=coverage.txt -covermode=atomic

.PHONY: verify
verify: clean license lint test

.PHONY: coverage
coverage: test
	bash <(curl -s https://codecov.io/bash) -t a5af28a3-92a2-4b35-9a77-54ad99b1ae00

.PHONY: clean
clean:
	-rm -rf bin
	-rm -rf coverage.txt
	-rm -rf *.tgz
	-rm -rf *.tgz
	-rm -rf *.asc
	-rm -rf *.sha512
	@go mod tidy &> /dev/null

# The source tarball is packed into a temporary directory and moved here: GNU tar exits 1 ("file changed as we read it")
# when the archive it writes is inside the directory being packed. Packing errors must fail the release.
release-src: clean
	tmp=$$(mktemp -d) && trap 'rm -rf "$$tmp"' EXIT && \
	$(TAR) -zcvf "$$tmp/$(RELEASE_SRC).tgz" \
	--exclude bin \
	--exclude .git \
	--exclude .idea \
	--exclude .DS_Store \
	--exclude dist \
	--exclude coverage.txt \
	--exclude $(RELEASE_SRC).tgz \
	--exclude $(RELEASE_BIN) \
	--exclude $(RELEASE_BIN).tgz \
	. && \
	mv "$$tmp/$(RELEASE_SRC).tgz" .

release-bin: build
	rm -rf $(RELEASE_BIN)
	mkdir $(RELEASE_BIN)
	cp -R bin $(RELEASE_BIN)
	cp -R dist/* $(RELEASE_BIN)
	cp -R CHANGES.md $(RELEASE_BIN)
	cp -R README.md $(RELEASE_BIN)
	$(TAR) -zcvf $(RELEASE_BIN).tgz $(RELEASE_BIN)
	rm -rf $(RELEASE_BIN)

# macOS bsdtar folds AppleDouble `._*` entries into the next entry when listing, so `-t` alone would hide them.
# `--options !mac-ext` lists every entry as it is. GNU tar has no such option and lists them anyway.
.PHONY: release-check
release-check:
	@list_opts=""; \
	if tar --version 2>/dev/null | grep -q bsdtar; then list_opts="--options !mac-ext"; fi; \
	for pkg in $(RELEASE_SRC).tgz $(RELEASE_BIN).tgz; do \
		if [ ! -f "$$pkg" ]; then echo "$$pkg was not created"; exit 1; fi; \
		entries=$$(tar $$list_opts -tzf "$$pkg") || { echo "$$pkg is not a readable tar.gz"; exit 1; }; \
		if echo "$$entries" | grep -E '(^|/)(\._|(\.DS_Store|__MACOSX)(/|$$))'; then \
			echo "$$pkg contains the macOS metadata files listed above"; \
			exit 1; \
		fi; \
	done

# release-check runs as a recipe step, not a prerequisite, so it always runs after both packages exist, even with -j.
# `-f` keeps the sub-make on this Makefile when it was given with `make -f`, since a sub-make does not inherit `-f`.
release: verify release-src release-bin
	@$(MAKE) -f $(firstword $(MAKEFILE_LIST)) --no-print-directory release-check
	gpg --batch --yes --armor --detach-sig $(RELEASE_SRC).tgz
	shasum -a 512 $(RELEASE_SRC).tgz > $(RELEASE_SRC).tgz.sha512
	gpg --batch --yes --armor --detach-sig $(RELEASE_BIN).tgz
	shasum -a 512 $(RELEASE_BIN).tgz > $(RELEASE_BIN).tgz.sha512

## Check that the status is consistent with CI.
check-codegen:
	$(MAKE) clean
	mkdir -p /tmp/swctl
	@go mod tidy &> /dev/null
	git diff >/tmp/swctl/check.diff 2>&1
	@if [ ! -z "`git status -s`" ]; then \
		echo "Following files are not consistent with CI:"; \
		git status -s; \
		exit 1; \
	fi

.PHONY: docker
docker: PUSH_OR_LOAD = --load
docker: PLATFORMS =

.PHONY: docker.push
docker.push: PUSH_OR_LOAD = --push
docker.push: PLATFORMS = --platform linux/386,linux/amd64,linux/arm64

docker docker.push:
	docker buildx create --use --driver docker-container --name skywalking_cli > /dev/null 2>&1 || true
	docker buildx build $(PUSH_OR_LOAD) $(PLATFORMS) --build-arg VERSION=$(VERSION) . -t $(HUB)/$(APP_NAME):$(VERSION) -t $(HUB)/$(APP_NAME):latest
	docker buildx rm skywalking_cli

.PHONY: install
install: clean
	$(BUILD_RULE)
	-cp $(OUT_DIR)/$(BINARY)-$(VERSION)-$(OSNAME)-* $(DESTDIR)/swctl

.PHONY: uninstall
uninstall: $(OSNAME)
	-rm $(DESTDIR)/$(PROJECT)/swctl
