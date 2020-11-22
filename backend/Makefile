VERSION ?= latest
IMG ?= commonpool/backend:$(VERSION)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

CONTAINER ?= $(shell command -v docker || command -v podman || exit 1)
SHELL=bash

# define standard colors
BLACK        := $(shell tput -Txterm setaf 0)
RED          := $(shell tput -Txterm setaf 1)
GREEN        := $(shell tput -Txterm setaf 2)
YELLOW       := $(shell tput -Txterm setaf 3)
LIGHTPURPLE  := $(shell tput -Txterm setaf 4)
PURPLE       := $(shell tput -Txterm setaf 5)
BLUE         := $(shell tput -Txterm setaf 6)
WHITE        := $(shell tput -Txterm setaf 7)
RESET := $(shell tput -Txterm sgr0)


.PHONY: test
test: fmt vet
	go test ./...

.PHONY: run
run: fmt vet
	go run main.go

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet: 
	go vet ./...

.PHONY: pack
pack: test
	$(CONTAINER) build . -t $(IMG)

.PHONY: push
push: 
	$(CONTAINER) push $(IMG)

.PHONY: up
up:
	@whoami
	@echo "> > > ${BLUE}Setting up commonpool backend environment...${RESET}"
	@if [ $$(command -v docker) ] && [ $$(command -v docker-compose) ] ; then \
		docker-compose up; \
	else \
		if [ "$$(podman pod ps | grep commonpool-dev | wc -l)" == "0" ] ; then \
			echo "> > > ${BLUE}Starting PostgreSQL${RESET}"; \
			podman play kube podman.yaml; \
		else \
			read -p "Development pod is already running. Re-create it? [yN]" yn; \
			case $$yn in \
				[Yy]* )podman pod rm commonpool-dev -f; podman play kube podman.yaml;; \
				""|[Nn]* ) exit;; \
				* ) echo "Please answer yes or no.";; \
			esac \
		fi \
	fi

.PHONY: down
down:
	@if [ $$(command -v docker) ] ; then \
		docker-compose down; \
	else \
		if [ "$$(podman pod ps | grep commonpool-dev | wc -l)" == "1" ] ; then \
			echo "> > > ${BLUE}Stopping PostgreSQL${RESET} "; \
			podman pod rm commonpool-dev -f; \
		fi \
	fi

.PHONY: .env-export
.env-export:
	sed -ne '/^export / {p;d}; /.*=/ s/^/export / p' .env > .env-export

.PHONY: start
start: up .env-export
	. .env-export && go run main.go