GO11MODULES=on
APP?=application
REGISTRY?=gcr.io/images
COMMIT_SHA=$(shell git rev-parse --short HEAD)



.PHONY: build
## build: build the application
build: clean
    @echo "Building..."
    @go build -o ${APP} main.go

.PHONY: run
## run: runs go run main.go
run:
	go run -race main.go

.PHONY: clean
## clean: cleans the binary
clean:
    @echo "Cleaning"
    @rm -rf ${APP}
    @go get -u -v golang.org/x/tools/cmd/goimports
    @gofmt -l -w -s .
    @goimports -l -w .

.PHONY: test
## test: runs go test with default values
test:
	go test -v -count=1 -race ./...

.PHONY: lint
## lint: lint the application code for problems
lint:
	@golangci-lint run

.PHONY: setup
## setup: setup go modules
setup:
	@go mod init \
		&& go mod tidy \
		&& go mod vendor

# helper rule for deployment
check-environment:
ifndef ENV
    $(error ENV not set, allowed values - `staging` or `production`)
endif

.PHONY: docker-build
## docker-build: builds the stringifier docker image to registry
docker-push: build
	docker build -t ${APP}:${COMMIT_SHA} .

.PHONY: docker-push
## docker-push: pushes the stringifier docker image to registry
docker-push: check-environment docker-build
	docker push ${REGISTRY}/${ENV}/${APP}:${COMMIT_SHA}

.PHONY: help
## help: Prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
