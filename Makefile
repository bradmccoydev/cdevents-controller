# Makefile for releasing cdevents-controller
#
# The release version is controlled from pkg/version

TAG?=latest
NAME:=cdevents-controller
DOCKER_REPOSITORY:=bradmccoydev
DOCKER_IMAGE_NAME:=$(DOCKER_REPOSITORY)/$(NAME)
GIT_COMMIT:=$(shell git describe --dirty --always)
VERSION:=$(shell grep 'VERSION' pkg/version/version.go | awk '{ print $$4 }' | tr -d '"')
EXTRA_RUN_ARGS?=

run:
	go run -ldflags "-s -w -X github.com/bradmccoydev/cdevents-controller/pkg/version.REVISION=$(GIT_COMMIT)" cmd/cdevents-controller/* \
	--level=debug --grpc-port=9999 --backend-url=https://httpbin.org/status/401 --backend-url=https://httpbin.org/status/500 \
	--ui-logo=https://raw.githubusercontent.com/bradmccoydev/cdevents-controller/gh-pages/cuddle_clap.gif $(EXTRA_RUN_ARGS)

.PHONY: test
test:
	go test ./... -coverprofile cover.out

build:
	GIT_COMMIT=$$(git rev-list -1 HEAD) && CGO_ENABLED=0 go build  -ldflags "-s -w -X github.com/bradmccoydev/cdevents-controller/pkg/version.REVISION=$(GIT_COMMIT)" -a -o ./bin/cdevents-controller ./cmd/cdevents-controller/*
	GIT_COMMIT=$$(git rev-list -1 HEAD) && CGO_ENABLED=0 go build  -ldflags "-s -w -X github.com/bradmccoydev/cdevents-controller/pkg/version.REVISION=$(GIT_COMMIT)" -a -o ./bin/cdeventscli ./cmd/cdeventscli/*

tidy:
	rm -f go.sum; go mod tidy -compat=1.19

vet:
	go vet ./...

fmt:
	gofmt -l -s -w ./
	goimports -l -w ./

build-charts:
	helm lint charts/*
	helm package charts/*

build-container:
	docker build -t $(DOCKER_IMAGE_NAME):$(VERSION) .

build-xx:
	docker buildx build \
	--platform=linux/amd64 \
	-t $(DOCKER_IMAGE_NAME):$(VERSION) \
	--load \
	-f Dockerfile.xx .

build-base:
	docker build -f Dockerfile.base -t $(DOCKER_REPOSITORY)/cdevents-controller-base:latest .

push-base: build-base
	docker push $(DOCKER_REPOSITORY)/cdevents-controller-base:latest

test-container:
	@docker rm -f cdevents-controller || true
	@docker run -dp 9898:9898 --name=cdevents-controller $(DOCKER_IMAGE_NAME):$(VERSION)
	@docker ps
	@TOKEN=$$(curl -sd 'test' localhost:9898/token | jq -r .token) && \
	curl -sH "Authorization: Bearer $${TOKEN}" localhost:9898/token/validate | grep test

push-container:
	docker tag $(DOCKER_IMAGE_NAME):$(VERSION) $(DOCKER_IMAGE_NAME):latest
	docker push $(DOCKER_IMAGE_NAME):$(VERSION)
	docker push $(DOCKER_IMAGE_NAME):latest

version-set:
	@next="$(TAG)" && \
	current="$(VERSION)" && \
	/usr/bin/sed -i '' "s/$$current/$$next/g" pkg/version/version.go && \
	/usr/bin/sed -i '' "s/tag: $$current/tag: $$next/g" charts/cdevents-controller/values.yaml && \
	/usr/bin/sed -i '' "s/tag: $$current/tag: $$next/g" charts/cdevents-controller/values-prod.yaml && \
	/usr/bin/sed -i '' "s/appVersion: $$current/appVersion: $$next/g" charts/cdevents-controller/Chart.yaml && \
	/usr/bin/sed -i '' "s/version: $$current/version: $$next/g" charts/cdevents-controller/Chart.yaml && \
	/usr/bin/sed -i '' "s/cdevents-controller:$$current/cdevents-controller:$$next/g" kustomize/deployment.yaml && \
	/usr/bin/sed -i '' "s/cdevents-controller:$$current/cdevents-controller:$$next/g" deploy/webapp/frontend/deployment.yaml && \
	/usr/bin/sed -i '' "s/cdevents-controller:$$current/cdevents-controller:$$next/g" deploy/webapp/backend/deployment.yaml && \
	/usr/bin/sed -i '' "s/cdevents-controller:$$current/cdevents-controller:$$next/g" deploy/bases/frontend/deployment.yaml && \
	/usr/bin/sed -i '' "s/cdevents-controller:$$current/cdevents-controller:$$next/g" deploy/bases/backend/deployment.yaml && \
	/usr/bin/sed -i '' "s/$$current/$$next/g" cue/main.cue && \
	echo "Version $$next set in code, deployment, chart and kustomize"

release:
	git tag $(VERSION)
	git push origin $(VERSION)

swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	go get github.com/swaggo/swag/gen@latest
	go get github.com/swaggo/swag/cmd/swag@latest
	cd pkg/api && $$(go env GOPATH)/bin/swag init -g server.go

.PHONY: cue-mod
cue-mod:
	@cd cue && go mod init github.com/bradmccoydev/cdevents-controller/cue
	@cd cue && go get k8s.io/api/...
	@cd cue && cue get go k8s.io/api/...

.PHONY: cue-gen
cue-gen:
	@cd cue && cue fmt ./... && cue vet --all-errors --concrete ./...
	@cd cue && cue gen
