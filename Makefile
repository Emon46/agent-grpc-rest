
REGISTRY ?= observoinc
# Image URL to use all building/pushing image targets
TAG      ?=  latest
IMG ?= $(REGISTRY)/control-agent:$(TAG)
# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.25.0


PLATFORM ?= linux_amd64
VERSION          = $(TAG)_$(subst /,_,$(PLATFORM))
# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: fmt vet ## Run tests.
	go test  -cover -race ./... -coverprofile cover.out
##@ Build

.PHONY: build
build: fmt vet ## Build manager binary.
	CGO_ENABLED=0 GOOS=linux go build -o bin/control-agent cmd/server/main.go

.PHONY: run
run: fmt vet ## Run a controller from your host.
	go run server/main.go

# If you wish built the manager image targeting other platforms you can use the --platform flag.
# (i.e. docker build --platform linux/arm64 ). However, you must enable docker buildKit for it.
# More info: https://docs.docker.com/develop/develop-images/build_enhancements/
.PHONY: docker-build
docker-build: build test ## Build docker image with the manager.
	docker build -t ${IMG} .

.PHONY: docker-push
docker-push: docker-build ## Push docker image with the manager.
	docker push ${IMG}


.PHONY: push-to-kind
push-to-kind: build docker-build
	@echo "Loading docker image into kind cluster...."
	@kind load docker-image $(IMG)
	@echo "Image has been pushed successfully into kind cluster."

.PHONY: clean
clean:
	rm -rf .go bin
	rm pb/*.go

gen:
	protoc --proto_path=proto proto/*.proto  --go_out=:pb --go-grpc_out=:pb

client:
	go run client/main.go -address 0.0.0.0:8081
server:
	go run cmd/server/main.go -port 8080