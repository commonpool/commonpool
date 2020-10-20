VERSION ?= latest
IMG ?= commonpool/backend:$(VERSION)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

SHELL=bash

test: fmt vet
	go test ./...

run: fmt vet
	go run main.go

fmt:
	go fmt ./...

vet: 
	go vet ./...

docker-build: test
	docker build . -t $(IMG)

docker-push: 
	docker push $(IMG)

compose-up:
	docker-compose up -d

compose-down:
	docker-compose down

.env-export:
	sed -ne '/^export / {p;d}; /.*=/ s/^/export / p' .env > .env-export

%: .env-export
	. .env-export && make -f secondary "$@"

start: compose-up .env-export
	. .env-export && go run main.go