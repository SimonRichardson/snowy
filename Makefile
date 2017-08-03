.PHONY: all
all:
	go get github.com/Masterminds/glide
	glide install
	$(MAKE) clean build

.PHONY: build
build: dist/documents

dist/documents:
	go build -o dist/documents github.com/trussle/snowy/cmd/documents

.PHONY: clean
clean: FORCE
	rm -rf dist/documents

FORCE:

.PHONY: integration-tests
integration-tests:
	docker-compose run documents go test -v -tags=integration ./cmd/... ./pkg/...
