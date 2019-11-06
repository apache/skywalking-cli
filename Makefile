# Example:
#   make build
#   make install

.PHONY: build install

build:
	GO_BUILD_FLAGS="-v" ./scripts/build
	./bin/swctl --version