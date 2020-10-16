VERSION ?= latest
IMG ?= commonpool/frontend:$(VERSION)

start: 
	npm run start

build:
	npm run build

build-prod:
	npm run build:prod

test:
	npm run test

lint:
	npm run lint

e2e:
	npm run e2e

e2e-watch:
	npm run e2e:watch

docker-build:
	docker build . -t $(IMG)

docker-push: 
	docker push $(IMG)