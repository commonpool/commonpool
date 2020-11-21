VERSION ?= latest
IMG ?= commonpool/frontend:$(VERSION)

CONTAINER=$(shell command -v docker || command -v podman || exit 1)

BUILD_CONTAINER=$(shell (command -v docker && echo build . -t $(IMG)) || (command -v podman && echo build . -t $(IMG)) || (command -v buildah && echo bud -t $(IMG) .))
PUSH_CONTAINER=$(shell (command -v docker && echo push $(IMG)) || (command -v podman && echo push $(IMG)) || (command -v buildah && echo push $(IMG)))

.PHONY: start
start:
	npm run start

.PHONY: build
build:
	npm run build

.PHONY: build-prot
build-prod:
	npm run build:prod

.PHONY: test
test:
	npm run test

.PHONY: lint
lint:
	npm run lint

.PHONY: e2e
e2e:
	npm run e2e

.PHONY: e2e-watch
e2e-watch:
	npm run e2e:watch

.PHONY: pack
pack:
	$(BUILD_CONTAINER)

.PHONY: push
push:
	$(PUSH_CONTAINER)
