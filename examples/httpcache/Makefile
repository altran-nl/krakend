.PHONY: all deps build

# This Makefile is a simple example that demonstrates usual steps to build a binary that can be run in the same
# architecture that was compiled in. The "ldflags" in the build assure that any needed dependency is included in the
# binary and no external dependencies are needed to run the service.

KRAKEND_VERSION=$(shell git describe --always --long --dirty --tags)
BIN_NAME=krakend_httpcache_example_${KRAKEND_VERSION}

all: deps build

deps:
	go get "github.com/altran-nl/krakend/config/viper"
	go get "github.com/altran-nl/krakend/proxy"
	go get "github.com/altran-nl/krakend/router/gin"
	go get "github.com/altran-nl/krakend/logging/gologging"
	go get "github.com/gin-gonic/contrib/cache"
	go get "github.com/gregjones/httpcache"

build:
	go build -a -ldflags="-X github.com/altran-nl/krakend/core.KrakendVersion=${KRAKEND_VERSION}" -o ${BIN_NAME}
	@echo "You can now use ./${BIN_NAME}"
