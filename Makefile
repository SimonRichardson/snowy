.PHONY: all
all:
	go get github.com/Masterminds/glide
	glide install
	$(MAKE) clean build

.PHONY: build
build: dist/documents

dist/documents:
	go build -o dist/documents github.com/trussle/snowy/cmd/documents

pkg/store/mocks/store.go:
	mockgen -package=mocks -destination=pkg/store/mocks/store.go github.com/trussle/snowy/pkg/store Store

pkg/repository/mocks/repository.go:
	mockgen -package=mocks -destination=pkg/repository/mocks/repository.go github.com/trussle/snowy/pkg/repository Repository

pkg/metrics/mocks/metrics.go:
	mockgen -package=mocks -destination=pkg/metrics/mocks/metrics.go github.com/trussle/snowy/pkg/metrics Gauge,HistogramVec
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
