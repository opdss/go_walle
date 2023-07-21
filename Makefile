GO_VERSION ?= 1.18.8
GOOS ?= linux
GOARCH ?= amd64
GOPATH ?= $(shell go env GOPATH)
NODE_VERSION ?= 16.11.1
COMPOSE_PROJECT_NAME := ${TAG}-$(shell git rev-parse --abbrev-ref HEAD)
BRANCH_NAME ?= $(shell git rev-parse --abbrev-ref HEAD | sed "s!/!-!g")
GIT_TAG := $(shell git rev-parse --short HEAD)
ifeq (${BRANCH_NAME},main)
TAG    := ${GIT_TAG}-go${GO_VERSION}
TRACKED_BRANCH := true
LATEST_TAG := latest
else
TAG    := ${GIT_TAG}-${BRANCH_NAME}-go${GO_VERSION}
ifneq (,$(findstring release-,$(BRANCH_NAME)))
TRACKED_BRANCH := true
LATEST_TAG := ${BRANCH_NAME}-latest
endif
endif
CUSTOMTAG ?=

FILEEXT :=
ifeq (${GOOS},windows)
FILEEXT := .exe
endif

.DEFAULT_GOAL := help
.PHONY: help
help:
	@awk 'BEGIN { \
		FS = ":.*##"; \
		printf "\nUsage:\n  make \033[36m<target>\033[0m\n"\
	} \
	/^[a-zA-Z_-]+:.*?##/ { \
		printf "  \033[36m%-17s\033[0m %s\n", $$1, $$2 \
	} \
	/^##@/ { \
		printf "\n\033[1m%s\033[0m\n", substr($$0, 5) \
	} ' $(MAKEFILE_LIST)

##@ Dependencies

.PHONY: build-dev-deps
build-dev-deps: ## 安装开发工具
	go install github.com/cosmtrek/air@latest
	go install golang.org/x/tools/cmd/stringer


.PHONY: build-web-deps
build-web-deps: ## 安装web依赖包
	cd web && npm ci --registry https://registry.npm.taobao.org

##@ Build

.PHONY: web
web: ## 编译构建web
	cd web && npm ci --registry https://registry.npm.taobao.org && npm run build

.PHONY: backend
backend: ## 编译构建api
	@echo "Running ${@}"
	./scripts/release.sh build main.go

.PHONY: install
install: web backend ## 编译构建本项目,先构建web，再构建api


