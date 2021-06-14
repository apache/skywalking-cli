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

OS = $(shell uname)

GO = go
GO_PATH = $$($(GO) env GOPATH)
GO_BUILD = $(GO) build
GO_GET = $(GO) get
GO_INSTALL = $(GO) install
GO_TEST = $(GO) test
GO_LINT = $(GO_PATH)/bin/golangci-lint
GO_LICENSER = $(GO_PATH)/bin/go-licenser
ARCH := $(shell uname)
OSNAME := $(if $(findstring Darwin,$(ARCH)),darwin,linux)
GOBINDATA_VERSION := v3.21.0
GO_BINDATA = $(GO_PATH)/bin/go-bindata
GO_BUILD_FLAGS = -v
GO_BUILD_LDFLAGS = -X main.version=$(VERSION)

PLATFORMS := windows linux darwin
os = $(word 1, $@)
ARCH = amd64

SHELL = /bin/bash

all: clean license deps lint test build

tools:
	mkdir -p $(GO_PATH)/bin
	$(GO_BINDATA) -v || curl --location --output $(GO_BINDATA) https://github.com/kevinburke/go-bindata/releases/download/$(GOBINDATA_VERSION)/go-bindata-$(OSNAME)-amd64 \
		&& chmod +x $(GO_BINDATA)
	$(GO_LINT) version || curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GO_PATH)/bin
	$(GO_LICENSER) -version || GO111MODULE=off $(GO_GET) -u github.com/elastic/go-licenser

deps: tools
	$(GO_GET) -v -t -d ./...

.PHONY: assets
assets: tools
	cd assets \
		&& $(GO_BINDATA) --nocompress --nometadata --pkg assets --ignore '.*\.go' \
			-o "assets.gen.go" ./... \
		&& ../scripts/build-header.sh assets.gen.go \
		&& cd ..

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p $(OUT_DIR)
	GOOS=$(os) GOARCH=$(ARCH) $(GO_BUILD) $(GO_BUILD_FLAGS) -ldflags "$(GO_BUILD_LDFLAGS)" -o $(OUT_DIR)/$(BINARY)-$(VERSION)-$(os)-$(ARCH) cmd/swctl/main.go

.PHONY: lint
lint: tools
	$(GO_LINT) run -v --timeout 5m ./...

.PHONY: test
test: clean lint
	$(GO_TEST) ./... -coverprofile=coverage.txt -covermode=atomic

.PHONY: build
build: deps windows linux darwin

.PHONY: license
license: clean tools
	$(GO_LICENSER) -d -licensor='Apache Software Foundation (ASF)' .

.PHONY: verify
verify: clean license lint test

.PHONY: fix
fix: tools
	$(GO_LINT) run -v --fix ./...
	$(GO_LICENSER) -licensor='Apache Software Foundation (ASF)' .

.PHONY: coverage
coverage: test
	bash <(curl -s https://codecov.io/bash) -t a5af28a3-92a2-4b35-9a77-54ad99b1ae00

.PHONY: clean
clean: tools
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
	--exclude .github \
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
docker:
	docker build . -t $(HUB)/$(APP_NAME):$(VERSION)

.PHONY: test-commands
test-commands:
	@if ! docker run --name oap -p 12800:12800 -p 11800:11800 -d -e SW_HEALTH_CHECKER=default -e SW_TELEMETRY=prometheus apache/skywalking-oap-server:8.4.0-es7 > /dev/null 2>&1;then \
		docker container stop oap; \
		docker container prune -f; \
		docker run --name oap -p 12800:12800 -p 11800:11800 -d -e SW_HEALTH_CHECKER=default -e SW_TELEMETRY=prometheus apache/skywalking-oap-server:8.4.0-es7; \
	fi
	VERSION=$(VERSION) ./scripts/test_commands.sh
	@docker container stop oap
	@docker container prune -f
