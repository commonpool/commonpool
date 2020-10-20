VERSION ?= latest
IMG ?= commonpool/backend:$(VERSION)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

SHELL=bash

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

.PHONY: docker-build
docker-build: test
	docker build . -t $(IMG)

.PHONY: docker-push
docker-push: 
	docker push $(IMG)

.PHONY: compose-up
compose-up:
	docker-compose up -d

.PHONY: compose-down
compose-down:
	docker-compose down

.env-export:
	sed -ne '/^export / {p;d}; /.*=/ s/^/export / p' .env > .env-export

.PHONY: start
start: compose-up .env-export
	. .env-export && go run main.go