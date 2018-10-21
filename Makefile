.DEFAULT_GOAL := help

# make configs taken from kubernetes
DBG_MAKEFILE ?=
ifeq ($(DBG_MAKEFILE),1)
    $(warning ***** starting Makefile for goal(s) "$(MAKECMDGOALS)")
    $(warning ***** $(shell date))
    $(warning ***** setting debug flags)
		DEBUG = true
else
    # If we're not debugging the Makefile, don't echo recipes.
    MAKEFLAGS += -s
		DEBUG = false
endif
# It's necessary to set this because some environments don't link sh -> bash.
SHELL := /bin/bash
# We don't need make's built-in rules.
MAKEFLAGS += --no-builtin-rules


# constants and envars
ELF_NAME = keybasectl
ELF_APPENVIRONMENT = dev
build_production: ELF_APPENVIRONMENT = production
ELF_VERSION ?= v0.0.1
OS = $(shell grep "^ID=" /etc/os-release | cut -d"=" -f2 )
BUILD_TIMEOUT = 120 # seconds
DEP_GRAPH_PATH := tmp/$(ELF_NAME)-$(strip $(ELF_VERSION)).deps.png
CONTAINER_IMAGENAME = 'stefancocora/$(ELF_NAME)'
CONTAINER_VERSION = $(ELF_VERSION)
CONTAINER_NAME := '$(ELF_NAME)'
CONTAINER_NOCACHE := 'nocache'
container_iterateimage : CONTAINER_NOCACHE = 'withcache'
CONTAINER_BUILD_ACTION := 'build'

# possible values: local/localRkt/localDocker/concourseCI
# each of those values tells the build script how to behave if inside/outside a build container
BUILD_ENV ?= local
OUTPUT_DIR = nothing
ifeq ($(BUILD_ENV), local)
	OUTPUT_DIR := bin/go
else ifeq ($(BUILD_ENV), localRkt)
	OUTPUT_DIR := bin/rkt
else ifeq ($(BUILD_ENV), localDocker)
	OUTPUT_DIR := bin/docker
else ifeq ($(BUILD_ENV), concourseCI)
	OUTPUT_DIR := bin/runc
endif


# https://github.com/River-Island/product-backbone/blob/master/Makefile
envars := ORDS_ENDPOINT=http://localhost:8100 \
	DB_SCHEMA=backbone \
	DB_USER=backbone \
	DB_PASS=mypassword \
	DB_HOST=localhost \
	DB_PORT=6432 \
	ENV=dev \
	VERSION=local-machine

.PHONY: build
build:					## Build the artifact using system go
ifeq ($(DEBUG),true)
	$(info elf appenvironment: $(ELF_APPENVIRONMENT))
	$(info elf version: $(ELF_VERSION))
	$(info debug: $(DEBUG))
endif

	@echo "--> Building ELF ..."
ifeq ($(OS),alpine)
	$(info detected OS: $(OS))
	timeout -t $(BUILD_TIMEOUT) util/build.sh build $(ELF_NAME) $(ELF_BUILD_ENV) $(ELF_APPENVIRONMENT) $(ELF_VERSION) $(OUTPUT_DIR) $(DEBUG) $(BUILD_ACTION)
else
	$(info detected OS: $(OS))
	timeout --preserve-status $(BUILD_TIMEOUT) util/build.sh build $(ELF_NAME) $(ELF_BUILD_ENV) $(ELF_APPENVIRONMENT) $(ELF_VERSION) $(OUTPUT_DIR) $(DEBUG) $(BUILD_ACTION)
endif

.PHONY: build_production
build_production:			## Build the artifact using system go
ifeq ($(DEBUG),true)
	$(info elf appenvironment: $(ELF_APPENVIRONMENT))
	$(info elf version: $(ELF_VERSION))
	$(info debug: $(DEBUG))
endif

	@echo "--> Building ELF ..."
