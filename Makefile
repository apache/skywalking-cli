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

HUB ?= docker.io/apache

RELEASE_BIN = skywalking-cli-$(VERSION)-bin
RELEASE_SRC = skywalking-cli-$(VERSION)-src

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
	@$(GO_GET) -v -t -d ./...

$(GO_LINT):
	@$(GO_LINT) version > /dev/null 2>&1 || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
$(LICENSE_EYE):
	@$(LICENSE_EYE) --version > /dev/null 2>&1 || go install github.com/apache/skywalking-eyes/cmd/license-eye@d38fe05

.PHONY: lint
lint: $(GO_LINT)
	$(GO_LINT) run -v --timeout 5m ./...
.PHONY: fix-lint
fix-lint: $(GO_LINT)
	$(GO_LINT) run -v --fix ./...

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

release-src: clean
	-tar -zcvf $(RELEASE_SRC).tgz \
	--exclude bin \
	--exclude .git \
	--exclude .idea \
	--exclude .DS_Store \
	--exclude dist \
	--exclude $(RELEASE_SRC).tgz \
	.

release-bin: build
	-mkdir $(RELEASE_BIN)
	-cp -R bin $(RELEASE_BIN)
	-cp -R dist/* $(RELEASE_BIN)
	-cp -R CHANGES.md $(RELEASE_BIN)
	-cp -R README.md $(RELEASE_BIN)
	-tar -zcvf $(RELEASE_BIN).tgz $(RELEASE_BIN)
	-rm -rf $(RELEASE_BIN)

release: verify release-src release-bin
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
