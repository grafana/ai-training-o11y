# When these variables are set in the Dockerfile they will not be overwritten
BUILD_VERSION ?= $(shell git describe --tags --long --always --abbrev=10)
BUILD_COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# These env variables should be in sync with the ones exported in the Dockerfile
GO_OPT= -ldflags "-X main.Branch=$(BUILD_BRANCH) -X main.Revision=$(BUILD_COMMIT) -X main.Version=$(BUILD_VERSION)"

.PHONY: exe
exe:
	GO111MODULE=on go build $(GO_OPT) -o ./dist $(BUILD_INFO) ./...
