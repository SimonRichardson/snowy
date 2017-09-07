PATH_SNOWY = github.com/trussle/snowy

.PHONY: all
all:
	go get github.com/Masterminds/glide
	glide install
	$(MAKE) clean build

.PHONY: build
build: dist/documents

dist/documents:
	go build -o dist/documents ${PATH_SNOWY}/cmd/documents

pkg/store/mocks/store.go:
	mockgen -package=mocks -destination=pkg/store/mocks/store.go ${PATH_SNOWY}/pkg/store Store

pkg/repository/mocks/repository.go:
	mockgen -package=mocks -destination=pkg/repository/mocks/repository.go ${PATH_SNOWY}/pkg/repository Repository

pkg/metrics/mocks/metrics.go:
	mockgen -package=mocks -destination=pkg/metrics/mocks/metrics.go ${PATH_SNOWY}/pkg/metrics Gauge,HistogramVec
	sed -i '' -- 's/github.com\/trussle\/snowy\/vendor\///g' ./pkg/metrics/mocks/metrics.go

pkg/metrics/mocks/observer.go:
	mockgen -package=mocks -destination=pkg/metrics/mocks/observer.go github.com/prometheus/client_golang/prometheus Observer

.PHONY: build-mocks
build-mocks: FORCE
	$(MAKE) pkg/store/mocks/store.go
	$(MAKE) pkg/repository/mocks/repository.go
	$(MAKE) pkg/metrics/mocks/metrics.go
	$(MAKE) pkg/metrics/mocks/observer.go

.PHONY: clean
clean: FORCE
	rm -f dist/documents

.PHONY: clean-mocks
clean-mocks: FORCE
	rm -f pkg/store/mocks/store.go
	rm -f pkg/repository/mocks/repository.go
	rm -f pkg/metrics/mocks/metrics.go
	rm -f pkg/metrics/mocks/observer.go

FORCE:

.PHONY: integration-tests
integration-tests:
	docker-compose run documents go test -v -tags=integration ./cmd/... ./pkg/...

.PHONY: documentation
documentation:
	go test -v -tags=documentation ./pkg/... -run=TestDocumentation_

.PHONY: coverage-tests
coverage-tests:
	docker-compose run documents go test -covermode=count -coverprofile=bin/coverage.out -v -tags=integration ${COVER_PKG}

.PHONY: coverage-view
coverage-view:
	go tool cover -html=bin/coverage.out

PWD ?= ${GOPATH}/src/${PATH_SNOWY}
TAG ?= dev
BRANCH ?= dev
ifeq ($(BRANCH),master)
	TAG=latest
endif

.PHONY: build-docker
build-docker:
	@echo "Building '${TAG}' for '${BRANCH}'"
	docker run --rm -v ${PWD}:/go/src/${PATH_SNOWY} -w /go/src/${PATH_SNOWY} iron/go:dev go build -o documents ${PATH_SNOWY}/cmd/documents
	docker build -t teamtrussle/snowy:${TAG} .

.PHONY: push-docker-tag
push-docker-tag: FORCE
	@echo "Pushing '${TAG}' for '${BRANCH}'"
	docker login -u ${DOCKER_HUB_USERNAME} -p ${DOCKER_HUB_PASSWORD}
	docker push teamtrussle/snowy:${TAG}

.PHONY: push-docker
ifeq ($(TAG),latest)
push-docker: FORCE
	@echo "Pushing '${TAG}' for '${BRANCH}'"
	docker login -u ${DOCKER_HUB_USERNAME} -p ${DOCKER_HUB_PASSWORD}
	docker push teamtrussle/snowy:${TAG}
else
push-docker: FORCE
	@echo "Pushing requires branch '${BRANCH}' to be master"
endif
