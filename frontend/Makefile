VERSION ?= latest
IMG ?= commonpool/frontend:$(VERSION)

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

.PHONY: docker-build
docker-build:
	docker build . -t $(IMG)

.PHONY: docker-push
docker-push:
	docker push $(IMG)