ifeq ($(OS),alpine)
	$(info detected OS: $(OS))
	timeout -t $(BUILD_TIMEOUT) util/build.sh build $(ELF_NAME) $(ELF_BUILD_ENV) $(ELF_APPENVIRONMENT) $(ELF_VERSION) $(OUTPUT_DIR) $(DEBUG) $(BUILD_ACTION)
else
	$(info detected OS: $(OS))
	timeout --preserve-status $(BUILD_TIMEOUT) util/build.sh build $(ELF_NAME) $(ELF_BUILD_ENV) $(ELF_APPENVIRONMENT) $(ELF_VERSION) $(OUTPUT_DIR) $(DEBUG) $(BUILD_ACTION)
endif

.PHONY: container_image
container_image:			## container:docker: Build a container image without using docker cache
	@echo "--> Building container image without caches..."
ifeq ($(DEBUG),true)
	$(info version: $(CONTAINER_BUILD_ACTION))
	$(info version: $(CONTAINER_VERSION))
	$(info nocache: $(CONTAINER_NOCACHE))
	$(info imagename: $(CONTAINER_IMAGENAME))
	$(info debug: $(DEBUG))
endif
	timeout --preserve-status 120s util/buildcontainer.sh $(CONTAINER_BUILD_ACTION) $(CONTAINER_VERSION) $(CONTAINER_NOCACHE) $(CONTAINER_IMAGENAME) $(DEBUG)

.PHONY: container_iterateimage
container_iterateimage:			## container:docker: Build a container image using docker cache
	@echo "--> Building container image ..."
ifeq ($(DEBUG),true)
	$(info containerbuildaction: $(CONTAINER_BUILD_ACTION))
	$(info version: $(CONTAINER_VERSION))
	$(info nocache: $(CONTAINER_NOCACHE))
	$(info imagename: $(CONTAINER_IMAGENAME))
	$(info debug: $(DEBUG))
endif
	timeout --preserve-status 120s util/buildcontainer.sh $(CONTAINER_BUILD_ACTION) $(CONTAINER_VERSION) $(CONTAINER_NOCACHE) $(CONTAINER_IMAGENAME) $(DEBUG)

.PHONY: deps
deps:					## Update the package dependencies
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

.PHONY: dep_graph
dep_graph:				## Generate a dep graph of build dependencies
	@echo "generating dep dependency graph ..."
ifeq ($(DEBUG),true)
	$(info $(DEP_GRAPH_PATH))
endif
	dep status -dot | dot -Tpng > $(DEP_GRAPH_PATH)
	eog $(DEP_GRAPH_PATH) &

.PHONY: run_docker_arch
run_docker_arch:			## Run ELF in a docker container
	$(info run the elf file inside a docker container === Nothing to see here yet)

.PHONY: run_rkt_arch
run_rkt_arch:				## Run ELF in a rkt container
	$(info run the elf file inside a rkt container === Nothing to see here yet)

.PHONY: run_rkt_gentoo
run_rkt_gentoo:				## Run ELF in a rkt container
	$(info run the elf file inside a rkt container === Nothing to see here yet)

.PHONY: gometalinter
gometalinter:				## golang: Run gometalinter against the src code
# $(shell command -v gometalinter || $(go get -u github.com/alecthomas/gometalinter && 	gometalinter --install >/dev/null))
	@echo "running gometalinter ..."
	gometalinter ./... --vendor --skip=vendor --exclude=".*_mock.*.go" --exclude="vendor.*" --cyclo-over=15 --deadline=2m --disable-all \
		--enable=errcheck \
		--enable=vet \
		--enable=deadcode \
		--enable=gocyclo \
		--enable=golint \
		--enable=varcheck \
		--enable=structcheck \
		--enable=maligned \
		--enable=vetshadow \
		--enable=ineffassign \
		--enable=interfacer \
		--enable=unconvert \
		--enable=goconst \
		--enable=gosimple \
		--enable=staticcheck \
		--enable=gas

# has to be the last make target
# all other targets should have a comment assigned
.PHONY: help
help:					## Show this help
	$(info -------------------------------------------------------------------)
	$(info $(ELF_NAME))
	$(info -------------------------------------------------------------------)
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
